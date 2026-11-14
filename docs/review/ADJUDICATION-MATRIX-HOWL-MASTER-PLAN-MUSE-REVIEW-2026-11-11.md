# ADJUDICATION MATRIX — HOWL MASTER PLAN vs INDEPENDENT MUSE REVIEW

**Date:** 2026-11-11  
**Type:** Read-only evidence adjudication  
**Status:** **ADJUDICATION ONLY — NOT MASTER-PLAN ADOPTION, NOT AN EDIT, NOT A FREEZE, NOT AN ACCEPTANCE, NOT R6, NOT F3-C AMENDMENT, NOT F3-D, NOT IMPLEMENTATION AUTHORIZATION**  
**Master plan:** `docs/review/HOWL-COMPLETE-PRODUCT-ARCHITECTURE-REQUIREMENTS-DEPENDENCY-MASTER-PLAN-2026-11-11.md`  
**Independent review:** `docs/review/REVIEW-INDEPENDENT-HOWL-COMPLETE-PRODUCT-ARCHITECTURE-MASTER-PLAN-2026-11-11.md`

---

## 1. Artifact verification

| Artifact | Expected SHA-256 | Observed SHA-256 | Identity |
|---|---|---|---|
| Master plan | `e39cc396fa20adf2656a1b4e650a8bdc6c9a49cb46c3ddf95309ff144229ef8c` | `e39cc396fa20adf2656a1b4e650a8bdc6c9a49cb46c3ddf95309ff144229ef8c` | **MATCH** |
| Independent review | `af493c0b6118ebede556d4927e8ab93b9378ce46c8298f3bb9e38554f6265e9d` | `af493c0b6118ebede556d4927e8ab93b9378ce46c8298f3bb9e38554f6265e9d` | **MATCH** |

Master-plan size: 703 lines, 83,354 bytes. Review size: 250 lines, 33,410 bytes. Both identities match the supplied values.

Neither artifact was edited. This adjudication did not run tests, builds, migrations, model loads, or application source changes. Evidence is static documentation/source inspection.

---

## 2. Executive adjudication

The independent review is **substantially valid**. Its core verdict — the master plan is a strong inventory but is not ready for adoption — is supported. It correctly identifies material requirement omissions, a terminology collision with safety consequences, incomplete per-chat schema mapping, a missing prompt-path gate, and completeness gaps.

The adjudication nevertheless makes four corrections to the review:

1. **W8 is not open.** The ledger records W8 as **VERIFIED** at `PROJECT_ISSUES_AND_EXPANSION.md:1864`. The earlier W8 `OPEN` line at `:1840` is superseded by the follow-up. The review’s F-07 and F-14 carry the stale status.
2. **W4 is not a current blocker in the form stated.** W4’s stated condition was that `main.go` did not set the WebView2 profile path. Current `main.go:64` sets `WebviewUserDataPath` under the isolation switch, and W5–W7 record the isolated native runtime implementation/verification. W4 is therefore **superseded in substance by W5–W7**, not a live “blocked” item. A formal W4 closure record may be missing, which is an evidence gap, not evidence that the blocker remains.
3. **F-13 is not a critical defect.** It is a **critical control confirmation**: “Critical” describes the consequence if protected boundaries were violated; the review itself says the plan complies and no reopening is warranted. The severity label should remain visible as a control consequence, but the finding is not a current defect and requires no plan correction.
4. **F-11 and F-15 require narrower corrections than written.** The dependency order should not be silently resequenced: preserve the accepted F3 order and add parallel/evidence guardrails. FND-15 lists a missing remote boundary as a gap; it does not actually design one. The correction is to constrain its wording, not to claim the plan authorized remote work.

**Disposition count:** 11 valid findings (including three control confirmations), 4 partially valid findings, 0 invalid, 0 unresolved. “Valid” does not mean every proposed correction is accepted verbatim; it means the finding identifies a real issue in the master plan or a real control that must be preserved.

**Overall verdict:** **DO NOT ADOPT the master plan unchanged.** Apply the evidence-only corrections, record the owner decisions, correct the review’s stale W4/W8 carry-forward in the next planning revision, and re-review the corrected plan. No protected boundary is reopened by this adjudication.

---

## 3. Complete F-01–F-15 disposition table

“Exact wording” below is the review’s exact finding heading. The detailed evidence and exact proposed correction for each finding follow in §4.

