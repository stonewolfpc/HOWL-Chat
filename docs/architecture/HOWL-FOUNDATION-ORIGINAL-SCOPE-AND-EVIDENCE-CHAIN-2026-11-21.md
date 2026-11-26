# Howl Foundation — Original Scope Recovery and Evidence Chain

## Read-only audit `HFS-AUDIT-2026-11-21-01`

| Field | Value |
|---|---|
| Date | **2026-11-21** (local, CST) |
| Type | Documentation discovery, historical scope recovery, evidence reconciliation |
| Baseline | **`HFR-CP-2026-11-21-01`** — retained, not overwritten |
| Related package | **`HF1-PLAN-2026-11-21-01`** — proposed; not authorized |
| Scope | Original Howl foundation work; implementation history; tests; remaining work |
| Production changes | **0** |
| Tests executed by this audit | **0** — inspection and recorded evidence only |
| Models run / corpus / benchmarks | **0** |
| Canonical MEMORY-09 files changed | **0** (four canonical paths hashed; all four match HFR’s recorded digests) |
| This file | New, untracked; **audit record**, not a requirements baseline |

**Evidence labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** · **Superseded**.

---

## 0. Method, source hierarchy, and limits

### 0.1 What was inspected

**HFS-01 (Observed).** This audit inspected:

- The owner-facing working issue and expansion records: `ISSUES.md`, `PROJECT_ISSUES_AND_EXPANSION.md`.
- Historical execution records: `STAGE_STATUS.md`, `T1_NOTES.md`, `TECH_DEBT.md`, `IMAGE_PAGE_WORKLOG.md`, and the F2/F3/W records referenced by the expansion document.
- The 2026-11-11 design and amendment package: `RECONCILED-ARCHITECTURE-PLAN-2026-11-11.md`, `INTEGRATION-AND-ACCEPTANCE-PLAN-2026-11-11.md`, `AMENDMENT-01-*`, `AMENDMENT-02-*`, and the four design notes.
- `docs/review/HOWL-COMPLETE-PRODUCT-ARCHITECTURE-REQUIREMENTS-DEPENDENCY-MASTER-PLAN-2026-11-11.md` — a later planning inventory that reproduces T0.1–T3.12 and explicitly reclassifies the Tier 1 ledger as `UNVERIFIED` despite the `STAGE_STATUS` closures.
- The later requirement record: `REQUIREMENTS-HOWL-CHAT-STATE-LIFECYCLE-2026-11-19.md` (HSL).
- The readiness audit `HFR-CP-2026-11-21-01` and the Phase 1 plan `HF1-PLAN-2026-11-21-01`.
- Current code and test names for chat state, prompt assembly, lorebook/memory, multi-character queue, resource management, process ownership, model intelligence, image loading, and audio.

No UI was rendered; no model or GPU was used; no test was run by this audit. The test results in §B are the results recorded by `HFR-CP-2026-11-21-01` and earlier execution records, not fresh runs.

### 0.2 Source hierarchy

**HFS-02 (Source-derived — two original-list generations exist).**

| Generation | Source | Role | What it contains |
|---|---|---|---|
| **V1** | `ISSUES.md` | Earlier owner-facing working list | Tiers 0–4; no T2.4/T2.5; no F2/F3 ledgers |
| **V2** | `PROJECT_ISSUES_AND_EXPANSION.md` | Later merged owner-facing working list | Tiers 0–4; adds T2.4/T2.5; adds F2/F3, W-series, deferred FUT-HW1/FUT-PROMPT1 |
| **Engineering appendices** | F2/F3/W sections in V2; `STAGE_STATUS.md`; review records | Implementation history | Approved storage/model work, ownership cutovers, defect repairs |
| **Later design** | 2026-11-11 architecture/amendment set | Design proposals and owner decisions | DWP-01–DWP-09; no implementation authorization |
| **Later requirement** | HSL, 2026-11-19 | Owner-stated Howl responsibility split | HSL-1–HSL-11; no implementation |
| **Audit/plan** | HFR, 2026-11-21; HF1, 2026-11-21 | Readiness evidence and proposed work package | A–E readiness; Phase 1 acceptance criteria |

**HFS-03 (Calculated — V2 is the later recovered list, but it silently narrows V1).** V2 is later, explicitly merges earlier notes, retains the T0–T3 wording, and adds T2.4/T2.5 and F2/F3 history. V1 is retained as the earlier wording, not silently merged. The differences are mostly additions:

- V2 adds T2.4 (speaker identity) and T2.5 (character/world/scenario image sets).
- V2 adds the F2, F3-A, F3-A-R1, F3-B, F3-C, W13–W28, FUT-HW1, and FUT-PROMPT1 records.
- V2 leaves the ordering law and most Tier wording intact.
- **V2 drops T0.5 (concurrent instances killing each other's model processes) without a recorded closure or rationale.** V1 has T0.5 at `ISSUES.md:69–86`; later `docs/asr-repair-design.md:10` still says “T2.2 and T0.5 remain open and unchanged.” The omission is therefore not evidence that T0.5 was fixed or deselected. This audit preserves it.
- V2’s “Documentation age” section says no line is assumed finished; this audit takes that seriously.

**HFS-04 (Observed — V2 is not a signed requirements baseline).** It is a working project document, not an owner-signed specification. Its statements are source-derived evidence, not proof that a stage closed. A closure claim is accepted only where the document records verification, and then at the evidence level that verification actually supports.

**HFS-05 (Calculated — original active count).** The recovered original product list contains **27 active items** as the union of V1 and V2:

- T0.1–T0.5: **5** release blockers (V2 lists only T0.1–T0.4; its omission of T0.5 is unreconciled).
- T1.1–T1.5: **5** core-mechanics items.
- T2.1–T2.5: **5** UI items (V1 had only T2.1–T2.3; V2 adds T2.4/T2.5).
- T3.1–T3.12: **12** additions.

Counts: **V1 = 25**, **V2 = 26** (25 + 2 additions − 1 removal), **union = 27**. This audit treats **27** as the recovered original scope; neither the 25-row V1 nor the 26-row V2 is complete by itself. Tier 4 contains **2 parked ideas**, not scheduled work. E.1–E.8 are parked expansion concepts, not active requirements. F2/F3/W and the deferred FUT-* entries are engineering or deferred ledgers, not members of the 27.

**HFS-06 (Calculated — “foundation” must not be narrowed to A–E).** The readiness audit’s A–E categories are an evidence organization, not the original list. They omit release blockers, UI work, image-page work, per-field generation, audio generation, ternary support, diffusion-version checking, and LLM-image-prompt research. A statement that “the foundation is ready” cannot be made from A–E alone.

### 0.3 Limits

**HFS-07 (Observed).**

- The original documents do not consistently distinguish “requirement,” “proposed solution,” and “implementation detail”; this audit labels each where the source permits.
- Authorship is not recorded per line. T0–T3 use an owner-facing first-person issue style; this audit calls them **user-directed** as a source-derived judgment, not a proven authorship record.
- F2/F3/W entries are engineering-led records; their status words are not independent verification.
- No historical test was rerun. A test name proves a test exists, not that it passes now.
- “No implementation found” means no implementation was found **in the inspected paths**, not that none exists anywhere.
- The audit does not reopen L-2, Decision A, G0, G1, or any settled MEMORY-09 decision.

---

## A. Original scope recovery

### A.1 Exact original list — Tier 0 (release blockers)

Source: `PROJECT_ISSUES_AND_EXPANSION.md:50–95`; earlier wording in `ISSUES.md:32–92`.

| ID | Original name and wording | Source | Type | Provenance | Later history / status |
|---|---|---|---|---|---|
| **T0.1** | **“My messages don't appear until restart.”** Symptom: a sent message does not show; restart proves it was saved. The document says test the JSON count first, then chat-switch rendering. | V2:50–63 | Defect / acceptance requirement | User-directed (source-derived) | `docs/review/ASSESSMENT-howl-application-reality-check-shared-foundations-2026-11-09.md:121–124` still records T0.1/T0.2 as open. No closure entry was found. Current code has a live user-message renderer and a duplicate-turn guard, but no same-session regression test. |
| **T0.2** | **“The model receives my message twice.”** Suspected separate from T0.1; likely prompt duplication even when saved once. Required evidence: a prompt dump or server log showing once per turn. | V2:65–69 | Defect / acceptance requirement | User-directed | Still recorded open at 2026-11-09. `recordUserTurn` now avoids a consecutive duplicate, but no prompt-dump test was found. |
| **T0.3** | **“Image model will not load.”** Black textless popup; no model load; no image-model-directory chooser. Three-part fix: validate stored paths at startup, add a chooser, give the dialog a real error state. | V2:71–81 | Defect + solution | User-directed; solution details are implementation-derived | Path preference/validation tests exist (`app-image-model-path_test.go:12–53`; `internal/image/engine/sdcpp/load_safety_test.go`). A full chooser/error-state acceptance was not found. |
| **T0.4** | **“Image preview is broken.”** The generated image does not preview; this blocks judging image-page layout (T2.3). Required evidence: prompt → image in the preview panel, correct aspect, no crop. | V2:83–89 | Defect / acceptance requirement | User-directed | `IMAGE_PAGE_WORKLOG.md:108–120` records a live check: 1024×1024 full-resolution image, correct MIME, displayed with `object-fit: contain`; the remaining softness was attributed to model output, not the preview path. That is a narrow live observation, not a permanent regression test. |

### A.2 Exact original list — Tier 1 (core mechanics)

Source: `PROJECT_ISSUES_AND_EXPANSION.md:96–143`.

| ID | Original name and wording | Source | Type | Provenance | Later history / status |
|---|---|---|---|---|---|
| **T1.1** | **“State is not owned properly.”** Tab changes must not unload the model or drop the connection; returning must preserve state. Samplers must persist visually across exit/reload. | V2:98–104 | Requirement + acceptance | User-directed | V2 ledger says `UNVERIFIED`; `STAGE_STATUS.md:102–140` records a 2026-10-31 live run. The later 2026-11-11 master plan still calls current ownership split, so the old live result is narrow historical evidence, not product-wide closure. |
| **T1.2** | **“VRAM arbitration between chat and image models.”** Check whether both fit; if not, a confirm dialog; unload, re-check freed VRAM, load; mirror the flow in reverse. | V2:106–115 | Requirement | User-directed | V2 ledger says `UNVERIFIED`; `STAGE_STATUS.md:142–186` records both directions with live VRAM readings and a stubbed harness confirm dialog. `internal/resource` still has no tests (`HFR-22`); the later master plan calls current resource authority incomplete. Historical live evidence and “no resource tests” are both true: behavior was observed, no automated regression exists. |
| **T1.3** | **“Think tags / reasoning handling.”** On/off/hidden with selectable intensity; hidden reasoning must not count against context; a thought must never be fed back next turn. | V2:117–124 | Requirement | User-directed | V2 ledger says `UNVERIFIED`; `STAGE_STATUS.md:187–243` records implementation and reasoning flags. `TestStripThinkBlocks` (`app_think_test.go:7`) and `TestThinkSplitter` (`app_think_stream_test.go:30`) cover stripping/splitting, not budget. No real token-budget test exists, so “does not count” is unverified. `T1_NOTES.md:40` says “OPEN” but is an earlier scratch note; the later STAGE_STATUS record supersedes it. |
| **T1.4** | **“Sampler defaults the user never has to touch.”** Family-specific defaults and correct prompt formatting; no UI card the backend ignores. | V2:126–132 | Requirement | User-directed | V2 ledger says `UNVERIFIED`; `STAGE_STATUS.md:244–286` records two families with measured sampler/template differences and no ignored UI cards. The later master plan calls current model/template evidence partial. `internal/modelintel` tests cover classification/caching, not model quality. |
| **T1.5** | **“Header says the wrong name.”** “Howl Chat - Tron Edition” → “Howl Chat”; bundle with T1.1. | V2:134–137 | One-line defect | User-directed | Current `frontend/dist/index.html:20` title is `HOWL Chat`; other pages use `HOWL Chat` variants. No dedicated test; source observation only. |

### A.3 Exact original list — Tier 2 (UI)

Source: `PROJECT_ISSUES_AND_EXPANSION.md:145–219`.

| ID | Original name and wording | Source | Type | Provenance | Later history / status |
|---|---|---|---|---|---|
| **T2.1** | **“Theme unification.”** Unify dark/light palettes across every page; fix light-theme overrides; verify system-light default. Required evidence: side-by-side screenshots of each page in both themes. | V2:147–156 | Requirement | User-directed | `STAGE_STATUS.md:426–535` records Stage A/B/C work; Stage C (723–790) measured five pages, not all nine, with a luminance heuristic and visual checks. It explicitly limits the gate (16 px+ text; not WCAG). Character/world/image pages were not in the quoted Stage C set. **Partial E2, not full acceptance.** |
| **T2.2** | **“Ghost elements and stale markup.”** Sweep every element lookup against the DOM; zero write-to-missing-ID, zero duplicate IDs, no orphaned markup. | V2:158–164 | Requirement | User-directed | `STAGE_STATUS.md:988–1100` records settings-surface and duplicate-function repairs; `STAGE_STATUS.md:1368–1487` records **18 never-created IDs over 30 reference sites**. **T2.2 remains OPEN**; some individual ghosts were fixed (e.g. `ig-history-status`). |
| **T2.3** | **“Image page layout — target Invoke AI organisation.”** Preview prominent, controls grouped, history organised; reference Invoke AI but do not chase parity; depends on T0.4. | V2:166–170 | Addition / UI requirement | User-directed; “Invoke AI” is a proposed reference | Image page repair worklog records functional repairs, not an acceptance against this layout requirement. The later master plan says T0.3/T0.4 and runtime acceptance remain unresolved. **Not verified.** |
| **T2.4** | **“Speaker identity in chat (thumbnail or initial).”** Thumbnail beside every speaker; fallback enlarged first initial on dark blue; fixed size; resized thumbnail, not full card image; prerequisite for multi-character legibility. | V2:172–192 | Addition / UI requirement | User-directed; detailed spec is implementation-derived | `chat-render.js:62–110` includes a character thumbnail when an image resolves in `addAssistantMessage`; without an image there is no avatar at all, and the later `updateBubbleWithCharacterInfo` path (lines 150–169) uses a generic 👤 placeholder, not the specified enlarged initial. No dedicated T2.4 test. **Partial E1.** |
| **T2.5** | **“Character / world / scenario image sets.”** Cards can hold several images with one default; worlds/scenarios likewise; find whether it was scrapped, overwritten, or broke. | V2:194–219 | Regression / feature restoration | User-directed | `character-images.js`, `scenario-form.js`, `context-panels.js`, and `character-configuration-integrity.test.js:305–341` show multi-image storage and preservation; the test does not assert thumbnail-default selection. Character tests exist; world/scenario coverage is partial. **E2 on the storage slice, no full restoration audit.** |

### A.4 Exact original list — Tier 3 (additions)

Source: `PROJECT_ISSUES_AND_EXPANSION.md:221–331`.

| ID | Original name and wording | Source | Type | Provenance | Later history / status |
|---|---|---|---|---|---|
| **T3.1** | **“Context handling — verify before optimising.”** First question: was the context flush ever implemented? If so, measure prefill; candidate `cache_prompt: true` + `--cache-reuse N`; full unloads unacceptable. | V2:225–234 | Verification + optional optimization | User-directed; cache proposal is a candidate solution | `truncateMessagesToFitContext` (`app-prompt.go:594–668`) uses **4 characters per token**, not a tokenizer; `TrimToMaxTokens` is a placeholder (`internal/backend/types/context.go:168`). No `cache_prompt`/`--cache-reuse` implementation found. **E0/E1; question still open.** |
| **T3.2** | **“Load speed — learn from LM Studio.”** Find why LM Studio loads faster and adapt what is portable. | V2:236–239 | Research / optimization | User-directed | No benchmark harness, baseline, or measured comparison found; `internal/bench` is an invariant/stress scanner, not a performance benchmark (`HFR-31`). **E0.** |
| **T3.3** | **“Scoped lorebook + memory system.”** Lore separable by world/scenario/character; primary/secondary keys; word-boundary matching; scan depth 2–3; recursion capped at 2 with no-recursion flag; fixed token budget with priority/insertion order; precedence character → scenario → world; dedupe; lore budget ~15–20%; memory per character, tied to who said it, natural-language triggers, hard context guard; AI-authored lore with specificity/frequency checks. | V2:241–266 | Addition; bundled lore + memory requirement | User-directed; mechanics are implementation-derived | Split below into T3.3-L (lorebook), T3.3-M (legacy memory), T3.3-A (AI-authored lore) to avoid false accounting. The original document treats them as one system; later decisions separate legacy lorebook/memory from MEMORY-09. |
| **T3.4** | **“Per-field AI character sheets.”** Every field gets its own generate button; model sees filled fields and returns only the target field; draft/accept/regenerate/lock per field; grammar-constrained JSON; per-kind schemas; no preset voices; exemplar dialogue. | V2:268–276 | Addition | User-directed; detailed design is implementation-derived | No field-generation binding, grammar-constrained JSON path, or per-field accept/lock UI was found in the inspected frontend/backend. **E0 in inspected scope.** |
| **T3.5** | **“Multi-character chat.”** 2+ characters fully separate, each with the appropriate system prompt, context flushed between speakers; root cause of voice bleed is flat context; needs structured transcript + per-character state. | V2:278–288 | Addition | User-directed; root-cause statement is an inference in the source | `character-queue.js` and `streamMultiCharacterResponses` (`chat-stream.js:244–320`) implement FIFO rotation and per-character calls. No test references `CharacterQueue` or `streamMultiCharacterResponses` (`HFR-17`); no contamination/flush test. **E1 orchestration.** |
| **T3.6** | **“Ternary models as first-class.”** Load ternary formats without new dependencies; must just work when the user points at the file. | V2:290–293 | Addition | User-directed | No ternary-specific loader path or test was found in the inspected Go/JS code. **E0 in inspected scope.** |
| **T3.7** | **“Audio and TTS end to end.”** Test audio model capability and TTS through the UI; needs headphones for verification. | V2:295–298 | Addition / verification | User-directed | Extensive audio packages/tests exist (ASR, TTS, manager, pipeline, helper). `qwen_asr_e2e_test.go` and TTS helper interop tests are present. No UI + headphone end-to-end acceptance was found; HFR records no real model run. **E2 components, E0 UI.** |
| **T3.8** | **“Audio generation.”** Newer small audio model, short clips fit 8 GB; play generated audio inside the app. | V2:300–303 | Addition | User-directed | `internal/audio` is ASR/TTS, not a music/audio-generation engine; no audio-generation UI was found. **E0 in inspected scope.** |
| **T3.9** | **“Diffusion version check.”** Confirm current Stable Diffusion release; update to latest patched version if one exists. | V2:305–308 | Research / maintenance | User-directed | No version-check implementation or recorded check was found; image code targets sd.cpp. **E0.** |
| **T3.10** | **“LLM writes the image prompt — research gated.”** Can a local LLM prompt the image generator well, at acceptable cost/quality/speed? Every part is a big if; research before building. | V2:310–315 | Research-gated addition | User-directed; flow is a proposal | No implementation or research result found in the inspected scope. **E0 / not authorized.** |
| **T3.11** | **“Compatibility for everyone else.”** Auto context/layering currently assume this machine; generalise as far as possible, ship beta, fix from reports. | V2:317–321 | Requirement / beta scope | User-directed | Hardware detection has NVIDIA/AMD/sysfs fallbacks; `FUT-HW1` remains deferred. No compatibility matrix/acceptance. **E0/E1.** |
| **T3.12** | **“Whole-app optimisation pass.”** Separate stage after everything above; look for waste everywhere. | V2:323–326 | Optimization stage | User-directed | No benchmark harness or whole-app measurement found. **E0.** |

### A.5 Tier 4 and expansion concepts — not active scope

**HFS-08 (Observed).** Tier 4 (`V2:333–344`) contains:

- **Reactive 3D character** — parked; enormous VRAM cost; likely incompatible with multi-character chat.
- **Multi-character with avatars** — sub-case of the same resource problem.

E.1–E.8 (`V2:346–527`) are explicitly “parked until foundations are solid”: model-aware image anchoring, multi-image vision, local-chat references, friction points, “fake strengths,” custom workspace layout, curtain vault, image staging boards. These are **not** active requirements and are not counted in the 27.

### A.6 F2/F3/W — engineering-led foundation work, not the original product list

**HFS-09 (Observed).** F2/F3/W are recorded in the same working document but are a different kind of object: owner-approved engineering requirements, architecture decisions, cutover work, and defect repairs. They must not be merged into the 27 or into MEMORY-09.

#### F2 — storage, content library, model intelligence

| ID | Original requirement / decision | Source | Status recorded |
|---|---|---|---|
| F2-A1 | Multiple user-controlled directories per content category; do not force moves | V2:568 | APPROVED; backend/settings implemented; browsers still need all roots |
| F2-A2 | Explicit read-only / read-write / save-destination roles | V2:569 | APPROVED; registry roles implemented; old save/import paths still need routing |
| F2-A3 | External library visibility without forced migration; show source, editability, destination | V2:570 | APPROVED; character discovery/source metadata implemented; edit/delete behavior incomplete |
| F2-A4 | SillyTavern character-library compatibility | V2:571 | APPROVED; JSON/PNG read-only discovery; import/export/editable operations future |
| F2-A5 | Documents/Pictures/Videos defaults through Windows known folders | V2:572 | APPROVED; implemented/unit-tested; redirected/OneDrive runtime confirmation required |
| F2-A6 | Separate comprehensive cleanup; Reset Settings is settings-only | V2:573 | APPROVED; cleanup separately named/previewed/confirmed |
| F2-B1 | Unified model directory registry across roots | V2:579 | APPROVED; image/model registry retained; single model-root manager UI/cross-kind listing incomplete |
| F2-B2 | Unified model identification across chat/image/audio/TTS | V2:580 | APPROVED; `internal/modelintel` exists; audio execution adapters/tokenizer/template inspection incomplete |
| F2-B3 | Cached identification with lightweight revalidation; no huge-file hashing | V2:581 | APPROVED; size+modtime cache; runtime/packaged behavior unverified |
| F2-B4 | Evidence-based provenance/template detection; no marketing-name facts | V2:582 | APPROVED; Verified/Source-derived/Inferred/Unknown/Conflicting model; tokenizer/template/online source inspection absent |
| F2-B5 | Separate modality workspaces; shared catalog does not merge UIs | V2:583 | APPROVED; capability classification; UI filtering incomplete |
| F2-B6 | Shared backend model intelligence/resource coordination | V2:584 | APPROVED as direction; metadata cache exists; resource consolidation later |
| F2-B7 | Music/audio expansion deferred; TTS contract frozen | V2:585 | DEFERRED / FROZEN |
| F2-C1–C8 | One authoritative owner; JSON mirrors untrusted; localStorage cache; chat folder format; media sidecars; create-only Save; explicit migration; disconnected roots retained | V2:625–644 | **DESIGN RECORDED**, no cutover/migration by F2-C |

#### F3 — ownership safety and lorebook cutover

| ID | Original requirement / decision | Source | Status recorded |
|---|---|---|---|
| F3-A1 | Remove automatic mirror writes from character/world/scenario/lorebook saves; leave legacy JSON untouched | V2:650 | IMPLEMENTED |
| F3-A2 | `ExportContentItem` only content-library write; structured result; non-authoritative envelope | V2:651 | IMPLEMENTED |
| F3-A3 | Stable export identity, safe segments, collision detection; never replace on normal export | V2:652 | IMPLEMENTED |
| F3-A4 | Replacement requires `RoleWrite` + explicit `replace=true`; `RoleSave` create-only | V2:653 | IMPLEMENTED |
| F3-A5 | Last-known character inventory for unavailable roots; stale not deleted/live | V2:654 | IMPLEMENTED |
| F3-A6 | Gallery containment incl. symlink/reparse checks; atomic image+sidecar commit | V2:655 | IMPLEMENTED |
| F3-A7 | Wipe uses configured roots, preflights, preserves external roots, reports failures | V2:656 | IMPLEMENTED |
| F3-A8 | Deletion/rename propagation deferred until authoritative ownership moves | V2:657 | DEFERRED |
| F3-A-R1-1–8 | Unsafe-target rejection, backup/restore, export identity, wipe staging/rollback, gallery containment; limits recorded | V2:667–685 | IMPLEMENTED / LIMITATIONS; final suite passed; audio flake unresolved |
| F3-B | Lorebook ownership cutover: canonical JSON file store, explicit import, read-back adapter, legacy store retained | V2:1214–1425 | Implemented/accepted in the record; native-window acceptance pending per registry note |
| F3-C | Store foundation, explicit import, LF-1 migration, chat consumer, single-writer audit, native A9/A10/A12 | V2:1426–1762 | Implemented with native evidence recorded; not a product-wide acceptance |

#### W13–W28 — ownership/editor and prompt-boundary repairs

Source: `PROJECT_ISSUES_AND_EXPANSION.md:1872–1931`; `docs/review/ENGINEERING-ISOLATED-NATIVE-RUNTIME-2026-11-10.md` §13–18.

| ID | Recorded defect / change | Recorded status | Evidence limit |
|---|---|---|---|
| W13 | Chat snapshot re-sent after canonical edit; `SetLorebooks` overlays canonical content | FIXED IN CODE — unit-tested, not natively verified | No model loaded |
| W14 | Post-cutover ownership boundaries; deleted entries dropped; disconnected roots retained; new lorebook can be enabled | FIXED IN CODE — unit-tested, not natively accepted | — |
| W15 | Private character lore leaked because `visibility`/unowned owner were not enforced | FIXED IN CODE — unit-tested, not natively accepted | — |
| W16 | Disconnected roots marked `Stale`, excluded from prompt | FIXED IN CODE — unit-tested, not natively accepted | — |
| W17 | No inference seam; added narrow `inferenceClient` and recording client | IMPLEMENTED + PROMPT-BOUNDARY TESTED IN-PROCESS | Not native GUI, not model inference |
| W18 | Character turns send frontend-composed prompt; `resolveCharacterLorePromptBlock` had no production call site | **OPEN** at that record | Superseded by W19/W20/W21 for code/native arming, but see W22–W25 |
| W19 | Frontend character composer calls bound `BuildCharacterLoreBlock` | FIXED IN CODE — binding unit-tested; native UI not exercised | — |
| W20 | Test-only `promptcapture` recorder; production compile-time no-op; composer fails loudly if resolver unavailable | IMPLEMENTED + NATIVE ARMING VERIFIED | GUI chat/character prompt exercise pending |
| W21 | Native ordinary-chat capture: 2 messages, probe present, forbidden marker absent, persisted | NATIVE ORDINARY-CHAT CAPTURE VERIFIED | Character/lorebook scenario pending |
| W22 | Canonical Lorebooks editor never wrote `ownerId`; replace stripped owner | **OPEN** — blocks native character acceptance | Needed owner control + preserve-on-replace |
| W23 | Character-scope Owner control distinct from Linked Characters; preserve omitted owner; honor explicit set/clear | FIXED IN CODE — store contract unit-tested; UI round-trip/native pending | — |
| W24 | Owner selector read on-disk characters while chat reads `CharacterStorage`; duplicate names/renames unresolved | FIXED IN CODE — identity-source mismatch; native pending | — |
| W25 | Native editor round-trip attempted; owner selector lists characters; coordinate drive could not complete second rename | **INCOMPLETE** | Needs element-addressable drive/UI harness |
| W26 | Stable-ownership proposal: reuse `character.id`; journaled name→ID migration | PHASE 1/2 COMPLETE — awaiting product decisions | Migration authority/legacy chats/ownerName/memory scope |
| W27 | Stable character ownership implemented; native Wails acceptance pass | NATIVE WAILS ACCEPTANCE PASS | No model inference performed |
| W27a | Cross-owner conflict collapse fixed; same-trigger different owners separated | FIXED + REGRESSION TESTS PASS | Native Aria/Bram prompt blocks |
| W28 | Memory identity containment; `owner_kind`; stable-ID capture; export/import same character ID; fail closed | GO TESTS + FRONTEND TEST + BUILDS PASS | Native memory checks via W28b seam; legacy ownership repair still a product decision |
| W28a/b/c | Export/import classification, native test seam, security audit | Tests pass; 7/7 native seam checks; security tests pass | Seam is test-only; snapshot is load-level, not live-model prompt evidence |

**HFS-10 (Calculated — W18–W23 are not a clean closure chain).** W18 was corrected by W19, but W22 exposed an authoring gap that reopened the native character path; W23 fixed the store contract, W25 could not complete the native UI round trip, W26/W27 eventually established stable IDs and native acceptance. A reader must not cite W18 as “open” without W19–W27, nor cite W23 as “native verified” — it was code-level only.

### A.7 Dates, historical status, and dependencies

**HFS-10a (Observed/Source-derived).** The T-item documents are working records, not dated specifications. The dates below come from the execution records; where no date exists, that is stated. Dependency arrows are **this audit’s calculated reading** of the original wording, not new requirements.

| Items | Date / historical status | Dependencies and later clarification |
|---|---|---|
| **T0.1–T0.4** | V2 marks all four `UNVERIFIED`; the 2026-11-09 assessment still lists T0.1/T0.2 open. T0.4 has a narrow 2026-11-04 live preview check, but no current regression. | T0.3 → T0.4 (preview cannot be judged without a loaded image model). T2.3 depends on T0.4. No later formal clarification found. |
| **T0.5** | V1 only; observed 2026-11-01; V2 drops it; `docs/asr-repair-design.md:10` still says open. | Independent of T0.1–T0.4; depends on process ownership (`internal/proc`). No closure record. |
| **T1.1** | `STAGE_STATUS` closed it 2026-10-31 on a live host; V2 ledger still says `UNVERIFIED`; current tree not rerun. | Foundation for T3.1/T3.5 state assumptions. HSL-1–HSL-11 (2026-11-19) later clarifies state lifecycle, much more broadly than T1.1. |
| **T1.2** | `STAGE_STATUS` closed it 2026-10-31 with live VRAM readings; V2 ledger `UNVERIFIED`; no automated resource tests now. | Depends on chat and image loaders; overlaps FUT-HW1 (lease authority). |
| **T1.3** | `STAGE_STATUS` closed it 2026-10-31; V2 ledger `UNVERIFIED`; no real token-budget test. | Depends on model reasoning/template support; T3.1 and T3.3 token budgets depend on it conceptually. |
| **T1.4** | `STAGE_STATUS` closed it 2026-10-31 with two families; V2 ledger `UNVERIFIED`; no full matrix. | Independent; modelintel later extends model identification. |
| **T1.5** | V2 `UNVERIFIED`; source title now `HOWL Chat`. | Independent one-line fix; no dependency. |
| **T2.1** | Stage A/B/C 2026-10-31–11-01; five pages measured; partial. | Cross-cutting; no hard dependency, but must be rechecked after UI changes. |
| **T2.2** | Inventory 2026-11-02; **OPEN** (18 ghosts/30 sites); later ASR note confirms still open. | Cross-cutting; T2.1/T2.3/T3.4 all add DOM surface and can add ghosts. |
| **T2.3** | No closure date; image-page repairs are separate. | **Depends on T0.4.** |
| **T2.4** | V2 addition; no dedicated test or closure. | Depends on character image data (T2.5); T3.5 needs speaker identity. |
| **T2.5** | V2 addition; character slice implemented, world/scenario partial. | Independent data/UI; T2.4 reads its thumbnail. |
| **T3.1** | V2 addition; question still open; context path is character-based. | Depends on T1.3 (hidden reasoning) and T3.3 (lore budget); HSL later adds state/branch context. |
| **T3.2** | V2 addition; no baseline/harness. | Independent research; T3.12 depends on its harness. |
| **T3.3-L** | V2 addition; resolver/store substantial; no full acceptance. | Depends on T2.2 (UI IDs), T3.1 (budget); F3-B/C later cut over storage; MEMORY-09 split later separates legacy lore. |
| **T3.3-M** | V2 addition; legacy memory exists; not MEMORY-09. | Depends on T3.3-L conversion and T1.1 state; HSL/PCB later expose identity/authority gaps. |
| **T3.3-A** | V2 addition; no implementation found. | Depends on model generation + T3.3-L; not authorized. |
| **T3.4** | V2 addition; no implementation found. | Depends on model/grammar support and T2.2 UI stability; not authorized. |
| **T3.5** | V2 addition; queue exists; no isolation test. | Depends on T1.1/T3.1/T3.3-M and T2.4; HSL later adds state/transition requirements. |
| **T3.6** | V2 addition; no implementation found. | Independent loader/model-format work; not authorized. |
| **T3.7** | V2 addition; packages/tests exist; no UI end-to-end. | Depends on audio manager and UI; T3.8 follows. |
| **T3.8** | V2 addition; no implementation found. | Depends on T3.7 audio infrastructure; not authorized. |
| **T3.9** | V2 addition; no version check found. | Independent image-stack maintenance; T3.10 depends on it. |
| **T3.10** | V2 addition; research-gated; no implementation. | Depends on T3.9 and T3.2; not authorized. |
| **T3.11** | V2 addition; partial detection; no matrix. | Cross-cutting; overlaps FUT-HW1 and T1.2. |
| **T3.12** | V2 addition; separate final stage; no harness. | After T3.2 and the other items; no implementation. |

**HFS-10b (Calculated — V1/V2 reconciliation).** V1 has **25** active items (T0.1–T0.5, T1.1–T1.5, T2.1–T2.3, T3.1–T3.12). V2 displays **26**: it adds T2.4/T2.5 and removes T0.5 (25 + 2 − 1 = 26). The union is **27**. V2’s removal of T0.5 has no recorded decision. This audit treats V2 as authoritative for T2.4/T2.5 and the F2/F3 appendices, but **not** as authority to erase T0.5.

---

## B. Evidence matrix — 27 recovered original items

**HFS-11a (Observed — V2’s omission).** This matrix includes **T0.5** because V1 defines it and a later document still calls it open. It is not an addition by this audit; it is a recovered item that V2 dropped without explanation.

**HFS-11 (Calculated).** Evidence levels are those defined in `HFR-CP-2026-11-21-01` §2: **E0** not established, **E1** implemented, **E2** working in a limited test, **E3** end-to-end verified. A historic live observation can be E2/E3 for that date, but does not become a current automated claim.

| Item | Original scope / source | Design status | Implementation evidence | Test evidence | End-to-end | Remaining work | Class / confidence |
|---|---|---|---|---|---|---|---|
| **T0.1** | Send must render same-session; V2:50–63 | No formal design; symptom→test | `addUserMessage` (`chat-render.js`); `recordUserTurn` (`app-chat.go:360`); persistence append | `chat-render-ownership.test.js` covers assistant identity, not user render; no same-session test | **No** | Reproduce JSON count + chat-switch render; add regression if still broken | Source-derived / **E0–E1**, medium |
| **T0.2** | Prompt must contain one user turn; V2:65–69 | Duplicate-turn guard in code | `recordUserTurn` skips consecutive identical user turn; `buildPromptMessages` includes list once | `TestBuildPromptMessagesOrder` checks order, not duplication; no prompt dump | **No** | Prompt-dump / server-log test; check frontend+backend double append | Observed code / **E1**, medium |
| **T0.3** | Image model must load, chooser/error state; V2:71–81 | Partial design in issue | `preferModelPath`, checkpoint validation, model discovery; image-directory chooser exists in `image-model-dirs.js` | `TestPreferModelPath*`; `TestValidateCheckpoint*`; no invalid-path dialog run | **No** | Exercise invalid stored path, chooser, visible error, VRAM climb | Observed partial / **E1**, medium |
| **T0.4** | Generated image preview; V2:83–89 | No formal design | Image generation/canvas/gallery path | `IMAGE_PAGE_WORKLOG.md:108–120` live check: full-res preview, correct MIME; not a test | **Narrow live observation** | Keep regression/visual evidence; distinguish generation softness | Observed / **E2 (old), current unverified** |
| **T0.5** | Concurrent instances must not kill each other’s children; `ISSUES.md:69–86` (V1; dropped from V2 without rationale; still called open at `docs/asr-repair-design.md:10`) | Ownership audit required; job-object approach | `internal/proc/job_windows.go` attach/resume; `D3` ownership | `TestSuspendedChildRunsOnlyAfterOwnership`; `TestStartFailsClosedWhenJobUnavailable` | **No** two-instance test | Two-instance acceptance; ensure `ReapOrphans` does not kill by path | Observed / **E2 primitive, E0 scenario** |
| **T1.1** | Model/state survives tab switch; samplers persist; V2:98–104 | Stage design | App/JS rehydration; guarded emits; settings store | `STAGE_STATUS.md:102–140` live run: same PID/VRAM, sampler restart | **Historical yes on host** | Rerun on current build; tab-switch + settings persistence regression | Observed historical / **E3 narrow** |
| **T1.2** | Chat/image VRAM arbitration; V2:106–115 | `CheckModelFit`/`FreeVRAMFor` | `app-vram-arbitration.go`; `internal/resource` advisory | `STAGE_STATUS.md:142–186` both directions; no automated test; `internal/resource` no tests | **Historical yes on host** | Automated fault/lease tests; current behavior after resource changes | Observed historical / **E2/E3 narrow; current E1** |
| **T1.3** | Think modes; hidden excluded; never fed back; V2:117–124 | Reasoning flags + stripping | `app-llama.go` flags; `http_client` reasoning sink; `app-think.go`; UI card | `TestStripThinkBlocks`; `TestThinkSplitter`; live reasoning card in STAGE_STATUS | **No two-turn measured prompt** | Real prompt/token test that hidden thought does not reappear or consume budget | Observed / **E2** |
| **T1.4** | Family defaults/formats; V2:126–132 | Family table + template resolution | `app-family-defaults.go`; `app-model-template.go`; `modelintel` | `STAGE_STATUS.md:244–286` two families measured; no automated template matrix | **Historical yes on host** | Real model matrix; regression for sticky templates | Observed / **E2/E3 narrow** |
| **T1.5** | Header name; V2:134–137 | Trivial | `frontend/dist/index.html:20` title `HOWL Chat`; other pages variants. Wordmark differs in capitalization from the literal `Howl Chat`, but the Tron suffix is gone. | None | **Source only** | Confirm the intended capitalization/header across pages | Observed / **E1** |
| **T2.1** | Theme unification all pages; V2:147–156 | Token vocabulary + light/dark | `theme.css`, `theme-light.css`, settings theme path | `STAGE_STATUS.md:723–790` five pages, luminance heuristic, visual buttons | **Partial** | Character/world/image pages; persisted per-page screenshots; contrast limits | Observed / **E2 partial** |
| **T2.2** | Zero ghost IDs; V2:158–164 | Sweep classified | Some fixes (`ig-history-status`, settings cleanup); 18 ghosts remain | `STAGE_STATUS.md:1368–1487` source inventory; no runtime test | **No** | Resolve 18 IDs/30 sites; zero duplicate IDs; decide manager modules | Observed / **E1; open** |
| **T2.3** | Invoke-AI-style image page; V2:166–170 | Reference only | Image page functional repairs, not layout acceptance | No acceptance test | **No** | Layout review against stated criteria; depends T0.4 | Source-derived / **E0–E1** |
| **T2.4** | Speaker thumbnail/initial; V2:172–192 | Detailed spec | `chat-render.js:62–110`: thumbnail if image; no avatar in `addAssistantMessage` when absent; later path uses generic 👤, not initial | No T2.4 test; e2e avatars in sidebar/scenario only | **No** | Initial fallback, fixed sizing, two-character screenshot test | Observed partial / **E1** |
| **T2.5** | Multi-image sets with default; V2:194–219 | Data/UI design | `character-images.js`; `scenario-form.js`; `context-panels.js`; arrays | `character-configuration-integrity.test.js:305–341`; e2e scenario image | **Partial** | World/scenario full round-trip; regression investigation | Observed / **E2 partial** |
| **T3.1** | Was flush implemented? Measure; cache reuse; V2:225–234 | No accepted design | `truncateMessagesToFitContext` char estimate; `TrimToMaxTokens` placeholder; no cache flags | No token/flush test; HFR-13/36 | **No** | Real tokenizer/budget; measure prefill; cache-reuse experiment | Observed / **E0–E1** |
| **T3.2** | Load-speed research; V2:236–239 | Research queue only | No benchmark harness; `internal/bench` invariants | HFR-31: no `Benchmark` function | **No** | Baseline load-time harness; compare LM Studio | Observed absence / **E0** |
| **T3.3-L** | Scoped lorebook mechanics; V2:241–266 | Lorebook resolver design | `internal/backend/lorebook`; canonical `internal/lorebookstore`; scopes/triggers/conflicts | `TestResolveExactPriorityAndScope`, `TestResolveAppliesCharacterBudget`, `TestResolveBudgetIsDeterministic`, `TestResolveSkipsHiddenVisibility`, conflict/visibility tests | **No full long-chat test** | Token-based lore budget (currently chars); scan depth/recursion parity; native multi-character | Observed / **E2** |
| **T3.3-M** | Per-character memory with trigger/context guard; V2:257–263 | Legacy memory design | `internal/backend/memory` CRUD/converter; frontend `MemoryStore`; localStorage + backend sync | `TestMemoryCaptureUsesStableIDAndFailsClosedOnAmbiguity`, `TestPromptOwnershipFilterRejectsLegacyAndForeignRecords`, `TestMemoryToLorebookEntryCarriesOwnerKind`; frontend memory tests | **No MEMORY-09 integration** | Frontend authority bypass; token guard; legacy ownership decision | Observed / **E2 narrow** |
| **T3.3-A** | AI-authored lore with specificity/frequency; V2:264–266 | No dedicated design found | No generator/specificity/frequency implementation found in inspected paths | None | **No** | Separate design/authorization; not part of HF1 | Source-derived absence / **E0** |
| **T3.4** | Per-field AI sheets; V2:268–276 | No design found | No field-generation binding/UI found | None | **No** | Design + authorization; not Phase 1 | Source-derived absence / **E0** |
| **T3.5** | Multi-character separation/flush; V2:278–288 | Queue design; no state design | `CharacterQueue`; `streamMultiCharacterResponses`; per-character prompts | No test references queue/stream (`HFR-17`); no contamination test | **No** | 2–3 character isolation, flush, abort/recovery, restart | Observed / **E1** |
| **T3.6** | Ternary first-class; V2:290–293 | No design found | No ternary path found | None | **No** | Research/implement later; not Phase 1 | Source-derived absence / **E0** |
| **T3.7** | Audio/TTS end to end; V2:295–298 | Audio architecture docs | `internal/audio` ASR/TTS/manager/pipeline; UI bindings | Many package tests incl. helper interop; `qwen_asr_e2e_test.go`; no UI/headphone run | **No** | Real UI + hardware check; distinguish package from product E3 | Observed / **E2 components, E0 UI** |
| **T3.8** | Audio generation; V2:300–303 | No design found | No generation engine/UI found | None | **No** | Research/authorize separately | Source-derived absence / **E0** |
| **T3.9** | Diffusion version check; V2:305–308 | Research only | No version-check code found | None | **No** | Verify current sd.cpp release; record update decision | Source-derived absence / **E0** |
| **T3.10** | LLM image prompt research; V2:310–315 | Research-gated | No implementation found | None | **No** | Research only; not authorized | Source-derived absence / **E0** |
| **T3.11** | Broad compatibility; V2:317–321 | FUT-HW1/architecture | Hardware detection fallbacks; no compatibility matrix | Audio hardware tests only | **No** | Compatibility matrix and beta evidence | Observed partial / **E0–E1** |
| **T3.12** | Whole-app optimization; V2:323–326 | Separate stage | No benchmark harness | None | **No** | Benchmark first; then optimize | Source-derived absence / **E0** |

### B.1 Cross-cutting findings

**HFS-12 (Observed — state identity is positional).** `appendChatMessage` writes `ID: fmt.Sprintf("%d", len(a.chatMessages)-1)` (`app-prompt.go:281–317`); `buildChatHistory` writes `len(history)` (`app-chat.go:373–418`); `ReplaceChatHistory` writes `i` (`app-prompt.go:327–389`). `ChatHistory` has an `ID` field, but it is not stable across edits, deletions, branches or reorder.

**HFS-13 (Observed — token budgeting is a character estimate).** `truncateMessagesToFitContext` (`app-prompt.go:594–668`) uses `contextSize * 4` characters and calls that a token fit. `TrimToMaxTokens` is explicitly a placeholder. `Message.Tokens` is initialized to `0` and never populated by a tokenizer (`internal/backend/types/message.go:84`). Lore `ContextBudget`/`MaxCharacters` is also measured in characters (`internal/backend/lorebook/lorebook.go:119–187, 429–450`; `app-prompt.go:881–897`). The UI tutorial calls the lore budget “tokens”; that wording is not implemented.

**HFS-14 (Observed — the tokenizer seam exists below the app boundary).** `llama.HTTPClient.Tokenize` calls llama-server `POST /tokenize` (`internal/backend/llama/http_client.go`), but `App.inferenceClient` (`app.go:37–46`) exposes only generation, context size, model info, health and loaded state — **not** `Tokenize`. Phase 1’s finding is correct: a real token-count path needs an interface extension or a tokenizer seam, and the test fakes must support it.

**HFS-15 (Observed — multi-character has no test).** `CharacterQueue` and `streamMultiCharacterResponses` are implemented but no test file references them; `HFR-17` records this. The queue passes prior responses as text and the character prompt reads frontend `MemoryStore`, so isolation is prompt-level, not backend-enforced.

**HFS-16 (Observed — legacy memory is not MEMORY-09).** `internal/backend/memory` is explicitly “automated lorebook entries captured during chat sessions” and converts records to `lorebook.Entry`; `frontend/dist/js/memory-store.js` stores `howl-chat-memories` in localStorage and syncs to backend `AddMemory`/`EditMemory`. `eval/mem1/m09` is a separate 1,951-line engine with no caller outside its package (`HFR-09`). They must not be added together.

---

## C. Historical milestone reconciliation

**HFS-17 (Observed — timeline).**

| Date | Milestone | What it establishes | Evidence limit |
|---|---|---|---|
| 2026-10-31 | T1.1 state ownership closed | Live model survives page switch; samplers persist | One host; harness stream limitation; not rerun |
| 2026-10-31 | T1.2 VRAM arbitration closed | Both directions live with VRAM readings | Confirm dialog stubbed in harness; no automated resource tests |
| 2026-10-31 | T1.3 think tags closed | Reasoning flags, stripping, UI card | No token-budget measurement; `T1_NOTES.md` stale “OPEN” |
| 2026-10-31 | T1.4 family defaults closed | Two families measured with different templates/samplers | No full model matrix; UI effective values not visually rechecked |
| 2026-11-01 | T2.1 Stage C | Five pages measured in both themes | 16 px+ heuristic; not WCAG; not all pages; screenshots not all persisted |
| 2026-11-02 | T2.2 ghost inventory | 18 never-created IDs / 30 reference sites | Source-only; 16 template lookups unclassified; open |
| 2026-11-04 | Audio mic/history tranche | PCM WAV capture, temp cleanup, composer target | No live Wails mic; no concurrent chat+ASR; no real Whisper |
| 2026-11-09 | F2/F3-A/FUT-HW1/FUT-PROMPT1 recorded | Storage/model ownership; deferred hardware/prompt requirements | Status words are project records, not fresh tests |
| 2026-11-10 | F3-B/F3-C, W13–W28 | Canonical lorebook store/import/chat consumer; stable ownership; prompt-boundary capture | Native acceptance varies by W item; W22 open, W25 incomplete |
| 2026-11-11 | Reconciled architecture/amendments | DWP design packages, single-owner boundaries; OD-01–OD-04 | Design only; no implementation authorization |
| 2026-11-12/13 | MEMORY-09 research/pre-registration | New memory architecture separated from legacy lorebook; G1 only | No G2; no production integration |
| 2026-11-19 | HSL requirement | Howl owns state lifecycle; MEMORY-09 must not infer transitions | Requirement only; no implementation; no acceptance |
| 2026-11-21 | HFR audit | A–E readiness snapshot; no E3; context placeholder | Focused tests; no model/GPU; omits original UI/release scope |
| 2026-11-21 | HF1 plan | Proposed chat state/context package | Not authorized; tokenizer seam unresolved |

**HFS-18 (Calculated — do not collapse historical claims).** The 2026-10-31 records are real evidence for their scenarios, but they do not cover the current tree. The later 2026-11-11 master plan reclassifies T1.1–T1.4’s ledger state as `UNVERIFIED` and says current ownership is still split; that is a status correction, not a refutation of the live observations. Conversely, HFR’s “no resource tests” does not erase the historical live VRAM evidence; it means the current repository lacks a regression harness. All three statements belong in the record.

**HFS-19 (Observed — W18–W23 chronology).** W18 identified frontend-composed character prompts; W19 moved lore resolution to `BuildCharacterLoreBlock`; W20 added a test-only inference recorder; W21 verified ordinary-chat native capture; W22 found the missing `ownerId` authoring path; W23 implemented owner control/preserve-on-replace at the store contract. Native character acceptance remained blocked through W24/W25 and was only later addressed by W26/W27. A document that cites W18 as the final state or W23 as native-verified is reading the chain incorrectly.

---

## D. Readiness audit reconciliation (`HFR-CP-2026-11-21-01`)

### D.1 Mapping A–E to the original list

**HFS-20 (Calculated — HFR categories are audit constructs).**

| HFR category | Original items it actually supports | Relevance | What HFR added or changed |
|---|---|---|---|
| **A. Chat/state machine** | T1.1 state ownership; T3.1 context; T3.5 structured transcript | **Partially relevant** | HFR’s “state machine” is largely HSL (2026-11-19), a later owner requirement, not an original Tier item. T1.1 was about tab/state persistence, not branches/revert. Memory authority appears only as a cross-cutting note (HFR-38), not as legacy-memory assessment. |
| **B. Multi-character orchestration** | T2.4 speaker identity; T3.3-M memory; T3.5 multi-character | **Relevant** | HFR correctly finds no orchestration test. It does not treat T2.4/T2.5 UI or T3.3-L lorebook scope as the original scope. |
| **C. Hardware-aware intelligence** | T1.2 VRAM arbitration; T3.11 compatibility; FUT-HW1 | **Relevant but broader** | HFR’s lease/telemetry/OOM category is FUT-HW1 plus audit engineering, not a single original work item. T1.2 had historical live evidence. |
| **D. Inference engine / Kobold-like capabilities** | T1.3 reasoning; T1.4 sampler/template; T3.6 ternary; T3.7 audio | **Relevant but renamed** | “Kobold-like capabilities” is audit framing. T3.7 is audio, not LLM inference; T3.6 is a loader requirement. |
| **E. Performance** | T3.2 load speed; T3.12 whole-app optimization | **Partially relevant** | E is a measurement dimension HFR introduced to test the foundation. It is not one of the recovered original items. |

### D.2 Original items omitted from HFR

**HFS-21 (Observed — omissions).**

HFR does not assess, or only mentions in passing:

- **T0.1–T0.4** release blockers (message render, duplicate prompt, image load, image preview). **T0.5 is worse:** V2 omits it, but V1 defines it and a later ASR design note still calls T0.5 open. HFR does not assess it.
- **T1.5** header name (trivial, but original).
- **T2.1–T2.3** theme, ghost markup, image-page layout.
- **T2.4–T2.5** speaker identity and image sets (T2.4 is mentioned under B only as context; T2.5 is absent).
- **T3.3-A** AI-authored lore.
- **T3.4** per-field AI character sheets.
- **T3.6** ternary models.
- **T3.8** audio generation.
- **T3.9** diffusion version check.
- **T3.10** LLM image-prompt research.
- **F2** storage/content/model-intelligence requirements.
- **F3-A/F3-B/F3-C** ownership safety and lorebook cutover, except incidental references.
- **W13–W28** ownership/editor and prompt-boundary history.
- **Tier 4/E** parked ideas (correctly excluded from active readiness, but should be named as excluded).

**HFS-22 (Calculated — omission is not a finding of failure).** HFR was a readiness audit for A–E and MEMORY-09, not a full product audit. The omission matters only if HFR is treated as the authoritative scope or as proof that the original foundation is complete. It is neither.

### D.3 Readiness categories introduced by the audit

**HFS-23 (Observed).** The following are audit/design constructs, not original user wording:

- “A. Chat/state machine” as a five-part subsystem with branch/revert/lifecycle requirements.
- “C. Hardware-aware intelligence” as a lease/telemetry/OOM authority.
- “D. Inference engine and Kobold-like capabilities.”
- “E. Performance” as a foundation subsystem.
- E0–E3 evidence levels.
- The nine integration scenarios in HFR §10.
- The A–E “no subsystem has E3” verdict.

They are useful and should be retained for evidence, but they do not replace T0–T3.

### D.4 Valid results retained

**HFS-24 (Observed).**

- Focused Go tests passed for chat history/ownership, prompt order/lore privacy, think stripping, modelintel, llama props, GGUF admission, loader rejection, and process ownership.
- `internal/resource` has **no test files**.
- Frontend suite is **15 pass, 1 fail** at `tests/frontend/character-configuration-integrity.test.js:409`; the failure is separate and still open.
- No model, GPU, telemetry-loss, OOM, KV-cache, or real multi-character run occurred.
- `eval/mem1/m09` is 1,951 lines with no caller outside its package.
- These results remain valid at the recorded commit/tree; this audit did not rerun them.

### D.5 Finding-level HFR map

**HFS-24a (Calculated — every HFR ID is accounted for).** Ranges below cover all HFR-01–HFR-54. A range does not promote its findings; it records that the findings in it support the same original item(s).

| HFR findings | Original item(s) supported | Relevance | Limit / drift |
|---|---|---|---|
| **HFR-01–02** | All 27, cross-cutting | **Relevant as verdict** | HFR’s “A–E” verdict is not the original list; it says nothing about T0/T2/T3.4–T3.12. |
| **HFR-03–05** | None directly; evidence method | **Outside original scope** | E0–E3 levels are audit constructs, not user requirements. |
| **HFR-06–08** | Cross-cutting test record | **Relevant** | Focused tests only; no model/GPU; result is a snapshot. |
| **HFR-09** | MEMORY-09 ledger only | **Outside original Howl scope** | The 1,951-line engine is not legacy lorebook or original T3.3. |
| **HFR-10–15** | T1.1; T3.1; T3.3-M; T3.5 | **Partial** | HFR-10/11 support history/ownership; HFR-12/13/14 are HSL/later state requirements; HFR-15 says no MEMORY-09 dependency for the state machine. |
| **HFR-16–20** | T2.4; T3.3-M; T3.5 | **Partial** | HFR-17 directly supports “queue exists but no test”; HFR-19 supports the legacy-memory authority gap; no T2.4/T2.5 UI assessment. |
| **HFR-21–25** | T1.2; T3.11; FUT-HW1 | **Partial** | HFR-22 correctly says no resource tests; it does not erase the historical T1.2 live evidence. HFR’s lease authority is FUT-HW1, not original T1.2. |
| **HFR-26–30** | T1.3; T1.4; T3.6; T3.7 | **Partial** | HFR-27 supports parsing/admission only; ternary and audio are not tested. “Kobold-like” is audit framing. |
| **HFR-31–35** | T3.2; T3.12 | **Relevant to the measurement gap** | HFR-31’s no-benchmark finding supports both; E is not an original item. |
| **HFR-36–41** | T3.1; T3.3; T3.5; cross-cutting | **Directly relevant** | HFR-36/37 support the token/identity findings; HFR-38/39 support legacy/MEMORY-09 separation; HFR-41 is unrelated frontend health. |
| **HFR-42** | T1.1; T3.1; T3.3; T3.5; T3.11 | **Partial** | Nine integration scenarios are audit constructs; several presuppose HSL/MEMORY-09. |
| **HFR-43** | All; gate discipline | **Methodology** | “Evidence, not a meeting” is not a product requirement. |
| **HFR-44–45** | MEMORY-09 ledger; separation | **Outside original scope** | Correctly preserves settled decisions; does not assess T0–T3 implementation. |
| **HFR-46–47** | Cross-cutting independence | **Relevant** | “A–E can progress” is planning advice, not original scope. |
| **HFR-48–49** | Checkpoint and next evidence | **Relevant as recommendation** | Its order omits T0 release blockers and T2 UI work. |
| **HFR-50** | Invariants | **Methodology** | No-E3-without-scenario is sound; it is not user wording. |
| **HFR-51** | Resume instruction | **Methodology** | Re-run tests; no implementation authorization. |
| **HFR-52–54** | Limits and evidence paths | **Relevant** | Static audit plus focused tests; no runtime; paths support traceability. |

**HFS-24b (Calculated — HFR’s biggest scope drift).** HFR’s A–E map is useful for the specific capabilities it names, but it replaces the owner’s priority order with a subsystem taxonomy. It also treats HSL (a 2026-11-19 requirement) as if it were part of the original foundation. The original ordering law — broken, UI, additions — is not represented.

---

## E. Phase 1 reconciliation (`HF1-PLAN-2026-11-21-01`)

### E.1 Mapping Phase 1 objectives to original scope

**HFS-25 (Calculated).**

| Phase 1 objective | Original grounding | Status of grounding |
|---|---|---|
| Stable state/message/turn identities | T1.1 state ownership; T3.5 structured transcript; **HSL-3** (later clarification) | **Grounded but partly later.** Stable message identity is an implementation subrequirement of T1.1/T3.5; `state_id`/`turn_id` are HSL/Phase-1 design specifics, not original Tier wording. |
| Non-destructive branches/restoration | Original “Next”/variations, regenerate; **HSL-4/5/6** (later) | **Grounded in later owner requirement.** The original Tiers did not define a branch model; T3.5’s structured transcript and the existing Next/regenerate defects make it a necessary subrequirement. |
| World/scenario/character composition | REG-1 protected relationship behavior; T3.3 scoping; T3.5 per-character prompts | **Directly grounded.** Phase 1’s stable-ID composition is a necessary condition for REG-1 and T3.3/T3.5. |
| Real token counting and budgeting | T1.3 hidden-thought budget; T3.1 flush; T3.3 fixed token budget; FUT-CHTOK1/FUT-CH2 | **Directly grounded.** The original documents ask for token-aware budgets and warn against false precision. |
| Predictable flush between character responses | T3.1; T3.5 | **Directly grounded.** T3.5 explicitly requires context flushed between speakers. |
| Preservation across model changes | T1.1 state ownership; T3.2 load speed; T3.11 compatibility | **Grounded as a necessary condition**, not an explicit original acceptance test. |
| Long/multi-character tests | Ordering law/stop rule; T3.5 | **Grounded.** T3.5 has no test today; Phase 1’s tests are the missing evidence. |

### E.2 Criteria that are grounded, derived, or introduced

**HFS-26 (Calculated).**

- **Grounded:** stable IDs; non-destructive branch/revert; composition; real token counting; flush; model-change survival; long/multi-character tests.
- **Derived subrequirements:** a `turn_id`; a `current_state_id`; parent/alternative links; explicit deletion sets; a tokenizer seam; fail-closed “unavailable” budget state; recomputation after model change.
- **Introduced by the plan/audit:** E3 acceptance levels; AC-1–AC-8 as named scenarios; “no character-count estimate as authority”; the requirement that tokenizer unavailability be reported as unavailable. These are reasonable evidence controls, but they are not quotations from the original Tiers.

### E.3 Mismatches and gaps

**HFS-27 (Observed — Phase 1 is not T3.3).** Phase 1 covers world/scenario/character composition but explicitly excludes MEMORY-09 integration and does not implement scoped lorebook mechanics or legacy memory. Therefore:

- Passing Phase 1 would not complete T3.3.
- Passing Phase 1 would not establish that legacy lorebook/memory is ready for MEMORY-09.
- T3.3-L and T3.3-M remain separate work/evidence.

**HFS-28 (Observed — tokenizer seam).** Phase 1 correctly identifies that `llama.HTTPClient.Tokenize` exists but `App.inferenceClient` does not expose it. The plan must extend the interface or add a seam; its tests must use a fake tokenizer for unit paths and a real model tokenizer for E3 token-count evidence. A tokenizer is model-specific; persisted counts cannot be treated as portable.

**HFS-29 (Observed — tokenizer availability).** The plan’s AC-8 correctly requires an honest `unavailable` state. It must not be weakened later into a character-count fallback; `truncateMessagesToFitContext` is exactly that fallback today.

**HFS-30 (Calculated — Phase 1 status unchanged).** Phase 1 remains **proposed, not authorized**. This audit does not promote it, change its acceptance criteria, or begin implementation.

---

## F. Separate memory-system accounting

### F.1 Three ledgers, never summed

**HFS-31 (Observed — mandatory separation).**

| Ledger | What it is | Where it lives | Completion/evidence today |
|---|---|---|---|
| **1. Original Howl foundation** | T0–T3 product work, F2/F3/W engineering work | This audit §A–B | Mixed; no overall E3; many original items unverified |
| **2. Legacy lorebook / memory** | Older automatic lorebook entries, natural-language triggers, localStorage memory, canonical lorebook store | `internal/backend/lorebook`, `internal/backend/memory`, `frontend/dist/js/lorebooks.js`, `memory-store.js`, `internal/lorebookstore` | Substantial code and tests; **not** MEMORY-09 |
| **3. MEMORY-09** | New memory architecture/evaluation | `eval/mem1/m09`, MEMORY-09 research/decision records | Engine unconnected; G1 fixture only; G2 blocked; no production capture |

**HFS-32 (Calculated — do not compute a combined percentage).** A legacy lorebook test is not MEMORY-09 progress; a MEMORY-09 fixture is not a Howl feature. A single “memory completion” percentage would be false.

### F.2 Legacy lorebook / memory — what exists

**HFS-33 (Observed).**

- `internal/backend/lorebook` implements scoped entries, exact/loose/semantic matching, primary/secondary triggers, priorities, conflict rules, cooldown/frequency, visibility, stale exclusion, and character/world/scenario routing.
- Budgets are **entry counts** and **character counts**, not model tokens (`budget.go`; `applyBudget` in `lorebook.go`).
- `internal/lorebookstore` provides canonical files, explicit import, stable IDs, owner migration, recovery and rollback, with tests.
- `internal/backend/memory` provides CRUD for automated memory records, auto-capture config, identity, and conversion to lorebook entries. `frontend/dist/js/memory-store.js` uses localStorage `howl-chat-memories` and syncs to backend bindings.
- Tests cover lorebook resolution/budget/visibility/conflicts, store import/owner migration/rollback, memory identity/ownership, and frontend memory ownership.
- **Limits:** character prompt reads frontend `MemoryStore` (bypasses backend authority); unowned/name-owned legacy records fail closed; lore budget is not token-based; no full long-conversation test; no MEMORY-09 provenance/capture.

**HFS-33a (Calculated — legacy ledger).**

| Question | Legacy lorebook / memory | Evidence |
|---|---|---|
| What was intended? | Scoped lore by world/scenario/character; natural-language primary/secondary triggers; priority/conflict/budget rules; per-character memory tied to the speaker; hard context guard (T3.3). | V2:241–266 |
| What exists? | A resolver with scopes, trigger modes, priority, visibility, conflict rules, character/world routing; a canonical file store; automated-memory CRUD + conversion; a frontend MemoryStore. | `internal/backend/lorebook`; `internal/lorebookstore`; `internal/backend/memory`; `frontend/dist/js/memory-store.js` |
| What remains useful? | Resolver/selection semantics; canonical file store and import/rollback; editor/viewer UI; owner/privacy enforcement patterns. | Tests and code paths in HFS-33 |
| What is unsuitable for the future architecture? | localStorage as memory authority; display-name/positional identity; character-count budgets; trigger-only retrieval as a memory model; no state/branch lifecycle; no production capture envelope. | HFR-13, HFR-19, HFR-38; T3.3 audit |
| What UI is reusable independently? | Lorebooks editor, character-memory viewer, context panels, memory candidate menu, story-tools surfaces. | Frontend modules; no dependency on MEMORY-09 schema |

**HFS-34 (Calculated — reusable UI, not reusable authority).** The Lorebooks editor, character memory viewer, context panels, and memory candidate menu can be reused as UI components. Their storage/trigger/ownership semantics cannot be assumed to carry into MEMORY-09. In particular:

- localStorage is not a durable authority;
- display names are not stable identity;
- trigger-only retrieval is not MEMORY-09 capture/classification;
- per-scope entry counts are not a token budget;
- no state/branch lifecycle exists.

### F.3 MEMORY-09 — what exists

**HFS-35 (Observed).**

- `eval/mem1/m09/` is 1,951 lines and has no caller outside its package (`HFR-09`).
- The production capture boundary is proposed, not implemented (`PCB-01…42`).
- G1 fixture contract passed for its fixture scope; production wiring is unproven.
- G2 classification remains blocked; no corpus/G2 experiment.
- The unresolved-class lifecycle remains blocked on an owner decision/L-2/Decision A amendment.
- MEMORY-09 is not a continuation of the legacy lorebook; it is a new architecture with its own requirements and evidence.

**HFS-36 (Calculated — no readiness implication).** Nothing in this audit supports a claim that Howl is ready for MEMORY-09 testing. The opposite is established: state identity, token budget, multi-character isolation, resource enforcement and real inference all remain unverified.

---

## G. Remaining-work register

**HFS-37 (Calculated — no new requirements).** These are the outstanding or unverified implications of the original items, not additions.

### G.1 Release blockers (Tier 0)

| Item | Remaining work | Verification needed | Dependency |
|---|---|---|---|
| T0.1 | Same-session render regression | One message, JSON count, chat-switch render | None |
| T0.2 | Prove one prompt occurrence | Prompt dump/server log | T0.1 test path |
| T0.3 | Invalid-path dialog + load; chooser exists, error/load path still unaccepted | Invalid path error; valid path VRAM rise | Image engine |
| T0.4 | Preserve preview evidence | Prompt→image, correct aspect/crop | T0.3 |
| T0.5 (V1; still open) | Two-instance ownership test | Second instance leaves first child alive | Process ownership |

### G.2 Core mechanics (Tier 1)

| Item | Remaining work | Verification needed |
|---|---|---|
| T1.1 | Rerun current tree | Tab switch/return; sampler restart |
| T1.2 | Automated resource/lease regression | Fit/deny/unload/recheck both directions; telemetry loss |
| T1.3 | Token-level hidden-reasoning check | Real prompt/token test; no reappearance turn 2 |
| T1.4 | Full model-family matrix | At least two families, current build |
| T1.5 | Visual title check | Page headers |

### G.3 UI (Tier 2)

| Item | Remaining work | Verification needed |
|---|---|---|
| T2.1 | Character/world/image pages; persisted screenshots | Both themes; contrast/visual |
| T2.2 | 18 ghosts/30 sites; zero write-to-missing-ID | Static sweep + runtime smoke |
| T2.3 | Layout acceptance against stated criteria | Image page review; depends T0.4 |
| T2.4 | Initial fallback, fixed sizing, two-character test | Screenshot/visual |
| T2.5 | World/scenario full round-trip | Multi-image save/load/default |

### G.4 Additions (Tier 3)

| Item | Remaining work | Verification needed |
|---|---|---|
| T3.1 | Real token budget/flush; measure prefill; cache experiment | Tokenizer-based budget; long chat; cache-reuse measurement |
| T3.2 | Load-time baseline | Reproducible harness; compare LM Studio |
| T3.3-L | Token-based lore budget; scan depth/recursion; long-chat test | Real tokenizer; 2–3 characters |
| T3.3-M | Backend authority; token guard; legacy ownership decision | Memory boundary test; no frontend bypass |
| T3.3-A | AI-authored lore design/implementation | Separate authorization |
| T3.4 | Per-field generation design/implementation | Separate authorization |
| T3.5 | Multi-character isolation/flush | 2–3 characters, conflicting facts, abort/restart |
| T3.6 | Ternary loader research/implementation | Real file, no new deps |
| T3.7 | UI + headphone end-to-end | Real ASR/TTS run |
| T3.8 | Audio-generation research | Separate authorization |
| T3.9 | Diffusion version check | Current release evidence |
| T3.10 | Image-prompt research | Research only |
| T3.11 | Compatibility matrix/beta | Multiple hardware profiles |
| T3.12 | Benchmark harness | Baseline before optimization |

### G.5 Engineering/deferred ledgers

| Area | Remaining work |
|---|---|
| F2-A | Route old save/import paths through library writer; external edit/delete; import/export adapters; redirected-folder runtime |
| F2-B | Model-root UI; audio/tokenizer/template inspection; packaged cache; UI capability filtering |
| F2-C | No implementation; decisions only |
| F3-A | A8 propagation deferred; comprehensive cleanup contract |
| F3-B/C | Native acceptance/limits; no product-wide E3 |
| W18–W28 | W22 open; W25 incomplete; stable-ID/native path landed later; no model inference |
| FUT-HW1 | Deferred; advisory allocation, no lease authority |
| FUT-PROMPT1 | Deferred; no matrix |
| FUT-MEM1/FUT-LORE1 | Deferred; distinct from MEMORY-09 |

---

## H. Recommended next audit or work package

**HFS-38 (Calculated — recommendation only).** The next bounded step should be **Phase 1A: stable identity and a real tokenizer seam**, not a broad “foundation” claim.

Why this subset:

- It directly serves T1.1, T3.1, T3.3 and T3.5.
- It supplies the missing prerequisite for HSL and for any future memory capture.
- It is narrow enough to verify: ID stability across save/load; token counts from a real tokenizer; an explicit unavailable state; no character-count fallback.
- It does not require MEMORY-09, model downloads, G2, BGE, or a new architecture.

It should be paired with a **separate T0 release-blocker re-verification** (T0.1/T0.2/T0.3/T0.4), because the readiness audit does not cover them and no current E3 exists. That is a verification package, not new scope.

**HFS-39 (Observed — do not start implementation).** This audit authorizes nothing. It recommends only:

1. Phase 1A identity + tokenizer seam, with real-tokenizer E3 evidence and an honest unavailable state.
2. A read-only T0/T1/T2 regression pass to replace historical/live claims with current evidence.
3. A separate W22/W25 native authoring closure before any character-lore acceptance claim.

**HFS-40 (Calculated — readiness statement).** Howl is **not** ready for MEMORY-09 testing. The original foundation has substantial code and several historical narrow successes, but no current end-to-end evidence across chat state, token budgeting, multi-character isolation, resource enforcement and real inference. MEMORY-09’s legacy/lorebook distinction must remain intact: existing lorebook/memory is useful prior work, not partial MEMORY-09.

---

## I. Verification and change record

**HFS-41 (Observed — this audit).**

| Check | Result |
|---|---|
| Production code/schema/API changed | **0** |
| Tests added/changed/run | **0** |
| Models run / corpus / benchmarks | **0** |
| Canonical MEMORY-09 files changed | **0** (four canonical paths hashed; all four match HFR’s recorded digests) |
| `HFR-CP-2026-11-21-01` changed | **No** |
| `HF1-PLAN-2026-11-21-01` changed | **No** |
| Scope count corrected from V2’s 26 to recovered 27 | **Yes — T0.5 restored from V1** |
| Pre-existing user changes touched | **No** |
| This file | New, untracked; documentation audit |
| Canonical hashes | r33 `dfb28b6003ca356f237d5e9bf127ba978d8a35e9f15b258b509eb1beeee5f158` · v9 `698254ca262c0874eb1bf87679959e69492106633e016757c21ad3375a4fe95b` · ledger `e5008020ff68bd21b41d03d09de4542f6fc2d0ecddfc1111cc5b48dae71dfc3c` · map `b4ba33ccbfe8f8c1bff292ef72f17595341454187f67ce1350c93cb001df34d4` |
| Canonical paths hashed | `docs/research/MEMORY-09-ARCH-STAGE-CONTRACT-2026-11-15.md` · `docs/research/MEMORY-09-ARCHITECTURE-v9-2026-11-17.md` · `docs/research/MEMORY-09-OWNER-DECISIONS-2026-11-17.md` · `docs/research/MEMORY-09-CONFORMANCE-MAP-2026-11-17.md` |

**HFS-42 (Observed — evidence paths).**

| Topic | Paths |
|---|---|
| Original list | `PROJECT_ISSUES_AND_EXPANSION.md`; `ISSUES.md` |
| Historical execution | `STAGE_STATUS.md`; `T1_NOTES.md`; `TECH_DEBT.md`; `IMAGE_PAGE_WORKLOG.md` |
| State/prompt | `app-chat.go`; `app-prompt.go`; `internal/backend/chat/` |
| Token placeholder | `internal/backend/types/context.go`; `internal/backend/types/message.go` |
| Tokenizer | `internal/backend/llama/http_client.go`; `app.go` (`inferenceClient`) |
| Multi-character | `frontend/dist/js/character-queue.js`; `chat-stream.js`; `chat-render.js` |
| Lorebook/memory | `internal/backend/lorebook/`; `internal/backend/memory/`; `internal/lorebookstore/`; `frontend/dist/js/memory-store.js` |
| Resource/process | `internal/resource/`; `internal/proc/job_windows.go`; `app-vram-arbitration.go` |
| Readiness | `docs/architecture/HOWL-FOUNDATION-READINESS-AND-MEMORY-09-DECISION-REGISTER-2026-11-21.md` |
| Phase 1 | `docs/architecture/HOWL-FOUNDATION-PHASE-1-CHAT-STATE-CONTEXT-INTEGRITY-2026-11-21.md` |

**End of audit.** The original Howl foundation is broader than HFR A–E. Its recovered list is T0.1–T3.12: **27 active items** when V1’s T0.5 is preserved, not V2’s 26. V2 adds T2.4/T2.5 but drops T0.5 without a closure; a later document still calls T0.5 open. T4/E are parked; F2/F3/W are engineering history; HSL is a later owner clarification; MEMORY-09 is a separate architecture. The strongest current evidence is narrow history/ownership/prompt-boundary testing and historical live observations. No subsystem has current E3, the context budget is still character-based, multi-character orchestration is untested, and the tokenizer is not exposed at the app seam. The legacy lorebook/memory system is real prior work, not MEMORY-09 progress.
