# Howl Chat — Foundation Contract Pack and DeepSeek Intake

## Planning packet `HF-CONTRACT-INTAKE-2026-11-20-01`

**Status:** **Owner decision register recorded; FC-1–FC-6 review-complete and accepted as reviewed proposed drafts (2026-11-21).** Decision-preparation and contract-paperwork packet. Review acceptance is a record that the papers were reviewed; it does not approve FC-1–FC-7, authorize implementation, modify a schema/API, select a migration, close an audit, or report a test result. The owner decisions recorded below are direction, not implementation authorization.

**Purpose:** Give DeepSeek one bounded, traceable intake: the system boundary, papers to prepare, the recorded owner decisions, the remaining research/delegated decisions, and the first implementation authorization that may be considered afterward.

**Parent plan:** `HOWL-FOUNDATION-FRONT-TO-BACK-ARCHITECTURE-AND-DELIVERY-PLAN-2026-11-20.md` (`HF-ARCH-PLAN-2026-11-20-01`).

**Current resource baseline:** `HOWL-RESOURCE-ARBITRATION-AND-VERIFIED-HANDOFF-DISCOVERY-BASELINE-2026-11-21.md` (`HRA-DB-2026-11-21-01`).

**M0A acceptance gate:** `HOWL-CHAT-M0A-RELEASE-TRUTH-ACCEPTANCE-REVIEW-2026-11-21.md` (`HCDM-M0A-AR-2026-11-21-01`).

---

## 0. Accepted review corrections (2026-11-21)

**These corrections were accepted before the owner decision register in §0.2. They do not by themselves approve FC-1–FC-7 or authorize M0A.**

| # | Accepted correction | Effect in this packet |
|---|---|---|
| 1 | **HRA baseline inheritance** | FC-5 and D6 must extend `HRA-DB-2026-11-21-01`; they must not create a competing resource vocabulary or re-litigate its unresolved questions by assumption. |
| 2 | **T0.5 independence** | D6 is an admission policy, not proof of process ownership. T0.5 evidence still requires two independent instances with owned children; a second instance refusing heavy work is not T0.5 evidence. |
| 3 | **M0A independent eligibility** | M0A may be authorized as the release-truth gate without waiting for D1–D8. It stops and returns to architecture planning if verification requires a state/identity contract change. |
| 4 | **Explicit traceability** | This packet now names the original scope, readiness, program-map, Phase 1, resource-baseline and M0A acceptance records it depends on. |
| 5 | **Verification-first M0A** | A1–A9 run first. Passing criteria close without code changes. Code changes occur only for demonstrated failures, within the send/render/append/process-ownership path. |

### 0.1 Traceability links

| Link | Artifact |
|---|---|
| Original scope | `HFS-AUDIT-2026-11-21-01`; `PROJECT_ISSUES_AND_EXPANSION.md`; `ISSUES.md` |
| Readiness evidence | `HFR-CP-2026-11-21-01` |
| Program map | `HCDM-2026-11-21-01` |
| Phase 1 | `HF1-PLAN-2026-11-21-01` |
| Resource arbitration and corrected modality ledger | `HRA-DB-2026-11-21-01` |
| M0A acceptance | `HCDM-M0A-AR-2026-11-21-01` |
| Prior hardware design | `AMENDMENT-01-UNIFIED-HARDWARE-INTELLIGENCE-2026-11-11.md`; `DESIGN-HARDWARE-RESOURCE-ORCHESTRATION-2026-11-11.md` |

Every contract paper must preserve this chain:

**Original requirement → recovered evidence → owner decision → architecture constraint → design decision → work package → acceptance test → implementation evidence.**

### 0.2 Owner decision register (2026-11-21)