| ID | Exact finding wording | Muse severity | Disposition | Master-plan correction? | Owner decision? | Primary axes |
|---|---|---|---|---|---|---|
| F-01 | “Mobile/remote classified "unestablished" despite binding ledger constraints” | High | **Valid** | Yes — reclassify MOB constraints; keep implementation unstarted | No for constraints; yes for any remote phase authorization | Requirements, dependencies, security |
| F-02 | “Content hot-swap vs model hot-swap share one word” | High | **Valid** | Yes — structural terminology split | No | Architecture, dependencies, protected boundary clarity |
| F-03 | “Per-chat persistence mapping is incomplete” | High | **Valid** | Yes — add union mapping before F3-E | Yes — durable target set | Requirements, architecture, dependencies |
| F-04 | “Provenance/seed/attachment classification is correct but needs a ledger anchor” | Medium | **Valid** | Yes — separate provenance families | Yes for adoption of new schema fields; no for annotation | Requirements, architecture |
| F-05 | “STAGE_STATUS closures vs current evidence handled correctly; carry explicitly” | High (if mishandled) / Low correction cost | **Valid — control confirmation** | No substantive correction; add standing rule | No | Release scope, verification discipline |
| F-06 | “RAG/live-prompt split confirmed; decision still missing” | High | **Valid** | Yes — make one-path decision an F3-E prerequisite | Yes | Architecture, dependencies, release scope |
| F-07 | “OnShutdown supersede confirmed; packaged lifecycle still open” | Medium | **Partially valid** | Yes for packaged-lifecycle caveat; **no** for “W4 blocked/W8 open” | No | Dependencies, release scope |
| F-08 | “Gallery atomicity history handled correctly; video gap correctly held” | Medium | **Valid — control confirmation + minor** | Minor — add video F3-G corollary | No | Architecture, dependencies |
| F-09 | “"Creative studio" needs adjudication, not dismissal” | Medium | **Valid** | Yes — separate F2-B5, E.6/E.8, and general studio | Yes only for general studio definition | Requirements, release scope |
| F-10 | “Avatar/3D classification is correct; desktop-interaction scope must stay explicit” | Medium | **Valid — control confirmation + minor** | Minor — make desktop-interaction recovery explicit | Yes if scope exceeds Tier-4 note | Requirements, architecture |
| F-11 | “Dependency order is sound; lifecycle placement and remote work need guards” | High | **Partially valid** | Yes — add parallel lifecycle guard and remote gate; do not reorder F3 | Yes only for remote phase; no for guard | Dependencies, release scope |
| F-12 | “Release/deferred taxonomy is sound; two scope-leak paths must be closed” | Medium | **Valid** | Yes — per-category gates; no-migration-by-inclusion | No | Release scope, architecture |
| F-13 | “Protected boundaries preserved; no reopening warranted on this evidence” | Critical (if violated) — plan complies | **Valid — critical control confirmation, not a defect** | No | No | Protected boundaries |
| F-14 | “Completeness: ledger coverage is broad but denominated too narrowly” | Medium | **Partially valid** | Yes for HW-C1–C3 and E.4/E.5; **correct** W4/W8 statuses | No | Requirements, dependencies |
| F-15 | “FND-15 must not plan a remote boundary the product has not authorized” | Low | **Partially valid** | Yes — narrow FND-15 wording to local/MOB constraints | No | Architecture, security |

---

## 4. Finding-by-finding adjudication

### F-01 — Mobile/remote classified “unestablished” despite binding ledger constraints

**Muse severity:** High. **Disposition: Valid.**

**Independent evidence.** `PROJECT_ISSUES_AND_EXPANSION.md:1783-1789` records MOB-1 through MOB-4 as binding product constraints: desktop remains fully usable without a phone, network, or companion; mobile is never required for local inference/creative work; a phone extends capability and never unlocks desktop function; Android is the test platform only when that phase is authorized. It also states that no mobile work is authorized or started and that any future remote API or pairing design must satisfy MOB-1–MOB-3. The master plan’s §4.8, S23, C-18, §8.2 row 16, §9.4, and §11.1 item 10 instead call Android/remote “unestablished” without carrying those constraints. The implementation absence is confirmed: no remote server, pairing, transport, or auth implementation was found; the only frontend “android” string is a user-agent check in `error-system.js`.

**Required correction.** Reclassify precisely:

> **Android/remote: established conditional product constraint (MOB-1–MOB-4); implementation unauthorized and unstarted. QR/pairing protocol, remote API, authentication, and remote-session behavior remain unestablished.**

Carry MOB-1–MOB-3 into §4.8, §8.2 row 16, §9.4, §12, and FND-15 as constraints on any future network work. Do not add a remote feature or implementation commitment.

**Owner decision.** None is needed to record the existing constraints. An explicit owner decision is required to authorize a remote phase, define a QR/pairing protocol, or approve remote API/authentication behavior.

**Axes.** Requirements, dependencies, security. **Confidence:** Confirmed. **Gap:** no remote implementation exists to test; the protocol remains undefined.

### F-02 — Content hot-swap vs model hot-swap share one word

**Muse severity:** High. **Disposition: Valid.**

**Independent evidence.** REG-1 at `PROJECT_ISSUES_AND_EXPANSION.md:700-722` protects world/scenario/character/lorebook reuse, independent usability, scoping, and content hot-swap. F3-B8 at `:1294` confirms that no attachment, scope, activation, or content hot-swap code was modified. The master plan uses “hot-swappable” for REG-1 in §10 and “hot-swap” for model-server replacement in §4.7/§7.3. These are disjoint: content activation does not require stream coordination or model rollback, and model replacement is not protected by REG-1.

**Required correction.** Make the split structural:

- **Content hot-swap (REG-1, protected):** activation/replacement of world, scenario, character, and lorebook context; no model lifecycle semantics.
- **Model hot-swap (unaccepted product behavior):** replacing the loaded inference server during an active conversation; requires request coordination, readiness/failure publication, and rollback. This is not protected by REG-1 and is not an accepted behavior.

Rename PROP-F3C-05 to “model hot-swap.” In §10, replace “hot-swappable” with “content hot-swappable.” Add the distinction wherever both are discussed.

**Owner decision.** None for terminology. Adopting model hot-swap as a product behavior remains separate and unauthorized.

**Axes.** Architecture, dependencies, protected-boundary clarity. **Confidence:** Strongly supported. **Gap:** no runtime hot-swap stream/rollback test exists.

### F-03 — Per-chat persistence mapping is incomplete

**Muse severity:** High. **Disposition: Valid.**

**Independent evidence.** The ledger already establishes deferred per-chat requirements: FUT-CH1 at `PROJECT_ISSUES_AND_EXPANSION.md:615` requires chat-scoped Author’s Note persistence; FUT-CH2 at `:616` requires per-chat context-capacity behavior before rollover; FUT-C1 at `:605` requires up to five composer cues persisted with the chat. The master plan lists all three but does not connect them to the F3-E schema decision. `ChatMeta` (`internal/backend/chat/types.go:50`) has no model field, and `ChatHistory` (`:72`) has only `hasImage`; model selection is global at `chat-model.js:169,271,393`. PROP-F3C-02 is therefore proposal-only, while FUT-CH1/FUT-CH2/FUT-C1 are already ledger requirements at the product level.

