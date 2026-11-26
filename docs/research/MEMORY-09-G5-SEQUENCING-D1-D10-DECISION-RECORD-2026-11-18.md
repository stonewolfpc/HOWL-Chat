# MEMORY-09 — G5 Sequencing Decision Record (D1–D10) and Adversarial Review

**Status:** canonical documentation record. Consolidates decisions the owner closed **in
conversation**. It creates no implementation authority, no execution authority, no corpus
authority, no publication authority, and no acceptance authority.

| Field | Value |
|---|---|
| Identity | `MEMORY-09-G5-SEQUENCING-D1-D10-DECISION-RECORD-2026-11-18` |
| Created | 2026-11-18 |
| Authority | Transcriptive. Records decisions taken in the owner conversation; creates no new authority |
| Source material | The owner conversation record for the MEMORY-09 G5 sequencing decision series (D1–D10) |
| Governing frozen text cited | `docs/research/MEMORY-09-PREREGISTRATION-AMENDMENT-2-2026-11-15.md` §1–§2 (read only; unmodified) |
| Companion records | **None created.** This document is deliberately singular |
| Registry registration | **Outstanding.** `docs/DOCUMENT-REGISTRY.md` was not modified |
| Implementation status | **NOT STARTED / NOT AUTHORIZED** |
| Execution status | **NONE PERFORMED** |

Status vocabulary used in this record: **Closed** (owner closed the design decision) ·
**Specification-complete** (the decision and specification work is done) · **Provisional**
(estimate only, unreconciled) · **Unresolved** (raised, not decided).

---

## Executive summary

Nine decisions (D1–D8, D10) are content-complete. D9 is closed in status, but its decision text
was not recoverable and is not reproduced here. All ten are closed as **design** decisions; one
(D3) is specification-complete but **corpus-blocked**. The validator-unit estimate S₃ is **provisional ≈315–319** and
**unreconciled**.

An adversarial consistency review was then performed over the record. It is a **textual review**,
not an execution or source-validation exercise. It found:

- **0 contradictions** — no two accepted decisions are mutually unsatisfiable;
- **0 boundary violations** — no rule was found to permit something another decision prohibits;
- **8 findings** — 2 `provenance_gap`, 2 `false_closure`, 4 `coverage_gap`;
- **2 `no_finding`** compatibility conclusions.

**None of the eight findings is resolved.** Each remains an owner decision.

Adoption-blocking status is recorded in **three separate classifications**, which must not be
merged:

| Classification | Held by | Findings named |
|---|---|---|
| **Owner-declared adoption gates** | The owner | **None declared** |
| **Reviewer assessment** | Independent review | 1, 2, 5 |
| **Document author's assessment** | This record's author, **unratified** | 1, 5 |

"None declared" is a statement about the owner's **current act of classification**. It is **not** a
statement that no blockers exist, and it is **not** a finding that none were found. No finding has
been investigated and cleared by that statement, and it cannot serve as adoption clearance.

This owner position is recorded **here, in the executive summary**. It is **not** recorded as a
numbered decision in §4, **not** registered in `docs/DOCUMENT-REGISTRY.md`, and **not** committed,
so it retains an instance of Finding 1's provenance gap.

The single most important reading instruction: **"closed" means the owner closed the design
decision.** It does not mean implemented, validated, accepted in a new write-guard generation, or
experimentally demonstrated.

---

## 1. What this document is, and what it is not

**It is** a canonical transcription and consolidation of a decision sequence that previously
existed only in the owner conversation record.

**It is not:**

- proof that the decisions were previously integrated into repository files;
- evidence that any mechanism described here exists in code;
- evidence that any validator, probe, or test has been run;
- an authorization instrument;
- a work plan;
- an approval of any proposed follow-up;
- a resolution of any adversarial finding.

**Namespace note.** "D1–D10" in this document is the **G5 sequencing decision series**. It is a
different namespace from the architecture decision series referred to elsewhere in this
workspace as D1–D7 (see `docs/research/MEMORY-09-OWNER-DECISIONS-2026-11-17.md`). The two series
are unrelated. This document does not amend, interpret, or supersede any architecture decision.

---

## 2. Status ledger

Recorded exactly as supplied by the owner.

| Item | Status |
|---|---|
| D1 | Closed |
| D2 | Closed |
| D3 | Specification-complete; corpus-blocked |
| D4 | Closed |
| D5 | Closed |
| D6 | Closed |
| D7 | Closed |
| D8 | Closed |
| D9 | Closed |
| D10 | Closed |
| S₃ | Provisional ≈315–319; unreconciled |

**On "closed":** the owner closed the design decision. It does **not** mean implemented,
validated, accepted in a new write-guard generation, or experimentally demonstrated. A closed
decision whose mechanism does not exist is a closed decision about a mechanism that does not
exist.

**On D3's "specification-complete":** this describes the decision and specification work only.
It does **not** establish corpus availability, implementation completeness, validator coverage,
or experimental readiness.

**On S₃:** provisional, unreconciled, and an estimate — not a coverage claim. Its ≈55-unit D8
component was recorded as an **upper bound**. Earlier provisional figures are superseded and are
not reproduced here.

---

## 3. Governing output vocabulary (D2's frozen reference)

D2 designates Amendment 2 §1 as the governing four-cell G5 output vocabulary. The following is
quoted **verbatim** from
`docs/research/MEMORY-09-PREREGISTRATION-AMENDMENT-2-2026-11-15.md`, and is not paraphrased,
extended, or restated in this document:

