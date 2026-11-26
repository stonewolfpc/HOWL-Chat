# FC-2 — Composition and Context Contract

## Contract paper `FC-2-2026-11-21-01`

| Field | Value |
|---|---|
| Date | **2026-11-21** (local, CST) |
| Type | Contract paper — review draft |
| Authority | D3 and D4 owner decisions; `HF-CONTRACT-INTAKE-2026-11-20-01` |
| Gate | M2 |
| Status | **PROPOSED — technology selection unresolved** |
| Implementation authorization | **None** |
| Tests/models/process runs | **0** |

**Labels:** **Owner-approved** · **Proposed** · **Calculated** · **Unresolved** · **Observed** · **Unestablished**.

**Identifier note:** Clauses suffixed with a letter (`FC2-05a`, `FC2-09a`, `FC2-05b`) are related clauses added without renumbering existing ones, so cross-references from other papers and the register remain valid. Open-question identifiers in §9 use the `FC2-D4-*` form; they identify questions, not clauses, and are cited as such.

---

## 1. Contract decision

**FC2-01 (Owner-approved — D4).** Howl must research, identify and adopt the best model-compatible tokenizer and prompt-construction approach for its requirements, regardless of implementation difficulty. Model-specific adapters may sit behind one consistent Howl interface.

**FC2-02 (Calculated — D4 is not a tokenizer-library choice).** The contract covers the complete pipeline: source collection, canonical resolution, privacy projection, model-specific formatting, tokenization, budget validation, reproducible accounting and fail-closed dispatch. Selecting a tokenizer library would satisfy only one stage.

**FC2-03 (Proposed).** The output of composition is a **PromptPlan**: a preflight artifact describing exactly what would be sent, what was included or excluded, why, and whether the request is valid.

**FC2-04 (Owner-approved — D3).** Composition resolves immutable canonical references and conversation state; it does not silently read mutable browser-local copies or “latest” objects.

---

## 2. PromptPlan artifact

**FC2-05 (Proposed — minimum fields).**

| Field | Obligation |
|---|---|
| `prompt_plan_id` | Immutable identity of this finalized plan artifact |
| `conversation_id` / `state_id` | FC-1 state being composed |
| `model_identity` | Model artifact and revision |
| `tokenizer_identity` | Exact tokenizer and version |
| `chat_template_identity` | Exact template and version |
| `adapter_identity` | Model-family formatting adapter |
| `composition_policy_identity` | Composition policy and version governing source selection, ordering, inclusion and trimming eligibility |
| `privacy_policy_identity` | Privacy and visibility policy governing this plan, with its version where the policy is independently versioned |
| `source_manifest` | Ordered sources and their authority |
| `segments` | System, world, scenario, character, lore, memory, history, tools, multimodal parts |
| `token_counts` | Per-segment and total counts from the actual tokenizer |
| `output_reserve` | Tokens reserved for generation |
| `context_limit` | Target model context limit |
| `included` | Segment IDs included |
| `excluded` | Segment IDs excluded, with reason |
| `transformed` | Summaries, projections or formatting changes, with provenance |
| `validation_status` | `valid`, `invalid` or `unknown` |
| `failure_reason` | Explicit reason when not valid |
| `content_hash` | Hash of the canonical reproducible plan content, excluding `prompt_plan_id` and `content_hash` themselves |

**FC2-05a (Proposed — plan immutability).** A finalized PromptPlan is immutable. A change to any input, model, tokenizer, template, adapter, composition policy or privacy policy produces a new `prompt_plan_id`; a plan is never edited in place to represent different composition rules. Two plans that differ in governing policy are not equivalent reproducibility inputs even when their state and canonical revisions match. A plan that cannot name its composition-policy identity, or the privacy policy governing it — with a version identifier where that policy is independently versioned — is not a valid plan (FC2-15).

**FC2-05b (Proposed — plan identity versus dispatch-attempt identity).** `prompt_plan_id` identifies the immutable plan artifact and is bound to its content hash and its composition-policy and privacy-policy identities. `request_id` is not a field of that artifact. A dispatch attempt is recorded separately, in an attempt record that references `prompt_plan_id` and carries its own `request_id`, so an unchanged plan retains its `prompt_plan_id` and hash across attempts while each attempt is identified on its own. Because the attempt lives outside the plan, recording an attempt cannot alter the plan's contents, its identity or its hash. This clause defines identity only. It does not establish retry eligibility or recovery behaviour, and no other clause of this contract does either: whether a further attempt may be made after a dispatch outcome, and how that outcome is recorded and recovered, remain unresolved. FC2-07 and FC2-17–FC2-20 describe composition, dispatch, explicit failure and state preservation; they do not decide retry eligibility or a recovery policy. A changed plan is never recorded under an earlier plan's identity, and a plan whose model, tokenizer, template or adapter identity has changed is recomputed rather than re-dispatched as if it were the same plan (FC2-16, FC2-23).

