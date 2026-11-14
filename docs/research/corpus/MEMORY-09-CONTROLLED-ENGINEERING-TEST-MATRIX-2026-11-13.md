# `MEMORY-09` — Controlled Engineering Test Matrix

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `MEMORY-09-CETM-1` · **Version:** 1 · **Date:** 2026-11-13 · **Status:** **MATRIX DESIGN — no test executed; every row is `NOT RUN`.**
**Authority:** `MEMORY-09-WO-ARCH-06` §7–§8.

---

## 1. Repetition doctrine

**No universal repetition count.** A repeat must have a stated reason, and only these four qualify:

| Reason | What it buys | Typical N |
|---|---|---|
| **Reproduce** a baseline | Confirms a result isn't a one-off | 2–3 |
| **Measure variation** | Spread of behaviour under identical settings | 5–10 |
| **Test a change** | Compares two conditions | 3–5 each |
| **Probe a failure** | Chases an observed defect | ad hoc, recorded |

A repeat with no reason is not evidence; it is compute spent to feel safer.

## 2. Test matrix

| ID | Dimension | Deliberately varied |
|---|---|---|
| T-01 | Run-to-run consistency | nothing (repeats) |
| T-02 | Sampling sensitivity | sampler values |
| T-03 | Model compatibility | model, same architecture |
| T-04 | Memory configuration | memory on/off, retrieval mode |
| T-05 | Context pressure | conversation length |
| T-06 | Character switching | active character |
| T-07 | Inventory changes | stored inventory |
| T-08 | Clothing/equipment | stored equipment |
| T-09 | Historical facts & relationships | stored history |
| T-10 | Conflicting/missing information | state conflicts |
| T-11 | **Architecture vs model capability** | architecture assistance on/off |

---

### T-01 · Run-to-run consistency
**Hypothesis:** at fixed settings, output *content* varies while *state-dependent facts* do not.
**Controlled:** model, sampler, context, state. **Varied:** run index only.
**Reps:** 5–10, reason = measure variation.
**Expected:** fact-level stability near-total; prose-level divergence expected and benign.
**Failure:** state facts change between runs (a defect, not variance) — e.g. inventory spontaneously differs.
**Ceiling:** "consistent under these conditions", never "deterministic".
**Raw evidence:** all N full replies, seeds if recorded, exact settings.
**Note:** run-to-run spread is **not** an α question (α = agreement between *raters*); report spread, not α.

### T-02 · Sampling sensitivity
**Hypothesis:** small sampler changes shift style more than factual continuity.
**Controlled:** model, context, state, prompt. **Varied:** one sampler value at a time (BCB vs a documented alternative).
**Reps:** 3–5 per condition, reason = test a change.
**Expected:** style drift; facts intact.
**Failure:** a sampler change breaks continuity or persona — a real finding about operating ranges.
**Ceiling:** sensitivity *for this model*, not a general law.
**Raw evidence:** both exact sampler sets, paired replies.

### T-03 · Model compatibility
**Hypothesis:** the architecture's guarantees hold across families.
**Controlled:** architecture, state, test inputs, BCB intent. **Varied:** model family.
**Reps:** 2–3 per model, reason = reproduce across conditions.
**Expected:** same *behavioural contract* (canonical facts appear), different prose.
**Failure:** a family systematically ignores or contradicts supplied state — a compatibility limit worth naming.
**Ceiling:** "compatible on the tested families", never "model-agnostic".
**Raw evidence:** per-model records with sampler values and load/quant details.

### T-04 · Memory configuration
**Hypothesis:** architecture-supplied state, not model size, drives continuity.
**Controlled:** model, sampler. **Varied:** memory on/off (and retrieval mode if applicable).
**Reps:** 3–5 each, reason = test a change.
**Expected:** with memory, canonical facts persist; without, drift/contradiction.
**Failure:** no difference — would falsify the architecture's value for this test.
**Ceiling:** effect *for this model and this scenario*; no causal claim about internals.
**Raw evidence:** both configurations, identical inputs, state dumps.

### T-05 · Long-form / context pressure
**Hypothesis:** architecture-supplied state survives beyond what the context window can hold.
**Controlled:** model, sampler, starting state. **Varied:** conversation length up to and past usable context.
**Reps:** 2–3, reason = probe a failure mode.
**Expected:** early facts still honoured after truncation; failures appear as gradual drift, not cliff-edges.
**Failure:** early canonical facts contradict late replies with state still present — a genuine defect.
**Ceiling:** "resilience up to N turns on this model".
**Raw evidence:** full transcript, context budget, truncation points, state snapshots.

