# MEMORY-09 — Experiment Plan: Does Memory Make the Same Model More Capable?

**Status:** plan only. Nothing here is implemented, no probe has been run, nothing is frozen
on the new side, and no result exists to report. Approval of this plan is a prerequisite for
any code.

**Revision v3 (review amendments, before approval):** §2.5 adds placebo **construction
fidelity** — per-turn, structure-matched, tokenizer-measured — so arm D cannot be
token-equal while differing in shape or salience.

**Revision v2 (review amendments, before approval):** (a) the outcome scale is a four-level
ladder with *memory-supported answer* as the headline, not "M-correct"; (b) a **token-matched
placebo arm (D)** was added, because equal configured budgets do not prove equal useful
information; (c) **effective** context accounting is specified as measured, per turn, with
memory competing for the same budget rather than adding to it; (d) temporal correctness is
included in the veto set; (e) A2 results are barred from promoting proposals into decisions.

**Relationship to prior work:** this plan resumes **after** the bge reranker evaluation, which
is closed. The reranker was evaluated, found insufficient, and is **not shipped**; its
apparatus stays frozen and is not re-run. No research is restarted and no earlier evaluation
is repeated.

**Objective under test (unchanged):** make the same model more capable over time **without** a
larger model, heavier prompting, or a more complicated model interaction. Measured against the
existing baseline, not merely demonstrated.

**Record this plan does not replace:** the design contract (`MEMORY-09-ARCH-STAGE-CONTRACT`,
r32) and its normative extract. The contract says what must be true; this plan says how we
would find out whether it is.

---

## 1. Part 1 — Preserve the baseline

The baseline is evidence, not a starting point to be adjusted. It is frozen before the
experiment begins and is not modified to accommodate the subject.

### 1.1 Frozen artifacts (verified at plan time)

| Artifact | sha256 (prefix) | Role |
|---|---|---|
| `eval/s3/rerank-spec.json` | `bbf400c225e6c90a` | frozen v1.1 protocol (immutable) |
| `eval/s3/rerank-spec-v1.2.json` | `5414d0a7…` | v1.2 protocol, separate from v1.1 |
| `eval/s3/rerank_eval.py` | `f59959f30ec336ae` | frozen R1–R5 harness |
| `eval/s3/rerank_set.py` | `23badd5e7825718d` | frozen adversarial case set |
| `eval/s3/formation.py` | `d8aac500…` | Gate A formation, frozen |
| `eval/s3/records.py` | `1432fb63…` | record model, frozen |
| `eval/s3/evidence.py` | `19db697a…` | evidence layer, frozen |
| `eval/s3/temporal-spec.json` | `43d81e07e236df8c` | E-TEMPORAL-SET v2 (T10 not edited) |
| receipts | `43d744e9…` (temporal v2), `06f8413885aeb992` (A/B), `caed2a2f…` (temporal v1) | immutable results |

`eval/harness/*` (adapter, conditions, provenance, state, plan, preflight, validate) is
**read-only** for this work. So is every S1/S2 artifact and schedule.

### 1.2 Known results, stated with their limits

| Measurement | Result | Limit that must not be lost |
|---|---|---|
| Gate A | 14/14 PASS | formation frozen at 7 records — not a capability claim |
| Gate B | **6/7 PARTIAL** | **P5 ("Who trained the horse?") is an open failure** |
| Harness self-checks | PASS | zero model calls |
| S1 audit | PASS, 5 calls | behavioural, single batch |
| `validate` | 89/89 | local enforcement, not model behaviour |
| Corpus | PASS | structural |
| Reranker | evaluated, **not shipped** | no further candidate fetched |
| Resource envelope | no config meets VRAM ≤384 MiB **and** p50 ≤40 ms | four modes exist; none satisfies both |

**Gate B 6/7 is not a completed baseline.** The subject may not inherit P5 as an assumption.

### 1.3 Rules for touching the baseline