**FC2-06 (Proposed).** The plan is an auditable execution artifact. It is not a workaround for last-minute trimming and not a promise that the model will behave well.

---

## 3. Composition pipeline

**FC2-07 (Proposed — ordered stages).**

| Stage | Obligation | Failure behavior |
|---|---|---|
| **Collect** | Gather candidate sources from FC-1 state, FC-4 canonical refs, memory/retrieval and modality inputs | Mark missing source; do not substitute |
| **Resolve** | Resolve exact canonical revisions and stable IDs | `unresolved`; fail closed if required |
| **Project** | Apply privacy and visibility rules | Exclude only by explicit policy; record exclusion |
| **Order** | Apply the model-family prompt order | Invalid order is an error |
| **Format** | Apply chat template and special-token rules | Adapter failure is explicit |
| **Tokenize** | Count with the actual target tokenizer | `unknown`; no estimate as authority |
| **Budget** | Compare total plus output reserve to context limit | Invalid; fail closed |
| **Account** | Record included/excluded/transformed segments and reasons | Incomplete plan is invalid |
| **Dispatch** | Send only a valid plan | No valid plan, no dispatch |

**FC2-08 (Proposed).** Source collection must cover system instructions, conversation history, world, scenario, character, lore, retrieved memory, tool definitions, structured payloads and multimodal inputs, including their token costs. Where a modality's token cost cannot be determined exactly by the target tokenizer, the plan records the accounting basis and marks the count as estimated or unknown; an estimate is never presented as an exact count (`FC2-D4-MULTIMODAL`).

**FC2-09 (Proposed).** Every segment names its source authority. A display name, browser-local record or snapshot copy is never silently promoted to canonical source.

**FC2-09a (Proposed — memory boundary scope).** FC-2 does not define, authorize or integrate the MEMORY-09 record system, and it does not resolve the delegated D8 decision that gates FC-7. Retrieved memory from sources authorized elsewhere may be composed as an input segment under FC2-08 and FC2-09; this paper neither prohibits such inputs nor grants them authority. For MEMORY-09 specifically, no integration rule, record applicability or authority is established here; that belongs to FC-7 once it exists.

---

## 4. Model-family adapters

**FC2-10 (Proposed).** The consistent Howl interface accepts a model identity and a source plan; model-specific adapters own:

- chat-template application;
- special-token handling;
- role/system-message placement;
- tool-call formatting;
- multimodal part formatting;
- tokenizer invocation;
- model-specific stop/control semantics needed for composition.

**FC2-11 (Proposed).** An adapter declares its supported model families, template versions and known limitations. An unsupported family is refused explicitly, not silently formatted with a generic template.

**FC2-12 (Unresolved — D4 research).** The adapter boundary, tokenizer invocation path and template source of truth require a compatibility matrix and benchmark evidence across the actual target model families. No technology is selected here.

---

## 5. Token accounting and budget

**FC2-13 (Owner-approved — D4, with a calculated extension).** Token accounting uses the actual target model's tokenizer and chat template. Character-count estimates, heuristic ratios and stale persisted counts are not authority.

**FC2-14 (Proposed — budget equation).**

`system + world + scenario + character + lore + memory + history + tools + modality + output_reserve <= context_limit`

The left side is the sum of every included segment named in FC2-05 and FC2-08, each counted by the actual target tokenizer, plus the output reserve; an excluded segment contributes zero and is recorded as excluded. The exact segment ordering and reserve policy are model/adapter-specific. The inequality and the need for a reserve are not optional. The aggregate governs: a per-segment ceiling never substitutes for it, and the output reserve is counted once per request rather than once per segment.

**FC2-15 (Proposed — unknown state).** If the tokenizer or template is unavailable, the plan status is `unknown`. Howl does not present an estimate as compliance and does not dispatch a budget-dependent request as if it were valid. The same applies to an included segment whose token cost cannot be determined: while any included segment's cost is estimated or unknown, the aggregate budget cannot be certified `valid`, and an estimate retained for observability cannot establish that the inequality in FC2-14 is satisfied (FC2-08).