**Required correction.** Add a union table to §4.7/§8.2:

| Item | Kind | Status |
|---|---|---|
| FUT-CH1 Author’s Note content | Deferred but real per-chat product requirement | Must be considered in F3-E schema |
| FUT-C1 composer cues | Deferred but real per-chat product requirement | Must be considered in F3-E schema |
| FUT-CH2 context-capacity behavior | Deferred product requirement; whether the setting itself is persisted per chat is not yet specified | Resolve before treating F3-E as complete |
| PROP-F3C-02 per-chat model preference | Proposal, not adopted | Separate from historical provenance |
| PROP-F3C-01/03/04 provenance, seed, attachments | Proposals, not adopted | Separate schema fields |

The F3-E target-set decision must address this union, not only the model pointer.

**Owner decision.** Yes — the durable chat target set, including whether the proposed model pointer joins the deferred per-chat requirements.

**Axes.** Requirements, architecture, dependencies. **Confidence:** Confirmed for source facts; planning inference strongly supported. **Gap:** FUT-CH2 persistence semantics are not fully specified.

### F-04 — Provenance/seed/attachment classification is correct but needs a ledger anchor

**Muse severity:** Medium. **Disposition: Valid.**

**Independent evidence.** The master plan correctly labels PROP-F3C-01/03/04 as proposed, not adopted. The distinction from existing provenance uses is missing. F2-B4 (`PROJECT_ISSUES_AND_EXPANSION.md:582`) covers evidence-based model/template detection; FUT-C5 (`:609`) covers sampler provenance/research; FUT-ST2 (`:610`) covers imported-field provenance. Image metadata records prompt/model/options/seed (`internal/image/types/types.go`), while chat history does not. `types.Message.Metadata` exists in a separate path but is not the persisted `ChatHistory` shape. The `-1` seed sentinel is global, not per-turn.

**Required correction.** Add a provenance taxonomy:

| Provenance kind | Source/status |
|---|---|
| Model metadata/template provenance | F2-B4, approved evidence model; not per-turn generation proof |
| Sampler/research provenance | FUT-C5, deferred |
| Imported-field provenance | FUT-ST2, deferred |
| Image-generation parameters | Existing image metadata; source-confirmed, not chat provenance |
| Per-turn generation provenance | PROP-F3C-01, proposal arising from the Phase-0 contract gap; not adopted |
| Explicit/random/zero/legacy-unknown seed state | PROP-F3C-03, proposal; not adopted |
| Durable image/attachment identity | PROP-F3C-04, proposal; not adopted |

State explicitly that adoption of one kind does not adopt another.

**Owner decision.** Yes for adopting new durable fields; no for the taxonomy annotation itself.

**Axes.** Requirements, architecture. **Confidence:** Confirmed for source shape facts; classification strongly supported. **Gap:** no per-turn provenance implementation or runtime proof.

### F-05 — STAGE_STATUS closures vs current evidence handled correctly; carry explicitly

**Muse severity:** High if mishandled / Low correction cost. **Disposition: Valid — control confirmation.**

**Independent evidence.** `STAGE_STATUS.md:102,142,187,244` records T1.1–T1.4 as closed on 2026-10-31. F3-C still records split chat/memory ownership and open defects ROOT-01, WIPE-02, IDENTITY-03, REACH-04, STARTUP-05, and ALIAS-06. The master plan already handles this correctly in §4.1, C-01, and §13 item 10: historical closure is not product-wide verification. There is no current critical defect in the plan on this point.

**Required correction.** No substantive correction. Add a standing rule:

> No `STAGE_STATUS` closure may be cited as verification for chat, memory, resource, or ownership scope without current source-or-runtime evidence.

**Owner decision.** No.

**Axes.** Release scope, verification discipline. **Confidence:** Confirmed. **Gap:** no new runtime evidence was produced.

### F-06 — RAG/live-prompt split confirmed; decision still missing

**Muse severity:** High. **Disposition: Valid.**

**Independent evidence.** `app-prompt.go:671` builds the live prompt directly from `App` state. `app.go:186` attaches RAG to `chat.Service`, and `internal/backend/chat/service.go:1053-1054` is the RAG attachment point. The live `buildPromptMessages` path does not demonstrably call RAG. The master plan records the split in S2/S7, C-02, FND-09/FND-10, and §11.1 item 5, but §8.3 places “complete context/memory/prompt foundations” after F3-D–F3-H and does not make one prompt path an explicit F3-E prerequisite.

**Required correction.** Add to §8.3 before F3-E:

> The owner must designate one live prompt assembly path and its owner for RAG, lore sync, tokenizer/context accounting, and character isolation. No third prompt path may be introduced. F3-E cannot be treated as complete while the live path and `chat.Service` remain separate owners of prompt-relevant state.

**Owner decision.** Yes — one live prompt path and its ownership boundary.

**Axes.** Architecture, dependencies, release scope. **Confidence:** Strongly supported. **Gap:** no runtime RAG effect was measured.

### F-07 — OnShutdown supersede confirmed; packaged lifecycle still open

**Muse severity:** Medium. **Disposition: Partially valid.**

**Independent evidence.** The core claim is correct: `main.go:53` registers `OnShutdown: app.OnShutdown`, so the 2026-11-09 “unregistered” finding is superseded. The master plan correctly says lifecycle is improved but not fully verified (S1, FND-16, C-20, §12). Packaged lifecycle, multi-instance behavior, every child type, and cross-process identity remain unproven.