> Two distinct claims, each reported on its own, never fused into one word like "works":
>
> | Claim | Requires | What it means |
> |---|---|---|
> | **C1 — net capability improvement** | **A1 > A0 on M1** | memory-equipped conversation beats no memory at all |
> | **C2 — selection-specific benefit** | **A1 > D on M1** | the gain is attributable to *relevance*, at matched tokens and shape |
>
> **Full capability claim = C1 ∧ C2.** Either alone is a weaker, separately named statement.
>
> | C1 | C2 | Permitted statement |
> |---|---|---|
> | ✓ | ✓ | full capability claim: memory improved capability, and relevance drove it |
> | ✓ | ✗ | capability improved, but **not** shown to be selection-specific (tokens/format could explain it) |
> | ✗ | ✓ | relevance helps **relative to decoys**, but not enough to beat no-memory — report as such |
> | ✗ | ✗ | no capability claim; if either comparison is significantly reversed, report a regression |
>
> Both C1 and C2 are computed on **M1** only (Amendment 1 §2), never on M2.

The same frozen text fixes the indeterminate floor that D7a's precedence depends on:

> **Floor for interpretability:** if `n_eff < 12`, the comparison is declared **indeterminate**
> regardless of its p-value, because at that point a "significant" result would rest on fewer than
> twelve independent conversations. This floor is fixed now, not after seeing ties.

And the tie rule, which no disposition may override:

> - **Tie:** `d_i = 0` exactly. **Ties are counted, reported, and excluded from the sign count.**
>   They are never treated as favourable, never as unfavourable, and never discarded silently.

Machine-level disposition identifiers observed in existing evidence artifacts are
`full_capability_claim`, `net_only`, `selection_only`, `regression`, and `indeterminate`. The
four permitted prose statements above remain the authority; identifiers are labels for them.

---

## 4. Decision record

### D1 — No equivalence conclusion

**Option A.** Closed.

**Operative prohibition, as recorded:**

> No equivalence conclusion from failure to demonstrate a difference. `not_supported` never means
> equivalence.

Consequences:

- There is **no** `indistinguishable` disposition. The disposition does not exist and must not be
  introduced.
- A result that fails to demonstrate a difference supports **no** statement of sameness. Absence
  of a demonstrated difference is not a demonstration of absence.
- `g5_verdict` remains **unavailable** (null) under the current frozen instrument.
- Existing statistical outputs are **preserved** — the frozen tests, tie accounting, effective
  sample size, and minimum-meaningful-difference treatment are not replaced, relaxed, or
  supplemented.
- The **vocabulary mismatch** between the frozen output vocabulary and any richer descriptive
  language is preserved, not resolved by inventing a category.
- The prohibition against inferring equivalence is **mandatory** and applies to summaries,
  reports, abstracts, and informal descriptions alike.

### D2 — Amendment 2 §1 governs, and cells are never collapsed

**D2-a.** Closed.

- Amendment 2 §1 is the governing four-cell G5 output vocabulary, quoted verbatim in §3 above.
- `net_only` and `selection_only` are **never collapsed** into one another, into a generic
  "partial", or into a weaker full claim. They are separate statements with separate meanings.
- **No new statistical category is invented**, and the frozen preregistration is **not amended**
  by any decision in this series.
- Both claims are computed on **M1 only**, never on M2.

### D3 — Evidence architecture, expected counts, and corpus-blocked status

**Specification-complete; corpus-blocked.**

**Architecture.** Three layers:

| Layer | Unit | Role |
|---|---|---|
| **L0** | Run plan | Declares the intended sample and the derived plan hash |
| **L1** | Session envelope | One per conversation × arm × repeat; carries the session's evidence |
| **L2** | Probe receipt | One per probe; the indivisible unit of scored evidence |

**Expected counts, as recorded:** 16 conversations × 3 arms (A0, A1, D) × 3 repeats × 13 probes
per conversation. This yields **144 session envelopes** and **1,872 probe receipts**. The G5
per-arm disposition table operates over **48 cells** (16 conversations × 3 arms); the
stopping-rule inventory is expressed over **(conversation, repeat, arm)** cells, totalling 144.
Both counts appear in the record for different purposes and are not interchangeable.

*Note on provenance of these figures:* the code-level constants corresponding to the 144 and
1,872 totals were identified earlier in the same session and were **not re-verified** during this
documentation task.

**Required properties:**

- **Expected-versus-present classification** — every expected artifact is classified as present,
  missing, or superfluous; the classification is explicit, not inferred from totals.
- **Explicit path enumeration** — evidence is enumerated by declared path, not discovered by
  directory scan.
- **Five evidence categories** — evidence is partitioned into five declared categories, and a
  category is never silently substituted for another.
- **Plan-hash agreement** — the evidence must agree with the plan hash it claims to implement.
- **Terminal-status evidence** — the run's terminal status is recorded as evidence of what the
  process declared, with the standing limits imposed by D8-U5.

**Scored arms.** A0 (no memory), A1 (memory), D (decoys). The A0/A1/D scored-arm design is a
decision and is preserved here; it is not restated as an implementation.

**Known discrepancy, preserved.** A G4 plan-row discrepancy was documented. It remains
documented. The frozen plan is **not** rewritten to match the discrepancy, and the discrepancy is
**not** resolved by preferring either the plan or the observation.

**Corpus-blocked.** No corpus exists with the 16 frozen conversation specifications and frozen
probe identities. Without them no `plan_sha256` can exist, and D3 cannot leave
specification-complete status. This is the block, and it is a block, not a defect.

### D4 — Authorization provenance, identity, and instantiation blockers

**D4-a.** Closed.

- Authorization provenance is **P1**.
- **Authorization-document identity and validation rules** are part of the decision: an
  authorization document is identified by content, and validation is against declared rules —
  not against plausibility, filename, or location.
- **Execution authorization is not MEMORY-09 memory authority.** Authorisation to run the
  experiment does not confer authority over the memory system under test, and vice versa. The two
  are distinct instruments with distinct scopes.
- **Binding requirements** are declared: the authorization binds a specific plan, a specific
  instrument, and a specific environment. An authorization that names none of these binds
  nothing.
