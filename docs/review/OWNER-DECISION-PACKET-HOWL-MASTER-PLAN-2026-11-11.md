# OWNER DECISION PACKET — HOWL MASTER PLAN

**Date:** 2026-11-11  
**Type:** Read-only owner-decision preparation  
**Status:** **PREPARATION ONLY — NOT AN OWNER DECISION, NOT MASTER-PLAN ADOPTION, NOT F3-C AMENDMENT, NOT R6, NOT F3-D, NOT IMPLEMENTATION AUTHORIZATION**  
**Basis:** `docs/review/ADJUDICATION-MATRIX-HOWL-MASTER-PLAN-MUSE-REVIEW-2026-11-11.md`  
**Basis SHA-256:** `a9fdaeefdc4ea57c28a706cc09e57a77566f3ed6150d2e4e305904de7b6049ed` — **MATCH**

No master plan, adjudication, source, contract, freeze, acceptance record, or protected artifact was modified.

---

## 1. Purpose and filtering rule

This packet isolates the decisions that genuinely require the project owner’s judgment before a corrected master plan is prepared. It does not ask the owner to decide technical implementation, to ratify facts already established by the ledger, or to reopen protected F3-C boundaries.

A matter is presented as an owner decision only when it materially changes one or more of:

- intended product behavior;
- committed release scope;
- privacy, security, or local/remote boundaries;
- compatibility commitments;
- a durable data/retention promise that users can rely on;
- acceptance of a new product capability rather than a mechanism.

Architecture mechanisms are routed to technical investigation. Existing approved directions are preserved. Unresolved evidence gaps are not turned into requirements.

### 1.1 Decision classification summary

| Class | Count | Items |
|---|---:|---|
| Owner input required before corrected plan | **1** | OD-01: durable chat truth set |
| Genuine owner decisions that can be deferred | **5** | OD-02 model hot-swap behavior; OD-03 remote/mobile phase; OD-04 general creative studio; OD-05 desktop-avatar interaction; OD-06 Story Tools scope |
| Technical investigations | **9** | T-01–T-09 |
| Editorial corrections requiring no owner decision | **10 groups** | COR-01–COR-10 |
| Existing decisions/protected boundaries | Preserved | F3-C, REG-1, F2-B5, F2-B7, Gate v2, A9/A10/A12/A13, D-STORAGE, V20, R6, F3-D |

**Why only one decision is “now.”** The corrected plan can carry every other open product question as explicitly deferred without misrepresenting it as approved. The chat truth set is different: it determines whether the next corrected plan describes the four F3-C proposals as product requirements, as unresolved proposals, or as rejected scope, and it determines the meaning of the F3-E schema dependency. Without that answer, the plan cannot state a stable product-level position.

The owner is **not** being asked to amend F3-C, authorize implementation, or accept a schema. OD-01 asks only what Howl should promise as product behavior; any F3-C amendment, migration, or implementation remains a separate authorization.

---

## 2. OD-01 — Durable chat truth set

**Status:** **DECISION REQUIRED BEFORE CORRECTED PLAN.**  
**Recommendation:** clearly labeled below; not approved.

### 2.1 Exact question

> Which, if any, of the following must be guaranteed product behaviors for Howl chats:  
> (a) per-turn model and engine provenance;  
> (b) explicit-versus-random seed state, including the actual seed used;  
> (c) durable image and attachment references;  
> (d) a per-chat continuation-model preference?

The owner may select all, a subset, none, or defer. The four are separable product behaviors even though they share a schema dependency.

### 2.2 Why this decision matters

The current chat screen can show content that the persisted chat does not fully own:

- The live assistant turn is produced without writing the actual model, engine, effective parameters, or seed into `ChatHistory`.
- `ChatHistory` stores only a `hasImage` boolean, not a durable image/attachment reference.
- Model selection is global (`howl-last-model-path`), not per chat.
- `ChatMeta` has no model field.
- The richer `types.Message.Metadata` type exists in a separate path and is not the persisted chat shape.