**Correction to the review.** The review’s proposed carry-forward says “W4 (blocked)” and “W8 (open).” Those statuses are stale:

- W4 at `PROJECT_ISSUES_AND_EXPANSION.md:1812` says native acceptance was blocked because `main.go` did not set `WebviewUserDataPath`. Current `main.go:64` sets it under the isolation switch, and W5–W7 record the isolated profile/runtime implementation and verification. **W4 is superseded in substance by W5–W7**; a formal closure record is missing, which is an evidence gap, not a live blocker.
- W8 at `:1840` says unwritable-root/content creation were unexplored, but the later follow-up at `:1864` records W8 as **VERIFIED**. **W8 is not open.**

**Required correction.** Keep “packaged lifecycle remains open” but do not import the stale W4/W8 labels. Record in §8.2/§12:

> W4: historical WebView2-profile blocker, superseded by W5–W7 isolation; formal closure record may be missing. W8: unwritable-root rejection verified; do not list as open. Remaining lifecycle gaps: packaged/multi-instance behavior, all child types, and F3-C native evidence.

**Owner decision.** No. The correction is evidentiary.

**Axes.** Dependencies, release scope. **Confidence:** Confirmed on OnShutdown and W4/W8 statuses. **Gap:** packaged lifecycle remains untested.

### F-08 — Gallery atomicity history handled correctly; video gap correctly held

**Muse severity:** Medium. **Disposition: Valid — control confirmation plus minor correction.**

**Independent evidence.** F2-C reported non-atomic image/sidecar writes, uncontained reads, and suppressed delete errors. F3-A6 at `PROJECT_ISSUES_AND_EXPANSION.md:669` records containment plus atomic image+sidecar commit. Current `internal/image/gallery/gallery.go:72` writes the image and sidecar atomically and removes the image if the sidecar fails; `ResolveFile` enforces gallery containment, and `Delete` returns real errors. The master plan’s C-19 correctly treats the old finding as historical and the repair as not runtime-proven. Video still has no gallery sidecar/catalogue/delete contract (S9, §6, F3-G).

**Required correction.** Keep C-19. Add to §13 item 7: “No video sidecar/catalogue/delete contract exists, so video attachments have no planning basis until F3-G.” No architectural change is required.

**Owner decision.** No.

**Axes.** Architecture, dependencies. **Confidence:** Strongly supported. **Gap:** no native gallery/video acceptance.

### F-09 — “Creative studio” needs adjudication, not dismissal

**Muse severity:** Medium. **Disposition: Valid.**

**Independent evidence.** The ledger north star calls Howl a world generator and inspiration tool. F2-B5 at `PROJECT_ISSUES_AND_EXPANSION.md:583` is **APPROVED**: separate modality-specific creative workspaces; the shared catalog does not merge chat and image interfaces; capability classification keeps chat/image/TTS distinct; the remaining follow-up is explicit capability filtering in UI selectors. E.6/E.8 are parked workspace-layout/staging-board concepts. The master plan calls “creative studio / general asset generation” unestablished without distinguishing F2-B5 from the parked concepts or from a genuinely new general studio.

**Required correction.** Replace the single unestablished row with three categories:

1. **F2-B5 modality-specific creative workspaces** — approved direction, deferred/partial implementation; carry the capability-filtering follow-up.
2. **E.6/E.8 custom layout and image staging boards** — parked, not authorized.
3. **General asset-generation/creative studio beyond F2-B5/E.6/E.8** — unestablished; owner must define before architecture.

Do not merge (3) into (1). Add F2-B5 to the approved/protected decision map.

**Owner decision.** Yes, only for category (3) and any expansion beyond F2-B5. No new decision is needed to preserve the approved F2-B5 direction.

**Axes.** Requirements, release scope. **Confidence:** Strongly supported. **Gap:** F2-B5 implementation/UI selector status was not runtime-verified here.

### F-10 — Avatar/3D classification is correct; desktop-interaction scope must stay explicit

**Muse severity:** Medium. **Disposition: Valid — control confirmation plus minor correction.**

**Independent evidence.** `docs/tier4-avatar.md` is explicitly parked and unscheduled: prebuilt body library, AI-generated skin, TTS-amplitude mouth, Tier B/C/D, external UltraTex, licensing caveat, unverified performance. Current 2D avatar rendering is source-confirmed in character/chat UI. The master plan separates 2D avatars from the parked 3D proposal and lists a “desktop character avatar as a persistent interactive surface” only as an optional future.

**Required correction.** Keep the Tier gating. Add an explicit sentence:

> Any desktop-avatar interaction scope beyond the Tier-4 note is unestablished and must be recovered from the owner; 2D avatar rendering is not evidence of that scope.

Do not promote T4.1/T4.2 or make 3D a core chat dependency.

**Owner decision.** Yes if a desktop interaction scope beyond the parked Tier-4 note is desired; no for preserving the current classification.

**Axes.** Requirements, architecture. **Confidence:** Strongly supported. **Gap:** no 3D performance or TTS-amplitude verification.

### F-11 — Dependency order is sound; lifecycle placement and remote work need guards

**Muse severity:** High. **Disposition: Partially valid.**

**Independent evidence.** The master plan’s §8.2 rows 1–7 correctly place F3-C disposition, root/path identity, and the durable schema decision before F3-D→F3-H; rows 9–12 correctly place model/prompt identity before admission/operations/provenance/memory isolation; media identity precedes attachments and avatars. The cycles in §8.4 are genuine. The weakness is timing: FND-16 lifecycle verification sits at row 8 / §8.3 step 5, after F3-D–F3-H, even though migration, wipe, and final-save evidence depend on a trustworthy shutdown boundary. The remote row is gated only by “separately authorized,” not by MOB-1–MOB-3.