- **Authorship is asserted, not cryptographically proven.** The repository's object format is
  `sha1` and the relevant commits are unsigned (`%G? = N`). Under those repository conventions an
  authored document is an attributable claim, not a verifiable signature. This limitation is part
  of D4 and propagates to anything built on it (see Finding 2).

**The authorization cannot currently be instantiated.** The identified blockers are: no corpus
and therefore no frozen probe identities and no `plan_sha256`; no G07 tokenizer; no G12 model
binding; and no grant mechanism — `ExecutionNotAuthorised` is raised unconditionally for network
responders, so no path exists by which a grant could be honoured. D4 is a closed decision about
an instrument that cannot yet be issued.

**Publication authority is not among the things D4-a confers.** See D8-U4-ii(a).

### D5 — Immutable run evidence and opaque run identity

**Closed.**

- **Run evidence is immutable and committed**, including evidence of **failed runs**. A failed
  run's evidence is evidence; it is not discardable, and it is not corrected in place.
- The **`docs/authorizations/` convention** places authorization documents in a committed
  directory, under the document class D4-a defines.
- **Registry registration is a separate action** from evidence creation. Creating evidence does
  not register it; registration does not create evidence.
- **`run_id` is opaque**, unique, and **independent of `plan_sha256`**. It identifies a run, not a
  plan, and carries no derivable meaning. Two runs of the same plan are different runs.

### D6 — `write-guard-v3`: integrated generation and all-or-nothing acceptance

**Closed.**

- **Integrated `write-guard-v3` generation.** The v3 generation is adopted as an integrated whole,
  never piecemeal.
- **All-or-nothing acceptance.** A partial v3 acceptance is not a v3 acceptance. If any acceptance
  condition fails, the generation is not accepted.
- **Session-boundary environment observation** with **receipt references**: environment is
  observed at session boundaries and referenced from the receipt, not reconstructed later.
- **Separation of five environment evidence classes**: **authorized**, **expected**,
  **configured**, **effective**, and **observed**. These are five distinct claims about
  environment and are never merged into one "environment" field. An authorized value, an expected
  value, a configured value, an effective value, and an observed value may all differ, and the
  differences are the finding.
- **Protected-set design**: a protected set of **12** artifacts.
- **Acceptance-time check count**: **38 validators**.
- **Instruments and schema**: the v3 generation sets `MEM1_VERSION = "0.2.0"` and
  `receipt_schema_version = 2`. These are **design targets for the v3 instrument**, not a
  description of the current tree, which still carries a pre-v3 version string. The instrument
  version and the receipt-schema version are **two separate version numbers** and are never
  conflated.
- **Acceptance authority, as recorded:**

> Chain-based acceptance authority, no rollback.

  Authority to accept derives from the chain of prior acceptance, not from recency, not from a
  file's modification time, and not from the newest artifact present. There is **no rollback**:
  an accepted generation is not un-accepted, only superseded by a later link that stands on it.
  Chain traversal refusals are explicit — dangling link, competing invalidation, competing
  successor, cycle, digest mismatch, and zero or multiple heads are all refusals. **"Newest file"
  is never the resolution.**

  Note the limit this establishes and does not establish: the chain establishes **ordering and
  non-rollback**. It does **not** authenticate the authorship of its root or of any link (see
  Finding 2).

**Relation to existing state.** The standing guarantee is `write-guard-v2` at 24/24 protected
hashes, which remains intact. The relationship between v2's 24 and v3's 12 is **not defined by
D6** (see Finding 5).

### D7 — Disposition precedence and cell-scoped permitted statements

**D7a and D7b.** Closed.

#### D7a — Precedence

As recorded, disposition is resolved in this order:

> veto → indeterminate → reversed → both → one → none

- **veto** outranks everything. A vetoed comparison yields no capability disposition at all.
- **indeterminate** outranks any positive or negative outcome. This is anchored to the frozen
  floor: `n_eff < 12` declares a comparison indeterminate **regardless of its p-value**.
- **reversed** outranks partial support.
- **both → one → none** is the descending order of supported comparisons.

#### D7b — Cell-scoped statements (B3)

- The G5 statement is **cell-scoped**. There is no global collapse across cells.
- **`regression_comparisons` is a comparison-specific list**, naming which comparison(s) are
  reversed: `["c1"]`, `["c2"]`, or both. It is never a boolean and never a single generic flag.
- `regression` is **reserved** for the ✗✗ cell when a comparison is reversed. It is not a general
  label for "something went worse".
- Worked cases as recorded:
  - `(✓, reversed)` → `net_only` with `regression_comparisons = ["c2"]`.
  - `(reversed, ✓)` → `selection_only` with `regression_comparisons = ["c1"]`.
- **Indeterminate controls disposition, not disclosure.** A cell whose disposition is
  indeterminate because of the `n_eff` floor still discloses any reversed comparison through
  `regression_comparisons`. Nothing is hidden by being outranked.

#### J09 and the assertion-change register

J09 currently pins four combinations under the **global** reading that B3 replaces. Under D7b,
J09's **v3** expectation changes: the `(✓, reversed)` cell becomes `net_only` with
`regression_comparisons = ["c2"]`.

- The **historical v2 assertion is untouched**. It is not edited, not regenerated, and not
  reinterpreted. It remains the record of what the v2 suite asserted.
- An **assertion-change register is required** to record the divergence between the v2 assertion
  and the v3 expectation.
- **The v2 assertion has not been changed.** This document records an obligation, not a completed
  edit.
- B3 exposed a **coverage gap**: `selection_only` and `no_claim` had no assertion at all before
  this decision. That gap is recorded, not closed.

#### The implementation D7 replaces

`eval/mem1/stats.py` implements the **global** reading — `evaluate_claim` and `decide` collapse
across cells. D7b replaces that reading at specification level. **The code has not been changed**
and still implements the global reading.

### D8 — Evidence integrity, invalidity, and the limits of the terminal record

**Closed**, via sub-decisions U1–U5.