The F3-C Phase-0 assessment records these as gaps that the frozen target set cannot preserve. The adjudication classifies them as **proposals arising from the contract gap**, not adopted requirements. The owner decision determines whether they become product requirements.

The consequence is user trust: after reload or model removal, Howl can show a reply but lose which model produced it, which seed was used, or which image was attached. That is materially different from losing an implementation detail.

### 2.3 Relevant historical requirements and evidence

| Evidence | What it establishes |
|---|---|
| F3-C Phase-0 assessment | Frozen target set cannot preserve per-turn model provenance, explicit/random seed state, durable attachment references, or per-chat continuation. These are contract gaps, not adopted fields. |
| `internal/backend/chat/types.go:50` | `ChatMeta` has identity/context fields but no model field. |
| `internal/backend/chat/types.go:72` | `ChatHistory` has text, role, sender, timestamp, `hasImage`, and optional `characterId`; no model/seed/attachment identity. |
| `app-prompt.go` live path | Builds and streams a reply without persisting effective model/seed provenance. |
| `frontend/dist/js/chat-model.js:169,271,393` | Model selection is global `localStorage`, not per chat. |
| `internal/backend/types/message.go` | Richer metadata type exists in a separate path; it is not the persisted `ChatHistory` shape. |
| `PROJECT_ISSUES_AND_EXPANSION.md:605,615,616` | FUT-C1, FUT-CH1, and FUT-CH2 already require deferred per-chat composer cues, Author’s Note persistence, and context-capacity behavior. |
| F3-C freeze and Gate v2 | No amendment, target-set change, acceptance, or implementation is authorized by this packet. |

The four proposed behaviors are **not** the same as F2-B4 model/template provenance, FUT-C5 sampler research, FUT-ST2 import-field provenance, or existing image-generation metadata. None of those existing requirements adopts per-turn chat provenance, per-turn seed state, durable chat attachments, or a per-chat model preference.

### 2.4 Options

#### Option A — Adopt all four as product requirements in principle

Howl would promise that a chat can restore the model/engine provenance of each assistant turn, the explicit/random seed state and actual seed, durable image/attachment references, and the last-selected continuation model.

**Architecture:** All four enter the product-level target set. A separate controlled F3-C amendment analysis would be required before F3-E finalizes the schema. The chat target set must include provenance, seed, attachment identity, and continuation preference. F3-F media identity and model identity/readiness become direct dependencies.

**Release scope:** Expands the intended chat data contract and may extend F3-C/F3-E/F3-F work. It does not itself authorize implementation or acceptance. It protects the “if it is in the chat screen, it should come back” expectation most completely.

**Dependencies:** F3-C Phase 0; F3-E; F3-F; model identity; all assistant paths; attachment identity; future model-switch behavior.

**Risk:** More schema and migration work; the F3-C target set must be amended through a separate authorized process.

#### Option B — Adopt visible continuity only: durable attachments and per-chat model preference

Howl would restore the image/attachment and the chat’s last-selected model, but would not promise historical per-turn model identity or seed state.

**Architecture:** Requires durable attachment identity and a per-chat model pointer. Avoids the full provenance object, but still touches F3-E/F3-F and model lifecycle. Historical model removal cannot be explained from the chat.

**Release scope:** Smaller than Option A; preserves the most visible continuity, but leaves “which model/seed produced this?” unresolved.

**Dependencies:** F3-E, F3-F, model identity/readiness, per-chat settings.

**Risk:** A later request for provenance/seed creates another schema migration.

#### Option C — Adopt durable attachment references only

Howl would guarantee that a visible image/attachment survives reload; model identity, seed state, and per-chat continuation remain deferred.

**Architecture:** Concentrates the schema change on attachment identity and the F3-F media contract. Avoids a model pointer and generation-provenance object.

