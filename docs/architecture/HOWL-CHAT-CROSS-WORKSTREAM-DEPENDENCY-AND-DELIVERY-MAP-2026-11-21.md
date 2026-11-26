# Howl Chat — Cross-Workstream Dependency and Delivery Map

## Program map `HCDM-2026-11-21-01`

| Field | Value |
|---|---|
| Date | **2026-11-21** (local, CST) |
| Type | Documentation and planning only |
| Starting baseline | **`HFS-AUDIT-2026-11-21-01`** |
| Related records | `HFR-CP-2026-11-21-01`; `HF1-PLAN-2026-11-21-01`; `HOWL-FOUNDATION-PROGRAM-DIRECTION-2026-11-21.md`; `HOWL-FOUNDATION-ORIGINAL-SCOPE-AND-EVIDENCE-CHAIN-2026-11-21.md`; `HOWL-CHAT-M0A-RELEASE-TRUTH-ACCEPTANCE-REVIEW-2026-11-21.md`; `HOWL-RESOURCE-ARBITRATION-AND-VERIFIED-HANDOFF-DISCOVERY-BASELINE-2026-11-21.md` |
| Production changes | **0** |
| Code/schema/API/migration changes | **0** |
| Tests executed by this map | **0** |
| Models / corpus / benchmarks | **0** |
| MEMORY-09 implementation/integration | **0** |
| Authorization | **None**; this map recommends, it does not authorize |
| Owner scope (corrected 2026-11-21) | Audio generation **complete**; TTS **deferred, required for chat finalization**; image generation **nearly complete, remaining work in scope**; vision **largely complete, gaps to assess**; video **incomplete, in scope**; audio capture/input separate; 3D mesh generation **removed**; 2D/3D visual-character rendering **removed** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** · **Superseded**.

**Scope correction (2026-11-21).** The modality ledger in `HRA-DB-2026-11-21-01` supersedes the older deferral/removal statements still present in earlier drafts of this map. Image generation is in scope; audio generation is complete; TTS is deferred but required for chat finalization; vision gaps are to be assessed; video is in scope; audio capture/input is its own capability. Historical statements are retained only as superseded history, not as current scheduling authority.

---

## 0. Purpose, source control, and limits

### 0.1 What this map is

**HCDM-01 (Observed — scope).** This map converts the recovered original scope and evidence into a **delivery program**. It is not another readiness audit, not a Phase 1 authorization, and not a MEMORY-09 plan.

It answers:

1. What Howl still needs.
2. What actually blocks release.
3. What can proceed independently.
4. What must be coordinated first.
5. What useful functionality can ship next.
6. What is required before MEMORY-09 can be tested in real Howl usage.

### 0.2 Sources and authority

**HCDM-02 (Source-derived).**

- `HFS-AUDIT-2026-11-21-01` is the scope-recovery baseline. It established the 27-item union: V1 has 25, V2 has 26, and `T0.5` is recovered from V1 because V2 dropped it without a closure.
- `PROJECT_ISSUES_AND_EXPANSION.md` remains the authoritative original owner-facing work ledger for T0–T3 and the F2/F3/W appendices.
- `HFR-CP-2026-11-21-01` remains the evidence snapshot. Its A–E categories are not the work list.
- `HF1-PLAN-2026-11-21-01` remains a **proposed** work package. This map reviews it; it does not authorize it.
- The 2026-11-21 program direction remains the program-level direction record.
- `HRA-DB-2026-11-21-01` is the current discovery baseline for resource arbitration, verified handoff and the corrected modality ledger.
- The owner’s current scope decisions supersede older roadmap entries where they conflict.

**HCDM-03 (Observed — no silent scope expansion).** This map itself creates no new product requirement. `HRA-DB-2026-11-21-01` records the owner-approved resource-arbitration requirement separately. Technical tasks marked **D-** are derived subtasks of an existing original item.

### 0.3 Owner scope controls applied

**HCDM-04 (Owner-approved — corrected direction, 2026-11-21).**

| Scope | Owner action | Planning consequence |
|---|---|---|
| Audio capture/input | **IN SCOPE — separate capability** | Own pipeline, lifecycle, resource envelope and failure modes. Not audio generation and not TTS. |
| Audio generation | **COMPLETE** | No remaining delivery claim here. Completion is an owner statement; link its evidence separately. |
| TTS | **DEFERRED, REQUIRED FOR CHAT FINALIZATION** | Not current scheduling; chat cannot be declared final without it. Re-entry needs a dated owner decision and its own work package. |
| Image generation | **NEARLY COMPLETE — REMAINING WORK IN SCOPE** | T0.3/T0.4/T2.3/T3.9/T3.10 are active/partial, not deferred. Close the remainder by evidence. |
| Vision | **LARGELY COMPLETE — GAPS TO ASSESS** | Inventory remaining gaps before claiming completion. Package presence is not acceptance. |
| Video | **INCOMPLETE — IN SCOPE** | Own capability contract, workspace/storage estimate and lifecycle boundary; not a still-image variant for resource accounting. |
| Ordinary image display, character portraits, image assets | **IN SCOPE** | Existing display/asset paths remain eligible. |
| 3D mesh generation | **REMOVED** | No milestone, research, dependency, infrastructure, UI, or release commitment. Historical evidence remains historical. |
| 2D/3D visual-character rendering | **REMOVED** | No visual-character workstream. The T2.4 speaker-identity marker and T2.5 image-set storage remain narrowly in scope as chat/UI requirements, not as a visual-character program. |

**HCDM-05 (Calculated — scope correction is not completion).** “Complete” and “nearly complete” are owner scope statements, not E3 acceptance. TTS remains a deferred final-scope requirement. Image-generation and vision remaining work stays in scope. The corrected ledger in `HRA-DB-2026-11-21-01` controls over the older deferral/removal wording.

### 0.4 Limits

**HCDM-06 (Observed).**

- No test was run for this map. Recorded results are cited as recorded, not revalidated.
- No model, GPU, packaged app, image engine, or audio device was exercised.
- File paths and test names are evidence of implementation, not proof of current runtime behavior.
- “No implementation found” means none was found in the inspected scope, not repository-wide absence.
- This map does not reopen L-2, Decision A, G0, G1, P-1/S-3, or any settled MEMORY-09 decision.
- This map does not alter `HFS-AUDIT-2026-11-21-01`, `HFR-CP-2026-11-21-01`, `HF1-PLAN-2026-11-21-01`, canonical r33/v9, the owner ledger, or the conformance map.

---

## A. Scope reconciliation

### A.1 The 27 original IDs — status after owner controls

**HCDM-07 (Calculated — 27 owner IDs, 30 accounting rows).** `T3.3` is split into lore, legacy memory, and AI-authored lore for accounting. `T3.7` is split into ASR/audio-input and TTS because TTS is deferred but required for chat finalization. Splitting is **not** a new requirement.

**HCDM-07A (Owner-approved).** Vision and video are capability-ledger entries, not new original T-IDs; recording them does not change the 27/30 counts.

Legend: **P0** release blocker or release-safety gate; **A** core architecture; **F** feature completion; **R** reliability/fault; **P** performance; **D** deferred; **X** removed; **C** completed.