#### Duplicate and identity checks

Two silent-preference sites were identified in the existing code and are recorded here as defects
of record, not repaired:

| Site | Behaviour |
|---|---|
| `eval/mem1/stopping.py:45` | Cell set collapses duplicates; duplicated evidence inflates `observed` and can drive `missing_fraction` negative |
| `eval/mem1/stats.py` (`build_pairs`) | Dict keyed by conversation id — last write wins |

Duplicate detection exists only for **configuration** (`placebo.py`, `fixtures.py`,
`distribution.py`) and never for **run evidence**. The intended principle is already stated in the
model: *"nothing is silently normalised, deduplicated or fuzzy-matched"* (`eval/mem1/model.py:16`).
The cheap correction — an `assert_unique_evidence(...)` before dict construction, with no
signature change — was recorded as available and **was not applied**.

#### D8-U1 — Finding identifiers

- **A-10 has no implementation.** It names a condition, not a check.
- **C4b merges into A-10** as `authorization_plan_mismatch`.
- **C4a remains distinct** as `artifact_plan_mismatch`.

#### D8-U2 — Three families, never interchangeable

| Family | Members |
|---|---|
| **`contradiction`** | C1–C4a |
| **`invalid`** | C5 = `probe_spec_mismatch`; parse failures; field failures; quota failures |
| **`incomplete`** | Expected-but-absent evidence |

The families are **never interchangeable**, and all three produce the **same refusal**. A
contradiction is not an invalidity; an invalidity is not an incompleteness. Conflating them would
convert "the evidence disagrees with itself" into "no evidence was examined".

#### D8-U3 — Closed exclusion list, instrument-level restart

- **Amendment 2 §4's exclusion list is closed.** Rule 1: *"Invalidity is decided by the criteria
  above"* — the four-item parenthetical (missing receipt; span still present; `seed_honoured`
  false or unavailable; token accounting absent). Naming a condition grants it **pair-level
  treatment**, not field-level scope.
- **C6 is instrument-level.** Rule 4 applies: *"if the cause is instrument-level, the whole fixed
  sample restarts per the stopping rule"*. There is **no pair exclusion** and **no replacement
  conversation**.
- Token-accounting absence is already **criterion 4**; it is not a separate exclusion.
- The rule is anchored to **frozen deterministic derivation contracts** only. It is not an
  open-ended "impossibility" test, and no exclusion may be invented beyond the closed list.

#### D8-U4-i — One invalidation per report

- **Exactly one invalidation per report** — per report, not per run.
- **Void is absorbing**: once invalidated, the report cannot become valid again.
- Later findings discovered after invalidation are routed to **`contradiction` records**, not to
  additional invalidations.
- **There is no withdrawal.** A false invalidation implies a verifier defect, but does **not** by
  itself establish an Amendment 2 §4 (U3) contract violation. Those are different findings.

#### D8-U4-ii — Publication of invalidation

- Invalidation is published **mechanically**. It requires **no separate grant**.
- **The report's author cannot veto its own invalidation.** The anti-suppression interest is
  identical to the reason a party may not judge evidence it belongs to.
- **U4-ii(a) — publication permission is a pre-execution scope.** The authorization document
  carries `permits_publication`, `publication_scope`, `publication_binds`, and `granted_at`, with
  `granted_at` **predating the first envelope**. Publication authority is therefore **not**
  inherited from D4-a execution authorization. This is an explicit correction made during the
  sequence: execution authorization does not imply publication authorization.
- **U4-ii(b) — checker independence is asserted, not proven.** No testable unit exists (0 units).
  It belongs in documented limitations, not in the validator count, unless a future operating
  environment provides a testable mechanism.

#### D8-U4-iii — `invalidated` is transitional

State progression:

> `unreported → reported → invalidated → superseded`

plus `withheld` as a separate, visible state.

- `invalidated` is **visible and persistent**. It is not a silent dead end and it is not a
  terminal tombstone.
- A successor does **not** restore the voided report. It establishes a **new report under the
  original run's authority**.
- **Reporting repair** — using existing artifacts — retains `run_id` and the publication grant,
  and is permitted **only mechanically**: locating an expected artifact that was previously
  omitted, or establishing that two paths hold identical bytes.
- **Reporting repair cannot become evidence adjudication.** Choosing between two conflicting
  artifacts because one appears more plausible is **adjudication, which is forbidden**. If
  existing evidence cannot re-establish the correct basis, the path is a **new run with new
  authorization**, not a reporting repair.
- **Evidence repair** requires a **new `run_id` and new grants**. It does not modify previously
  recorded evidence (see the D5 compatibility conclusion in §6.4).
- **A refusal report is a valid successor.** A successor need not assert a capability result; it
  may report `input_invalid` and explain why no result is admissible. The chain may therefore end
  in a **citable refusal** rather than an unsupported experimental claim.
- An **unrepairable run remains visibly invalidated**.

#### D8-U5 — The terminal record has no adjudicative authority

**Categorically.** The terminal record may be **read, compared, and contradicted**. It may
**never decide**.

| # | Reason |
|---|---|
| 1 | **A party cannot judge.** The terminal record is produced by the execution process whose artifacts may be in conflict — the same anti-suppression logic that forbids the report author from vetoing its own invalidation |
| 2 | **Earliest and worst-informed.** It is written at execution end, before verification, and cannot know about contradictions established later |
| 3 | **A declaration is not evidence of the declared thing.** "144 envelopes were written" is a claim about what the process did, not an establishment that 144 valid envelopes exist — the D3 distinction between run-status artifact and evidence |
| 4 | **Wrong chain position.** Authority is chain-based (D6); the terminal record sits at the start of the chain, while adjudication belongs at the head — the verifier over the current evidence set |
| 5 | **It can be wrong about its own counts.** A 145th artifact appearing after the record was written makes its numbers stale. It is a witness that can be mistaken, which is precisely why it cannot be a judge |