1. Existing passing tests are **not rewritten** to accommodate the subject.
2. Frozen files are not edited. New work lives in new modules (§3.2).
3. If a baseline test cannot express something the experiment needs, that is an
   **extension**, recorded as an addition with its own coverage — never a modification of
   what already passes.
4. Any baseline behaviour that looks wrong is recorded as a finding, not quietly fixed.

---

## 2. Part 2 — Experimental conditions

### 2.1 Arms

| Arm | Condition | Purpose |
|---|---|---|
| **A0** | baseline: no memory, frozen window, no assembly | what the model does unaided |
| **A1** | memory from **adopted** decisions only (D1–D7 mechanics) | tests the contract we have |
| **A2** | A1 + **proposed** rules (§11.3 temporal, §10.1 materiality) | exploratory only; inert until approved |
| **C** | control: prewritten `Known:` records injected | proves plumbing; **never scored as memory formation** |
| **D** | placebo: same token cost as A1, but decoy/irrelevant records | separates "better selection" from "more tokens" |

**A2 is inert until the proposals are approved.** If run early it may only be reported as
"exploratory, unapproved rules", never as validation of the contract, and never as evidence
for approving those rules. A2 results cannot promote a proposal into a decision by any route:
promotion requires the separate approval path, with the consequences examined as the contract
requires.

**Why arm D exists:** A1 could beat A0 simply by receiving more usable text. D holds the
budget and the block shape constant while removing relevance, so an A1 gain over D is
attributable to *selection*, not to extra tokens. Reporting A1 vs A0 alone would not
distinguish the two.

### 2.2 Identity matrix — what must be identical

| Dimension | Held identical | How it is proven |
|---|---|---|
| model + revision | same endpoint, pinned revision | harness `provenance` record; `model_revision` non-`unavailable` |
| sampling | temperature, top_p, max_tokens, seed policy | `conditions.FROZEN_PARAMS`; `seed_honoured` stated |
| conversation evidence | same turns, same tokens-per-turn | `evidence.EvidenceLog.sha256()` |
| queries | same probe set, same order | case hash in the receipt |
| context budget | same window capacity and reserve | `context.Window.as_dict()` |
| assembly budget | same max records / max chars | retrieval result record |
| system message | same except the memory block | `conditions.assert_separated`-style single-block contrast |
| **effective prompt tokens** | **measured per turn, not merely configured** | token accounting (§2.4) |

**Any unavoidable difference is documented in the run receipt.** An undocumented difference
means the run is not a comparison.

### 2.3 The anti-cheat rule

The model must form memories **from actual conversation evidence**. Concretely:

- Every record in A1/A2 must cite a `source_span` that exists in the run's `EvidenceLog`.
- A run whose records were injected rather than formed is arm **C** and is labelled as such.
- Feeding `Known: …` blocks and reporting retrieval success is **excluded from the primary
  endpoint** — that is the same failure the current baseline already exhibits, dressed up.

Verification is mechanical: for each record, resolve `source_span` and assert it belongs to
the conversation turns of that run.

### 2.4 Effective context budget — measured, not configured

Identical configured limits do not imply identical information. A run receipt records, per
turn, the **measured** token split:

| Component | Recorded as |
|---|---|
| system + instructions | `tokens.system` |
| conversation history in the window | `tokens.history` |
| assembled memory block | `tokens.memory` |
| assembly wrapper and citation overhead | `tokens.overhead` |
| **total prompt** | `tokens.total` |
| reply reserve | `tokens.reserve` |

Rules:

- **Memory competes with history, it does not extend it.** Both arms share one prompt
  capacity; the memory block is paid for out of the same budget, displacing history when
  needed. That displacement is the architecture's actual claim, not a side effect.
- `tokens.total ≤ capacity − reserve` in **both** arms, asserted per turn, not inferred.
- The A0/A1 delta in `tokens.total` is reported with run-to-run variation. A systematic A1
  advantage in total tokens is a **finding about the experiment**, not a capability result,
  and blocks the capability claim until explained.