| ID | Original item | Current status | Class | Actual remaining work |
|---|---|---|---|---|
| **T0.1** | User messages do not appear until restart | Open; no recorded closure; current render/guard exists | **P0** | Reproduce JSON count, same-session render, chat-switch render; fix if still broken; add regression |
| **T0.2** | Model receives the user message twice | Open; duplicate guard exists but no prompt-dump evidence | **P0** | Prompt dump/server log showing exactly one occurrence; remove any double append |
| **T0.3** | Image model will not load | Image generation nearly complete; remaining loading-path work in scope | **F / R** | Validate stored paths, chooser and load path; close by evidence |
| **T0.4** | Generated image preview is broken | Image generation nearly complete; remaining preview work in scope | **F / R** | Verify preview end-to-end; no “deferred” disposition |
| **T0.5** | Concurrent instances kill each other’s model processes | Original defect; current job-object code appears to replace path-based reaping; **no two-instance test** | **P0 / R** | Two packaged instances; second must not kill first; abrupt exit kills only its own job; closure record |
| **T1.1** | State is not owned properly | Historical 2026-10-31 closure; V2 ledger `UNVERIFIED`; current tree not rerun | **A / P0-gate** | Current-tree tab switch, model survival, sampler persistence; decide whether historical evidence is accepted |
| **T1.2** | VRAM arbitration between chat and image models | Historical closure; resource tests absent; both sides remain in scope | **A / R** | Chat/image fit, reload, release verification; superseded by the cross-modality arbitration baseline `HRA-DB-2026-11-21-01` for architecture, not for closure |
| **T1.3** | Think tags / reasoning handling | Historical closure; stripping tests exist; hidden-thought budget unverified | **A / F** | Two-turn prompt/token test; hidden thought never reappears; hidden not counted |
| **T1.4** | Sampler defaults the user never has to touch | Historical closure; two-family evidence only | **F / A** | Current-tree model-family matrix; template/sampler application |
| **T1.5** | Header says the wrong name | Tron suffix gone; current title `HOWL Chat`; capitalization differs from literal `Howl Chat` | **C** | Confirm intended capitalization/page headers; no blocker |
| **T2.1** | Theme unification | Partial; Stage C measured five pages, not all pages; 16 px+ heuristic | **F / R** | Character/world/image-display pages; both themes; persisted screenshots; contrast limits |
| **T2.2** | Ghost elements and stale markup | **OPEN**; 18 never-created IDs over 30 reference sites | **P0 / R** | Resolve or remove 18 ghosts; zero write-to-missing-ID; duplicate-ID sweep; runtime smoke |
| **T2.3** | Image page layout — Invoke AI reference | Image generation remaining work in scope | **F** | Layout acceptance against the worklog; no blanket deferral |
| **T2.4** | Speaker identity in chat (thumbnail or initial) | Partial; image thumbnail exists; no enlarged-initial fallback; no dedicated test | **F** | Initial fallback, fixed sizing, two-speaker test. This is a chat identity marker, not a visual-character program |
| **T2.5** | Character/world/scenario image sets | Partial; character storage/preservation exists; thumbnail-default and world/scenario coverage incomplete | **F** | Multi-image round-trip and default selection; preserve ordinary image assets |
| **T3.1** | Context handling — verify before optimizing | Open; character-count estimate; no cache-reuse measurement | **A / R** | Real tokenizer, real budget, flush, prefill measurement; no character fallback |
| **T3.2** | Load speed — learn from LM Studio | Open; no benchmark harness | **P** | Reproducible load-time baseline; compare; optimize only after measurement |
| **T3.3-L** | Scoped lorebook mechanics | Substantial resolver/store; budgets in characters; no full long-chat acceptance | **A / F** | Token budget, scan depth/recursion parity, long-chat and multi-character acceptance |
| **T3.3-M** | Per-character legacy memory with triggers/context guard | Legacy memory exists; not MEMORY-09; frontend-local authority | **F / A** | Backend authority, token guard, legacy ownership decision; keep separate from MEMORY-09 |
| **T3.3-A** | AI-authored lore with specificity/frequency checks | No implementation found | **F** | Design and authorization first; not part of Phase 1 |
| **T3.4** | Per-field AI character sheets | No implementation found | **F** | Field-targeted generation, draft/accept/lock, grammar-constrained JSON; separate design |
| **T3.5** | Multi-character chat | Queue/stream exist; no contamination or flush test | **A / F** | 2–3 characters, private-context isolation, shared-scene handling, abort/restart |
| **T3.6** | Ternary models as first-class | No ternary loader path found; later master plan says planned | **F** | Research/loader evidence; no new dependency; separate authorization |
| **T3.7-A** | ASR/audio-input portion of T3.7 | Audio packages/tests exist; no UI/headphone end-to-end | **F / R** | Real microphone/UI transcript, availability/failure reporting; no TTS |
| **T3.7-T** | TTS portion of T3.7 | **Deferred, required for chat finalization** | **D / final** | No current scheduling; chat cannot be declared final without it; re-entry needs a dated owner decision |
| **T3.8** | Audio generation | **Complete** by owner | **C** | No remaining delivery claim; link evidence separately |
| **T3.9** | Diffusion version check | Image generation remaining work in scope | **F** | Verify current diffusion release and patched version; no blanket deferral |
| **T3.10** | LLM writes the image prompt — research gated | Image generation remaining work in scope | **F** | Research gate remains; scope disposition is active, not deferred |
| **T3.11** | Compatibility for everyone else | Partial hardware detection; no compatibility matrix | **A / P** | Beta compatibility matrix; hardware profiles; explicit unsupported paths |
| **T3.12** | Whole-app optimisation pass | Not started; no benchmark harness | **P** | Benchmark first; separate final stage |
| **T4.1/T4.2** | Parked 3D/visual-character concepts | **Removed** | **X** | None. Historical evidence only; no visual-character workstream |

### A.2 How many items actually remain?

**HCDM-08 (Calculated).**

- **27** owner IDs.
- **30** accounting rows after splitting `T3.3` and `T3.7`.
- **4 P0 items/gates:** T0.1, T0.2, T0.5, T2.2.
- **2 completed rows:** T1.5 and T3.8.
- **1 deferred row:** T3.7-T, which remains a required chat-finalization capability.
- **0 image-generation rows deferred:** T0.3, T0.4, T2.3, T3.9 and T3.10 remain in scope.
- **23 active/partial rows** across architecture, feature completion, reliability, and performance.
- **2 removed concepts outside the 30:** T4.1/T4.2 and visual-character rendering.

**HCDM-09 (Calculated — 27 is not 27 unfinished).** The union count is a recovery count, not a backlog count. T1.5 and T3.8 are complete by owner statement. T3.7-T is deferred but required for chat finalization. T0.3/T0.4/T2.3/T3.9/T3.10 remain in scope. Several historical items are partial rather than wholly unfinished.

### A.3 Duplicate and supersession reconciliation

**HCDM-10 (Calculated — do not double-count).**

| Overlap | Reconciliation |
|---|---|
| **T3.3-L** vs **F3-B/F3-C** | F3-B/B-C built and accepted the canonical lorebook store/import/chat consumer. That is enabling infrastructure, not a second lorebook feature. Remaining T3.3-L work is token budget, scanning/recursion parity, and end-to-end acceptance. |
| **T3.3-M** vs **MEMORY-09** | Legacy memory/lorebook is prior Howl functionality. It is **not** MEMORY-09 and is not partial MEMORY-09. Keep separate ledgers. |
| **T1.2** vs **FUT-HW1 / HRA-DB** | T1.2 is current chat/image arbitration behavior. `HRA-DB-2026-11-21-01` is the first-class cross-modality arbitration and verified-handoff architecture baseline; it is not T1.2 completion. Do not count the architecture baseline as T1.2 closure, or T1.2 as architecture approval. |
| **T3.5** vs **HSL/Phase 1** | T3.5 is the original multi-character requirement. HSL is a later state-lifecycle clarification. Phase 1 is a proposed work package. They overlap at context/identity but are not duplicates. |
| **T3.7-A** vs audio packages | Existing ASR/TTS code is implementation evidence. T3.7-A is the UI/end-to-end ASR requirement. TTS is deferred but required for finalization; package presence is not product acceptance. |
| **T0.3/T0.4/T2.3/T3.9/T3.10** | All are image-generation-adjacent and remain in scope. Do not maintain five separate release blockers for one near-complete capability, but do not drop the remaining work either. |
| **W22/W25 vs W27** | W27 records a native stable-ownership acceptance pass; W22 was open and W25 incomplete. The later pass may supersede parts of W22/W25, but that must be explicitly reconciled. Do not claim closure from the later record alone. |
| **F2-B2 tokenizer inspection** vs **T3.1 token accounting** | Model-intelligence tokenizer/template inspection is an input. T3.1 owns the actual context-budget path. Do not build two tokenizer authorities. |
| **T4/E** | Parked or removed concepts. They are not current delivery items and must not reappear as implied commitments. |