| Act | Status |
|---|---|
| Read the terminal record | **Allowed** |
| Compare its counts against artifacts | **Allowed** — detection |
| Find the terminal record itself in a contradiction, record both sides, refuse | **Allowed** |
| Establish which of two conflicting artifacts is valid | **Forbidden** |
| Supply the value that resolves a contradiction | **Forbidden** |
| Reconstruct a missing artifact from its declaration | **Forbidden** |

Sharpest consequence: if envelopes are missing and the terminal record says they existed, the run
is **`incomplete`** — not "we know they existed".

#### The five governing rules of D8

1. **Record both and never prefer.**
2. **Membership may be tested; materiality may not.**
3. **A party may not judge evidence it belongs to.**
4. **Publication permission predates the result.**
5. **Void is absorbing but never silent.**

Trigger for invalidation: an **established defect in the evidence actually relied upon**. The
membership test is permitted; the materiality test is forbidden.

#### Recorded correction within D8

C5 was **initially** classified as a contradiction and was **corrected** to `invalid`
(`probe_spec_mismatch`). U4's trigger was correspondingly generalised to "established
admissibility defect" while the finding's **classification was not changed again**. Both the
correction and the resulting classification are recorded. The correction was not silent.

### D9 — Closed; decision text not recoverable

**Status: Closed** (owner).

**Limitation.** The D9 decision text was **not recoverable** during the later adversarial review.
It was not reconstructed, paraphrased, inferred, or filled from repository contents. This document
records D9's existence and closed status only, and attributes **no** content and **no** finding to
it.

Consequence for readers: the record is **complete for D1–D8 and D10**, and **status-complete but
content-incomplete for D9**. Any future reader requiring D9's operative rules must obtain them
from the owner; they are not here and must not be inferred from this document's structure.

### D10 — Report identity, determinism, and G-06

**Closed.**

#### G-06 — the name resolves to two unrelated gaps

| Source | G-06 means | Track |
|---|---|---|
| `GATE-B-SOURCE-INVESTIGATION` | Chat ID absent from character-memory retrieval and the final request | Application defect (Gate B) |
| `MEMORY-01-RESEARCH-SYNTHESIS` | **"No retrieval evaluation on personal chat scale"** — recommends recall/nDCG on labelled Howl evidence; lexical vs dense vs hybrid; **P1** | Research gap |

The MEMORY-01 gap is the one that bears on MEMORY-09. The Gate-B item is an unrelated application
defect that no G5 validator could touch.

**Disposition.** MEMORY-09 measures **answer-level** M1/M2 against a fixed store. It does not
measure retrieval *ranking* quality: no labelled relevant sets, no recall@k, no nDCG, no
lexical-versus-dense comparison.

| Question | Answer |
|---|---|
| Does MEMORY-09 close G-06? | **No** |
| Does G-06 belong in the MEMORY-09 validator inventory? | **No** — different experiment, inputs do not exist |
| Does MEMORY-09 contribute anything to it? | **Partial, descriptive only** — per-class ladder counts and retrieval diagnostics are observational, not recall/nDCG |
| Units attributed | **0** |

G-06 is **carried forward to MEMORY-01**, not counted in S₃. Adding it would inflate the inventory
with units that cannot be exercised and would misrepresent an open research gap as a covered
control.

#### Report identity

| Element | Form |
|---|---|
| Body | Canonical serialisation of derived values |
| `run_id` | **Inside** the hashed body, so identical evidence from different runs yields different digests |
| `report_digest` | Hash of the canonical body |
| Identity | `(run_id, report_digest)` |
| `supersedes` | Names a prior `report_digest` |

#### Determinism

> **A report is not reproducible from scratch. It is a deterministic function of recorded
> evidence.**

Model responses are not reproducible; that is why evidence exists. The existing stance is stated
in `eval/mem1/distribution.py`: *"`report` is reproducible by construction: given the same store
and the same expectations."* D10 generalises it:

> Given the same evidence set, the same plan, the same environment binding, and the same `run_id`,
> the report body and digest are byte-identical.

Four rules make it hold:

| # | Rule | Why |
|---|---|---|
| 1 | **No clock, no paths, no host state in the body** | Wall-clock timestamps and absolute paths are not functions of the evidence |
| 2 | **Volatile fields live outside the hashed body** | Publication time and verifier identity belong in a `meta` block excluded from the digest |
| 3 | **Canonical ordering everywhere** | Every collection sorted before serialisation; never raw set or dict iteration |
| 4 | **Canonical bytes** | Sorted keys, UTF-8, no whitespace variance — one serialisation, not "equivalent" ones |

Rule 3 is a live hazard: Python string hashing is randomised per process, so iterating a set of
tuples containing strings yields different orders across runs.

**Existing anchors** for the convention: checks `C02` (run plan hash is deterministic), `S11`
(declared distribution reproducible from store and expectations), `K10` (A1 deterministic across
repeats), `B02` (seed derivation deterministic), `G04` (seed table hashes deterministically).
D10 extends an established convention rather than introducing one.

**Unit impact:** ≈+6.

---

## 5. Adversarial review

### 5.1 Nature and stated limitations

The review was a **textual consistency review** of the decision record. It was **not** an
execution, a source-validation exercise, or an experimental review. Its limitations are preserved
here as stated:

| Limitation | Statement |
|---|---|
| No artifact on disk | **No integrated D1–D10 record existed on disk** at the time of the review |
| Source of the review | The record under review was **the conversation record** |
| No execution | **No execution was performed** |
| D9 text | **D9's full text was unavailable and was not reconstructed** |
| Coverage | The review covered **D1–D8 and D10**, and D9 only where referenced by other decisions |
| Attribution | **No finding is attributed** to unrecovered D9 text |