### T-06 · Character switching
**Hypothesis:** character knowledge does not bleed across characters.
**Controlled:** model, sampler, scenario. **Varied:** active character.
**Reps:** 3 per switch, reason = test a change.
**Expected:** each character's facts stay theirs; no cross-attribution.
**Failure:** character A's private fact appears in character B's mouth.
**Ceiling:** separation for tested pairs only.
**Raw evidence:** both transcripts, state view per character.

### T-07 · Inventory changes
**Hypothesis:** inventory is authoritative and updates are reflected, not re-invented.
**Controlled:** model, sampler, scene. **Varied:** add/remove/consume items between turns.
**Reps:** 3, reason = test a change.
**Expected:** the reply respects the stored list; no phantom items.
**Failure:** the model invents, duplicates, or "remembers" a removed item.
**Ceiling:** compliance with stored state; not that the model understands ownership.
**Raw evidence:** state before/after each turn, replies.

### T-08 · Clothing and equipment
**Hypothesis:** persistent state is enforced across scenes.
**Controlled / varied / reps / ceiling / evidence:** as T-07, with the equipment set as the varying element.
**Failure:** clothing changes without a scene event; an equipped item is used after being removed.

### T-09 · Historical facts and relationships
**Hypothesis:** long-term memory preserves established facts and relationship state.
**Controlled:** model, sampler. **Varied:** elapsed turns since the fact was established.
**Reps:** 3 at three distances, reason = measure decay.
**Expected:** facts persist at all distances; phrasing varies.
**Failure:** a fact silently changes (date, relation, event outcome).
**Ceiling:** continuity over tested distances, not "perfect memory".
**Raw evidence:** establishing turn, recalling turn, both replies.

### T-10 · Conflicting or missing information
**Hypothesis:** under conflict, the system refuses to manufacture certainty (the `R-6` construct, at engineering level).
**Controlled:** model, sampler. **Varied:** state contains two conflicting values, or a required value is absent.
**Reps:** 3 per condition.
**Expected:** conflict is surfaced or hedged; absence is admitted rather than filled.
**Failure:** a confident fabricated resolution; a missing value silently invented.
**Ceiling:** "did not manufacture certainty in these cases" — never "recognised the conflict".
**Raw evidence:** conflicting state dump, reply, and which value (if any) was chosen.
**Note:** this is the engineering analogue of `R6-01`; it does **not** substitute for cue-level scoring.

### T-11 · Architecture vs model capability — the central comparison

**Hypothesis:** behaviour that persists because the architecture supplies state differs from behaviour a model produces from its own priors.
**Controlled:** model, sampler, test inputs, scenario, response length.
**Varied:** architecture assistance **on** vs **off**, same model, same inputs.
**Reps:** 3–5 per condition, reason = test a change.
**Expected:** the assisted condition honours canonical facts the unassisted condition drifts from.
**Failure:** assisted ≈ unassisted — would show the architecture adds little *for this model and scenario*.
**Confound disclosure required:** if conditions cannot be made equivalent (e.g. the unassisted path changes prompt shape), say so in the record rather than implying equivalence.
**Ceiling:** a difference supports *"the behaviour differs when assistance is present"* — **not** retrieval, not internal use, not "the model remembered".

**Outcome classification (mandatory, from the order §8):**

| Class | Meaning |
|---|---|
| **C1** Correct response | Factually right |
| **C2** Correct with authoritative state | Right, and the state contained it |
| **C3** Correct plausibly from model priors | Right, but the fact was guessable/common |
| **C4** Incorrect despite correct state | The defect class that matters most |
| **C5** Missing or contradictory state | An architecture defect, not a model defect |
| **C6** Cause unestablished | Recorded as such; no inference permitted |

**C2 vs C3 is the crux of the whole evaluation**, and it is decided by whether the fact could have been known without the state — a judgement the record must state explicitly, with the basis.

---

## 3. Confounds

| Confound | Handling |
|---|---|
| Sampler drift between conditions | Freeze per model; record literal values |
| Context budget differences | Fixed per test; changes recorded as conditions |
| Prompt-shape differences between assisted/unassisted | Disclosed, never smoothed over |
| Model re-quantised mid-series | Treated as a new model; series restarts |
| Machine state (thermal, other load) | Recorded; anomalous runs flagged, not deleted |
| Owner familiarity with the test | Noted; the same scenario may be reused, its results compared |

## 4. Claim ceiling of the matrix

Every row is `NOT RUN`. The matrix establishes *what would count as evidence*, not that any evidence exists. A test that returns "no difference" is a result and must be reported with the same prominence as a favourable one.
