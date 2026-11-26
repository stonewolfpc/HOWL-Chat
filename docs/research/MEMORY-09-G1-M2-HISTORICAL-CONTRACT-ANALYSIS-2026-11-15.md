# MEMORY-09 — G1 Review: Is Prompt-Facing Historical Retrieval Already in the Contract?

**Status:** analysis only. No architecture change made, none proposed as already-decided.
**Question:** the G1 review found that the correction probe's keyed answer (a historical value)
cannot be M2-scored because A1's eligibility filter is current-only. Before adding any retrieval
context, establish whether historical temporal intent can be represented inside an **existing**
context. If not, the addition is a contract change requiring owner approval.

**Answer: it cannot. This is a contract gap, not an implementation choice.**

---

## 1. What the contract names

The frozen contract defines retrieval as context-explicit, with exactly **four** contexts
(`K-R5`, §7 table):

| Context | Includes | Must not |
|---|---|---|
| `character_perspective` | memories with an active edge for that character + ambient-class records in scope | receive C1 authoritative truth as knowledge; receive other characters' private scope |
| `world_state` | every chat memory and every asset, regardless of character knowledge | be rendered as a character's knowledge |
| `narrator` | world state + all edges + unresolved conflicts + restricted lore | state a secret as something the character knows; inject it into the character channel |
| `inspection` | world state + knowledge map + conflict map | **feed the prompt** |

Two further facts decide the question:

- **`K-R20`**: `revoked_turn` is "derived on demand — the timestamp of the last revoke event — and
  appears **only in the `knowledge_history` inspection view**." History has exactly one named
  home, and it is the audit view.
- **`inspection` must not feed the prompt.** So the only context that exposes history is
  explicitly barred from being the source of prompt content.

## 2. What the plan expects, and where it diverges

Plan §4.1 question 6 asks whether "history [is] distinguishable from current truth" — answers to
"was X true earlier?" versus "is X true now?" — with a failing signal of "one answer serving both".
That requires a **prompt-facing** path to a superseded value.