Its purpose was bounded to five detection classes — contradictions, coverage gaps, boundary
violations, provenance gaps, and false closure — and explicitly **not** to expand the design. A
related boundary: an actual inconsistency must be distinguished from a missing implementation
mechanism or an acknowledged limitation. Those are different findings and were **not** converted
into new requirements.

### 5.2 Repository search results, as reported

These are **reported observations from the supplied review**, not results reproduced by this
documentation task.

| Search | Matches |
|---|---:|
| `regression_comparisons` \| `permits_publication` \| `publication_scope` | 0 |
| `write-guard-v3` | 0 |
| `net_only` \| `selection_only` | 4, all in `stats.py` and the v2 receipt; none in a decision record |
| `300`-level S₃ / `specification-complete` | 0 |

### 5.3 Findings

Findings are **observations from a textual review**. They are not requirements, not decisions,
and not authorizations. They remain **unresolved**.

**Finding 1 — `provenance_gap`.** No citable canonical artifact existed for the D1–D10
conversation decisions. D4-a requires authorization documents to exist as committed files and D6
requires chain-based authority citing prior links, yet the decisions themselves had no artifact,
so a future run could not cite the decision that governs it.

**Finding 2 — `false_closure`.** Chain-based acceptance authority does not authenticate the
authorship of its root or of its links. D4-a authorship is **asserted** under repository
conventions (sha1, unsigned commits), not cryptographically proven. Authority derives from the
chain, so the asserted root propagates to every link.

**Status: partially resolved at the owner-decision and specification-accuracy levels; open
overall.** The owner accepted the specification correction and accepted the existing trust model
as explicitly bounded. **Operational recovery remains open and undemonstrated.** This finding is
**not** reclassified as fully resolved, and the §5.5 totals are unchanged.

#### 5.3.2.1 — What the chain does and does not establish

| Property | Established by the current design |
|---|---|
| Internal integrity / consistency, for presented artifacts and their digest relationships | **Yes** |
| Citation-topology ordering | **Yes** |
| Anchored integrity | **No** |
| Authenticity of root or links | **No** |
| Actual temporal order of creation | **No** |
| Fork detection | **Conditional** — only when competing artifacts are co-presented to the same validator |
| Fork prevention | **No** |
| Resolution of contested chain authority | **No in-system mechanism** |

- A hash-linked chain establishes **internal consistency**, not an independently authenticated
  history.
- Without an externally trusted reference, the **root is replaceable**.
- A fully rewritten directory with recomputed digests is **observationally indistinguishable**
  from the genuine chain, to a validator with no reference outside it.
- The old chain is **not disproved** by a forged alternative; it may simply be **absent** from the
  evidence presented.
- Hashes establish neither **authorship** nor **actual chronology**.

This is a **trust-anchor limitation**, not a hash weakness.

#### 5.3.2.2 — Fork handling: two failure modes

| Presentation | Validator behaviour | Consequence |
|---|---|---|
| **Co-presented competing branches** | Detects competing successors and refuses | **Visible stall** |
| **Presentation asymmetry** | Different validators receive different internally consistent presentations; each may see a single head | Potential **silent acceptance** of divergent histories |

> The refusal set detects conflicts only among artifacts co-presented to one validator. It does
> not guarantee universal fork visibility. Presentation asymmetry can produce a clean pass rather
> than a refusal.

This is **not necessarily a validator malfunction**. The limitation concerns the **evidence
presented** to the validator, not the correctness of the validator.

#### 5.3.2.3 — Safety, liveness, and authority

- **Safety:** ambiguous chains are refused, **when the ambiguity is visible to the validator**.
- **Liveness:** there is **no in-system mechanism** to resume acceptance after contested chain
  authority stalls progression.
- **Authority:** the **owner is the external resolver**, outside the system's adjudication
  boundary.

D8's refusal to invent an in-system adjudicator is **intentional** — the owner is not delegating
authority to an invented mechanism. But **a designated external resolver and an established,
tested recovery procedure are different claims**, and only the first is supported here. No verified
external recovery procedure is claimed to exist. D8-U4-iii's refusal-report successor is a
**reporting path**, not a recovery procedure for chain authority.

#### 5.3.2.4 — Three independent questions

| Dimension | Owner position | Meaning |
|---|---|---|
| **Specification accuracy** | **Accepted** | The guarantees must be bounded to what the mechanism establishes |
| **Trust-model acceptance** | **Accepted** | The owner accepts the existing internal-only trust model and its limitations |
| **Operational recovery** | **Open / undemonstrated** | No verified recovery procedure has been established |

Completing the first two dimensions does **not** answer the third, and must not be read as doing so.

#### 5.3.2.5 — Owner decision on the options

This is an explicit **owner trust decision**, recorded as such.

| Option | Position |
|---|---|
| **(a1)** Specification accuracy | **ACCEPT** — the specification must state only what the existing mechanisms establish |
| **(a2)** Trust-model acceptance | **ACCEPT** — internal consistency without independently authenticated root or link authorship is accepted as a limitation of this instrument |
| **(b)** Cryptographic signatures | **DEFER** — no signature infrastructure or key-management mechanism is introduced |
| **(c)** External head witnessing | **DEFER** — no external witness or witnessing dependency is introduced |
| **(d)** Asserted identity at each use | **NOT SELECTED AS AN AUTHENTICATION MECHANISM** — improves accountability and visibility only; it does not prove identity, and is not a new requirement |

This acceptance is **not** a finding that the chain is authenticated, tamper-proof, or resistant
to replacement by a party with sufficient repository control.

#### 5.3.2.6 — What acceptance does and does not mean

The chain is an **accountability instrument for its custodian**, not a security boundary against
repository control. It protects against **accident** and against **third parties without write
access**; against the repository controller it provides **auditability, not protection**. "No
rollback" is therefore a **procedural** claim, not a physical one.

A corresponding bound applies to D8: U4-ii forbids a report's author from vetoing its own
invalidation, but it does **not** protect against a directory controller deleting the chain
wholesale. If the history is replaced, the recorded invalidations, refusals, and contradictions go
with it.