**Release scope:** Narrowest visible-trust improvement. Model switching remains global; historical turns do not identify their model; seed state remains unowned.

**Dependencies:** F3-F; F3-E attachment field; media catalogue; unavailable-reference semantics.

**Risk:** The chat can still lose model/seed meaning; a later provenance decision causes another migration.

#### Option D — Defer all four; retain current behavior

Howl keeps global model selection, `hasImage`-only chat history, and no per-turn model/seed provenance. The master plan records the known loss and the later migration cost.

**Architecture:** No F3-C target-set change for these behaviors now. F3-E can proceed on the frozen target set. The plan must state that model provenance, seed mode, durable attachments, and per-chat continuation are not product guarantees.

**Release scope:** Smallest immediate scope and fastest path to a corrected plan. It knowingly leaves a visible/persisted-data gap.

**Dependencies:** None added. A future decision would require a new migration and compatibility analysis.

**Risk:** The product can show a reply or image that cannot be fully reconstructed; model removal/switch explanations are unavailable. This is the status quo, not a neutral free option.

#### Option E — Other subset

The owner may choose any combination, for example provenance + attachments but not per-chat model, or per-chat model + attachments but not seed. The plan must then map each selected behavior separately; none is bundled by implication.

### 2.5 Architectural and release comparison

| Behavior | Option A | Option B | Option C | Option D |
|---|---|---|---|---|
| Per-turn model/engine provenance | Yes | No | No | No |
| Explicit/random seed + actual seed | Yes | No | No | No |
| Durable image/attachment reference | Yes | Yes | Yes | No |
| Per-chat continuation model | Yes | Yes | No | No |
| F3-C target-set change needed | Yes, via separate amendment analysis | Yes, partial | Yes, attachments only | No |
| F3-E schema impact | Full | Moderate | Moderate | Frozen target set |
| F3-F media dependency | Yes | Yes | Yes | Existing boolean only |
| Later migration risk | Lowest for these four | High for omitted provenance/seed | High for model/seed | Highest/known |

### 2.6 Recommendation — not an approved decision

**Recommendation: Option A, all four as product requirements in principle, with no implementation authorized.**

Reasons:

1. The four behaviors are the documented product gaps between the chat screen and persisted chat.
2. Adopting only visible attachments or continuity would still leave model/seed history unowned and likely create another migration.
3. The plan’s own trust principle and the F3-C Phase-0 assessment point to the complete set, not a UI-only patch.
4. “Adopt in principle” is not “amend F3-C” and not “implement.” It gives the corrected plan a stable product position while preserving every gate. The owner would still need a separate F3-C amendment analysis and authorization before any target-set change.

If the owner does not want to expand F3-C scope now, **Option D with an explicit migration-cost statement** is the honest alternative. It should not be described as preserving the current chat without loss.

### 2.7 Dependencies and affected decisions

- **F3-C Phase 0:** remains open; no amendment is authorized by this packet.
- **F3-E:** schema cannot be treated as final until the owner’s selection is known.
- **F3-F/F3-G:** attachment identity and video/media references depend on the selected attachment option.
- **F3-D:** remains blocked.
- **Model lifecycle:** per-chat continuation and provenance depend on model identity and readiness; seamless model hot-swap is not implied.
- **Per-chat deferred requirements:** FUT-CH1/FUT-CH2/FUT-C1 remain separate deferred requirements and are not automatically adopted by selecting Option A.
- **No owner decision is implied for** one prompt path, resource state machine, root identity, wipe semantics, or F3 sequence. Those remain technical/F3-C governance matters.

---

## 3. Decisions that can be deferred without blocking the corrected plan

These are genuine product or scope decisions, but the corrected plan can carry each as explicitly open without misrepresenting it as approved. None should be decided by default, silence, or inclusion.

### OD-02 — Model hot-swap behavior

**Exact question.** Should Howl support changing the loaded model during an active conversation?