**Required correction.** Do not resequence the accepted F3 order. Add a parallel guard:

> No migration, wipe, final-save, or native lifecycle evidence is accepted unless the shutdown boundary it exercised is recorded and verified. FND-16 verification may run in parallel with F3 work, but cannot be deferred until after evidence that depends on it.

Gate §8.2 row 16 on MOB-1–MOB-3 plus an explicit owner requirement. Make the resource-cycle breaker explicit: requested → loading → ready → failed → unloading → unloaded, with release evidence, not another registry map.

**Owner decision.** Yes for authorizing a remote phase or any actual F3 resequencing. No for the parallel guard, which preserves the existing order.

**Axes.** Dependencies, release scope. **Confidence:** Strongly supported as planning judgment. **Gap:** the actual shutdown boundary has not been packaged-runtime tested.

### F-12 — Release/deferred taxonomy is sound; two scope-leak paths must be closed

**Muse severity:** Medium. **Disposition: Valid.**

**Independent evidence.** §9 correctly separates release obligations, shared foundations, deferred features, and optional futures. The ledger ordering law and T3 gating are preserved. Two weaknesses remain: §9.1’s final bullet (“current ... functionality to the extent independently verified”) is too vague per category, and §9.5 permits architecture influence without explicitly forbidding migration/backfill/release-gate changes by inclusion.

**Required correction.** Add per-category gate references to §9.1 (chat→F3-E; characters→F3-D; images/videos→F3-F/G; lorebooks→F3-B; worlds/scenarios→F3-C). Extend §9.5:

> Reserving a field, identity, state, or boundary creates no migration, no backfill, no compatibility promise, and no release-gate change. Inclusion is not adoption.

**Owner decision.** No.

**Axes.** Release scope, architecture. **Confidence:** Strongly supported. **Gap:** none material.

### F-13 — Protected boundaries preserved; no reopening warranted on this evidence

**Muse severity:** Critical (if violated) — plan complies. **Disposition: Valid — critical control confirmation, not a defect.**

**Independent evidence.** The master plan preserves F3-C freeze, F3-C Phase 0 open, R6 unauthorized, acceptance not granted, F3-D blocked, six D-ROOT defects open, D-STORAGE blocked, V20 blocked, A9/A10/A12/A13 closed, TTS/F2-B7 frozen, REG-1 protected, Gate v2 V1–V17 adopted with V18–V24 proposed, and the local-first boundary. Its §16 protected-artifact hashes match the recorded values. No reopening is warranted on the evidence examined.

**Severity adjudication.** This is **not a current critical defect**. “Critical (if violated)” describes the consequence of a boundary breach, not the plan’s present state. The review expressly says the plan complies and no reopening is warranted. The correct classification is **critical control confirmation**. Do not silently relabel the control as low-risk; preserve “Critical (if violated)” while recording “current finding: no violation.”

**Required correction.** None. Preserve the boundaries and the hash record.

**Owner decision.** No.

**Axes.** Protected boundaries. **Confidence:** Confirmed. **Gap:** no runtime evidence was produced, and none is needed to confirm documentary preservation.

### F-14 — Completeness: ledger coverage is broad but denominated too narrowly

**Muse severity:** Medium. **Disposition: Partially valid.**

**Independent evidence.** HW-C1–HW-C3 at `PROJECT_ISSUES_AND_EXPANSION.md:1772-1774` are deferred validation targets: a hardware compatibility matrix, graceful low-end degradation, and measurements on the named machine. The master plan mentions future hardware measurement only generically and does not name them. E.4/E.5 are recovered in §4.3 but not carried as planning guardrails into error/UI sections. The review’s W4/W8 statuses are stale, as established under F-07.

**Correction to the review.** The HW-C portion is valid. The W4/W8 portion is not: W4 is superseded by W5–W7/current isolation, and W8 is **VERIFIED** at `:1864`, not open. The T2.9.5B F1/F2 audio carry-forwards were not independently re-read in this adjudication; they must not be treated as established constraints until their exact source definitions are verified.

**Required correction.** Add HW-C1–C3 to §8.2/§12 with deferred status; cite E.4/E.5 where UI/error planning touches them; record W4 as historical/superseded and W8 as verified; list T2.9.5B F1/F2 only if verified later. Do not add an audio implementation obligation.

**Owner decision.** No for carrying the constraints; verification ownership is separate.

**Axes.** Requirements, dependencies. **Confidence:** Strongly supported for HW-C; **partially supported** because the review’s W4/W8 statuses are incorrect. **Gap:** T2.9.5B F1/F2 not independently verified.

### F-15 — FND-15 must not plan a remote boundary the product has not authorized

**Muse severity:** Low. **Disposition: Partially valid.**

**Independent evidence.** FND-15 currently says there is no implemented per-connection disclosure, offline switch, connection history, or “remote pairing boundary,” and lists remote access among dependents. That is a gap statement, not an implemented remote design. But the wording can imply a remote-shaped foundation before MOB-1–MOB-3 and any owner authorization.

**Required correction.** Narrow FND-15:

> Current scope is local-first and local-only: model downloads, document imports, and local engine access. No remote API, pairing, authentication, or remote-session boundary is designed or authorized. If a remote phase is later authorized, MOB-1–MOB-3 are mandatory acceptance constraints.

Move remote work to §9.4 optional futures. Do not delete the local-first/offline requirement.

**Owner decision.** No for the wording; yes for any future remote authorization (same decision as F-01).