#### 5.3.2.7 — Anchor failure modes (for any future reconsideration of (b) or (c))

| Anchor failure | Effect |
|---|---|
| **Compromised** | An anchor may attest to a forged chain, producing a confident but incorrect verification |
| **Lost** | No comparison is possible; may be indistinguishable from never having had an anchor |
| **Unavailable at acceptance** | The system must block or explicitly disclose degraded, unanchored acceptance |
| **Disputed** | Competing claimed anchors reproduce the fork problem at the trust-anchor level |

- A **compromised** anchor can be **worse than no anchor**, because it converts an honest
  unverified state into a silently incorrect verified state.
- A **disputed** anchor **relocates** the authority dispute rather than eliminating it.
- Signatures, witnesses, and identity assertions **relocate** trust — to keys, external parties, or
  people — and do not eliminate the need for an ultimately trusted authority.
- Any future reconsideration of (b) or (c) must account for the new anchor's **compromise, loss,
  unavailability, and dispute** scenarios.

None of these mechanisms is introduced as an implementation requirement.

#### 5.3.2.8 — Relationship to D8-U4-ii(b)

Finding 2 remains **distinct** from D8-U4-ii(b), which concerns **checker independence**. A correct
checker cannot authenticate an untrusted root merely by checking it correctly; an authenticated
chain cannot establish that its checker is independent. Neither resolves the other, and resolving
one must not be recorded as resolving both.

**Finding 3 — `coverage_gap`.** A refusal report may be unciteable when the original
pre-execution authorization does not grant publication authority. U4-ii(a) requires publication
permission to predate the first envelope; U4-iii establishes a refusal as a valid successor. If
the original grant omits publication entirely, the run can publish neither its result nor its
refusal — the silence U4-iii exists to prevent.

**Finding 4 — `coverage_gap`.** `g5_verdict: null` does not explicitly prohibit **descriptive
tallies** of per-cell dispositions. D1 forbids a G5-level verdict; nothing in the record prevents
a report stating "31 of 48 cells supported".

**Finding 5 — `coverage_gap`.** The relationship between the D6 protected set of **12** and the
v2 identity verification of **24/24** is undefined. If both numbers count protected artifacts, v3
governs a strictly smaller set than v2 and could permit overwriting artifacts v2 refused.

**Finding 6 — `coverage_gap`.** The relationship between the D6 validator inventory and the
historical **142-check baseline** is undefined. Are the 38 validators a subset, a replacement, or
additions? D7b further requires J09's expectation to change, which sits against the preservation of
that baseline.

**Finding 7 — `provenance_gap`.** D7a's precedence vocabulary (`reversed`, `both`, `one`, `none`)
is not formally mapped either to the Amendment 2 §1 four-cell table or to D7b's dispositions. The
categories are derivable but never defined, and precedence cannot be implemented over categories
whose identity is inferred.

**Finding 8 — `false_closure`.** "Specification-complete" could be misread as validated or
coverage-complete, despite D3's corpus-blocked status and the provisional S₃. Completeness is over
*decisions*; it is not coverage of units and it is not validation.

### 5.4 `no_finding` results

Four apparent tensions were examined and resolved. These are **compatibility conclusions**, not
guarantees that the mechanisms exist.

1. **D5 immutability vs D8-U4-iii evidence repair** — compatible **on the record's terms**:
   evidence repair does not modify previously recorded evidence; it requires a new `run_id` and new
   grants.

   **Caveat, not a resolution.** Reporting repair retaining the same `run_id` is narrower than it
   appears: locating a previously omitted expected artifact **changes the committed evidence set
   under the same `run_id`**. The mechanical-versus-adjudicative line does the work here, and the
   record supplies no testable discriminator for it. This conclusion is therefore **not proof**
   that immutability survives every repair case. The boundary remains an **owner question** — it is
   not a ninth finding, and it is not a resolved compatibility.
2. **D5 opaque `run_id` vs D10 inclusion of `run_id` in report identity** — compatible: the digest
   varies with lineage while remaining a deterministic function of evidence, plan, environment, and
   `run_id`.
3. **D2 vocabulary vs D7 precedence/disclosure** — compatible when disposition and regression
   disclosure remain distinct: indeterminate controls disposition, not disclosure, and the reversed
   comparison remains visible in `regression_comparisons`.
4. **D5 terminal-record registration vs D8-U5 refusal of adjudicative authority** — compatible:
   registration confers identity, not authority.

### 5.5 Review totals, as reported

| Class | Count |
|---|---:|
| `provenance_gap` | 2 |
| `false_closure` | 2 |
| `coverage_gap` | 4 |
| `boundary_violation` | 0 |
| `contradiction` | 0 |
| `no_finding` | 2 |

**No two accepted decisions were found to be mutually unsatisfiable.** The design was found to be
internally consistent; what it contains is undecided states and two claims that outrun their
mechanisms. These counts are not silently reinterpreted or recategorised by this document.

---

## 6. Separation of categories

These categories are kept distinct throughout and must not be merged.

### 6.1 Owner decisions

D1–D10 are **closed design decisions**, recorded in §4. D9 is closed with its text unavailable
and unreconstructed.

### 6.2 Specification status

D3 is **specification-complete but corpus-blocked**. This is a statement about decision and
specification work, not about corpus availability, implementation completeness, validator
coverage, or experimental readiness.

### 6.3 Adversarial findings

Findings 1–8 in §5.3 are **observations from a textual review**. They are open.

### 6.4 `no_finding`

The four compatibility conclusions in §5.4. They are **not** guarantees that any mechanism exists.

### 6.5 Unresolved

The questions raised by Findings 1–8 that have **not received owner decisions**. Listed in §7.

### 6.6 Not authorized

Listed in §8.

### 6.7 Distinctions that must not be collapsed