### A.4 Engineering/deferred ledger — what still applies

**HCDM-11 (Observed/Source-derived).**

| Ledger | Remaining applicable work | Disposition |
|---|---|---|
| **F2-A** | Route old save/import paths through the library writer; external edit/delete; import/export adapters; redirected-folder runtime | Active supporting work, not a separate product milestone |
| **F2-B** | Model-root UI; capability filtering; tokenizer/template inspection; packaged cache; **TTS/audio portions deferred** | Active only for chat/model-intelligence scope |
| **F2-C** | One authoritative owner, JSON mirrors, media sidecars, create-only save, migration | Design only; no cutover authorized |
| **F3-A** | A8 propagation deferred; comprehensive cleanup contract | Implemented safety with limits |
| **F3-B/F3-C** | Native acceptance/limits; no product-wide E3 | Enabling infrastructure; not a second T3.3 |
| **W13–W28** | W22/W25 reconciliation; no model inference | Ownership/prompt-boundary history; close only with evidence |
| **FUT-HW1** | Lease/telemetry/OOM authority | First-class architecture baseline `HRA-DB-2026-11-21-01`; implementation not authorized. T1.2 is not this. |
| **FUT-PROMPT1** | Model/template matrix | Deferred; T1.4 is the current sampler/template requirement |
| **FUT-MEM1/FUT-LORE1** | Structured memory and lore intelligence | Deferred; distinct from MEMORY-09 |

---

## B. Release-blocker register

### B.1 Known or owner-reported release blockers

**HCDM-12 (Calculated — release blockers are user-visible or safety-critical).** A release claim requires these to be resolved or explicitly accepted with evidence.

| ID | Exact blocker | Evidence today | Severity | Consequence | Dependencies | Verification needed |
|---|---|---|---|---|---|---|
| **R-01 / T0.1** | Sent message does not appear until restart | Owner-reported in V1/V2; no closure record; current renderer and duplicate-turn guard exist, but no same-session regression | **High** | Core chat appears broken; user may resend and duplicate work | Chat send/render path; T0.2 diagnosis shares it | One message → JSON count → same-session bubble → chat-switch bubble |
| **R-02 / T0.2** | Model receives the same user message twice | Owner-reported; duplicate guard exists; no prompt dump | **High** | Doubled prompt, wasted context, distorted replies | Shares prompt/render path with R-01 | Prompt dump/server log: exactly one occurrence per turn |
| **R-03 / T2.2** | 18 ghost IDs over 30 reference sites | `STAGE_STATUS.md:1368–1487` source inventory; no runtime proof all are benign | **High/Medium** | Silent control failures, dead buttons, misleading settings; release quality cannot be claimed | UI IDs across settings/manager surfaces; T2.1/T2.4/T2.5 add DOM | Source sweep + runtime smoke; zero write-to-missing-ID in exercised paths |
| **R-04 / DEF-FE-01** | Frontend character-configuration test fails | `tests/frontend/character-configuration-integrity.test.js:409`: actual `custom_voice_not_in_the_fixed_select`, expected `''` | **Medium** | Test suite is not green; possible voice-config state bug | Character config/voice model | Triage, fix or explicitly classify as a known test-only issue; no green-suite claim |

### B.2 Release-safety gates — not known active defects

**HCDM-13 (Calculated — do not promote a fix to verified release safety).**

| ID | Gate | Evidence today | Severity if regressed | Current classification | Verification needed |
|---|---|---|---|---|---|
| **R-05 / T0.5** | Concurrent instances must not kill each other’s model processes | Original path-based `ReapOrphans` is gone; current code uses a per-process job object with `CREATE_SUSPENDED` + `AssignProcessToJobObject`; tests cover suspension/ownership and fail-closed start. `TestD3ChatServerSpawnOwnership` is gated on a model and is not a two-instance test. | **High** | **Release-safety verification gate**, not a confirmed current defect | Two independent packaged instances; second starts, first’s child survives; abrupt exit kills only its own job; explicit ownership report |
| **R-06 / T1.1** | Tab switch must not unload the model or lose state | 2026-10-31 live run in `STAGE_STATUS.md:102–140`; V2 ledger remains `UNVERIFIED`; current tree not rerun | **High** if broken | **Release evidence gate** | Current-build tab switch, model survival, sampler persistence, restart |
| **R-07 / T1.3** | Hidden reasoning must not count and must not reappear | Stripping/splitting tests exist; no real token-budget or two-turn prompt test | **Medium** | **Core correctness gate** | Real tokenizer + two-turn prompt evidence |
| **R-08 / T1.2** | Chat/image VRAM arbitration | Historical live evidence; no resource tests; image side deferred | **Medium** for chat; **deferred** for image | **Core reliability**, not a current release blocker for the chat-only path | Chat model load/reload/release; image side only on re-entry |

### B.3 What is not a release blocker

**HCDM-14 (Calculated).**

- T3.1 token budgeting is a **core architectural completion** and a prerequisite for long-conversation/MEMORY-09 evidence, not a reason to block every release if short-chat use is honest about limits.
- T3.5 multi-character isolation is a major feature and a MEMORY-09 readiness dependency, but chat can ship without claiming multi-character correctness.
- T3.3-L/T3.3-M/T3.3-A, T3.4, T3.6, T3.11, T3.12 are feature/architecture/performance work.
- TTS is deferred but required for chat finalization; it is not a gate for unrelated chat functionality but cannot be omitted from final scope.
- Image generation, vision and video remain in scope; their remaining work must be scheduled by evidence, not waved through as “nearly complete.”
- Audio generation is complete. 3D mesh generation and 2D/3D visual-character rendering are removed and are not gates.

### B.4 T0.5 — actual severity and relationship to release safety

**HCDM-15 (Observed — current mechanism).** `internal/proc/job_windows.go` no longer contains a path-based orphan reaper. `Start` launches the child suspended, assigns it to a per-process job object, and only then resumes it. `CloseJob` fires `KILL_ON_JOB_CLOSE` for that process’s attached children. `KillTree` targets a known PID.

**HCDM-16 (Inferred — original mechanism likely fixed).** A per-process job object should prevent one Howl instance from killing another instance’s child; the old executable-path sweep is absent. This is a strong code-level explanation, not proof of the packaged two-instance behavior.

**HCDM-17 (Unestablished — release safety).** `TestSuspendedChildRunsOnlyAfterOwnership` and `TestStartFailsClosedWhenJobUnavailable` prove the ownership handshake and fail-closed startup. `TestD3ChatServerSpawnOwnership` proves the chat spawn path attaches under a job, but it is model-gated and sequential, not a two-independent-instance test. There is no recorded evidence that two packaged Howl instances coexist without interference.

**HCDM-18 (Calculated — classification).** T0.5 is **not a known current release blocker**, but it is a **release-safety gate**. Severity if the fix is incomplete is high: one instance or harness can kill another’s model child, producing unexplained model loss and corrupting confidence in resource ownership. It should be closed by evidence before a release that claims safe concurrent use.

---

## C. Dependency graph

### C.1 Readable graph

**HCDM-19 (Calculated — critical path).**

```text
Release truth
  T0.1 same-session render ─┐
  T0.2 one prompt occurrence ├──► M0 release gate
  T0.5 two-instance ownership┘
  T2.2 ghost/ownership sweep ──► M0B UI release gate

Conversation foundation
  stable state/message/turn IDs
        │
        ├──► non-destructive branch/revert/restore
        │
        ├──► real token counts ──► context budget ──► lore budget
        │                                      └──► hidden-thought accounting
        │
        ├──► stable composition (world/scenario/character)
        │
        └──► multi-character isolation ──► real multi-character use
                    │
                    └──► MEMORY-09 readiness (capture/provenance)

Model/runtime
  model lifecycle + resource ownership ──► model swap ──► long-conversation use
        │
        └──► context recomputation on model change

Feature lanes
  lorebook completion ──► authoring/ownership ──► legacy memory boundary
  UI/theme/speaker/image-set completion
  ASR end-to-end (TTS deferred)
  performance measurement ──► optimization

Separate
  MEMORY-09 capture/classification/lifecycle ──► integration only after Howl readiness
```