**Axes.** Architecture, security. **Confidence:** Strongly supported. **Gap:** no remote requirement or protocol exists.

---

## 5. Critical classification reconciliations

| Question | Adjudicated position | Evidence / consequence |
|---|---|---|
| Content hot-swap vs model hot-swap | **Content hot-swap** is REG-1-protected activation/reuse of world, scenario, character, and lorebook context. **Model hot-swap** is unaccepted replacement of the inference server during a conversation, with unresolved stream/readiness/rollback behavior. They must never share an unqualified label. | REG-1 `PROJECT_ISSUES...:700-722`; F3-B8 `:1294`; master plan §4.7, §7.3, §10; no model hot-swap acceptance |
| Product-level requirement vs F3-C contract adoption | A ledger requirement is not made adopted by F3-C, and an F3-C proposal is not made a product requirement by appearing in the plan. The plan must keep the three states distinct: **ledger requirement**, **contract proposal**, **implemented behavior**. | FUT-CH1/CH2/C1 are deferred but real; PROP-F3C-01–05 are proposals; current global model selection is implemented behavior |
| Per-chat persistence union | FUT-CH1 (Author’s Note), FUT-C1 (composer cues), and FUT-CH2 (context-capacity behavior) are deferred per-chat product requirements. PROP-F3C-02 (model preference), PROP-F3C-01/03/04 (provenance/seed/attachment) are proposals. The F3-E schema decision must cover the union without promoting the proposals. | `PROJECT_ISSUES...:605,615,616`; `ChatMeta` has no model field; global `howl-last-model-path` |
| Provenance families | Model/template provenance (F2-B4), sampler provenance (FUT-C5), import-field provenance (FUT-ST2), image-generation metadata, and per-turn generation provenance (PROP-F3C-01) are different. Seed mode (PROP-F3C-03) and attachment identity (PROP-F3C-04) are separate again. | No field may borrow another family’s authority; all proposals remain unadopted |
| MOB-1–MOB-4 and remote/QR | MOB-1–MOB-3 are binding conditional constraints. MOB-4 names Android when a remote phase is authorized. Remote implementation is unauthorized/unstarted. QR/pairing protocol is unestablished. | `PROJECT_ISSUES...:1783-1789`; no implementation found |
| F2-B5 vs E.6/E.8 vs creative studio | F2-B5 is an **approved** direction for separate modality-specific creative workspaces and capability filtering. E.6/E.8 are **parked** layout/staging concepts. A general asset-generation studio is **unestablished** and needs an owner definition. | F2-B5 `:583`; E.6/E.8; master plan §4.8 must split them |
| STAGE_STATUS closures | T1.1–T1.4 closures are historical. They are not product-wide verification for chat/memory/resource scope while F3-C defects remain open. No accepted work should be reopened merely because the closure is old. | `STAGE_STATUS.md:102,142,187,244`; F3-C defects; plan C-01 |
| HW-C1–C3, W4, W8 | HW-C1–C3 are deferred/not-started validation targets and must be carried. W4’s old blocking condition is superseded by W5–W7/current `main.go:64`; a formal closure record may be missing. W8 is **VERIFIED** at `:1864`, not open. | `PROJECT_ISSUES...:1772-1774,1812,1840,1864`; `main.go:64` |
| FND-15 | Current scope is local-only. It may record the absence of a remote boundary, but must not design one or imply authorization. Any future remote boundary is gated by MOB-1–MOB-3 and an owner decision. | FND-15; MOB constraints; no remote implementation |
| F-13 severity | F-13 is a **critical control confirmation**, not a current critical defect. “Critical (if violated)” is the control-consequence label; the review says the plan complies. Do not change it to a low-severity defect or silently erase it. | F3-C freeze/gates; §16 hashes; no reopening |

---

## 6. Ten-correction traceability matrix

| Correction | Review §10 correction | Findings | Classification | Adjudicated action |
|---|---|---|---|---|
| COR-01 | Reclassify MOB-1–MOB-4 as established conditional constraints; gate remote work on MOB-1–MOB-3 | F-01 | **Evidence apply; owner gate for remote** | Replace “unestablished” with “constraints established, implementation unauthorized/unstarted”; keep QR/pairing unestablished. |
| COR-02 | Split content hot-swap from model hot-swap structurally | F-02 | **Evidence apply** | Rename REG-1 to “content hot-swap”; rename PROP-F3C-05 and all model references to “model hot-swap”; define both. |
| COR-03 | Add per-chat persistence union map | F-03 | **Evidence apply for mapping; owner for target set** | Add FUT-CH1/FUT-CH2/FUT-C1 vs PROP-F3C-02 table; require decision before F3-E. |
| COR-04 | Anchor provenance/seed/attachment proposals to the Phase-0 gap, distinct from ledger provenance | F-04 | **Evidence apply; owner for adoption** | Add provenance taxonomy; no promotion. |
| COR-05 | Add STAGE_STATUS non-verification rule; keep one-path prompt decision as F3-E prerequisite | F-05, F-06 | **Evidence apply; owner for prompt path** | Add rule; move the one-path decision into the pre-F3-E gate. |
| COR-06 | Carry W4/W8 and HW-C1–C3; cite E.4/E.5; note audio carry-forwards | F-07, F-14 | **Evidence apply for HW-C/E.4/E.5; history correction for W4/W8; audio unresolved** | Carry HW-C1–C3 and E.4/E.5. **Do not carry “W4 blocked/W8 open.”** Record W4 as superseded by W5–W7/current isolation, W8 as verified; verify T2.9.5B F1/F2 before using them. |
| COR-07 | Keep C-19; add video-catalogue corollary | F-08 | **Evidence apply** | Add F3-G video corollary; no media acceptance claim. |
| COR-08 | Adjudicate F2-B5 vs E.6/E.8 vs unestablished studio; keep avatar Tier gating | F-09, F-10 | **Evidence apply; owner for general studio and desktop-interaction scope** | Split three categories; carry approved F2-B5; preserve Tier 4. |
| COR-09 | Resequence/parallel-guard lifecycle; gate remote design on MOB; specify resource state machine | F-11 | **Evidence/planning apply; owner only for actual remote or F3 reorder** | Do not reorder F3; add parallel shutdown guard, MOB gate, and state machine. |
| COR-10 | Harden §9.1 with per-category gates; add no-migration-by-inclusion; constrain FND-15 | F-12, F-15 | **Evidence apply; protected** | Add gates and no-migration rule; narrow FND-15 to local/MOB constraints. |

