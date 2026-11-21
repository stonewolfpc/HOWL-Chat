# FC-3 — Multi-Character Isolation Contract

## Contract paper `FC-3-2026-11-21-01`

| Field | Value |
|---|---|
| Date | **2026-11-21** (local, CST) |
| Type | Contract paper — review draft |
| Authority | D5 and D7 owner decisions; `HF-CONTRACT-INTAKE-2026-11-20-01` |
| Gate | M3 |
| Status | **PROPOSED — not approved** |
| Implementation authorization | **None** |
| Tests/models/process runs | **0** |

**Labels:** **Owner-approved** · **Proposed** · **Unresolved** · **Observed** · **Unestablished**.

---

## 1. Contract decision

**FC3-01 (Owner-approved — D5).** A speaker change causes a hard flush of transient prompt context. Each character response is composed from durable state, with the active character's identity and visibility rules.

**FC3-02 (Owner-approved — D5).** Keyed runtime-context reuse is permitted only after adapter-specific evidence proves it cannot carry foreign private context. No reuse is authorized by this paper.

**FC3-03 (Owner-approved — D7).** TTS must be local-only and support unique, user-adjustable voices per character in multi-character chat. Voice identity and character configuration are not transient prompt context and must survive the hard flush.

**FC3-04 (Proposed).** Isolation is a property of the composed prompt and execution context, not of a display-name check, a character-card label or a global mutable transcript.

---

## 2. Execution context

**FC3-05 (Proposed — per-response context).**

| Field | Obligation |
|---|---|
| `character_id` | Stable identity; not display name |
| `character_revision` | Exact authored revision used |
| `conversation_id` | Stable conversation identity |
| `state_id` | Exact FC-1 conversation state being composed |
| `scenario_id` / `world_id` | Active composition, stable IDs |
| `visibility_scope` | Public/shared, character-private, user-private or system-only |
| `voice_profile_id` | D7 voice assignment; persisted configuration |
| `private_context_refs` | Character-scoped facts/history |
| `shared_scene_refs` | Explicitly shared transcript/scene material |
| `prompt_plan_id` | FC-2 plan for this response |
| `queue_class` | Interactive/background and fairness class |

**FC3-05a (Proposed — execution metadata is not prompt content).** The execution-context fields are not automatically prompt segments. `voice_profile_id` is a persisted configuration reference used by the TTS path, and `queue_class` is scheduling metadata; neither is model-visible input merely because it appears in the execution context. Prompt content is determined by FC-2 composition, not by the presence of a field here.

**FC3-06 (Proposed).** The active character is selected by stable ID from FC-1 state. A display name may appear in the prompt as a label; it never selects a character, voice or private context.

**FC3-07 (Proposed).** Private context is filtered before tokenization and before any adapter formatting. Removing it after prompt construction is not isolation evidence.

---

## 3. Hard flush semantics

**FC3-08 (Proposed — what a hard flush is).** A hard flush clears the transient prompt working set between speaker changes. It does not delete or rewrite durable conversation history, memory records, character configuration or voice assignment.

**FC3-09 (Proposed — what survives).**

| Survives hard flush | May be cleared |
|---|---|
| Conversation/state/message IDs | Transient prompt working set |
| Authored character revision | Character's prior runtime context |
| Voice profile assignment | Adapter/KV cache reuse, if any |
| Shared scene/transcript selected for inclusion | Transient copies of private context not selected for this character |
| Retention/deletion records | Temporary formatting buffers |

**FC3-09a (Proposed — flush is not deletion).** A hard flush may clear only transient copies, working-set references and temporary representations. Underlying durable private-context records remain governed by their own retention and deletion authority; a hard flush is never a deletion mechanism and does not remove authored or durable records.

**FC3-10 (Proposed — reconstruction).** After a flush, the next response is reconstructed from FC-1 state and FC-2 composition. The system does not rely on a surviving KV cache or a hidden mutable transcript.

**FC3-11 (Calculated — D5/D7 boundary).** Flushing private context must not flush voice identity. A character's voice is configuration referenced by stable ID, not a fact inside the prompt working set.

---

## 4. What may cross a speaker boundary

**FC3-12 (Proposed — default rule).** Only explicitly shared scene/transcript material may cross a speaker boundary. Character-private facts, private history and private memory do not.

**FC3-13 (Unresolved — `FC3-SHARED-TRANSCRIPT`).** The precise public/shared transcript definition is not decided. Candidate scopes include the visible multi-character scene, the user-visible transcript, or an explicit subset. This paper does not select one.

**FC3-14 (Proposed).** Prior responses included as shared scene context are labelled, attributed and bounded. They do not become an unbounded global history available to every character.

**FC3-15 (Proposed).** A shared scene reference is not a licence to include another character's private retrieved memory. Retrieval results carry visibility provenance.