**Why it matters.** Content hot-swap under REG-1 is already protected and does not answer this question. Model replacement is a different behavior: current `LoadModel` can stop the inference server while a stream is active, and there is no rollback after the old server is gone. The F3-C Phase-0 assessment calls model hot-swap adjacent product behavior, not an F3-C target-set requirement. It must not be conflated with protected content hot-swap.

**Options.**

- **A — No seamless model hot-swap.** The user must stop or finish generation before switching. The switch is explicit, the old stream is cancelled, readiness/failure is reported, and the UI never says loaded before ready. This is the safest bounded behavior.
- **B — Cancel-first safe switch.** Allow a switch during a conversation only after cancelling the active stream, with an explicit confirmation and no claim of preserving partial output. Readiness and failure publication are required.
- **C — Seamless hot-swap.** Queue or restart generation across the replacement, with readiness/failure/rollback guarantees. This is the largest behavior and architecture commitment.
- **D — Defer the decision.** Keep the current manual path and record it as unsafe/unaccepted; do not present it as a product feature.

**Recommendation (not approved).** Option A or B now: no seamless switch, explicit stop/cancel, and no false loaded state. Treat C as a separate future capability. This preserves content hot-swap while avoiding an unproven runtime promise.

**Consequences.** A/B are smaller and compatible with current manual loading, but still need lifecycle verification and truthful status. C requires stream coordination, rollback, model-state machine, and likely chat provenance/continuation work. Deferring leaves the current behavior as an explicit limitation, not a feature.

**Release/architecture.** No immediate F3-C target-set change. Seamless hot-swap would add FND-16/FND-05/FND-06 work and affect per-chat model restoration. It is not blocking the corrected plan.

### OD-03 — Remote/mobile phase

**Exact question.** Should Android/remote access remain an optional future phase, be removed from the roadmap, or be authorized for discovery now?

**Why it matters.** MOB-1–MOB-3 are already binding constraints: desktop stays fully usable without a phone, network, or companion; mobile never becomes required; connecting a phone extends capability but never unlocks desktop function. MOB-4 names Android as the test platform when a phase is authorized. No remote API, pairing protocol, authentication model, or QR contract exists. This is a product/security boundary, not a technical detail.

**Options.**

- **A — Retain conditional future, no work.** Keep MOB-1–MOB-4 as constraints; no remote implementation or protocol design. Corrected plan treats remote as optional and unauthorized.
- **B — Authorize discovery only.** Investigate a remote API/QR pairing/auth model without implementation. This still needs a separate authorization and security review.
- **C — Remove remote/mobile from the roadmap.** Keep desktop-only scope and do not reserve remote interfaces.
- **D — Defer the decision.** Same as A for planning, but explicitly revisitable.

**Recommendation (not approved).** Option A/D: retain the constraints, authorize nothing, and keep QR/pairing unestablished. The local-first product does not need a remote decision to prepare the corrected plan.

**Consequences.** A/D leave FND-15 local-only and avoid speculative security work. B opens a new design/security stream. C removes a future option and should not be chosen merely to simplify planning.

**Release/architecture.** No current release impact. Any remote work would require MOB-1–MOB-3 acceptance, a remote security boundary, and an explicit owner authorization.

### OD-04 — General creative studio

**Exact question.** Beyond the approved F2-B5 modality-specific workspaces and the parked E.6/E.8 concepts, is a general asset-generation or creative-studio capability a product direction?

**Why it matters.** F2-B5 is already approved: separate modality-specific creative workspaces; a shared catalog does not merge chat and image interfaces; capability filtering is required. E.6/E.8 are parked layout/staging-board ideas. The mission label “creative studio” is not itself a requirement. If it is treated as one, it could silently absorb approved-but-unimplemented work or create a new product surface.

**Options.**