### 6.1 Classification of the corrections

**Can be made from existing evidence:** COR-01 (classification), COR-02, COR-03 (mapping), COR-04 (taxonomy), COR-05 (STAGE rule), COR-06 (HW-C1–C3, E.4/E.5, and corrected W4/W8 statuses), COR-07, COR-08 (F2-B5/E.6/E.8 split), COR-09 (guard wording), COR-10 (scope/FND wording).

**Requires further source/history investigation:** the formal W4 closure record; the exact T2.9.5B F1/F2 definitions; the current implementation/UI state of F2-B5 capability filtering; the precise FUT-CH2 persistence semantics.

**Requires an explicit owner decision:** the durable chat target set (including whether proposed fields join FUT-CH1/CH2/C1); one live prompt path; any remote phase/QR protocol authorization; the definition of a general creative studio beyond F2-B5/E.6/E.8; any desktop-avatar interaction scope beyond the Tier-4 note; and any reordering of the F3 sequence.

**Must remain unresolved or protected:** F3-C Phase 0, the six D-ROOT defects, D-STORAGE, V20, R6, F3-D, TTS governance, A9/A10/A12/A13 closure, Gate V18–V24 proposals, and every accepted freeze boundary. The corrections above do not resolve them.

---

## 7. Owner-decision register

These are decisions the evidence cannot make. Existing master-plan decisions are marked “preserved” and are not manufactured anew.

| ID | Decision | Why it cannot be decided by evidence | Current status |
|---|---|---|---|
| OD-01 | Authorize or decline a remote/mobile phase; if authorized, define the remote API and QR/pairing protocol | The ledger establishes only conditional constraints and an Android test-platform designation; it does not establish a remote feature, protocol, auth model, or session contract | **OPEN — no remote work authorized** |
| OD-02 | Define whether a general creative studio/asset-generation capability exists beyond F2-B5 and E.6/E.8 | F2-B5 is approved; E.6/E.8 are parked; the general studio is not in the ledger | **OPEN — general scope undefined** |
| OD-03 | Decide the durable F3-E chat target set across FUT-CH1/FUT-CH2/FUT-C1 and the PROP-F3C proposals | The ledger requires some per-chat behavior; the F3-C target set does not contain the proposed fields. Whether to amend F3-C, plan a later migration, or defer is an owner choice | **OPEN — affects F3-C Phase 0 / F3-E** |
| OD-04 | Designate the one live prompt-assembly path and its owner for RAG, lore, tokenizer/context, and character isolation | Both `App` and `chat.Service` have credible claims; no runtime evidence establishes one as live-path owner | **OPEN — architecture split** |
| OD-05 | Decide whether any F3 sequence reordering is permitted | The review’s lifecycle concern can be met with a parallel guard; changing the accepted order is a separate authorization | **OPEN ONLY IF REORDERING** |
| OD-06 | Define any desktop-avatar interaction scope beyond the parked Tier-4 note | The Tier-4 note is a parked concept, not an interaction specification; 2D rendering is not evidence of 3D/desktop scope | **OPEN — optional future** |
| OD-07 | Resolve canonical root/path identity, wipe semantics, startup degradation, and the six D-ROOT defects | These are existing F3-C owner decisions with explicit defect records and closure authorities | **OPEN — preserved from master plan §11.1** |
| OD-08 | Preserve or change TTS governance; define story-tools scope | F2-B7 freezes TTS; story-tools committed scope is not established | **OPEN — preserved** |
| OD-09 | Decide media/attachment identity and any F3-F/F3-G acceptance scope | The media catalogue and unavailable-reference semantics are not established | **OPEN — preserved** |
| OD-10 | Adopt, reject, or require another revision of the corrected master plan | This adjudication does not adopt the plan; the review recommends re-review after correction | **OPEN — next authority step** |

**No decision above is granted by inclusion in this table.** OD-07/OD-08/OD-09 restate existing open decisions; they are not new requirements.

---

## 8. Additional evidence required