### C.2 Dependency table

**HCDM-20 (Calculated).** Types: **H** hard, **S** soft, **R** runtime, **V** verification.

| From | To | Type | Why |
|---|---|---|---|
| T0.1 | T0.2 diagnosis | **S/H** | Both may share the send/prompt path; a render fix must not mask a duplicate append |
| T0.1/T0.2 | Phase 1 identity work | **S/H sequencing** | Phase 1 rewrites prompt/state paths; doing it on an unverified send path creates avoidable rework |
| T0.5 | Release safety claim | **V** | Code fix is present; two-instance evidence is missing |
| T2.2 | UI feature acceptance (T2.1/T2.4/T2.5) | **S/V** | New UI adds IDs; ghost sweep must be rerun after UI changes |
| Stable IDs | Branch/revert/restore | **H** | Non-destructive state needs stable identity |
| Stable IDs | Context composition | **H** | Composition must bind to stable world/scenario/character/message IDs |
| Stable IDs | MEMORY-09 readiness | **V/H** | Memory provenance cannot depend on list position |
| Tokenizer seam | Real token counting | **H** | `App.inferenceClient` does not expose `Tokenize` |
| Real token counting | Context budget | **H** | No authoritative budget without model tokens |
| Context budget | Lore budget | **H** | Lore must consume a token share, not an unrelated character count |
| Context budget | Multi-character working sets | **H** | Flush policy needs a real budget |
| Context budget | Hidden-thought exclusion | **V** | Hidden mode cannot be proven not to consume budget without measurement |
| Stable composition | Multi-character prompts | **H** | Each character needs the correct world/scenario/character block |
| Lorebook canonical store | Lore runtime | **R** | Store/import exists; runtime still consumes it |
| Legacy memory boundary | Multi-character privacy | **R/H** | Character prompt must not read another character’s local memory |
| Model lifecycle | Model-change preservation | **R** | Unload/reload/switch must not alter history or IDs |
| Resource ownership | Model lifecycle | **H** | Load/swap needs one authority and process ownership |
| T1.2 chat arbitration | Image generation resource contract | **S/R** | Both chat and image sides remain in scope; cross-modality arbitration is `HRA-DB-2026-11-21-01` |
| W22/W25 | Character-lore acceptance | **V/H** | Authoring ownership must be proven before private-lore claims |
| T3.3-L | T3.3-A | **S** | AI lore needs a stable lore model; not a hard blocker to design |
| T3.3-M | MEMORY-09 | **None directly** | Legacy memory is not MEMORY-09; only shared state/provenance concepts matter |
| T3.7-A | TTS | **None** | TTS deferred; do not chain ASR to it |
| Performance harness | T3.12 | **H** | No optimization without baseline |
| Phase 1 evidence | MEMORY-09 real-workload test | **V** | Foundation must be honest before memory testing |

### C.3 Dependency rules

**HCDM-21 (Calculated).**

- Do not call T3.1 a hard prerequisite for every feature. Short chat can work with an honest limited budget.
- Do not call T3.5 a prerequisite for all releases. It is a prerequisite for claiming multi-character correctness and for meaningful memory testing with multiple characters.
- Do not call T0.5 a proven defect. It is a safety gate with a likely code fix and missing acceptance.
- Do not call image generation, vision, video or TTS dependencies of chat merely because they share model/resource code. Their resource contracts are handled by `HRA-DB-2026-11-21-01`; TTS remains deferred but is a chat-finalization requirement.

---

## D. Delivery milestones

### D.0 Sequencing principles

**HCDM-22 (Calculated).**

1. **Release truth before architecture.** T0.1/T0.2 are owner-visible chat failures. Phase 1 identity work should not be built on an unverified send/prompt path.
2. **Stable identity before context.** Token budgets, branching, multi-character and memory provenance all need stable IDs.
3. **Real tokenizer, no fallback.** A character estimate may remain a diagnostic, never authoritative accounting.
4. **Vertical slices.** A slice is complete only when UI, backend, persistence and runtime behavior work together.
5. **No placeholder architecture.** If a contract is unresolved, stop and resolve it; do not introduce a temporary shape that becomes permanent.
6. **Scope boundaries stay explicit.** TTS is deferred but required for chat finalization. Image generation, vision and video remain in scope. No capability is hidden inside a chat milestone, and no older deferral label overrides the corrected ledger.
7. **Accounting rows are not projects.** The 30 rows are a completeness ledger, not a one-row-per-sprint schedule. A milestone may cover several rows only where they share a contract and can be accepted together; a row must not be split into a project merely to keep the table moving.

### D.1 M0A — Release Truth and Process Safety

**HCDM-23 (Calculated — recommended first).**

| Field | Definition |
|---|---|
| Scope | T0.1, T0.2, T0.5 |
| User-visible result | A sent message appears immediately; the model receives it once; a second Howl instance does not kill the first’s model child |
| Prerequisites | Current-tree reproduction; prompt dump; packaged two-instance harness or explicit manual test |
| Acceptance | One message → one saved copy → same-session bubble; prompt/log shows one occurrence; two independent instances coexist; abrupt exit kills only the exiting instance’s child |
| Regression risks | Duplicate guard masking a second append; render fix hiding a persistence bug; test harness itself reaping children |
| Outside | Phase 1 IDs/branches/token budget; image generation; TTS; MEMORY-09 |
| Evidence | JSON count, prompt dump, screenshots/logs, two-instance PID/job evidence |

**Acceptance gate:** M0A criteria A1–A9 are defined in [`HOWL-CHAT-M0A-RELEASE-TRUTH-ACCEPTANCE-REVIEW-2026-11-21.md`](HOWL-CHAT-M0A-RELEASE-TRUTH-ACCEPTANCE-REVIEW-2026-11-21.md). The review adds the async frontend append race and the two-instance stop conditions; M0A should not be authorized without them.

### D.2 M0B — Ownership and Ghost Truth

**HCDM-24 (Calculated — can follow M0A or run in a separate UI lane).**

| Field | Definition |
|---|---|
| Scope | T2.2; T1.5 capitalization check; `DEF-FE-01` triage; frontend/backend ownership inventory for chat, lorebook, character and legacy memory. Bounded local fixes only; no state/identity redesign. |
| User-visible result | No dead controls from missing IDs; no duplicate-ID writes; header name is intentional; one documented authority map for the surfaces the next milestone touches |
| Prerequisites | Source inventory of 18 ghosts/30 sites; owned file map; M0A send/prompt ownership established |
| Acceptance | Zero write-to-missing-ID in exercised surfaces; duplicate-ID sweep clean; runtime smoke of settings/managers; title check; ownership map names the authoritative writer for each inspected surface; `DEF-FE-01` triaged |
| Regression risks | Fixing IDs by adding hidden dummy elements; broad HTML edits colliding with T2.1/T2.4/T2.5; ownership inventory turning into a premature state rewrite |
| Outside | State/identity redesign; theme completion; lorebook cutover; image-generation page (in scope but not M0B) |
| Evidence | Static sweep, runtime console/network check, screenshots, ownership inventory, defect triage |

### D.3 M1 — Conversation Identity and Restore

**HCDM-25 (Calculated — Phase 1A/1B).**

| Field | Definition |
|---|---|
| Scope | Stable state/message/turn IDs; legacy unresolved; non-destructive branch/alternative; revert/restore; persistence/restart; stable composition IDs |
| User-visible result | Regenerate/branch no longer destroys prior replies; restart restores the same current state; world/scenario/character composition survives restart |
| Prerequisites | M0A release truth; contract C1/C2/C3 agreed; no contradictory ID authority |
| Acceptance | Append → save → restart → branch → revert → restart: IDs unchanged; old material retained; deletion names its set; composition restored by stable IDs; no synthetic IDs for legacy records |
| Regression risks | Migration rewriting legacy chats; frontend/back-end ID mismatch; branch UI selecting wrong state |
| Outside | Token budget; multi-character isolation; lorebook/content migration; MEMORY-09 integration |
| Evidence | Persistence test, branch/revert inventory, restart test, legacy unresolved case |