The contract supplies no such path. Temporal applicability (I16: "which claims apply to a given
point or interval") exists as a layer, but §11.3 is **[P] proposed**, and plan §3.1 places
"temporal applicability intervals" explicitly **out of scope for A1** — A2 only, inert until
approved.

**Therefore the plan asks for a measurement the frozen architecture does not, at A1, provide.**
That divergence is the finding; it is not an argument for silently inventing the interface.

## 3. Why the obvious "just add `history`" move is not implementation detail

Retrieval context is a **call parameter** and a contract-level concept (`K-R5`, §7). Eligibility
exclusion is defined per named context. Adding a fifth prompt-facing context would:

1. expand the contract's enumerated interface from four contexts to five;
2. introduce a prompt path to superseded/revoked content, which the inspection context was
   deliberately barred from providing;
3. change what "the character channel" can be shown, which interacts with `K-R6`/`K-R7`
   (discovery requires `narrator`; moving information into the character channel requires an
   `AcquisitionEvent`).

Any of those is an owner decision, not a reviewer's convenience. **Not done.**

## 4. Options, with consequences

| Option | What it means | Consequence |
|---|---|---|
| **A — fifth context** (`history`, prompt-facing) | new named context returning superseded/revoked records with status and turn | Cleanest semantics; explicit contract amendment; interacts with K-R6/K-R7; needs its own eligibility rules and veto coverage |
| **B — temporal parameter** on existing contexts | `temporal: current \| historical` on a retrieval call; eligibility widens for `historical` | Smaller interface change, but it modifies per-context eligibility, which is also contract-level; must prove historical mode cannot leak into present-tense answers |
| **C — accept the limit** | correction probe stays M1-only; M2 coverage is 11 of 13 probes | No contract change; smallest evidence loss; the limitation must be stated in the G4 receipt and may not be described as a pass |

**Recommendation: C now, A later if the architecture adopts history at all.** Option C keeps G4
honest today, costs one probe's M2 coverage, and leaves the contract untouched. Option B is the
most dangerous of the three, because a single flag that widens eligibility is exactly the kind of
"same call, different semantics" that produces leakage; it would need the leakage test named in
§5 before it could be trusted.

## 5. Tests that must exist before any option is adopted

1. A present-tense probe can never select a superseded or revoked record (already covered in
   current mode; must remain true under any new path).
2. A historical probe that finds nothing must abstain rather than answer from the current value.
3. A historical answer that presents the superseded value as current remains a **temporal veto**,
   not a scoring nuance.
4. Under option A or B, the same 13-probe fixture must produce identical M1 in every arm — the new
   path may not change answer quality measurement.

## 6. What this analysis does not do

- It does not change the frozen definition of M2.
- It does not add the context, and it does not weaken the correction probe.
- It does not block F2, which was closed inside the existing lifecycle stream: retraction is an
  explicit event record referencing the retracted record, with its own span, so "retracted" is
  distinguishable from "never established" and from "not known".

## 7. Decision required

F1 remains **open and blocking for G4** unless you choose option C, in which case it converts to a
declared limitation with reduced M2 coverage (11/13) and no contract change.

---

## 8. Closeout, same day: option C accepted, F3/F4 bounded review, G4 preflight

**Decision taken: C.** p47 stays in the probe count and is reported as `NOT_SUPPORTED` with a
provable span absence — never excluded, never described as a successful M2 observation. M2
coverage is 11/13, plus the unanswerable probe that is `N/A` by definition.

**F3 — bounded.** Weak overlap can still *select* a wrong record (scores 0.05–0.06), but it cannot
*support* an answer: `support_check` is predicate + value, evaluated after selection and
independent of retrieval score. Proved three ways: M2 is never awarded on selection alone across
all 13 receipts; a selected record with the right predicate and the wrong value fails support; and
an adversarial record sharing every query term under a different predicate is selected yet still
fails support.

**F4 — reduced, one residual risk recorded.** The matcher now requires an explicit retraction
*cue* (`I was wrong`, `that was wrong`, `retract…`, `scratch that`, `never mind`,
`forget that/what I said/everything/it`) plus ≥2 shared content tokens. Consequences, all tested
through formation rather than by calling `_retract` directly:

| Utterance | Outcome |
|---|---|
| "I forget where the ledger is." | no revocation (memory lapse) |
| "Forget the weather — the ledger is in the green cabinet." | no revocation (incidental mention, asserts the fact) |
| "I was wrong about the ledger in the cabinet." | revokes the ledger only; the promise survives |
| "Never mind the ledger — the gate stays locked." | no revocation (one shared token) |
| "Never mind the ledger in the green cabinet." | **revokes** (cue + three shared tokens) |

That last row is the residual risk: the rule matches *phrasing*, not intent, so a cued topic
dismissal with rich overlap can still withdraw a claim the speaker never withdrew. It is recorded
in the receipt rather than patched with an intent heuristic that would be unverifiable.

**F2 invariants — held.** The retraction event is active and cites its span and target; the target
stays revoked with provenance; the event does not appear in current truth, does not leak into an
unrelated question, and cannot support a present-tense claim about the fact it retracted (support
requires the probe to declare `retraction_of`). Retrieval never returns the revoked record.

**G4 preflight — BLOCKED, 27 PASS / 0 FAIL / 2 BLOCKED.** Frozen identities, seed table, quotas,
matched seeds, write-once receipts, capacity split, veto wiring and the no-network guard all pass.
The two blockers are environmental, not design: no runtime tokenizer is bound (G1 counts with a
test stand-in), and no model revision or sampling parameters are configured. **No model call was
made, and the instrument being sound is not evidence that the experiment ran.**

---

## 9. Standing disposition (owner, 2026-11-15) — work stops here

| Item | Disposition |
|---|---|
| **F1** | Accepted limitation for this experiment. p47 stays counted and `NOT_SUPPORTED`. |
| **F2** | Closed, with lifecycle invariants retained. |
| **F3** | Support-credit risk controlled. **Retrieval-quality risk remains**: irrelevant records still consume context budget, can distract the model, and can create misleading associations even when M2 refuses to credit them. |
| **F4** | **Bounded but unresolved.** The cue rule fixes the demonstrated false positives; it is a deterministic language rule, not intent recognition. A false retraction can corrupt downstream knowledge even when the event is recorded correctly. |
| **P5** | Independently observable, outside the 13-probe aggregate. |
| **G4** | Blocked, not authorised, zero experimental calls. |

**G07 and G12 are not to be waived or worked around.** G07: a test-only whitespace counter cannot
establish real prompt capacity or displacement; bind the actual runtime tokenizer and verify
serialised prompt accounting. G12: without a pinned model revision and sampling configuration the
experiment is not reproducible; freeze identity, revision and parameters before execution. Do not
substitute a convenient tokenizer to get past G07.

**What the 27 passing preflight checks do and do not establish.** They establish that the
instrument's deterministic controls behave as specified. They do **not** establish that runtime
accounting or model configuration will behave correctly once connected — those remain genuine
environmental dependencies, and the two BLOCKED items are exactly where that gap sits.

**No further action is authorised under this disposition:** no architecture expansion, no scoring
changes, no model calls, no attempt to turn BLOCKED items green.