---

## 5. Queue, fairness and cancellation

**FC3-16 (Proposed — queue rules).**

- Each response request names its character, conversation state and priority class.
- Cancellation of one response does not cancel another character's durable turn or rewrite shared history.
- A cancelled or failed response leaves an explicit state; it does not become a silent empty assistant turn.

**FC3-17 (Proposed).** Interruption of TTS playback cancels audio output only. The committed assistant text remains durable. A new response may begin without deleting the interrupted turn.

**FC3-18 (Unresolved).** Fairness policy, maximum queue depth, preemption, multi-character turn ordering and whether any no-starvation guarantee is required need an owner decision. They are not fixed by the D5 hard-flush decision, and FC3-16 does not establish a starvation bound.

---

## 6. D7 voice identity and TTS boundary

**FC3-19 (Owner-approved — D7).** Distinct characters must have distinct voices in multi-character conversations. The user must be able to adjust each character's voice independently.

**FC3-20 (Proposed).** Voice assignment is bound to `character_id` and a versioned `voice_profile_id`. The assignment is persisted independently of prompt context and survives:

- speaker switching;
- hard flush;
- app restart;
- conversation restoration;
- model change, where the voice implementation remains available.

**FC3-21 (Proposed).** TTS is invoked from a committed assistant turn, never as a pre-commit side effect. A TTS failure does not roll back the text turn. A text cancellation does not silently leave an orphan audio job.

**FC3-22 (Proposed — local-only).** The TTS path must not depend on a cloud service. If the local voice engine or voice profile is unavailable, the result is an explicit unavailable/failed state, not a cloud fallback or another character's voice.

**FC3-23 (Owner-approved — D7).** TTS is not complete until the actual multi-character experience supports distinct, user-adjustable voices and correct preservation across switching and restoration. UI selectors alone are not acceptance.

---

## 7. Error and privacy behavior

**FC3-23a (Proposed — failure obligations).**

| Failure | Required behavior |
|---|---|
| Missing character ID | Refuse; do not resolve by display name |
| Missing character revision | `degraded`/`unresolved`; no latest fallback |
| Missing voice profile | Explicit unavailable; no silent substitute |
| Private context leak risk | Fail closed; do not dispatch |
| Queue cancellation | Preserve durable turns; report interrupted audio/text separately |
| Adapter reuse without evidence | Disabled; hard flush |
| Shared transcript ambiguity | Use only the explicit scope; unresolved scope is not expanded silently |

**FC3-24 (Proposed).** A privacy failure is not repaired by deleting evidence. The failed composition is recorded with a redacted reason.

---

## 8. Acceptance design — proposed, not run

1. Two characters with conflicting private facts produce prompts containing only their own private context.
2. A speaker switch causes a hard flush; the next prompt is reconstructed from durable state.
3. Character A's voice assignment remains A after switching to B and back.
4. A voice-profile change for one character does not change another character's voice.
5. Interrupting A's TTS cancels audio without deleting A's committed text.
6. Restart and conversation restoration preserve character IDs and voice assignments.
7. A missing voice profile produces unavailable, not a different voice.
8. No cloud TTS request occurs.
9. A display-name collision does not select the wrong character or voice.

---

## 9. Open questions

| ID | Question | Status |
|---|---|---|
| `FC3-SHARED-TRANSCRIPT` | What is the public/shared transcript? | **Unresolved** |
| `FC3-PRIVATE-VISIBILITY` | What exactly is private to a character? | **Unresolved** |
| `FC3-KEYED-REUSE` | What adapter evidence permits keyed reuse? | **Unresolved; not authorized** |
| `FC3-FAIRNESS` | What queue/fairness policy applies? | **Unresolved** |
| `FC3-VOICE-VERSION` | How are voice-profile revisions bound and migrated? | **Unresolved** |
| `FC3-TTS-CANCELLATION` | What is the exact audio cancellation/retention contract? | Proposed, not approved |

---

## 10. Traceability and non-authorization

| Link | Artifact |
|---|---|
| Owner decisions | D5, D7 in `HF-CONTRACT-INTAKE-2026-11-20-01` |
| State identity | FC-1 (this set) |
| Prompt composition | FC-2 (this set) |
| Voice storage | FC-4 (this set) |
| TTS and two-character end-to-end acceptance | FC-6 `D7-E2E-01` (this set) |
| Multi-character readiness | `HF1-PLAN-2026-11-21-01` |

**FC3-25 (Observed — no authorization).** This paper does not implement multi-character routing, TTS, voice cloning, voice selection, queue changes or context reuse. It does not select a shared-transcript definition or a fairness policy.

**End of FC-3.** Hard flush protects private context; it does not erase voice identity. Isolation is proven by per-character composition, not by UI labels.
