# REVIEW — Independent Review of HOWL Complete Product Architecture Master Plan

**Status:** **INDEPENDENT REVIEW ONLY — NOT APPROVAL, NOT ADOPTION, NOT FREEZE, NOT ACCEPTANCE, NOT IMPLEMENTATION, NOT R6, NOT F3-D.**
**Date:** 2026-11-11
**Reviewed artifact:** `docs/review/HOWL-COMPLETE-PRODUCT-ARCHITECTURE-REQUIREMENTS-DEPENDENCY-MASTER-PLAN-2026-11-11.md`
**Master-plan SHA-256:** `e39cc396fa20adf2656a1b4e650a8bdc6c9a49cb46c3ddf95309ff144229ef8c`
**Master-plan size:** 703 lines; 83,354 bytes.
**Scope:** documentation-only challenge review. No source, test, build, migration, model-load, destructive-test, application-data, authority-record, freeze, acceptance, P1a, R6, or F3-D action was performed or authorized. No existing file was modified. Static source inspection only; no runtime execution was performed for this review.
**Method:** verified artifact hash; read the complete 703-line plan; independently inspected the cited ledger, avatar note, model-ecosystem note, shared-foundations assessment, F2-C cutover architecture, F3-C Phase-0 assessment, freeze/gate records, STAGE_STATUS, TECH_DEBT, and current source for material claims. DeepSeek summaries were not relied on exclusively.

---

## 1. Overall assessment and readiness verdict

**The master plan is a strong inventory draft but is NOT ready for adoption as the planning baseline.**

It correctly preserves every protected boundary, correctly refuses to promote deferred features, correctly withholds WORKING/VERIFIED status, and correctly identifies the expensive architectural decisions (chat serialization, root identity, media identity, model identity/readiness, prompt ownership, operation ownership). Its evidence taxonomy (source-confirmed vs runtime-confirmed vs test-confirmed vs inferred) is sound and consistently applied.

It nevertheless contains **two material requirement misclassifications, one terminology collision with safety consequences, and several completeness/dependency gaps** that would corrupt planning if adopted unchanged:

1. Mobile/remote work is classified "unestablished / not a recovered requirement" when the ledger records binding conditional constraints (MOB-1–MOB-4).
2. "Creative studio / general asset generation" is dismissed as unestablished without adjudicating its relationship to the approved F2-B5 modality-workspace direction and the ledger's north-star world-generator statement.
3. Content hot-swap (REG-1, protected) and model hot-swap (unaccepted product behavior with missing stream coordination/rollback) share one word and adjacent table rows; the plan separates them but the separation must be made structurally unmistakable before schema and lifecycle planning proceeds.

**Verdict: DO NOT ADOPT. Require the corrections in §10, then re-review. No freeze, acceptance, R6, F3-D, migration, or implementation follows from this review or from the plan.**

---

## 2. Executive summary

- **What the plan gets right:** subsystem inventory breadth (S1–S24); foundation map (FND-01–FND-16); data-ownership map; Tiers 0–2 / Tier 3 / FUT / parked classification; F3-C proposal-only labeling (PROP-F3C-01–05); contradiction register (C-01–C-22); dependency order preserving F3-D→F3-H; scope-protection rule (§9.5); protected-decision table (§10); evidence-gap register (§12).
- **What must be corrected:** MOB constraints reclassification; creative-workspace adjudication; hot-swap terminology split; per-chat persistence mapping (FUT-CH1/FUT-CH2/FUT-C1 vs per-chat model preference); Wails/isolated-runtime and hardware-validation targets carried explicitly; release-scope precision per category; remote-boundary limits on FND-15.
- **What was independently confirmed:** global model selection with no per-chat field; history stores only `hasImage`; `ChatMeta` has no model field; `Lorebooks []interface{}` untyped; `SetRAGService` wired while live prompt path bypasses it; `OnShutdown` now registered; STAGE_STATUS T1 closures vs open defects; gallery atomicity history; TTS freeze; freeze/acceptance/gate preservation; §16 hash preservation.
- **What remains unresolved:** owner decisions in plan §11.1; story-tools scope; TTS governance; media identity; startup degradation policy; one prompt path; resource authority; unestablished-feature recovery from the owner.