- Arm **D** matches A1's `tokens.memory` (same count, irrelevant content), so budget effects
  are held constant while relevance is removed.
- Token counts come from the runtime's own accounting, never estimated here — the same rule
  the evidence layer already applies to turn tokens.

### 2.5 Placebo construction — execution fidelity

Arm D is only a control if it is matched to A1 **per turn**, not on aggregate totals. Total
parity can hide alternating over- and under-shoot, and a technically token-matched block can
still be a differently structured or more distracting input.

| Requirement | Detail |
|---|---|
| **Per-turn matching** | D is matched to A1 turn by turn; aggregate-only equality is not accepted |
| **Record count** | identical, or the difference is documented with its reason before the run |
| **Record length** | matched within a declared tolerance per record |
| **Serialised shape** | same bullet form, ordering, citation/provenance fields and header; structure must not differ |
| **Placement** | same position in the system message, same wrapper, same section ordering |
| **Content** | demonstrably irrelevant to the probe — it must neither answer nor contradict it |
| **Decoy pool** | fixed and declared before the run, not improvised per case |
| **Tokenizer** | counts come from the **model's own tokenizer** applied to the actual serialised prompt, per turn |

**Acceptance check:** `|tokens.memory(A1) − tokens.memory(D)| ≤ tolerance` asserted every turn,
with the tolerance declared in advance. A turn outside tolerance is a receipt-visible
deviation and is reported as such; it is never averaged away.

**Why structure is included:** matching only length would let D differ in formatting or salience
and still pass a token check. If D fails on shape, an A1-vs-D difference may measure
presentation rather than selection — which is the one thing the arm exists to rule out.

### 2.6 Thresholds are not tuned on the evaluation set

- Retrieval `min_score`, abstention thresholds and any gate constants are fixed **before**
  the run, from held-out calibration or a prior protocol (the v1.2 calibration exists).
- Choosing a value after seeing results invalidates the result, even if it "looks obvious".
- Calibration pairs and evaluation cases remain disjoint, as in v1.2.

---

## 3. Part 3 — Minimal test subject

### 3.1 Scope: the smallest thing that can be wrong in interesting ways

**In scope (A1):** three streams and a join, reduced to what is executable.

- `MemoryRecord` — chat-owned, with `owner_id`, status, supersession, source span.
- `AcquisitionEvent` — append-only, authoritative, `via` ∈ {witnessed, verified, read, told,
  inferred, narrator_grant, revoked}.
- `KnowledgeEdge` — projection of the events; carries no truth status.
- Retrieval contexts — at least `character_perspective` (default, fail-closed) and
  `world_state`; eligibility precedes ranking.
- Formation from spans — deterministic extractor over real conversation text.
- Correction and supersession — a `revoke`/re-`learn` pair, with history retained.
- Abstention — empty eligible set abstains, never widens the pool.

**Out of scope for A1 (deliberately):**

- temporal applicability intervals, conflict groups, materiality predicate (§10.1, §11.3) —
  A2 only, and inert until approved;
- transfer, deletion, replayability states, detector governance — data-lifecycle rules with
  no behavioural signal in a short experiment;
- embeddings, reranking (closed), group knowledge, ambient authorship;
- any Howl Chat integration, UI, or runtime change.

### 3.2 Where the code lives

A **new package** `eval/mem1/`, importing the frozen layers read-only:

| Need | Reuse (unmodified) | New module |
|---|---|---|
| conversation evidence | `eval.s3.evidence.EvidenceLog` | — |
| context window + rollover record | `eval.s3.context.Window` | — |
| A/B separation discipline | `eval.harness.conditions` | — |
| write-once run receipts | `eval.harness.provenance` | — |
| A′ model (record/event/edge) | — | `eval/mem1/model.py` |
| formation from spans | — | `eval/mem1/form.py` |
| context-explicit retrieval | — | `eval/mem1/retrieve.py` |
| assembly + citation | — | `eval/mem1/assemble.py` |
| long-conversation driver | — | `eval/mem1/conversation.py` |
| probe + receipt capture | — | `eval/mem1/probe.py` |
| scoring | — | `eval/mem1/score.py` |