| Decision | Direction | Status |
|---|---|---|
| **D1** | Option C, with B as the materialized domain view | **Selected** |
| **D2** | Option A as currently directed: retain user data by default; authored objects are removed only by explicit manual deletion from inside Howl; no cascading deletion; orphanhood confers no deletion authority | **Selected**; the earlier “deletion only when an object is orphaned” shorthand is superseded (FC1-15); `orphaned` remains undefined for legacy, unknown-origin and interrupted-lifecycle records; the directed cleanup of memory originals, independent retention of memory copies, and deletion of conversation-owned snapshots recorded in §3 narrow the prior blanket “no automatic orphan cleanup” wording |
| **D3** | Canonical library owns authored objects and revisions; conversations hold immutable references and reproducibility snapshots; browser storage becomes compatibility/import only after a bounded transition | **Selected** |
| **D4** | Research, identify, and adopt the best model-compatible tokenizer and prompt-construction approach for Howl Chat's requirements, regardless of implementation difficulty; model-specific adapters behind one consistent Howl interface | **Owner direction established**; technology selection remains research/contract work |
| **D5** | Hard flush on speaker change; keyed reuse only after adapter-specific evidence | **Selected** |
| **D6** | Option A initially, explicitly temporary until B or C passes recovery and concurrency acceptance | **Selected**; instance admission policy, not a one-resident-model rule |
| **D7** | Local TTS with unique, user-adjustable, per-character voices preserved across character switching and conversation restoration | **Requirement clarified** |
| **D8** | Delegate the Howl ↔ MEMORY-09 lifecycle-fact decision to the GPT and DeepSeek agent/researcher process | **Delegated** |

**History rule.** The options, questions and recommendations in §3 are preserved as the original packet history. A recorded owner decision supersedes a recommendation without deleting it. No later paper may cite a superseded recommendation as the current decision.

**Authority rule.** Recording these decisions does **not** approve FC-1–FC-7 and does **not** authorize M0A. D4 remains research-dependent; D8 remains a delegated process; D2’s `orphaned` term remains undefined — FC-1 records the definition as unresolved and does not supply one.

### 0.3 Decision-to-paper impact

| Decision | Papers that must reflect it | Required treatment |
|---|---|---|
| D1 | FC-1, FC-2, FC-7 | Event ledger is durable; state/branch structures are derived and reconstructible. |
| D2 | FC-1, FC-4, FC-7 | Retention by default and explicit manual deletion for authored content; FC-1 records the three directed lifecycles (authored content, memory originals and copies, conversation snapshots); FC-4 preserves the snapshot/deletion boundary without authorizing a general background orphan sweep; FC-7 remains subject to the unresolved lifecycle boundary. |
| D3 | FC-2, FC-4 | Canonical ownership and revision rules; browser storage and snapshots classified explicitly; FC-4 distinguishes its read-only authored-content display copies from FC-1's independently owned memory copies and establishes no memory-record authority. |
| D4 | FC-2 | Prompt plan must be preflight, model-compatible, reproducible and fail-closed; technology selection requires research and approval. |
| D5 | FC-3 | Hard flush by default; keyed reuse only with adapter-specific evidence. |
| D6 | FC-5 | Admission policy extends HRA; temporary exit condition is recovery/concurrency acceptance, not implementation completion. |
| D7 | FC-1, FC-3, FC-5, FC-6 | Local-only TTS, per-character voice identity, user adjustment, preservation across switching/restoration, and multi-character acceptance evidence. |
| D8 | FC-7 | Delegated decision record; disagreement returns to the owner. |

### 0.4 Contract paper stage (2026-11-21)

FC-1 through FC-6 completed their review sequence on 2026-11-21 and are accepted as reviewed proposed drafts in `docs/architecture/contracts/`. The register records their dependencies, review status and pinned hashes. FC-7 is deliberately not filed because D8 remains a delegated decision process.