### D.4 M2 — Context Budget and Composition Integrity

**HCDM-26 (Calculated — Phase 1C/1D).**

| Field | Definition |
|---|---|
| Scope | `Tokenize` seam; real counts; system/context/history/generation reserve; over-budget receipt; tokenizer-unavailable state; composition order; model-change recomputation |
| User-visible result | Context usage is honest; the app does not silently pretend a character count is a token count; switching models recomputes the budget |
| Prerequisites | M1 stable IDs; `App.inferenceClient` tokenizer seam; C4 agreed; model lifecycle interface C6 |
| Acceptance | Real tokenizer count; over-budget conversation with inclusion/omission receipt; unavailable tokenizer reports unavailable; model swap recomputes; no character fallback as authority |
| Regression risks | Interface fake drift; persisted token counts reused after model change; truncation rewriting authoritative history |
| Outside | Multi-character flush; lore generation; memory capture; performance claims |
| Evidence | Tokenizer test, prompt receipt, model-swap test, unavailable-tokenizer test |

### D.5 M3 — Multi-Character Isolation and Flush

**HCDM-27 (Calculated — Phase 1E).**

| Field | Definition |
|---|---|
| Scope | 2–3 characters; per-character working set; private lore/memory boundaries; shared scene; flush; abort/recovery/restart |
| User-visible result | Multiple characters speak in one scene without leaking private facts or voices; stopping one response does not corrupt the others |
| Prerequisites | M2 budget; M1 IDs; C5 agreed; legacy-memory boundary decision (not MEMORY-09) |
| Acceptance | Two characters with conflicting private facts and different histories; each prompt contains only its own private context; shared scene available; authoritative history unchanged; abort/restart |
| Regression risks | Frontend `MemoryStore` bypass; prior responses over-included; queue state lost on restart |
| Outside | MEMORY-09 capture/classification; AI-authored lore; visual-character work |
| Evidence | Two/three-character scenario, prompt inspection, restart, abort |

### D.6 M4 — Model Lifecycle, Resource Arbitration and Hardware Safety

**HCDM-28 (Calculated — architecture-dependent; not authorized).**

| Field | Definition |
|---|---|
| Scope | One resource authority; model load/unload/swap; cross-modality admission; verified handoff; logical reservation vs physical allocation; drain/release/readiness evidence; telemetry loss; failed load recovery; process and cross-instance ownership; T1.2; T1.3/T1.4; T3.6; T3.11 |
| User-visible result | Selecting is not loading; a handoff is explicit; failure never looks ready; unsaved work survives; the app does not claim resources it cannot verify |
| Prerequisites | `HRA-DB-2026-11-21-01`; `DQ-HRA-*` decisions; C6/C12; T0.5 evidence; M1 IDs for state preservation; M2 token recalculation |
| Acceptance | Selection-only change; atomic reservation; safe coexist/transition; drain; verified/partial/unknown release; workload-specific readiness; failed handoff recovery; work preservation; cross-instance policy; no history loss |
| Regression risks | Two resource authorities; stale ledger; busy-model kill; unknown-as-free; cross-instance double-spend; modality-specific side effects |
| Outside | TTS implementation (deferred); visual-character work; performance optimization; image/video feature completion (their resource contracts are in scope here) |
| Evidence | Fault-injected resource tests, handoff traces, reservation/commit races, release evidence, model swap, PID/job checks, cross-instance test |

**Planning gate:** M4 may not be authorized from this map alone. `HRA-DB-2026-11-21-01` must first resolve the resource model, lifecycle protocol, ownership/concurrency rules, failure/recovery behavior, evidence requirements and future-modality contracts. A text-only implementation slice is not authorized by this record.

### D.7 M5 — Lore, Authoring, and UI Completion

**HCDM-29 (Calculated — split into M5A/M5B).**

**M5A — Lore and authoring ownership**

| Field | Definition |
|---|---|
| Scope | T3.3-L runtime completion; T3.3-A design/implementation only if separately authorized; T3.4 separately authorized; W22/W25 reconciliation; F3-B/C limits; T2.5 storage |
| User-visible result | Lore entries resolve predictably; owner/privacy rules work; authors can see and correct what is scoped where |
| Prerequisites | M1 composition IDs; M2 token budget for lore; C7/C8; W22/W25 evidence |
| Acceptance | Token-based lore budget; scan depth/recursion; owner/privacy; canonical store round-trip; no frontend authority bypass |
| Regression risks | Duplicate lore authorities; owner-by-name; legacy records rewritten; F3 store reimplementation |
| Outside | MEMORY-09; AI-authored lore if not separately authorized; image generation |
| Evidence | Lorebook tests, long-chat prompt, owner acceptance, store round-trip |

**M5B — UI completion**

| Field | Definition |
|---|---|
| Scope | T2.1 theme completion; T2.4 speaker identity marker; T2.5 image sets/defaults; T1.5 header check |
| User-visible result | Consistent readable themes; clear speaker identity; multi-image sets persist with a chosen default |
| Prerequisites | M0B ownership/ghost truth; M1 stable composition; C3/C9 |
| Acceptance | Both themes across remaining pages; initial fallback; two-speaker layout; multi-image save/load/default |
| Regression risks | Ghost IDs reintroduced; visual changes colliding with T2.2; image-set data loss |
| Outside | Image generation; visual-character rendering; image-page redesign |
| Evidence | Screenshots, UI tests, storage round-trip |

### D.8 M6 — Performance and Compatibility

**HCDM-30 (Calculated — after runtime is stable).**

| Field | Definition |
|---|---|
| Scope | T3.2 load-time baseline; T3.12 whole-app optimization; T3.11 compatibility matrix |
| User-visible result | Measured load-time improvements and honest hardware compatibility, not intuition-driven changes |
| Prerequisites | M4 stable model lifecycle; reproducible harness; named hardware/models |
| Acceptance | Baseline, metric, variance, before/after; compatibility profiles; explicit unsupported paths |
| Regression risks | Optimizing an unstable path; benchmark gaming; hardware-specific assumptions |
| Outside | Correctness fixes; image-generation performance; TTS |
| Evidence | Benchmark records, hardware matrix, no speed claim without measurement |

### D.9 M7 — ASR End-to-End

**HCDM-31 (Calculated — independent feature lane; TTS excluded).**

| Field | Definition |
|---|---|
| Scope | T3.7-A ASR/audio-input; capability status; microphone/UI transcript; failure handling |
| User-visible result | Voice input works through the real composer with truthful availability/errors |
| Prerequisites | Existing audio packages; no TTS dependency; C9 status/error contract |
| Acceptance | Real microphone capture → transcript in composer; unavailable engine reports unavailable; failed capture leaves working state |
| Regression risks | Implicit TTS path; shared temp/ownership; audio manager displacing chat model |
| Outside | TTS (deferred but required for finalization); audio generation (complete); remote ASR unless separately authorized |
| Evidence | Real UI/audio run, capability log, failure case |

### D.10 M8 — MEMORY-09 Readiness Gate (not a product milestone)

**HCDM-32 (Calculated — separate workstream).**

| Field | Definition |
|---|---|
| Scope | Verify that Howl has the prerequisites for meaningful memory testing; no capture/classification integration |
| Prerequisites | M1 IDs, M2 budget, M3 multi-character, M4 persistence/model swap, M8 contract |
| Acceptance | Stable source identity; branch/deletion semantics; context budget; multi-character isolation; restoration; failure reporting; real-workload harness |
| Outside | MEMORY-09 capture, classification, lifecycle, retrieval, G2, BGE, integration |
| Evidence | Readiness checklist signed against actual Howl behavior |