- **A — No general studio; carry F2-B5 only.** Keep E.6/E.8 parked and the general studio unestablished. This is the narrowest honest position.
- **B — Define a general studio now.** The owner specifies the product surface, scope, and expected outputs; architecture can then reserve catalog/job/media dependencies.
- **C — Fold E.6/E.8 into the studio.** This promotes parked concepts and requires explicit review; do not do it by implication.
- **D — Defer.** Keep the categories separate and revisit later.

**Recommendation (not approved).** Option A/D: preserve F2-B5 as approved direction, keep E.6/E.8 parked, and leave a general studio unestablished. This avoids scope invention while preventing amnesia about the approved workspace direction.

**Consequences.** A/D require no new architecture; B would add a product surface and media/job requirements; C would promote parked layout concepts. None changes current release gates.

### OD-05 — Desktop-avatar interaction

**Exact question.** Is a persistent interactive desktop character avatar required beyond the parked Tier-4 reactive-avatar note?

**Why it matters.** The Tier-4 note is a parked concept with a prebuilt body library, AI-generated skin, TTS-amplitude mouth, Tier B/C/D, external UltraTex, licensing caveat, and unverified performance. Current 2D avatars are not evidence of a 3D or desktop-interaction scope. The TTS dependency is unresolved.

**Options.**

- **A — Keep parked.** No desktop-interaction requirement; Tier 4 remains unscheduled and optional.
- **B — Define Tier B only.** Approve an app-side projected portrait or limited 2D/3D presentation without full reactive body animation.
- **C — Authorize full Tier C/D.** Commit to external/3D avatar work and its resource/performance costs.
- **D — Defer.** Preserve the parked note and decide later.

**Recommendation (not approved).** Option A/D. Keep 3D/desktop interaction parked and explicitly non-core; do not let it shape F3-E or media foundations as a required dependency.

**Consequences.** A/D preserve the current scope and avoid unverified performance/resource commitments. B/C would create new media, TTS, WebView, and hardware dependencies. None is required for the corrected plan.

### OD-06 — Story Tools committed scope

**Exact question.** Is Story Tools a committed release surface, and if so, what specific capabilities define it?

**Why it matters.** `story-tools*.js` and the page exist, but the adjudication records the committed scope as unknown. A page or controls are not a requirement. Treating every visible control as scope would expand release commitments without evidence.

**Options.**

- **A — Leave scope unknown and non-committed.** Do not use the page as evidence of release scope; keep it as an existing surface with undefined commitment.
- **B — Define a bounded release scope now.** The owner names the committed capabilities and acceptance boundary.
- **C — Remove or de-emphasize the surface.** This is a product decision and must not be made by planning inference.
- **D — Defer.** Preserve the unknown status and revisit before release planning.

**Recommendation (not approved).** Option A/D. Keep the scope unknown; the corrected plan can state that no committed Story Tools scope is established. This is not blocking.

**Consequences.** A/D avoid inventing features. B adds release scope and acceptance work. C changes the visible product surface and requires explicit owner intent.

---

## 4. Technical investigations — no owner decision required

These are evidence, architecture, or implementation-ownership questions. They can change how a requirement is satisfied, but they do not change intended product behavior, privacy, compatibility, or release commitments. They should be resolved by DeepSeek/source investigation under existing requirements.