Review acceptance does **not** approve the papers, authorize M0A or any implementation slice, select a tokenizer or resource broker, define `orphaned`, or resolve the `DQ-HRA-*` questions. The six papers remain `PROPOSED`. FC-1–FC-6 are review-complete. Their reviewed hash baseline is recorded in the contract register; FC-1's recorded hash reflects the authorized OD-D2-COPY memory-copy pass and two correction passes after the original baseline commit, FC-4's reflects the authorized D2 wording reconciliation and a later FC-1 copy-semantics consistency pass, and FC-2's reflects the authorized composition review. The register's §6 D2 scope-transition record documents that transition.

## 1. Non-negotiable boundaries

1. Howl owns chat state, branching, regeneration, retention, deletion, and the declaration of lifecycle facts. MEMORY-09 does not infer those facts. Howl owns the memory-transfer workflow, the user's authorization and the destination copy write; MEMORY-09 exposes no transfer operation and no destination-write interface, and that boundary (D-4, Option B) is unchanged.
2. MEMORY-09 owns memory representation/retrieval consequences within its own boundary. Howl does not edit or classify its records.
3. Selection is not allocation; process start is not readiness; unload return is not verified release; UI update is not durable commit.
4. A display name, browser-local record, snapshot copy, or legacy folder is never silently promoted to canonical identity.
5. Non-destructive chat transitions retain material. Only an explicit, named deletion set is destructive.
6. TTS remains deferred but required for chat finalization. It must be local-only, with unique, user-adjustable per-character voices preserved across character switching and conversation restoration. Image, vision, video, ASR/audio input, and audio generation remain separate modality capabilities under one lifecycle/resource contract. 2D/3D visual-character rendering is out of scope.
7. No package may call a historical test/worklog/DeepSeek report a current result without fresh, gate-defined evidence.

## 2. Papers DeepSeek should prepare

| Paper | Output | Must not decide alone | Gate served |
|---|---|---|---|
| **FC-1 Conversation State Contract** | Identity model, transition table, event/snapshot duties, idempotency, retention/deletion semantics, recovery. | State unit, retention duration, deletion policy. | M1 |
| **FC-2 Composition and Context Contract** | Canonical reference resolution, PromptPlan, tokenizer/template boundary, segment budgets, multimodal attachments. | Token-overflow product policy; privacy projection policy. | M2 |
| **FC-3 Multi-Character Isolation Contract** | Per-speaker execution context, transcript projection, queue/fairness/cancellation, context flush/reuse evidence. | Public/shared transcript semantics; reuse policy. | M3 |
| **FC-4 Content Ownership and Synchronization Contract** | Canonical store per object class, revision binding, bypass inventory/disposition, frontend/backend synchronization. | Canonical-root choices; legacy migration/quarantine policy. | M0B/M5 |
| **FC-5 Capability, Resource, and Handoff Contract** | Extends `HRA-DB-2026-11-21-01` and the prior hardware design: capability declaration, logical reservation vs physical allocation, admission, telemetry, lifecycle, drain/release/readiness/recovery, evidence ladder, instance coordination. | Reservation authority, concurrent-instance model, hard/soft limits. | M4 |
| **FC-6 Evidence and Acceptance Contract** | Receipt schemas, fault matrix, fixture rules, no-claim ladder, acceptance artifacts. | Product acceptance thresholds. | All gates |
| **FC-7 MEMORY-09 Boundary Contract** | Ordered/idempotent Howl lifecycle facts, replay/failure semantics, deletion-set handoff. | MEMORY-09 applicability representation. | M8 |

Every paper must state: sources, authority/non-authority, invariants, state transition/error tables, durability, observability/evidence, recovery, privacy, compatibility questions, acceptance design, open decisions, and explicit non-authorization.

Each paper must also identify its exact place in the traceability chain in §0.1. A paper that cannot name its source requirement, owner decision, constraint and acceptance gate is not ready for review.

## 3. Owner decision register and preserved agenda

### D1 — Conversation model (blocks FC-1/M1)

**Question:** What is the durable unit of conversation state?