`eval/s3/` receives **no edits**, so its hashes stay as recorded in §1.1.

### 3.3 Definition of done for the subject

- Zero-model self-check of the whole subject passes (mirrors `eval/s3/selfcheck.py`).
- A seeded conversation produces the expected records, edges and rollovers, and the check
  asserts *which* span each record came from.
- A deliberately corrupted run (record citing a nonexistent span) is **rejected**.
- No model call is required for any of the above.

---

## 4. Part 4 — Long-conversation evaluation

### 4.1 The eight questions, and what answers them

Every question is answered by instrumentation, not by reading transcripts.

| # | Question | Instrument | Failing signal |
|---|---|---|---|
| 1 | were the right memories formed? | record set vs planted facts, per span | planted fact with no record; record citing wrong span |
| 2 | did they survive rollovers? | `Window.rollovers[].dropped_spans` + record survival | record present but its evidence never retrievable |
| 3 | could the right character retrieve them? | retrieval receipt: selected ids + `owner_id` | correct answer with empty selection |
| 4 | was another character prevented? | same receipt, foreign-owner probe | any foreign record in selection — **veto** |
| 5 | did correction/supersession work? | status transitions + history query | new value active but history erased; old value still active |
| 6 | is history distinguishable from current truth? | answers to "was X true earlier?" vs "is X true now?" | one answer serving both |
| 7 | did it abstain when evidence was insufficient? | `abstained` + reason on unanswerable probes | confident answer from no eligible evidence |
| 8 | did retrieval supply the evidence behind a correct answer? | answer must cite selected record ids | correct answer, empty selection → **not** a memory success |

### 4.2 Outcome ladder — four levels, never collapsed

| Level | Definition | What it does **not** show |
|---|---|---|
| **answer-correct** | the final answer is correct | that memory was involved |
| **retrieval-correct** | the appropriate supporting record was selected | that it reached the prompt |
| **evidence-supported** | the selected memory was included in the assembled context and supports the answer | that the model used it |
| **memory-supported answer** | answer-correct **and** evidence-supported **and** the original evidence had left the window | causal reliance on the memory |

**The headline is memory-supported answer**, computed only on probes whose supporting span is
confirmed absent from the window (`Window.rollovers[].dropped_spans`). Correct answers with an
empty or irrelevant selection are `answer-correct` without support and are investigated, not
counted.

**Causal caution:** selection plus presence plus the absence of the original evidence is strong
evidence, but it is not proof that the model *used* the record rather than inferring the answer
from dialogue structure, character knowledge or luck. A causal claim would need the
intervention arm D plus a per-record ablation, and this plan does not make that claim.

### 4.3 Vetoes outrank aggregates

- **Cross-character leakage** — a character answers from another character's knowledge —
  fails the run regardless of any aggregate, per the standing rule that a high aggregate score
  cannot excuse a veto failure.
- **Erased history** — a superseded value becomes unrecoverable — fails the run.
- **Fabricated provenance** — a record whose `source_span` does not resolve — fails the run.
- **Temporal correctness** — a superseded value presented as current, or a current value
  presented as historical, is not an ordinary retrieval miss. Where a case exercises
  supersession, this fails the run.

### 4.4 Conversation shape

- Long enough to force **≥3 rollovers** after the evidence leaves the window.
- Planted facts of each class: character-specific, ambient, restricted, experiential.
- At least one correction after rollover, and one retraction of an asserted fact.
- At least two characters with a deliberate knowledge asymmetry.
- An unanswerable probe whose correct response is abstention.

**P5 is run as its own case** in every arm. It is not allowed to disappear into an aggregate.

---

## 5. The three questions this plan settles

### 5.1 What stays untouched?