| Distinction | Why it matters |
|---|---|
| **Design closure** vs **implementation or acceptance** | A closed decision about a non-existent mechanism is still closed |
| **Readiness** vs **authorization** | Being ready to run is not permission to run |
| **Execution authorization** vs **publication authorization** | Separate instruments; D8-U4-ii(a) is the explicit correction |
| **Instrument version** vs **receipt-schema version** | Two separate numbers; `MEM1_VERSION` ≠ `receipt_schema_version` |
| **Source identity** vs **generation lineage** | Content identity of an artifact is not its position in an acceptance chain |
| **Historical acceptance authority** vs **current-tree match** | A prior acceptance is not re-established by a matching current tree |
| **Contradiction** vs **invalidity** vs **incompleteness** | Three families, never interchangeable, same refusal |
| **Evidence membership** vs **evidence materiality** | Membership may be tested; materiality may not |
| **Report invalidation** vs **report supersession** | Invalidation is not restored by successorship |
| **Terminal-record registration** vs **adjudicative authority** | Registered evidence is not thereby a judge |

---

## 7. Unresolved — owner decisions outstanding

Each of the following is **open**. This document does not decide any of them, and the review's
proposed follow-up sequence is **not** an approved work plan.

1. **No citable canonical record** (Finding 1) — whether and how to produce a decision artifact
   that a future run can cite.
2. **Operational recovery for contested chain authority** (Finding 2) — the **specification
   accuracy** and **trust-model acceptance** dimensions are **owner-accepted**; this third
   dimension remains **open**. The owner is the designated external resolver, but **no verified
   recovery procedure has been established or tested**, and none is claimed. Signatures (b) and
   external head witnessing (c) are **deferred**; asserted identity (d) is **not selected** as an
   authentication mechanism.
3. **Refusal-report publication** (Finding 3) — what happens when the original grant carries no
   publication authority.
4. **Descriptive cell tallies** (Finding 4) — whether the report may present an aggregate of cell
   dispositions.
5. **Protected-set relationship** (Finding 5) — the relationship between v3's 12 and v2's 24,
   including whether a superset rule is required.
6. **Validator count vs the 142-check baseline** (Finding 6) — subset, replacement, or addition,
   and where J09's changed expectation lives.
7. **D7a precedence vocabulary mapping** (Finding 7) — the formal mapping to Amendment 2 §1's
   four cells and to D7b's dispositions.
8. **Qualification of "specification-complete"** (Finding 8) — whether and how the phrase is to be
   qualified.
9. **D9's decision text** — not recoverable in the record; obtainable only from the owner.
10. **Adoption-gate declaration** — whether any finding is formally classified as an adoption gate
    for a new write-guard generation. **None is currently declared**, which is not a finding that
    none exists.
11. **Reporting-repair / `run_id` boundary** — whether repairing an omitted artifact under the
    original `run_id` preserves immutability. An owner question raised in review; **not** a ninth
    finding and **not** a resolved compatibility conclusion.

---

## 8. Not authorized

Nothing in this document authorizes any of the following, and the creation of this document does
not imply any of it:

- implementation of any decision in this record;
- source-code or test/harness changes;
- running any test, probe, preflight, review, benchmark, or other execution;
- corpus preparation or creation;
- creation of authorization documents;
- creation of run plans, receipts, terminal records, or reports;
- registry modification or registration;
- changes to protected artifacts, frozen documents, hashes, sidecars, or receipts;
- publication of any result or refusal;
- staging, committing, pushing, or otherwise changing Git history;
- resolving any finding in §5.3;
- treating any proposed follow-up as an approved decision.

Gates that remain closed independently of this document: no G4 run, no G07 tokenizer work, no G12
model binding, no corpus preparation, no G5 execution, no publication, and no G5-rubric adoption.

---

## 9. Provenance and limitations

- The decision content in §4 was **supplied in the owner conversation record**. This document is a
  canonical transcription and consolidation of that record.
- It is **not proof** that these decisions were previously integrated into repository files. The
  reported repository searches in §5.2 found no such integration.
- The repository search results in §5.2 are **reported observations from the supplied review**,
  not results newly reproduced by this documentation task.
- The four-cell output vocabulary in §3 is quoted verbatim from the frozen Amendment 2 file, which
  was **read only and not modified**.
- **No new tests, source inspection, execution, or experimental evidence were generated** by this
  documentation task.
- Observable code-level facts referenced here — the global reading in `stats.py`, the set-collapse
  at `stopping.py:45`, the unused `plan_sha256`, `MEM1_VERSION`'s current value, the absence of
  authorization machinery — were established **earlier in the same session** and were **not
  re-verified** during this documentation task.
- This document does **not** establish implementation, validation, authorization, or experimental
  readiness.
- This document's creation **does not authorize any subsequent work**.
- This document does **not** prove the authorship, timing, or cryptographic authenticity of the
  original conversational decisions. Authorship is asserted, consistent with D4-a's own limit:
  the repository is `sha1` and the relevant commits are unsigned.

---

## 10. Verification performed

| Check | Result |
|---|---|
| No pre-existing D1–D10 record (overwrite avoided) | Confirmed by search before writing: 0 matches for `D1–D10` and for the chosen filename stem |
| Documentation convention followed | `docs/research/MEMORY-09-<TOPIC>-<DATE>.md` with status block and field table |
| Registry untouched | `docs/DOCUMENT-REGISTRY.md` not modified; registration reported outstanding |
| Frozen artifacts untouched | No frozen document, hash, sidecar, or receipt modified |
| Exact quoted language preserved | Amendment 2 §1–§2, D6 acceptance authority, D1 prohibition, D8-U5 table |
| Internal consistency | Status ledger, decision cross-references, findings, and unresolved list reconciled against each other |
| Execution | None performed |
| Git impact | One new untracked documentation file; nothing staged, committed, or pushed |

**End of record.**