---

## E. Parallel workstream map

### E.1 Safe parallel lanes

**HCDM-33 (Calculated).** Parallelism is allowed only after interface ownership is agreed. File overlap alone is not permission to merge behavior.

| Lane | Scope | Required interfaces | Likely overlap | Integration point | Conflict avoidance | Pre-merge tests |
|---|---|---|---|---|---|---|
| **P0 — Release truth** | T0.1/T0.2/T0.5 | Send/prompt event contract; process ownership | `app-prompt.go`, `app-chat.go`, `chat-render.js`, `chat-stream.js` | Chat send and model spawn | One owner for prompt/render files; no Phase 1 rewrite until M0A closes | JSON count, prompt dump, two-instance |
| **P1 — State/context** | Phase 1A–1E | C1–C5 | `app-chat.go`, `app-prompt.go`, `internal/backend/chat`, `types`, `chat-actions.js`, `chat-session.js` | Prompt composition and persistence | Sequence after M0A; contract-first; no parallel edits to the same state types | ID/branch/restart, token budget, multi-character |
| **P2 — Model lifecycle** | M4 | C4/C6 | `app-model.go`, `app-llama.go`, `app-vram-arbitration.go`, `internal/resource`, `internal/proc` | Model load/swap and prompt assembly | Freeze tokenizer/resource interface before P1 budget work; one authority | Resource fault tests, model swap, process ownership |
| **P3 — Lore/authoring** | M5A | C3/C7/C8 | `app-lorebooks-store.go`, `internal/lorebookstore`, `lorebooks.js`, `character-form.js`, `app-memory.go` | Composition and character authoring | Freeze stable IDs first; do not create a second lore store | Lorebook/store/owner tests, long-chat prompt |
| **P4 — UI/theme** | M5B | C3/C9 | HTML/CSS, `chat-render.js`, `character-images.js`, `scenario-form.js` | Display and authoring surfaces | Run after M0B ownership/ghost truth or coordinate ID changes | Theme screenshots, speaker/image-set UI |
| **P5 — ASR** | M7 | C9; audio manager boundary | `internal/audio/asr`, `app-audio.go`, `audio-settings.js`, `audio-service.js` | Composer input and capability status | No TTS changes; do not schedule implicit synthesis | Real mic/UI, unavailable engine |
| **P6 — Performance** | M6 | Benchmark contract | No production edits first; later hot paths | Measured runtime | Start after P2/P3 stabilize | Baseline/variance; no speed claim without measurement |
| **P7 — MEMORY-09** | Separate M1–M8 | C10 only | `eval/mem1/m09`, MEMORY-09 docs | Capture boundary later | No production integration; no Howl requirement invented from memory | Memory fixtures only until gate |

### E.2 Unsafe combinations

**HCDM-34 (Observed/Calculated).**

- **P1 + P0 without sequencing:** both touch prompt/render and can produce two “fixes” for the same symptom.
- **P1 + P2 without a tokenizer/resource contract:** two token-count or model-state authorities.
- **P3 + P1 without stable IDs:** lorebook/authoring writes identity by name or position again.
- **P4 + P0B without a ghost-ID owner:** UI changes reintroduce missing IDs while the sweep is running.
- **P5 + TTS work:** TTS is deferred but required for finalization; implicit synthesis or shared lifecycle changes would violate the current scheduling boundary.
- **P6 + unstable P2/P3:** benchmarks measure a moving target.

---

## F. Architectural contract register

**HCDM-35 (Calculated — agreement before implementation; not authorization).**

| ID | Contract | Must agree | Consumers | Current state | Blocks |
|---|---|---|---|---|---|
| **C1** | Message/state/turn identity | ID assignment, uniqueness, persistence, legacy-unresolved rule | Chat, persistence, Phase 1, MEMORY-09 later | Positional IDs today; Phase 1 proposes stable IDs | M1, branch/revert, memory provenance |
| **C2** | Branch/revert/deletion semantics | What is retained, what is current, what deletion names | Chat UI, persistence, context composition | No state model; destructive paths exist | M1, restore, MEMORY-09 history |
| **C3** | Composition order and stable IDs | World/scenario/character IDs; prompt order; missing-ID behavior | Prompt assembly, authoring, multi-character | Names and IDs mixed; composition order exists | M1/M3/M5 |
| **C4** | Token accounting | Tokenizer seam; per-message counts; model-change recompute; unavailable state | Context, lore, multi-character, UI status | Character estimate; no `App` tokenizer | M2, M3, memory readiness |
| **C5** | Multi-character working set | Per-character private context; shared scene; flush; prior-response labels | Chat stream, character queue, lore/memory | Queue exists; no isolation test | M3 |
| **C6** | Model lifecycle/resource ownership | Load/unload/swap, lease/admission, telemetry confidence, process/job ownership | Model, resource, prompt, UI | Advisory allocation; job object in code; no two-instance test | M4, model swap |
| **C7** | Persistence authority | Chat/history authority; legacy unresolved; no silent rewrite | Chat, storage, migration | Split App/chat.Service and legacy paths | M1, M5 |
| **C8** | Lore/authoring ownership | Canonical store, owner IDs, W22/W25, no frontend authority | Lorebook, character, prompt | Store exists; authoring gap remains | M5A |
| **C9** | UI/backend error and status | Unavailable/loading/failed states; no silent fallback | All UI, model/audio/resource | Mixed status behavior | M0B/M4/M7 |
| **C10** | MEMORY-09 capture boundary | Stable source identity, provenance, branch/deletion, context receipt | Howl/ MEMORY-09 integration | Proposal only; not authorized | Any real memory test |
| **C11** | Capability-scope boundary | Image generation, vision and video are in scope; TTS is deferred but required for finalization; audio generation is complete; no capability is silently collapsed into a text-model special case | Program, resource, UI | Corrected 2026-11-21; older deferral wording superseded | Prevents hidden commitments and stale scope |
| **C12** | Resource arbitration and verified handoff | Admission, logical vs physical reservation, drain, release evidence, readiness proof, cross-instance coordination, evidence ladder | Resource authority, lifecycle adapters, chat, modalities, UI | `HRA-DB-2026-11-21-01`; unresolved `DQ-HRA-*` | Any claim of safe multi-modality or concurrent-instance resource use |
| **ADR-001** | [`ADR-001 — Foundational Backend Rebuild`](ADR-001-FOUNDATIONAL-BACKEND-REBUILD-2026-11-22.md) — decision record, not a contract | Five roots; configuration authority and its routing record; persistence lifecycle; runtime resolution and evidence levels; ownership classes and external-write grants; D1–D5 | PathService, ConfigAuthority, content stores, runtime resolution | Accepted as decision record 2026-11-22; no implementation authorization | Nothing. Informs M0A–M8; closes no C1–C12 item |
| **ADR-002** | [`ADR-002 — Storage Locations, Picker Navigation and Content Routing`](ADR-002-STORAGE-AND-NAVIGATION-POLICY-2026-11-22.md) — decision record, not a contract | Per-category locations; picker navigation state; browser surfaces are source-only; N1–N19; refines ADR-001 D1 for chats | PathService, UI pickers, content routing | Accepted as decision record 2026-11-22; no implementation authorization | Nothing. Informs M0A–M8; closes no C1–C12 item |

**HCDM-36 (Observed — interface conflict risk).** C1, C3, C4, C5, C6 and C7 are the high-collision contracts. They cross `app-prompt.go`, `app-chat.go`, `chat-stream.js`, model/resource code and persistence. They must be agreed before parallel implementation, not discovered through merge conflicts.

**HCDM-36a (Observed — decision records sit outside C1–C12).** The two rows labelled `ADR-001` and `ADR-002` are **decision records**, not contracts, and they are deliberately not numbered C13/C14: they close no contract, authorize no implementation, and do not amend FC-1–FC-6 or any frozen artifact. They are registered here so that a worker planning this map finds the detailed decisions rather than rediscovering them. Where a decision record and a contract appear to disagree, neither controls by silence — the discrepancy is recorded for owner resolution.