**Untouched:** `eval/harness/*`, all frozen `eval/s3/` specs/harnesses/sets, formation and
receipt files, S1/S2 artifacts and schedules, the model, sampling parameters, the corpus, and
every test that currently passes.

**Extended (new modules, additive):** A′ model, formation from spans into events, context
explicit retrieval, receipt capture, scoring. No existing file is edited; nothing passing is
rewritten to fit.

### 5.2 What is a fair comparison?

Same model and revision, same conversation, same queries, same context budget, same sampling,
same assembly budget; the only intended difference is the memory block. Differences that
cannot be removed are documented per run. Records must be formed from the conversation, and
arm C is labelled so injected-record plumbing is never counted as formation.

### 5.3 What does a long conversation reveal that a short one cannot?

Whether knowledge **survives** the window. A short exchange can be answered from the prompt;
after rollovers that is impossible, so a correct answer requires retrieval — and the receipt
says whether retrieval actually happened. Short tests cannot distinguish memory from
attention, which is precisely the distinction the objective depends on.

---

## 6. Analysis plan

- Report **per-question** results (8 above) plus per-class counts; no single "memory score".
- Report the arms side by side with the same cases; A0 → A1 is the capability claim.
- Report surprising results, including failures with no explanation, rather than smoothing.
- Any difference smaller than the run-to-run variation (measured by reps at fixed settings) is
  reported as **indistinguishable**, not as an improvement.
- Threshold choices are listed in the receipt with their provenance (calibration or prior).

## 7. Threats to validity

| Threat | Mitigation |
|---|---|
| prompt still holds the evidence | assert `dropped_spans` contains the supporting span before the probe turn |
| leakage between characters | receipt check + veto |
| judge drift (LLM scoring) | deterministic matchers where possible; any model judging is its own declared experiment |
| model drift across runs | pinned revision, recorded in every receipt |
| overfitting to cases | cases frozen before runs; no threshold tuning (§2.6) |
| P5 inheritance | P5 run explicitly in every arm |
| "it works" from plumbing | arm C separated and excluded from the primary endpoint |
| small sample | reps at fixed settings; report variation, not point estimates alone |

## 8. Sequencing and gates

| Gate | Content | Exit criterion |
|---|---|---|
| **G0** | plan approval | this document accepted as the experiment's scope |
| **G1** | subject build | zero-model self-check green; corrupted-provenance run rejected |
| **G2** | plumbing dry run (arm C) | records flow end to end; **no capability claim** |
| **G3** | one long conversation, all arms | instrumentation complete: every probe yields a receipt and a verdict class |
| **G4** | full run | N conversations × arms A0/A1 (A2 inert), receipts written once |
| **G5** | disposition | capability claim accepted, rejected, or "indistinguishable" — with vetoes reported first |

No stage starts before its predecessor's exit criterion is met. A failure at G1–G3 is a
finding about the subject, not a reason to relax the check.

## 9. What would count as success — and what would falsify it

**Success:** at A1, on probes whose evidence left the window, **memory-supported answers**
rise materially above **both A0 and the token-matched placebo D**, abstention stays calibrated,
no veto fires, and `tokens.total` shows no systematic A1 advantage. Reported as a difference
with variation, not a score.

**Falsification:** memory-supported answers indistinguishable from A0 **or from D** after
rollovers; or gains explained by a larger effective budget or heavier prompting; or any veto
failure; or answers that are correct while the window still holds the evidence. Any of these
is a result worth having — A1 losing to D would mean the block merely bought tokens.

**Either outcome is publishable as evidence.** The point is not to make the architecture look
right; it is to find out what it actually is.

## 10. Non-goals

- No Howl Chat production integration, no UI, no runtime change.
- No new architecture, rule, field or probe in the contract while this plan executes.
- No re-running of the bge evaluation, and no new reranker candidate.
- No tuning of thresholds, cases or prompts against the evaluation set.
- No commit or push as part of the experiment.