| ID | Investigation | Why it is technical, not an owner product decision |
|---|---|---|
| T-01 | Formal W4 closure record; confirm W4 is superseded by W5–W7/current isolation | W4 is a native-runtime evidence/status question. The product already has isolation behavior; this is record reconciliation. |
| T-02 | Exact T2.9.5B F1/F2 audio carry-forward definitions | The adjudication did not establish these as constraints. Verify before using; no audio scope change follows. |
| T-03 | Current F2-B5 capability-filtering implementation/UI state | F2-B5 is already approved. The remaining question is implementation completeness, not whether to adopt the direction. |
| T-04 | FUT-CH2 persistence semantics: per-chat setting vs derived behavior | The ledger already requires per-chat behavior; only the representation is unresolved. |
| T-05 | One live prompt path: `App` vs `chat.Service` ownership of RAG, lore, tokenizer, and isolation | This is architecture ownership. The product requirement that configured RAG/lore actually affect live replies should not be weakened; the mechanism is technical. |
| T-06 | Resource state machine and release evidence | Requested/loading/ready/failed/unloading/unloaded and actual release are implementation/evidence questions, not a new product promise. |
| T-07 | Media/attachment identity and unavailable-reference semantics | F3-F/F3-G contract design. The owner decision is only the product promise in OD-01; the identity mechanism is technical. |
| T-08 | Hardware matrix HW-C1–C3 and graceful low-end measurement | These are deferred validation targets, not a choice between product behaviors. |
| T-09 | F3-C native/lifecycle evidence and packaged acceptance | Existing F3-C gate work; no new owner decision and no reopening. |

**Do not escalate T-01–T-09 to the owner merely because they are unresolved.** Escalate only if an investigation reveals a change to product behavior, privacy, compatibility, or a release commitment.

---

## 5. Editorial corrections — no owner decision required

These can proceed in a corrected master-plan revision without owner input. They correct classification, terminology, traceability, and status; they do not adopt any feature or change a gate.

| Correction | Action | Findings |
|---|---|---|
| COR-01 | Reclassify Android/remote as MOB-1–MOB-4 conditional constraints, implementation unstarted; QR/pairing unestablished | F-01 |
| COR-02 | Rename content hot-swap (REG-1) vs model hot-swap (unaccepted behavior) | F-02 |
| COR-03 | Add per-chat persistence union map; no promotion of proposals | F-03 |
| COR-04 | Add provenance taxonomy; keep per-turn/seed/attachment proposals separate | F-04 |
| COR-05 | Add STAGE_STATUS non-verification rule and prompt-path gate wording | F-05/F-06 |
| COR-06 | Carry HW-C1–C3 and E.4/E.5; correct W4/W8 statuses; leave T2.9.5B unverified | F-07/F-14 |
| COR-07 | Keep C-19; add video F3-G corollary | F-08 |
| COR-08 | Split F2-B5 vs E.6/E.8 vs general studio; preserve avatar Tier gating | F-09/F-10 |
| COR-09 | Add parallel lifecycle guard, MOB remote gate, resource state machine; do not reorder F3 | F-11 |
| COR-10 | Add per-category release gates, no-migration-by-inclusion, and local/MOB FND-15 wording | F-12/F-15 |

**Important:** COR-03, COR-04, and COR-05 can add mapping and wording from existing evidence, but they do not answer OD-01 or adopt the proposed schema fields. The corrected plan must leave the product decision visibly open until the owner answers.

---

## 6. Existing accepted decisions and protected boundaries — do not reopen

| Item | Status | Packet treatment |
|---|---|---|
| F3-C architecture freeze | **GRANTED; NOT ACCEPTANCE** | Preserve; no amendment |
| F3-C Phase 0 | **OPEN** | Preserve; OD-01 may affect future analysis but does not close it |
| F3-C acceptance | **NOT GRANTED** | No implication |
| R6 | **UNAUTHORIZED** | Do not draft, imply, or authorize |
| F3-D | **BLOCKED** | No downstream authorization |
| Six D-ROOT defects | **OPEN** | Existing F3-C disposition process; not part of this packet |
| D-STORAGE | **BLOCKED** | Preserve |
| V20 | **BLOCKED pending F3-F** | Preserve |
| A9/A10/A12/A13 | **CLOSED** | Do not reopen |
| Gate v2 V1–V17 | **ADOPTED** | Preserve |
| Gate v2 V18–V24 | **PROPOSED — NOT ADOPTED** | Preserve |
| REG-1 | **PROTECTED** | Content hot-swap only; model hot-swap is not protected |
| F2-B5 | **APPROVED** | Carry as approved direction; no new owner decision |
| MOB-1–MOB-4 | **CONSTRAINT RECORDED** | Carry; no remote implementation |
| TTS / F2-B7 | **FROZEN** | No new TTS feature or contract change |
| Local-first boundary | **CURRENT PRODUCT PRINCIPLE** | No remote/cloud work by inference |