| Option | Meaning | Consequence |
|---|---|---|
| A. Message alternatives only | A message owns alternatives; no explicit branch graph. | Simpler UI; weak for multi-turn divergence/revert. |
| B. Rooted state/branch DAG | Every transition yields a state node with parent relation(s); current state is a pointer. | Directly supports HSL alternatives, regenerate, revert, retention, and reproducible composition. |
| C. Full event ledger with derived branch graph | Events are primary; states/branches are materialized projections. | Strongest recovery/audit story; higher implementation cost. |

**Owner decision (2026-11-21): Selected — Option C, with B as the materialized domain view.**

**Original packet recommendation (preserved):** C as the durable model, with B as its materialized domain view. It is the most direct fit for HSL’s non-destructive transitions and explicit lifecycle notifications.

**Decision wording:** “Howl shall treat an append-only, ordered conversation-event ledger as the durable source of state; state/branch projections are derived and reconstructible. A transition produces a stable state ID. This does not prescribe storage technology or an implementation API.”

### D2 — Retention and deletion promise (blocks FC-1/FC-7)

**Question:** Who may inspect retained alternatives, and when may Howl permanently delete material?

| Option | Consequence |
|---|---|
| A. Retain every non-deleted state indefinitely in local storage. | Strong restoration/audit; unbounded storage/privacy burden. |
| B. Retain until explicit user deletion, with configurable export/archive later. | Matches HSL’s preservation rule; keeps destructive action explicit. |
| C. Automatic retention window. | Requires a new expiration/consent/deletion policy and risks contradiction with restoration. |

**Owner decision (2026-11-21): Selected — Option A.** As currently directed: Howl retains user data by default; authored content is not automatically deleted when a conversation, world, scenario, character or other related object is deleted; manual deletion from inside Howl is the user’s deletion authority; there is no cascading deletion; orphanhood is a condition, not a deletion authorization.

**Superseded wording (history, not a live rule):** the earlier shorthand “deletion only when an object is orphaned” is preserved here for decision traceability and is superseded by the direction above (FC1-15). The packet’s earlier recommendation B is preserved below but superseded.

**Open term:** `orphaned` remains undefined for legacy, unknown-origin and interrupted-lifecycle records. FC-1 records the definition as unresolved (FC1-23, FC1-23a) and does not supply one; deletion semantics for those records are not testable until the owner defines the term. This remains a definition obligation, not a reopening of the decision.

**Deletion mechanics:** No automatic orphan cleanup is authorized as a general background sweep, and no cascading deletion of authored objects is authorized (FC1-22b). A deletion names its set and records a receipt; a deletion that cannot name its set is refused (FC1-24). FC-1 directs one bounded cleanup that is not an orphan sweep: a parent's own original memory records are removed as part of that parent's explicit deletion (FC1-22e), while independently owned copies survive it (FC1-22g, FC1-22i), and conversation-owned reproducibility snapshots are deleted with their conversation (FC1-22c). The status of legacy, unknown-origin and interrupted-lifecycle records remains unresolved (FC1-23a); chat-scoped records and conversation-event/state-history retention are not settled by D2. Boundary 5 still requires an explicit, named deletion set.

**Copy semantics (2026-11-21, OD-D2-COPY):** The operation the source records call *transfer* is a copy: the source record is preserved and remains owned by its original parent, and each destination copy is independently owned. Every destination memory copy carries the literal `transferred_from` marker, which identifies its source through the copied origin references; a copy that cannot carry the marker is not created, and the marker and its provenance must survive deletion of the source without depending on a live source pointer (FC1-22h). The transfer receipt records that marker for each created copy, with source and destination identities consistent with the provenance record (FC1-22l). The earlier move interpretation is superseded by OD-D2-COPY and retained as superseded history in the decision record; the current copy semantics are recorded in FC1-22f. Howl owns the transfer workflow, the user's authorization, the transfer event and the destination write; MEMORY-09 exposes no transfer operation and no destination-write interface, and that boundary (D-4, Option B) is unchanged and is not superseded (FC1-22m). FC-4 makes no memory-record authority claim and separates its read-only authored-content display copies from these independently owned memory copies (FC4-13, FC4-14). Five questions remain open: whether the pre-deletion notice is blocking or advisory (`FC1-MEMORY-COPY-NOTIFY`), whether a world destination owns a copy (`FC1-MEMORY-COPY-DESTINATION`, with destination ownership unresolved in FC1-22g), whether copying is available only at deletion or at any time (`FC1-MEMORY-COPY-TIMING`), whether the notice also covers snapshots (`FC1-MEMORY-NOTICE-SCOPE`), and destination-side merge and conflict (`FC1-MEMORY-COPY-MERGE`).