| Evidence | Needed to resolve | Current limitation |
|---|---|---|
| Formal W4 closure record | Whether W4 is explicitly closed, not merely superseded by W5–W7/current isolation | Ledger has W4 BLOCKED and later W5–W7 implementation/verification; no explicit W4 closure line was found |
| W8 confirmation | Whether the later VERIFIED line is the controlling status | `:1840` says OPEN; `:1864` says VERIFIED. The later record controls, but the earlier line should be marked historical in a corrected plan |
| T2.9.5B F1/F2 definitions | Whether audio carry-forwards are real planning constraints | Not independently re-read in this adjudication; do not use until verified |
| F2-B5 implementation state | Whether capability filtering is implemented, partial, or only approved direction | Ledger says capability classification exists; UI selector filtering remains a follow-up |
| FUT-CH2 persistence semantics | Whether context capacity is a per-chat persisted setting or a derived behavior | Ledger requires per-chat behavior; field-level persistence is not specified |
| Chat durable schema | Whether PROP-F3C-01/02/03/04 join the F3-E target set | No adopted fields; `ChatMeta`/`ChatHistory` lack them |
| One live prompt path | Which of `App`/`chat.Service` owns RAG, lore, tokenizer, and isolation | Two source paths; no runtime proof |
| Resource state machine | Whether requested/loading/ready/failed/unloading/unloaded transitions and release are verified | Current resource manager is advisory; no release evidence |
| F3-C gate evidence | Current Phase 0, six defects, D-STORAGE, V20, F3-D positions | Preserved as open/blocked; this adjudication did not run the gate |
| Remote requirement/protocol | Whether any remote feature is wanted and what security boundary applies | MOB constraints exist; no remote feature/protocol exists |
| Creative-studio definition | Whether a general studio exists beyond F2-B5/E.6/E.8 | No dedicated requirement/page/service found |
| Hardware matrix | HW-C1–C3 measurements and graceful low-end behavior | Deferred/not started |

---

## 9. Protected-boundary confirmation

| Boundary | Adjudicated status |
|---|---|
| F3-C architecture freeze | **PRESERVED.** Freeze is not acceptance, verification, or implementation authorization. |
| F3-C Phase 0 | **OPEN.** Not closed by this adjudication. |
| R6 | **UNAUTHORIZED.** No R6 work or text. |
| F3-C acceptance | **NOT GRANTED.** |
| F3-D | **BLOCKED.** |
| D-ROOT / six defects | **OPEN.** No bulk closure, no canonical root selected. |
| D-STORAGE | **BLOCKED.** |
| V20 | **BLOCKED pending F3-F.** |
| A9/A10/A12/A13 | **CLOSED.** No reopening on this evidence. |
| Gate v2 V1–V17 | **ADOPTED.** V18–V24 remain **PROPOSED — NOT ADOPTED**. |
| TTS / F2-B7 | **FROZEN.** No new TTS feature or contract change. |
| REG-1 | **PROTECTED.** Content hot-swap only; model hot-swap is not protected by REG-1. |
| Local-first boundary | **PRESERVED.** No remote/cloud architecture is authorized. |

**Protected F3-C artifact hashes verified unchanged:**

- `ADJUDICATION-PACKAGE...`: `85e1968bac4eb70f489bbd7dc89057debd17e7b3bc2336aa69c036a4352f83e8`
- `OWNER-DECISION-REQUEST...`: `ed8bf2e2326a2ec74838dc0c8cd4de3332b8807d2d0a2647a18cfda060bc9c7c`
- `INVENTORY-F3C-ROOT...`: `343104b528818df699d752c419770fdc1bfd7fe71bc1f249e0dc6847f308c6b8`
- `INVESTIGATION-F3C-DROOT...`: `494539b9dde27a78c9c4c091131a93dce99c1debd81d11fea911e6d411e2e6a6`
- `DISPOSITION-REQUEST...`: `ce32ac0bedbe216b55c065be073b2165209e0ebad4cb435897b829336d29d2c5`

No protected file was edited by this adjudication.

---

## 10. Recommended next action

1. **Do not adopt the master plan unchanged.** The review’s central verdict is sustained, with the corrections above.
2. **Prepare a corrected master-plan revision later, not in this task.** Apply the evidence-only corrections: MOB classification, hot-swap terminology, per-chat union map, provenance taxonomy, prompt-path gate, F2-B5 split, HW-C/E.4/E.5 carriage, FND-15 narrowing, and §9.1/§9.5 hardening.
3. **Correct the review’s stale W4/W8 carry-forward in that revision.** Record W4 as historical/superseded by W5–W7 (with a formal-closure evidence gap) and W8 as verified, not open.
4. **Do not reorder F3.** Add the parallel shutdown-boundary guard; require an explicit owner decision before any actual resequencing.
5. **Send only the owner-decision register to the owner.** Do not ask the owner to ratify evidence-only classifications that the ledger already establishes.
6. **Keep F3-C and all protected boundaries untouched.** No amendment, acceptance, R6, F3-D, migration, or implementation follows from the corrected plan without separate authorization.
7. **Re-review the corrected plan** against the two hashes and the ten-correction matrix before any adoption decision.

---

## 11. Adjudication self-check

- [x] Both supplied artifact hashes were verified exactly.
- [x] The complete 250-line review, all 15 findings, and all 10 corrections were read.
- [x] F-01–F-15 each has an exact heading, severity, disposition, independent evidence, correction, owner status, axis impact, and confidence/gap.
- [x] Content hot-swap and model hot-swap are separated.
- [x] Product requirements, contract proposals, and implemented behavior are separated.
- [x] Per-chat persistence is mapped across FUT-CH1/FUT-CH2/FUT-C1 and PROP-F3C-02.
- [x] Provenance, seed, and attachment proposals are separated from F2-B4/FUT-C5/FUT-ST2 and image metadata.
- [x] MOB-1–MOB-4 are incorporated as conditional constraints; remote implementation and QR remain unestablished.
- [x] F2-B5 is distinguished from E.6/E.8 and from a general unestablished creative studio.
- [x] HW-C1–C3, W4, W8, and FND-15 use their source definitions; W4/W8 stale statuses are corrected.
- [x] STAGE_STATUS closures and every accepted/protected boundary are preserved.
- [x] F-13 is classified as a critical control confirmation, not a current critical defect.
- [x] No artifact was edited; no source/test/build/migration/model action was performed.

**Status:** complete read-only adjudication. The master plan remains unadopted. The review remains a valid independent review with the W4/W8 correction noted. The next action is a corrected planning revision and the owner decisions listed in §7—not implementation.