**FC2-16 (Proposed — model change).** Persisted token counts are invalidated when the model, tokenizer or template changes. Counts are recomputed; old counts are never reused as authority.

---

## 6. Overflow and failure

**FC2-17 (Owner-approved — D4, explicit failure, with a calculated extension).** A request that cannot be constructed safely fails explicitly. Howl does not silently send incomplete, invalid or over-budget context.

**FC2-18 (Proposed — no silent trimming).** Automatic trimming is not authorized by this contract. Any trimming requires a separate, explicit policy that names eligible segments, records the receipt and preserves the authoritative history. The old D4 A/B/C/D choice is not reinstated here.

**FC2-19 (Proposed).** An invalid plan returns a structured failure: which source, which segment, which limit, which adapter and what the user can do next. A generic “context too long” message is insufficient.

**FC2-20 (Proposed).** A failed composition does not delete, rewrite or reorder durable conversation state. User work remains.

---

## 7. Reproducibility and D1/D3 binding

**FC2-21 (Proposed).** A PromptPlan is reproducible when the same state, canonical revisions, adapter/tokenizer/template versions, composition-policy identity and privacy-policy identity produce the same plan hash.

**FC2-22 (Proposed).** The plan records exact canonical revision IDs and exact conversation state ID. A snapshot may reference the plan, but the plan does not become canonical conversation state.

**FC2-23 (Proposed).** If a referenced revision is missing, the plan is `unresolved`; it does not silently use the latest revision. If the template or tokenizer version changed, the plan is recomputed, not replayed as if identical.

---

## 8. Privacy and observability

**FC2-24 (Proposed).** Privacy projection happens before tokenization where possible, and the plan records what was projected or excluded. A private character segment must not appear in another character's plan; FC-3 owns the isolation rules.

**FC2-25 (Proposed).** The plan artifact may contain sensitive content. Storage, logging and export must distinguish the full plan from a redacted receipt. A content hash alone does not make the underlying text safe to expose.

**FC2-26 (Proposed).** Observability records plan identity, status, token totals, model/tokenizer/template identities, validation reason and dispatch decision. It does not claim model quality or response correctness.

---

## 9. Open research and decisions

| ID | Question | Status |
|---|---|---|
| `FC2-D4-ADAPTER` | Which model-family adapters are required first? | Research |
| `FC2-D4-TOKENIZER` | Which tokenizer implementation/path is authoritative per family? | Research |
| `FC2-D4-TEMPLATE` | Where do templates come from and how are they versioned? | Research |
| `FC2-D4-MULTIMODAL` | How are image/audio/video token costs represented? | Research |
| `FC2-D4-TOOLS` | How are tool definitions and structured payloads counted? | Research |
| `FC2-D4-TRIM` | Is any automatic trimming policy permitted? | Unresolved; not authorized |
| `FC2-D4-BENCH` | What compatibility benchmark proves a recommendation? | Proposed |

---

## 10. Acceptance design — proposed, not run

1. A plan for a known model/template produces a token total that matches the target tokenizer.
2. A missing tokenizer produces `unknown`, not an estimate.
3. An over-budget plan is invalid before dispatch; no model call occurs.
4. A model/template change invalidates old counts and recomputes.
5. A missing canonical revision produces `unresolved`; no latest fallback.
6. Included/excluded/transformed segments are auditable by ID and reason.
7. An unsupported model family is refused by the adapter, not silently formatted.
8. A failed composition leaves conversation state unchanged.
9. A prompt plan hash is stable for identical inputs, canonical revisions, model, adapter, tokenizer, template, composition-policy identity and privacy-policy identity.
10. A redacted receipt does not expose private segment content.

---

## 11. Traceability and non-authorization

| Link | Artifact |
|---|---|
| Owner decisions | D3, D4 in `HF-CONTRACT-INTAKE-2026-11-20-01` |
| State and snapshots | FC-1 (this set) |
| Canonical authority | FC-4 (this set) |
| Character isolation | FC-3 (this set) |
| Acceptance | FC-6 (this set) |
| Existing tokenizer seam | `HF1-PLAN-2026-11-21-01` §5.4 |

**FC2-27 (Observed — no authorization).** This paper does not select a tokenizer, prompt library, adapter implementation or template source. It does not authorize implementation or change any protected record.

**End of FC-2.** The prompt plan is a complete, model-specific, preflight account of the request — not a tokenizer choice and not a last-minute repair surface.