---

## 3. Findings register

### F-01 — Mobile/remote classified "unestablished" despite binding ledger constraints
- **Severity:** High
- **Plan sections/IDs:** §4.8 ("Unestablished ... not an authoritative project requirement"); S23; C-18; §8.2 row 16; §9.4; §11.1 item 10; §12 gap row.
- **Supporting artifact:** `PROJECT_ISSUES_AND_EXPANSION.md:1779-1790` — "Desktop-complete architecture requirement ... binding product constraint": MOB-1 (desktop fully usable with no phone/network/companion), MOB-2 (mobile never required for inference/creative work), MOB-3 (phone extends capability, never unlocks desktop function), MOB-4 (Android is the mobile test platform when that phase is authorized). "No mobile work is authorized or started. Any future remote API or pairing design must satisfy MOB-1–MOB-3."
- **Independent evidence:** `code_search` for `MOB-4|mobile test platform` returns only the ledger constraint row; `code_search` for remote/server/pairing/QR implementation returns no server, pairing, transport, or auth implementation; the only frontend "android" match is a user-agent check in `error-system.js:221`. Implementation absence confirmed; constraint presence also confirmed.
- **Why it matters:** "No implementation" is not "no requirement." Adopting the plan as written would erase a binding conditional constraint and let future remote/pairing design proceed without MOB-1–MOB-3 as acceptance criteria. This is exactly the product-vs-ledger error the mission warned against: absence from the F3-C contract is not proof of absence from the product requirements.
- **Correction:** reclassify Android/remote as **established conditional constraint, implementation unauthorized and unstarted**: desktop-complete is mandatory; any companion is optional and dependency-free; Android is the named test platform when/if the phase is authorized. Keep QR/pairing protocol itself unestablished. Carry MOB-1–MOB-4 into the dependency register as constraints on FND-15 and any future network work.
- **Confidence:** Confirmed.

### F-02 — Content hot-swap vs model hot-swap share one word
- **Severity:** High
- **Plan sections/IDs:** PROP-F3C-05 ("ADJACENT PRODUCT BEHAVIOR, NOT AN F3-C TARGET-SET REQUIREMENT"); C-05/C-07/C-08; §10 REG-1 row ("reusable, independently usable, scoped, and hot-swappable"); §7.3 hot-swap row.
- **Supporting artifacts:** `PROJECT_ISSUES_AND_EXPANSION.md:700-722` — REG-1 protects world/scenario/character/lorebook reuse, independent usability, scoping, and content hot-swap; `PROJECT_ISSUES_AND_EXPANSION.md:1294` (F3-B8: no attachment/scope/activation/hot-swap code modified); `ASSESSMENT-F3C-PHASE-0...` hot-swap adjacency; `app-model.go` manual load path (plan §7.3).
- **Independent evidence:** source search confirms manual load/unload exists; plan §7.3 and C-07 record no in-stream coordination and no rollback after the old server stops. REG-1 text contains no model-switching requirement; it constrains content relationships. The two "hot-swaps" are disjoint: content activation vs model-server replacement during an active conversation.
- **Why it matters:** A planner or implementer reading "hot-swappable" in the REG-1 row and "hot-swap ... not safe" in C-07 can conclude either that model switching is protected or that content swap needs stream coordination. Both misreadings cause wrong gates: one authorizes unsafe model behavior, the other over-constrains content work.
- **Correction:** rename throughout: **content hot-swap (REG-1, protected)** vs **model hot-swap (unaccepted product behavior)**. State the model-switching safety rule explicitly: no stream coordination plus no rollback equals no product behavior; per-chat continuation and provenance retrofit (C-08/C-09, PROP-F3C-01–04) are prerequisites, not decorations.
- **Confidence:** Strongly supported (source for manual path not re-traced line-by-line here; stream/rollback absence per cited assessment and plan, consistent with prior audits).