**Existing root, wipe, startup, and media decisions are not reopened here.** They remain in their existing F3-C/F3-F governance processes, with their own owners and closure authorities. This packet cannot preempt them.

---

## 7. Recommended path

1. **Owner answers OD-01 only.** Select the product-level chat truth set, or explicitly defer it and accept the stated migration cost. Do not treat the recommendation as approved.
2. **DeepSeek completes T-01–T-09** without changing product scope. Resolve the W4 record, F2-B5 state, prompt-path ownership, and resource/media implementation evidence.
3. **Prepare the corrected master plan.** Apply COR-01–COR-10; keep OD-02–OD-06 explicitly deferred; preserve all protected boundaries.
4. **Re-review the corrected plan.** Use the two artifact hashes, the ten-correction matrix, and the OD-01 answer. No adoption follows automatically.
5. **Do not authorize F3-C amendment, implementation, migration, F3-D, R6, or acceptance from this packet.** Those require separate decisions.

---

## 8. Decision register

### 8.1 Requiring owner input now

| ID | Decision | Recommendation (not approved) |
|---|---|---|
| OD-01 | Select the durable chat truth set: per-turn provenance, seed state, durable attachments, per-chat continuation model — all, subset, or none | All four in principle; separate F3-C amendment analysis and no implementation authorization |

### 8.2 Can be deferred without blocking the corrected plan

| ID | Decision | Default planning treatment |
|---|---|---|
| OD-02 | Model hot-swap behavior during active conversation | Record as unaccepted; recommend stop/cancel-first, no seamless hot-swap |
| OD-03 | Remote/mobile phase and QR/pairing | Retain MOB constraints; authorize no remote work; QR unestablished |
| OD-04 | General creative studio beyond F2-B5/E.6/E.8 | Carry F2-B5; keep E.6/E.8 parked; general studio unestablished |
| OD-05 | Desktop-avatar interaction beyond Tier 4 | Keep parked; no desktop-interaction commitment |
| OD-06 | Story Tools committed scope | Leave unknown/non-committed; do not infer from page controls |

### 8.3 Technical investigations still required

T-01 W4 closure record; T-02 T2.9.5B F1/F2; T-03 F2-B5 UI state; T-04 FUT-CH2 semantics; T-05 one prompt path; T-06 resource state machine; T-07 media identity; T-08 hardware matrix; T-09 F3-C native/lifecycle evidence.

### 8.4 Corrections that can proceed without owner input

COR-01 MOB classification; COR-02 hot-swap terminology; COR-03 per-chat union map; COR-04 provenance taxonomy; COR-05 STAGE/prompt rule; COR-06 HW-C/E.4/E.5 and corrected W4/W8; COR-07 video corollary; COR-08 F2-B5/E.6/E.8 split; COR-09 lifecycle/remote/resource guard; COR-10 release gates/no-migration/FND-15.

### 8.5 Must not be reopened

F3-C freeze/Phase 0/acceptance; R6; F3-D; six D-ROOT defects; D-STORAGE; V20; A9/A10/A12/A13; Gate v2 V1–V17 and proposed V18–V24; REG-1 content hot-swap; F2-B5; MOB-1–MOB-4; TTS/F2-B7; local-first boundary.

---

## 9. Status

This packet is complete preparation only. It does not decide OD-01–OD-06, adopt the master plan, amend F3-C, authorize implementation, migration, testing, builds, R6, F3-D, or any new requirement. The corrected master plan can be prepared after the evidence-only corrections and the OD-01 answer, while OD-02–OD-06 remain visibly deferred.