---

## G. MEMORY-09 readiness relationship

### G.1 Howl capabilities needed regardless of MEMORY-09

**HCDM-37 (Calculated).**

- Stable state/message/turn identity.
- Non-destructive branch/revert/restore and explicit deletion.
- Correct world/scenario/character composition.
- Real tokenizer-based context budget.
- Multi-character isolation and flush.
- Predictable model lifecycle and process ownership.
- Persistence across restart and model change.
- Truthful availability/error reporting.

These are Howl product requirements. They are not created by MEMORY-09.

### G.2 Capabilities specifically needed for meaningful memory testing

**HCDM-38 (Calculated).**

| Capability | Why memory needs it | Current evidence |
|---|---|---|
| Stable source identity | Provenance must name the message/state that produced a memory | Positional IDs; Phase 1 proposed |
| Branch/revert/deletion semantics | Memory must know what happened to the source material | No state model |
| Context receipt | Memory must know what the character could actually see | Character estimate only |
| Multi-character boundaries | Memory access must not cross characters or private facts | No isolation test |
| Persistence/restoration | Memory must survive restart/model change | History persistence partial |
| Failure reporting | Capture must not invent completeness when unavailable | Split authority/local fallback |
| Production capture boundary | Facts must enter once, with origin and identity | `PCB-01…42` proposal only |

### G.3 MEMORY-09 work that stays in its own workstream

**HCDM-39 (Observed).**

- M1 capture boundary implementation.
- M2 autonomous classification, G2, corpus.
- M3 unresolved-class lifecycle and L-2/Decision A questions.
- M4 acquisition/lifecycle event semantics.
- M5 ordering/authority/provenance/reconstruction/integrity.
- M6 conformance execution.
- M7 engine/harness integration.
- M8 Howl ↔ MEMORY-09 contract.
- G2/BGE/model/corpus authorization.

**HCDM-40 (Calculated — separation rule).** No MEMORY-09 task is a prerequisite for unrelated Howl UI, lorebook, ASR, performance or model-format work. No legacy lorebook/memory completion is counted as MEMORY-09 progress.

### G.4 Readiness gate for real memory testing

**HCDM-41 (Calculated — gate, not implementation).**

A real-workload MEMORY-09 test is meaningful only when:

1. Message/state/turn IDs are stable across save/load/branch/restart.
2. Branch/revert/delete behavior is explicit and non-destructive.
3. The character working set has a real token budget and an inclusion/omission receipt.
4. Multiple characters have isolated private context.
5. Model changes preserve history and recompute context.
6. Failure/unavailable states are truthful.
7. A production capture boundary is separately designed and authorized.
8. The owner can run a real creative workload, not only a fixture.

Until then, MEMORY-09 can proceed on its own architecture/experiments, but Howl is **not ready** for real-workload memory testing.

---

## H. Deferred and removed scope

**HCDM-42 (Observed — owner controls).**

| Scope | Disposition | Historical evidence preserved | Current delivery effect |
|---|---|---|---|
| **Audio capture/input** | **IN SCOPE — separate capability** | `internal/audio/`; audio settings/tests | Own pipeline/lifecycle/resource envelope; not generation or TTS |
| **Audio generation** (T3.8) | **COMPLETE** | Historical issue/expansion records; owner statement | No remaining delivery claim; evidence must be linked separately |
| **TTS** (T3.7-T) | **DEFERRED, REQUIRED FOR CHAT FINALIZATION** | Existing TTS packages, helper tests, design records | No current scheduling; cannot be omitted from final scope; re-entry needs explicit owner decision |
| **Image generation** (T0.3, T0.4, T2.3, T3.9, T3.10) | **IN SCOPE — NEARLY COMPLETE** | `IMAGE_PAGE_WORKLOG.md`; image engine/tests; preview observation | Remaining work stays in scope; close by evidence, not by a blanket deferral; ordinary image display/portraits/assets remain |
| **Vision** | **IN SCOPE — LARGELY COMPLETE, GAPS TO ASSESS** | Vision projector/input paths and tests | Inventory remaining gaps before claiming completion |
| **Video** | **IN SCOPE — INCOMPLETE** | Video mode paths in the image workflow | Own capability contract, workspace/storage and lifecycle; not a still-image variant |
| **3D mesh generation** | **REMOVED** | Historical issue/expansion records | No milestone or implied roadmap commitment |
| **2D/3D visual-character rendering** | **REMOVED** | Historical records only | No visual-character workstream; T2.4 is limited to a speaker-identity marker |

**HCDM-43 (Calculated — capability boundaries).** Image generation, vision and video remain in scope; their remaining work must be closed by evidence. TTS is deferred but required for chat finalization. Audio capture/input, audio generation and TTS are separate capabilities. Existing image display, character portraits and image assets remain in scope. No boundary may be used to silently drop remaining work or to schedule it without authorization.

**HCDM-44 (Observed — no hidden commitment).** These items must not reappear as:
- a prerequisite for a chat milestone;
- an implied “later” phase inside M4/M5;
- a release gate for unrelated functionality;
- a reason to keep an unused resource/UI path active.

---

## I. Open questions and risks

**HCDM-45 (Calculated — only sequencing/planning blockers).**

| ID | Question / risk | Why it matters | Current evidence | Disposition |
|---|---|---|---|---|
| **Q-01** | Is T0.1 still reproducible on the current build, and is T0.2 a separate prompt duplicate? | Determines whether M0A is a fix or a verification-only close | Owner report; no current regression | Investigate first; do not assume fixed |
| **Q-02** | Does the current job-object path actually prevent two packaged instances from interfering? | Release-safety claim; T0.5 closure | Code suggests yes; no two-instance test | Build the M0A two-instance evidence |
| **Q-03** | Are T1.1–T1.4 historical closures accepted as current, or must each be rerun? | Release/evidence scope; avoids stale claims | 2026-10-31 live records; V2 `UNVERIFIED` | Rerun the minimal current-tree checks; do not rely on old dates |
| **Q-04** | Do W22/W25 remain open after W27’s native ownership pass? | Character-lore authoring acceptance; private lore | W22 open/W25 incomplete; W27 later native pass | Reconcile records before claiming M5A |
| **Q-05** | Is `DEF-FE-01` a product defect or a test-fixture defect? | Release gate honesty; voice configuration | Test fails; HFR says separate | Triage; do not hide |
| **Q-06** | Does ASR remain active while TTS is deferred but required for finalization? | M7 scope; audio dependency | Owner separates audio capture/input, audio generation and TTS; TTS is deferred but final-scope | Treat ASR as active; TTS as deferred but required for chat finalization; confirm only if it changes sequencing |
| **Q-07** | What exactly remains in the nearly complete image-generation scope? | Prevents dropping remaining work or claiming completion from package presence | Owner says remaining image work stays in scope; T0.3/T0.4/T2.3/T3.9/T3.10 | Inventory and close by evidence; display/assets stay |
| **Q-08** | Are T2.4/T2.5 in scope without reopening visual-character work? | Owner removed visual-character rendering | T2.4 is speaker identity; T2.5 is image-set storage/display | Keep them narrowly; no visual-character program |
| **Q-09** | Do two Howl instances coordinate reservations, or is concurrent heavy work explicitly unsupported? | T0.5 prevents process interference but not logical double-spend | Current job-object ownership is per instance; no shared reservation ledger | Resolve in `HRA-DB-2026-11-21-01` before claiming concurrent-instance resource safety |

**HCDM-46 (Calculated — settled items not reopened).** L-2, Decision A, G0, G1, P-1/S-3, r33/v9, the owner ledger and the conformance map are not open questions here. The only MEMORY-09 questions are sequencing and readiness, not settled semantics.

---

## J. Recommended next bounded work package

### J.1 Recommendation

**HCDM-47 (Calculated — recommendation only).** Review the M0A acceptance criteria, then authorize **M0A — Release Truth and Process Safety** before Phase 1 implementation if those criteria are accepted.