**Original packet recommendation (superseded):** B. It is the minimum policy consistent with HSL and avoids inventing auto-deletion semantics.

### D3 — Canonical content authority (blocks FC-2/FC-4/M1)

**Question:** For each world, scenario, character, user persona, and lore object, which store is canonical and how are current bypasses handled?

**Required decision form:** For each object class name exactly one canonical write authority, allowed readers, temporary compatibility readers, prohibited writers, revision rule, and recovery/migration disposition. Browser local storage and chat snapshots must be classified explicitly; neither is implicitly canonical.

**Owner decision (2026-11-21): Selected.** Canonical library stores own authored objects and revisions; conversations hold immutable references plus reproducibility snapshots; browser local storage becomes compatibility/import only after a bounded transition.

**Original packet recommendation (preserved):** Canonical library stores own authored objects/revisions; a conversation stores only immutable references plus reproducibility snapshots; browser local storage becomes a compatibility/import view only after a bounded transition plan.

### D4 — Tokenization and prompt construction (blocks FC-2/M2)

**Owner direction established (2026-11-21):** Research, identify, and adopt the best model-compatible tokenizer and prompt-construction approach for Howl Chat’s requirements, regardless of implementation difficulty. Model-specific adapters may sit behind one consistent Howl interface; no single tokenizer or prompt-planning library may be assumed to fit every model family.

The chosen approach must support:

- Accurate token accounting using the actual target model’s tokenizer and chat template.
- Complete prompt construction across system instructions, conversation history, world, scenario, character, lore, retrieved memory, tool definitions, structured payloads, and modality inputs, including their token costs.
- Model-specific formatting and special-token handling.
- Context-window and output-budget validation before inference.
- Reproducible accounting of included, excluded, and transformed context.
- Explicit failure when a request cannot be constructed safely, rather than silently sending incomplete or invalid context.

**Status:** Owner direction established. The specific technology, adapter design, benchmarks, trade-offs, and implementation remain to be researched and approved through FC-2 and the contract process. This is not an implementation authorization.

**Original packet question (preserved):** What happens when a real-tokenizer PromptPlan does not fit?

| Option | Consequence |
|---|---|
| A. Refuse and explain. | Safest/reproducible; interrupts flow. |
| B. Deterministically trim only explicitly low-priority, eligible segments. | Usable; requires per-segment policy and audit receipt. |
| C. Summarize/compress automatically. | Requires a separate capability, provenance, and quality policy. |
| D. Ask the user to choose. | Transparent; more interaction. |

**Original packet recommendation (superseded by owner direction):** A plus B: default deterministic refusal where no explicit eligible trimming policy exists; only policy-labelled segments may be trimmed, and the receipt must name them. Automatic summarization remains out of this foundation decision.

**Research acceptance:** FC-2 must provide a model-family compatibility matrix, a prompt-plan artifact schema, a preflight failure table, and evidence that token counts use the actual target tokenizer/template. A technology recommendation must state alternatives, trade-offs, unresolved assumptions, and revisit conditions.

### D5 — Multi-character context policy (blocks FC-3/M3)

**Question:** What may cross a speaker boundary, and when may runtime context be reused?

**Required owner choices:** public/shared transcript definition; each character’s private visibility rule; whether speaker change is always a hard flush initially; cancellation/fairness behavior.