### F-03 — Per-chat persistence mapping is incomplete
- **Severity:** High
- **Plan sections/IDs:** §4.7 PROP-F3C-02 (per-chat continuation preference, proposed/not adopted); §7.4 per-chat row; C-08; §8.2 rows 3–5; §11.1 item 1.
- **Supporting artifacts:** Ledger FUT-CH1 (chat-scoped Author's Note with per-chat persistence), FUT-CH2 (single-character chats retain configured context capacity), FUT-C1 (up to five chat message cues persisted with chat); Phase-0 assessment on `howl-last-model-path` global storage; `frontend/dist/js/chat-model.js:169,271,393` (global key read/write).
- **Independent evidence:** `code_search howl-last-model-path` confirms global `localStorage` use; `read internal/backend/chat/types.go:50-78` confirms `ChatMeta` carries identity/scope fields but **no model field**, and `ChatHistory` carries only `hasImage`. Per-chat preference therefore has nowhere to persist regardless of the F3-C decision.
- **Why it matters:** The plan correctly labels per-chat model preference "proposed — not adopted," but a reader can infer that *no* per-chat persistence is required. The ledger already requires per-chat persistence for Author's Note content, context capacity behavior, and composer cues. The schema decision before F3-E must cover the union, not just the model pointer.
- **Correction:** add an explicit mapping: adopted deferred per-chat persistence (FUT-CH1 content, FUT-CH2 capacity behavior, FUT-C1 cues) vs proposed per-chat model preference (PROP-F3C-02) vs proposed provenance/seed/attachment fields. Decide the durable chat target set once (plan §11.1 item 1) against that union; record any accepted later-migration cost explicitly.
- **Confidence:** Confirmed for source facts; Strongly supported for the planning inference.

### F-04 — Provenance/seed/attachment classification is correct but needs a ledger anchor
- **Severity:** Medium
- **Plan sections/IDs:** §4.7 PROP-F3C-01/03/04; §7.4 provenance/seed/attachment rows; C-03/C-09/C-10/C-11; §11.1 item 1.
- **Supporting artifacts:** F3-C Phase-0 assessment (per-message metadata gap; `Metadata` richer type in `internal/backend/types/message.go:31-43` not the persisted shape; `InferenceOptions.Seed` default `-1`; `currentSeed` global-only; zero-seed validity); ledger FUT-C5/FUT-ST2 provenance language (sampler provenance, import field provenance — not per-turn generation provenance).
- **Independent evidence:** `read types.go` confirms `ChatHistory` shape (`id, sender, text, role, timestamp, hasImage, characterId?`); `code_search hasImage` confirms boolean-only persistence across `app-chat.go`, `app-prompt.go`, `chat-actions.js:62`, `chat-session.js:146`. Seed `-1` defaults confirmed in `app.go:223` and sampler/image defaults; per-turn seed mode absent per assessment.
- **Why it matters:** The plan's "proposed — not adopted" labels are correct and must survive. The residual risk is a future claim that provenance/seed/attachment were "never product requirements" because F3-C did not adopt them. The ledger's provenance language is about sampler/import evidence, not per-turn generation identity — the distinction must be written down so neither side borrows the other's authority.
- **Correction:** annotate PROP-F3C-01/03/04 as **product-schema proposals arising from the Phase-0 contract gap**, distinct from ledger FUT-C5/FUT-ST2 provenance items. Do not promote; do not erase.
- **Confidence:** Confirmed for source/shape facts; Strongly supported for classification.

### F-05 — STAGE_STATUS closures vs current evidence handled correctly; carry explicitly
- **Severity:** High (if mishandled) / Low correction cost
- **Plan sections/IDs:** §4.1 T1.1–T1.4 rows; C-01; §11.1 item 2; §13 item 10.
- **Supporting artifacts:** `STAGE_STATUS.md:102,142,187,244` (T1.1–T1.4 closed 2026-10-31); F3-C defect inventory (ROOT-01, WIPE-02, IDENTITY-03, REACH-04, STARTUP-05, ALIAS-06 open).
- **Independent evidence:** STAGE_STATUS closure lines verified by search; defect-open status reproduced from plan §2.3 and consistent with cited defect disposition request (hash-verified, not content-audited here).
- **Why it matters:** The most likely planning corruption is treating a 2026-10-31 stage closure as product-wide verification. The plan correctly treats closures as historical and state ownership as partial. This must survive adoption.
- **Correction:** keep C-01 and §13 item 10 verbatim; add a standing rule that no STAGE_STATUS closure is cited as verification for chat/memory/resource scope without current source-or-runtime evidence.
- **Confidence:** Confirmed.

### F-06 — RAG/live-prompt split confirmed; decision still missing
- **Severity:** High
- **Plan sections/IDs:** S2; S7; C-02; FND-09/FND-10; §8.2 rows 5, 12; §11.1 item 5; §13 item 4.
- **Supporting artifacts:** `ASSESSMENT-...-SHARED-FOUNDATIONS-2026-11-09.md:32,273`; `app-prompt.go:671 buildPromptMessages`; `app-chat.go:176`; `app.go:186 SetRAGService`; `internal/backend/chat/service.go:1053-1054`.
- **Independent evidence:** `code_search buildPromptMessages|SetRAGService` confirms the live `App` prompt builder, the `chat.Service` RAG attachment point, and the assessment's bypass claim. No RAG call inside the live path was demonstrated in the cited assessment; this review did not re-trace every call site.
- **Why it matters:** RAG/lore can be configured yet never affect replies — a silent capability failure. Every prompt/memory/lore decision downstream (FND-09/10, FUT-LORE1, multi-character isolation) depends on which path owns assembly.
- **Correction:** require the §11.1 item 5 decision (one live prompt path; RAG/lore/tokenizer/isolation owner) before F3-E and before any RAG/lore-capability claim. Forbid a third prompt path.
- **Confidence:** Strongly supported.

### F-07 — OnShutdown supersede confirmed; packaged lifecycle still open
- **Severity:** Medium
- **Plan sections/IDs:** §1 lifecycle bullet; S1; FND-16; C-20; §12 gaps.
- **Supporting artifacts:** `main.go:53` (`OnShutdown: app.OnShutdown`); `ENGINEERING-F1-LIFECYCLE-PROCESS-OWNERSHIP-2026-11-10.md:71,147`; 2026-11-09 assessment (unregistered finding).
- **Independent evidence:** `code_search OnShutdown` and `read main.go` context confirm registration; F1 record confirms the repair plus residual limits (packaged close, multi-instance, alias/cross-process, every child type).
- **Why it matters:** Repeating the stale "unregistered" finding would misdirect lifecycle work; claiming "fixed" would over-authorize packaged/native acceptance.
- **Correction:** keep "improved but not fully verified" and C-20 exactly; carry W4 (WebView2 profile path blocked) and W8 (unwritable-root/rejected-creation paths open) from the Wails/isolated-runtime records into §8.2/§12 so packaged acceptance cannot be inferred from registration alone.
- **Confidence:** Confirmed.

### F-08 — Gallery atomicity history handled correctly; video gap correctly held
- **Severity:** Medium
- **Plan sections/IDs:** S8/S9; FND-08; §6 images/videos rows; §10 F2-C5/F3-A6 rows; C-19; §13 item 7.
- **Supporting artifacts:** `ARCHITECTURE-F2C-...-2026-11-09.md:99,109,132,134` (non-atomic pair; uncontained read; suppressed delete); `ENGINEERING-F3A-...-MIRROR-SEMANTICS...:18-19,77-81` (temp-file+rename; sidecar-failure image removal; `ResolveFile` containment; truthful deletes); amendment notes on F3-F ownership (`gallery.go:94-128` writeFileAtomic without media root/index/journal).
- **Independent evidence:** not re-traced to `gallery.go` lines here; the two-record sequence (old defect → later repair) is consistently cited across independent records. Plan correctly refuses to treat the old finding as current fact and refuses to treat the repair as runtime-proven.
- **Why it matters:** Media attachment planning (PROP-F3C-04, F3-F/F3-G) depends on whether PNG+sidecar is already atomic. Over-claiming atomicity would let attachments assume durability the store does not provide.
- **Correction:** keep C-19; add the video corollary explicitly in §13 item 7: no video sidecar/catalogue/delete contract exists, so video attachments have no planning basis until F3-G.
- **Confidence:** Strongly supported.

### F-09 — "Creative studio" needs adjudication, not dismissal
- **Severity:** Medium
- **Plan sections/IDs:** S24; §4.8 creative-studio row; §7.4 creative-studio row; §8.2 row 16; §9.4; C-18; §11.1 item 10.
- **Supporting artifacts:** Ledger north star ("world generator and inspiration tool"); F2-B5 (separate modality-specific creative workspaces; shared catalog does not merge chat/image interfaces); E.6/E.8 (workspace/staging-board parked concepts).
- **Independent evidence:** F2-B5 verified in ledger text (§F2 table); no dedicated "creative studio" page/service/owner found in plan's source survey (consistent with plan). The gap is interpretive, not evidentiary: whether the mission's "creative-studio direction" refers to F2-B5 workspaces, E.6/E.8 boards, or a new asset-generation scope.
- **Why it matters:** "Unestablished — do not treat as scope" is safe against invention but unsafe against amnesia: a real approved workspace-separation direction (F2-B5) could be starved because an adjacent vague label was dismissed.
- **Correction:** adjudicate explicitly: (a) F2-B5 modality workspaces — approved direction, deferred implementation; (b) E.6/E.8 layout/staging concepts — parked; (c) any general asset-generation studio beyond (a)+(b) — unestablished, owner must define before architecture. Do not merge (c) into (a).
- **Confidence:** Strongly supported.

### F-10 — Avatar/3D classification is correct; desktop-interaction scope must stay explicit
- **Severity:** Medium
- **Plan sections/IDs:** S22; §4.3 T4.1/T4.2; §4.8 avatar rows; §7.4 3D row; §8.2 row 15; §9.3/§9.4; C-17.
- **Supporting artifacts:** `docs/tier4-avatar.md` (status parked; skeleton-once/skin-per-character; MakeHuman/CC0; TTS-amplitude mouth; Tier B/C/D; UltraTex external; licensing flag).
- **Independent evidence:** read avatar-note head (80 lines) confirms parked status, body library, and pipeline structure; TTS-governance dependency per C-17 consistent with TTS freeze (F-13).
- **Why it matters:** 2D avatar rendering (source-confirmed) is routinely conflated with the 3D reactive-avatar proposal. The plan correctly separates them; adoption must not blur the line when media/image/TTS dependencies are discussed.
- **Correction:** keep parked/Tier-gated classification; require any "desktop avatar" interaction scope beyond the Tier-4 note to be recovered from the owner (plan §11.1 item 10) rather than inferred from 2D rendering or the mission label.
- **Confidence:** Strongly supported.

### F-11 — Dependency order is sound; lifecycle placement and remote work need guards
- **Severity:** High
- **Plan sections/IDs:** §8.1 graph; §8.2 rows 1–16; §8.3 steps 1–9; §8.4 cycles.
- **Independent analysis:** rows 1–7 correctly put F3-C disposition, root/path identity, and the durable schema decision before F3-D→F3-H; rows 9–12 correctly put model identity/prompt before resource admission/operation ownership before provenance/memory isolation; media completion after identity/catalogue is correct. The acknowledged cycles (resource registration-before-readiness; App/chat.Service/persistence/frontend projections; lore cache/snapshot timing; media producer/owner/consumer) are genuine and correctly named.
- **Why it matters:** Two sequencing traps remain. First, FND-16 lifecycle verification sits at step 8 while every migration, wipe, and final-save test needs a trustworthy shutdown boundary — deferring it invites rework or untrustworthy stage evidence. Second, row 16 + step 8.3(8) ("remote/QR ... if separately authorized") could be read as permission to design remote architecture once any authorization exists, bypassing MOB-1–MOB-3 and F-01.
- **Correction:** (a) move lifecycle/shutdown verification earlier or run it in parallel with guardrails: no migration/wipe/final-save evidence is accepted without the shutdown boundary it exercised; (b) gate row 16 on MOB-1–MOB-3 satisfaction plus an explicit owner requirement, not on a generic future authorization; (c) make the §8.4 resource-cycle breaker explicit: a verified model-state machine (requested→loading→ready→failed→unloading→unloaded) with release evidence, not another registry map.
- **Confidence:** Inferred (planning judgment on a sound structure).

### F-12 — Release/deferred taxonomy is sound; two scope-leak paths must be closed
- **Severity:** Medium
- **Plan sections/IDs:** §9.1–§9.5; §7.4; §4.7.
- **Supporting artifacts:** Ledger ordering law (Broken→UI→Additions→Ideas; T3 gated on Tiers 0–2); FUT tables; F2-B7 freeze; REG-1.
- **Independent evidence:** ledger ordering law and T3 gating confirmed in `PROJECT_ISSUES_AND_EXPANSION.md` head; plan §9.5 states the correct rule (future capability may influence architecture, never authorize implementation/migration/gate change).
- **Why it matters:** §9.1's "current ... functionality to the extent independently verified" is too vague to plan against per category, and §7.4's "architectural change if frozen without it" column invites reserving fields that later get read as approval.
- **Correction:** (a) append gate references per §9.1 category (e.g., chat→F3-E gate; media→F3-F/G; characters→F3-D); (b) strengthen §9.5 with "no migration by inclusion": reserving a field/identity/state creates no migration, no backfill, and no release-gate change.
- **Confidence:** Strongly supported.

### F-13 — Protected boundaries preserved; no reopening warranted on this evidence
- **Severity:** Critical (if violated) — plan complies
- **Plan sections/IDs:** §2.3; §10; §11 C-04/C-05/C-21; §13 item 10; §16.
- **Supporting artifacts:** `FREEZE-F3C-ARCHITECTURE-2026-11-11.md:9,164` (freeze granted; not acceptance; no P1a/source/test authorization); `F3C-MASTER-...:30,1091`; `GATE-F3C-...-V2-2026-11-11.md:12,314,318`; ledger F2-B7; `STAGE_STATUS`/F3-B records on A9/A10/A12/A13; V20 blocked pending F3-F.
- **Independent evidence:** freeze/acceptance/gate language verified by search; TTS freeze confirmed in ledger F2-B7 row; §16 hashes for all five protected F3-C artifacts re-verified byte-identical in this review (ADJUDICATION `85e1968...`, OWNER-REQUEST `ed8bf2e2...`, INVENTORY `343104b5...`, INVESTIGATION `494539b9...`, DISPOSITION-REQUEST `ce32ac0b...`).
- **Why it matters:** Any baseline that silently reopens accepted work or converts freeze into acceptance would invalidate the F3-C measurement contract.
- **Correction:** none — preserve. Do not reopen A9/A10/A12/A13, D-STORAGE, V20, TTS governance, or Gate V18–V24 on the evidence in this review.
- **Confidence:** Confirmed.

### F-14 — Completeness: ledger coverage is broad but denominated too narrowly
- **Severity:** Medium
- **Plan sections/IDs:** §2.1; §4.1–§4.8; §12; §14.
- **Supporting artifacts:** ledger HW-C1–HW-C3 (hardware compatibility targets, deferred/not started); W1–W8 Wails/isolated-runtime decisions (W4 blocked WebView2 profile path; W5–W7 implemented/verified; W8 open); E.4/E.5 friction/fake-strength guardrails; T2.9.5B audio carry-forwards (F1/F2) referenced in decision registers.
- **Independent evidence:** HW-C1–C3 and W1–W8 blocks read in ledger §§Hardware/Wails (2026-11-10); plan §12 covers packaged-run and measurement gaps generically but does not name HW-C1–C3 or W4/W8 as planning constraints.
- **Why it matters:** A baseline that omits named deferred validation targets and blocked native-acceptance items will schedule stage work that cannot be evidenced (e.g., compatibility claims without HW-C1 matrix; native runs without W4).
- **Correction:** carry HW-C1–HW-C3, W4 (blocked), and W8 (open) explicitly into §8.2/§12 with owners; cite E.4/E.5 as design guardrails where UI/error planning touches them; note T2.9.5B F1/F2 carry-forwards as out-of-scope-but-adjacent so audio planning does not assume them closed.
- **Confidence:** Strongly supported.

### F-15 — FND-15 must not plan a remote boundary the product has not authorized
- **Severity:** Low (architectural hygiene)
- **Plan sections/IDs:** FND-15; §7.4 Android/remote row; S23.
- **Why it matters:** FND-15 lists "remote pairing boundary" among dependents while the remote requirement is (correctly per F-01 correction) conditional-only. Foundations should not acquire remote-shaped interfaces before the owner establishes the requirement.
- **Correction:** limit FND-15 dependents to implemented local paths (model downloads, document imports) plus MOB-1–MOB-3 constraints; move any remote/pairing boundary to §9.4 optional futures gated on F-01.
- **Confidence:** Inferred.

---

## 4. Requirements completeness review

The plan recovers the named ledger requirements across Tiers 0–2, Tier 3, Tier 4/E-series, F2-B/FUT, FUT-*/deferred register, FUT-HW1/HW2, FUT-PROMPT1, F2-C, F3-A/R1, F3-B, and F3-C proposal deltas, plus the avatar and model-ecosystem notes. Tier qualifiers ("nothing starts until Tiers 0–2 are closed") and deferred/parked labels are faithfully preserved. No invented requirements were found; unestablished items are explicitly marked rather than filled.

Corrections required (from findings): reclassify MOB-1–MOB-4 as established conditional constraints (F-01); adjudicate creative-workspace scope against F2-B5/E.6/E.8 (F-09); carry HW-C1–HW-C3, W4, W8, E.4/E.5 guardrails, and audio carry-forwards explicitly (F-14); keep story-tools scope unknown pending owner input rather than inferring from controls/tutorials (plan already does — preserve).

Historical sequencing is respected: ledger ordering law, T3 gating, F3-D→F3-H order, and deferred-slot assignments are preserved. Formal requirements are distinguished from proposals/experiments with one exception class: the F-02 terminology collision and the F-03 union-of-per-chat-persistence mapping, both corrected above.

---

## 5. Product-requirement vs contract-adoption classification review

The plan's core discipline — "not adopted into F3-C" does not mean "not a product requirement" — is correct and must survive adoption. Verified applications:

- **Per-chat continuation (F-03):** per-chat model preference is correctly proposal-only; ledger per-chat persistence (FUT-CH1/FUT-CH2/FUT-C1) is correctly deferred-but-real. Adopt the union mapping.
- **Per-turn provenance (F-04):** correctly proposal-only; ledger provenance items are a different kind (sampler/import evidence). Keep the kinds separate.
- **Seed state (F-04):** `-1`-means-random global behavior is source-confirmed; explicit/random/zero/legacy-unknown semantics are correctly proposal-only. Preserve zero-validity and no-`-1`-sentinel rules as proposal content, not current fact.
- **Durable attachment references (F-04/F-08):** `hasImage`-only persistence is source-confirmed; durable references are correctly proposal/F3-F work. Missing references must remain unavailable, never substituted (plan states; preserve).
- **Model hot-swap (F-02):** manual load path exists; stream-coordinated, rollback-safe hot-swap is not an accepted behavior. Adjacent-product-behavior labeling is correct provided the REG-1 content-hot-swap split is made explicit.
- **Android/remote (F-01):** correction required — conditional product constraints exist; implementation does not.
- **Creative-studio direction (F-09):** adjudication required — approved workspace separation vs parked concepts vs unestablished general scope.

---

## 6. Shared-foundation and dependency review

FND-01–FND-16 cover lifecycle, ownership, identity, model capability, engines, resources, provenance, media, memory/context, prompts, errors, settings, UI projection, import/export, security boundaries, and process ownership. The ownership map (§6) and dependency graph (§8.1–§8.2) are coherent: identity before references; schema before F3-E; F3 order preserved; model identity/prompt before admission/operations before provenance/isolation; media identity before attachments/avatars.

Independent challenges sustained: the RAG/live-path split (F-06) with a required one-path decision; the resource registration-before-readiness cycle (F-11) requiring a verified state machine; the media producer/owner/consumer cycle requiring F3-F/F3-G identity first; lifecycle verification timing (F-11 correction); remote-boundary hygiene (F-15). No circular dependency was found that the plan conceals — §8.4 names them — but the breakers must be proved in later design, not assumed from the diagram.

---

## 7. Release-scope and future-compatibility review

§9 correctly separates release obligations (Tiers 0–2, F3-A/R1 repairs, F3-B accepted scope, F3-C open work in stage order) from deferred implementation (T3, FUT-*, T4, audio expansion, video completion, AI-assisted creation) and optional futures (Tier C/D avatars, remote/QR, creative studio, AI-as-user, vault, multi-user/cloud — none authorized). The §9.5 rule is the correct adoption guard.

Required hardening (F-12): per-category gate references in §9.1; explicit "no migration by inclusion" alongside "may influence architecture." Future-compatibility analysis (§7.4) is valuable precisely because it names the rework cost of freezing the chat schema, root identity, media identity, model identity, prompt ownership, and operation model prematurely — adopt that warning, not the fields.

---

## 8. Protected-boundary review

Verified preserved: F3-C freeze (granted; not acceptance; no P1a/source/test authorization); master contract control; Gate v2 V1–V17 adopted with V18–V24 proposed; F3-C Phase 0 open; R6 unauthorized; F3-C acceptance not granted; F3-D blocked; D-ROOT six defects open; D-STORAGE blocked; V20 blocked; A9/A10/A12/A13 closed; TTS/F2-B7 frozen; REG-1 protected; local-first boundary. §16 preservation hashes re-verified byte-identical. No protected artifact was modified by the plan or by this review. No reopening is recommended on the evidence examined.

---

## 9. Evidence gaps and unresolved owner decisions

Plan §12 gaps are adopted as review findings: static-only evidence; no packaged/native run; no model load/inference; no concurrency/race testing; no alias-matrix, two-process, reload-with-attachment, provenance/seed, hot-swap stream/rollback, OOM/VRAM-matrix, audio/TTS end-to-end, gallery runtime, 3D performance, Android/studio requirement, or exhaustive-corpus evidence. Unknown remains unknown.

Unresolved decisions (plan §11.1, preserved): chat durable target set; canonical root/path identity; wipe semantics; startup degradation; one live prompt path; resource authority; media identity; TTS governance (unchanged); story-tools scope; unestablished-feature requirements. Plus this review's additions: MOB reclassification acceptance (F-01); workspace adjudication (F-09); hot-swap terminology split (F-02); per-chat union mapping (F-03); lifecycle timing and remote gating (F-11); §9 hardening (F-12); HW-C/W4/W8 carriage (F-14).

---

## 10. Specific corrections required before adoption

1. Reclassify MOB-1–MOB-4 as established conditional constraints; gate any remote work on MOB-1–MOB-3 (F-01).
2. Split content hot-swap (REG-1) from model hot-swap structurally throughout (F-02).
3. Add the per-chat persistence union map: FUT-CH1/FUT-CH2/FUT-C1 (deferred-real) vs PROP-F3C-02 (proposed) (F-03).
4. Anchor PROP-F3C-01/03/04 as Phase-0 schema proposals distinct from ledger provenance items (F-04).
5. Add the STAGE_STATUS non-verification rule (F-05); keep the one-path prompt decision as an F3-E prerequisite (F-06).
6. Carry W4/W8 and HW-C1–HW-C3 explicitly; cite E.4/E.5 guardrails; note audio carry-forwards (F-07/F-14).
7. Keep C-19 and add the video-catalogue corollary (F-08).
8. Adjudicate F2-B5 vs E.6/E.8 vs unestablished studio scope (F-09); keep avatar Tier gating with explicit desktop-interaction recovery (F-10).
9. Resequence or parallel-guard lifecycle verification; gate remote design on MOB constraints; specify the resource state-machine breaker (F-11).
10. Harden §9.1 with per-category gates and §9.5 with no-migration-by-inclusion (F-12); constrain FND-15 to local paths (F-15).

---

## 11. What this review does not authorize

- It does not adopt the master plan as a baseline.
- It does not amend, freeze, accept, or authorize implementation for F3-C or any stage.
- It does not authorize R6, F3-D, migration, cleanup, model work, native runs, or any source/test/build/data action.
- It does not close any defect, reopen any accepted work (A9/A10/A12/A13), adopt Gate V18–V24, unblock V20/D-STORAGE/F3-D, or change TTS governance.
- It does not establish Android/remote, creative-studio, desktop-avatar-interaction, or story-tools scope requirements.
- It does not treat absence of evidence as proof of absence, and it resolves no contradiction by silence: F-01–F-15 record the evidence and the owner decision each needs.