**Scope:** T0.1, T0.2, T0.5. Exact acceptance criteria: [`HOWL-CHAT-M0A-RELEASE-TRUTH-ACCEPTANCE-REVIEW-2026-11-21.md`](HOWL-CHAT-M0A-RELEASE-TRUTH-ACCEPTANCE-REVIEW-2026-11-21.md) A1–A9.

**Why this package:**

- T0.1/T0.2 are the only owner-reported chat failures that directly break normal use. A release cannot be credible while the user cannot trust that a message appears once and reaches the model once.
- T0.5 has a likely code fix but no two-instance evidence. A release-safety claim needs the packaged behavior, not a code reading.
- It is small, testable, and does not introduce architecture.
- It prevents Phase 1 from rebuilding prompt/state code on top of an unverified send path.
- It creates the prompt-dump and two-instance harness that later Phase 1 and model-lifecycle work will reuse.

**Prerequisites:**

- Reproduce T0.1/T0.2 on the current build; capture JSON count and prompt/server log.
- Identify the single owner of the send/render/prompt path.
- Build a two-instance packaged test or an explicitly documented manual equivalent; do not use a pattern-based kill against the developer’s own process tree.
- Freeze `app-prompt.go`/`chat-render.js` edits for the duration so the fixes are not raced by Phase 1.

**Acceptance:**

1. One sent message produces one saved user turn and one same-session bubble.
2. The model prompt contains the user turn exactly once.
3. Chat switch/return does not duplicate or hide the message.
4. Two independent Howl instances coexist; starting the second leaves the first’s model child alive.
5. Abrupt exit kills only the exiting instance’s owned child/job.
6. Evidence is recorded; no “looks fixed” closure.

**Outside M0A:** stable IDs, branches, token budgets, multi-character, image generation, TTS, MEMORY-09.

### J.2 What comes next

**HCDM-48 (Calculated).**

1. **M0B** — ownership and ghost truth (parallel UI lane after M0A ownership is clear).
2. **M1** — Phase 1A/1B identity, branch/revert, restore.
3. **M2** — Phase 1C/1D tokenizer and context budget.
4. **M3** — Phase 1E multi-character isolation.
5. **M4/M5/M6** — model lifecycle, lore/authoring/UI, performance.
6. **M7** — ASR/audio-input end-to-end; TTS remains deferred but required for chat finalization.
7. **M8** — MEMORY-09 readiness gate; memory work continues separately.
8. **HRA** — resource arbitration and verified handoff architecture gates; no implementation slice authorized.

### J.3 Phase 1 status

**HCDM-49 (Observed — unchanged).** `HF1-PLAN-2026-11-21-01` remains **proposed, not authorized**. Phase 1 remains the correct first **architectural** package, but it should be split into M1–M3 and preceded by M0A. Its real-tokenizer requirement is preserved: no character-count fallback may become authoritative. Its valid findings on identity, branching, composition, flushing, model-change preservation, and long/multi-character verification remain in force.

**HCDM-50 (Calculated — final state).** Howl is not ready for MEMORY-09 real-workload testing. The program should first make chat truthful, then make state/context honest, then make multi-character and model lifecycle reliable, while delivering usable UI and authoring improvements in bounded slices. Image generation, vision and video remain in scope; TTS is deferred but required for finalization; audio generation is complete; 3D mesh and visual-character work are removed; the separate MEMORY-09 workstream must not distort that sequence.

---

## K. Verification and change record

**HCDM-51 (Observed — this map).**

| Check | Result |
|---|---|
| Production code/schema/API/migration changed | **0** |
| Tests added/changed/run | **0** |
| Models/corpus/benchmarks | **0** |
| MEMORY-09 implementation/integration | **0** |
| `HFS-AUDIT-2026-11-21-01` changed | **No** |
| `HFR-CP-2026-11-21-01` changed | **No** |
| `HF1-PLAN-2026-11-21-01` changed | **No** |
| Program direction changed | **No** |
| Canonical r33/v9/ledger/map changed | **No** |
| Pre-existing user changes touched | **No** |
| This file | New, untracked; program map |
| Resource arbitration baseline | `HRA-DB-2026-11-21-01` added; implementation not authorized |
| Canonical hashes | r33 `dfb28b6003ca356f…` · v9 `698254ca262c0874…` · ledger `e5008020ff68bd21…` · map `b4ba33ccbfe8f8c1…` |

**HCDM-52 (Observed — source paths).**

| Topic | Paths |
|---|---|
| Original scope | `PROJECT_ISSUES_AND_EXPANSION.md`; `ISSUES.md` |
| Scope audit | `docs/architecture/HOWL-FOUNDATION-ORIGINAL-SCOPE-AND-EVIDENCE-CHAIN-2026-11-21.md` |
| Readiness | `docs/architecture/HOWL-FOUNDATION-READINESS-AND-MEMORY-09-DECISION-REGISTER-2026-11-21.md` |
| Phase 1 | `docs/architecture/HOWL-FOUNDATION-PHASE-1-CHAT-STATE-CONTEXT-INTEGRITY-2026-11-21.md` |
| Program direction | `docs/architecture/HOWL-MEMORY-09-PROGRAM-DIRECTION-2026-11-21.md` |
| Product inventory | `docs/review/HOWL-COMPLETE-PRODUCT-ARCHITECTURE-REQUIREMENTS-DEPENDENCY-MASTER-PLAN-2026-11-11.md` |
| Resource arbitration | `docs/architecture/HOWL-RESOURCE-ARBITRATION-AND-VERIFIED-HANDOFF-DISCOVERY-BASELINE-2026-11-21.md`; `AMENDMENT-01-UNIFIED-HARDWARE-INTELLIGENCE-2026-11-11.md`; `DESIGN-HARDWARE-RESOURCE-ORCHESTRATION-2026-11-11.md` |
| Modality scope | `HRA-DB-2026-11-21-01` §2; `PROJECT_ISSUES_AND_EXPANSION.md`; `IMAGE_PAGE_WORKLOG.md` |
| Process ownership | `internal/proc/job_windows.go`; `d3_chat_ownership_test.go`; `internal/proc/job_windows_test.go` |
| Chat/prompt | `app-chat.go`; `app-prompt.go`; `internal/backend/chat/` |
| State types | `internal/backend/types/context.go`; `internal/backend/types/message.go` |
| Tokenizer | `internal/backend/llama/http_client.go`; `app.go` |
| Multi-character | `frontend/dist/js/character-queue.js`; `chat-stream.js`; `chat-render.js` |
| Lorebook/memory | `internal/backend/lorebook/`; `internal/lorebookstore/`; `internal/backend/memory/`; `frontend/dist/js/memory-store.js` |
| Image page history | `IMAGE_PAGE_WORKLOG.md`; `PROJECT_ISSUES_AND_EXPANSION.md` |
| ASR/TTS | `internal/audio/`; `app-audio.go`; `docs/asr-repair-design.md` |

**HCDM-53 (Observed — proposed link, not applied).** The program-direction record `docs/architecture/HOWL-MEMORY-09-PROGRAM-DIRECTION-2026-11-21.md` does not link this map. Proposed insertion: add this file to §2.1. It was **not applied** because the program-direction record is a dated checkpoint and this map does not modify existing records.

**HCDM-54 (Observed — scope correction recorded).** The corrected modality ledger and the resource-arbitration discovery baseline are recorded in `HRA-DB-2026-11-21-01`. Older deferral/removal statements in this map are superseded for current scope. No implementation is authorized by that baseline or by this correction.

**End of map.** The delivery program is: release truth → conversation identity/restore → real token budget → multi-character isolation → model lifecycle/resource arbitration → lore/authoring/UI completion → measured performance → ASR. Image generation, vision and video remain in scope; TTS is deferred but required for chat finalization; audio generation is complete; 3D mesh and visual-character work are removed. Phase 1 remains proposed and is split, not weakened. MEMORY-09 remains separate and cannot be tested meaningfully until the Howl readiness gate is met.