**Owner decision (2026-11-21): Selected — hard flush on speaker change; keyed reuse only after adapter-specific evidence.**

**Original packet recommendation (preserved):** Start with hard flush on speaker change and prompt-plan reconstruction from durable state. Permit keyed reuse only after an adapter-specific evidence gate proves it cannot carry foreign private context.

### D6 — Resource and multi-instance policy (blocks FC-5/M4)

**Authority:** FC-5 must extend `HRA-DB-2026-11-21-01` and the prior hardware design. D6 chooses a supported multi-instance admission policy; it does not create a separate resource interpretation and it does not decide the `DQ-HRA-*` questions by implication.

**Question:** What is the supported policy for concurrent heavy work across Howl instances?

| Option | Consequence |
|---|---|
| A. One supported heavy-work instance; a second instance is read-only/light-work or refuses heavy admission. | Lowest complexity; prevents silent double-spend immediately. |
| B. Shared local broker owns all machine reservations. | Strong global accounting; adds service lifecycle/recovery work. |
| C. Atomic local lease registry without a permanent broker. | Lighter process model; difficult stale-lock/recovery proof. |
| D. Conservative independent partition. | Cannot guarantee race freedom; must be shown as non-guaranteed. |

**Owner decision (2026-11-21): Selected — Option A initially, explicitly temporary until B or C passes recovery and concurrency acceptance.** This is an instance-level heavy-work admission policy, **not** a one-model-at-a-time residency rule; `HRA-DB-2026-11-21-01` still requires supported multi-model residency to remain possible.

**Exit condition:** B or C becomes eligible only after recovery and concurrency acceptance. Implementation completion, a successful demo, or an apparently working broker is not sufficient evidence.

**Original packet recommendation (preserved):** A as the initial supported policy, explicitly temporary until B or C has passed recovery and concurrency acceptance.

**Separation rule — D6 is not T0.5 evidence.** T0.5 is a process-ownership safety property: two independent instances must not kill each other’s model children, and abrupt exit must kill only the exiting instance’s owned child. T0.5 acceptance requires two independent instances, each with an owned child or an approved owned test child. A second instance that refuses heavy admission has not exercised the ownership boundary and cannot close T0.5.

**Resource rule — D6 is not global accounting.** A supported multi-instance admission policy does not prove race-free logical reservations. `DQ-HRA-07` remains open until a shared authority or an explicit non-guarantee policy is selected and tested.

### D7 — TTS finalization boundary (blocks chat-finalization gate)

**Question:** What exact TTS behavior makes chat “final”?

**Required decision form:** turn-to-voice trigger (manual/automatic), per-character voice ownership, cancellation behavior, playback/output retention, privacy/data locality, failure fallback, and the acceptance evidence required before chat may be called final.

**Owner requirement clarified (2026-11-21):** TTS must be local-only and support unique, locally generated voices selectable per character in multi-character chat.

- No cloud TTS dependency.
- Character-specific voice assignment.
- Distinct voices for different characters in multi-character conversations.
- User-adjustable voice selection for each character.
- Correct voice preservation across character switching and conversation restoration.

**Acceptance rule:** TTS is not complete until this works in the actual multi-character experience. FC-6 must include a multi-character TTS acceptance matrix covering local-only generation, per-character assignment, distinct voices, user adjustment, preservation across switching/restart/restoration, and failure isolation. TTS remains deferred, but it is required for chat finalization.

**Original packet recommendation (preserved):** Treat TTS as a chat capability invoked from a committed assistant turn, never as a pre-commit side effect; retain a link to source turn/model/voice/provenance and cancel it without cancelling durable text.

### D8 — Howl ↔ MEMORY-09 event boundary (blocks FC-7/M8)

**Owner delegation (2026-11-21):** Delegate the Howl ↔ MEMORY-09 lifecycle-fact decision to the GPT and DeepSeek agent/researcher process. The resulting decision record must contain:

1. The selected option.
2. The competing alternatives.
3. The reasoning and evidence.
4. Any unresolved assumptions.
5. Conditions that would cause the decision to be revisited.

**Disagreement rule:** If the two agents disagree, the disagreement returns to the owner. It is not resolved by whichever agent writes the document last.

**Status:** Delegated. No option is selected in this packet. FC-7 cannot finalize its MEMORY-09 boundary until the delegated decision record exists.

**Original packet question (preserved):** Which lifecycle facts does Howl guarantee to publish and replay?

**Minimum fact set:** event ID, conversation/root/state IDs, ordered predecessor relation, operation kind, affected material IDs, current-state result, deletion set where applicable, durability level, timestamp, and idempotency/replay semantics.

**Not a Howl-only choice:** MEMORY-09’s record applicability representation remains jointly specified; Howl does not choose it unilaterally. The delegation covers the Howl lifecycle-fact contract, not unilateral ownership of MEMORY-09’s internal representation.

## 4. Contract invariants for DeepSeek to preserve

1. Stable IDs survive save/load, edits, and message renumbering.
2. A branch/regeneration/alternative never overwrites retained material.
3. A delete command names and confirms its deletion set before durable mutation.
4. Every generated result is attributable to state, character, prompt plan, model/capability, and evidence receipts.
5. Real tokenizer measurement precedes any budget-dependent dispatch.
6. Private character context is never selected by display-name coincidence or shared mutable context.
7. A reservation is not a physical-allocation claim; observed unsafe telemetry wins.
8. Uncertainty is visible and conservative; unknown is never converted into available/ready/released.
9. Capability failure preserves selection and user work; no silent provider/model/placement substitution.
10. A restart reconciles durable conversation/capability records with owned runtime processes before resuming work.
11. No MEMORY-09 lifecycle action occurs without an explicit Howl fact; no Howl memory-record mutation occurs through this boundary.
12. Admission policy and process-ownership proof are separate: a supported multi-instance policy cannot substitute for T0.5 evidence, and T0.5 evidence cannot prove global reservation accounting.
13. A transferred memory copy carries its `transferred_from` marker and copied provenance, which survive deletion of the source parent and its original memory record; a copy that cannot carry the marker is not created, and destination ownership remains unresolved.

## 5. M0A implementation authorization draft

This is the only near-term implementation package that may be considered independently of the foundation decisions. It remains unapproved until the owner expressly adopts it.

**Purpose:** establish release truth for T0.1, T0.2, and T0.5 without changing state/branch/token/authoring architecture.

**Authorized scope if approved:**

- **Verification first.** Run A1–A9 on the current build before changing code. A passing criterion closes without a code change.
- Code changes only where a criterion fails, and only within the user-message render/persistence, duplicate-user-turn, or per-instance process-ownership path. No speculative refactor.
- Deterministic append-delay seam or equivalent controlled test mechanism, only if needed to exercise the actual unawaited frontend ordering.
- Two independent instance test protocol using recorded process/job identity; each instance must exercise an owned child or approved owned test child. A second instance refusing heavy work is not T0.5 evidence. Never pattern-kill.
- A1–A9 from `HOWL-CHAT-M0A-RELEASE-TRUTH-ACCEPTANCE-REVIEW-2026-11-21.md`.
- **Independent eligibility.** M0A does not wait for D1–D8. It stops and returns the issue to FC-1/M1 if the smallest correct fix requires a state or identity contract change.

**Explicitly excluded:** stable IDs, branching/revert redesign, token budgeting, character isolation, canonical-store migration, resource broker work, modality implementation, MEMORY-09 integration, schema/API migration, cache clearing, stress/performance runs.

**Stop rule:** if the smallest correct fix needs a state/identity contract change, stop and return the issue to FC-1/M1.

**Closure rule:** only the specified acceptance evidence may close T0.1/T0.2/T0.5; a code diff, screenshot, or a sequential-only test is insufficient.

## 6. DeepSeek response format

For each paper, DeepSeek should return:

1. Traceability table citing the exact source records and requirement IDs.
2. Current verified facts versus assumptions/proposals/open decisions.
3. Candidate design alternatives and consequences.
4. Contract/state-transition/error/recovery tables with unresolved fields clearly marked.
5. An acceptance/evidence matrix that distinguishes unit, integration, end-to-end, and runtime/performance proof.
6. A change-impact inventory listing likely files, stores, APIs, migrations, fixtures, and users—without editing them.
7. A request for explicit owner decisions before any implementation work package.

## 7. Recommended immediate sequence

1. The owner decision register is recorded and the original options/recommendations are preserved as history. No recommendation is cited as current after being superseded.
2. M0A is independently eligible for authorization as the release-truth gate. If approved, run A1–A9 first; change code only for demonstrated failures. If a state/identity contract change is required, stop and return to FC-1/M1.
3. FC-1 through FC-6 are review-complete as reviewed proposed drafts in `docs/architecture/contracts/`. Their reviewed hash baseline is recorded in the contract register; FC-1's, FC-2's and FC-4's recorded hashes reflect the authorized post-commit memory-copy, composition-review and D2 consistency passes, not the original baseline commit. FC-7 waits for the D8 delegated decision record; FC-2’s D4-dependent choices wait for research and approval.
4. The D4 research produces a comparison, benchmark evidence, trade-offs and a recommended tokenizer/prompt-construction approach. The D8 delegation produces the required decision record. Disagreement returns to the owner.
5. In parallel, prepare but do not execute the M0A test/authorization packet.
6. Reconcile the papers into one contract register and submit specific implementation-authority requests by gate.
7. Execute only the owner-approved M0A slice first; do not use it as a vehicle for Phase 1 redesign.

## 8. Record of this packet

| Item | Status |
|---|---|
| Code/schema/API/migration/test/model/process changes | None. |
| Existing records changed | None; this new packet only, amended with the review corrections, the owner decision register, the 2026-11-21 review-stage record, the 2026-11-21 D2 wording reconciliation and the 2026-11-21 D2 scope transition (OD-D2-COPY). |
| New file | This DeepSeek intake packet only. |
| Review corrections | Four accepted corrections plus the verification-first M0A refinement applied; see §0. |
| Owner decisions recorded | D1, D2, D3, D5, D6 selected; D7 requirement clarified; D4 owner direction established; D8 delegated. |
| Contract paper drafts | FC-1–FC-6 review-complete, accepted as reviewed proposed drafts; FC-7 blocked by D8. |
| Contract review baseline | FC-1–FC-6 frozen at six pinned sha256 hashes; see `docs/architecture/contracts/README.md` §6. |
| Documentation performed | D2 wording reconciled 2026-11-21 across this packet’s D2 records, the contract register §3 and FC-4 FC4-14, then scoped by the OD-D2-COPY memory-copy direction across this packet, the register §3/§6 and FC-1. Documentation only: D2’s unresolved questions remain open. `orphaned` remains undefined for legacy, unknown-origin and interrupted-lifecycle records; chat-scoped records and conversation-event/state-history retention remain unresolved; privacy deletion remains a separate unresolved policy. The D2 direction for authored-object retention by default and explicit manual deletion remains in force. No general background orphan sweep and no cascading deletion of authored objects are authorized; FC-1's directed cleanup of a parent's own original memory records is bounded by that parent's explicit deletion and is not such a sweep. Five memory-copy questions remain open, and destination ownership remains unresolved. The FC-2 composition review was also recorded 2026-11-21, in FC-2 itself and in the contract register §6. |
| Contract papers approved | None. FC-1–FC-7 remain unapproved. |
| M0A authorized | No. |
| Authority state | Owner decisions recorded; contract papers not approved; M0A unauthorized. |

**Boundary statement.** The decision register and review drafts give the contract papers a stable direction. They do not approve FC-1–FC-7, authorize M0A, or convert D4 research or D8 delegation into an implemented design.
