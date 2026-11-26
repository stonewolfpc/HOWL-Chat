# Phase 1 — Discover · Tranche 02 (2026-11-21 · contract pack and MEMORY-09 G1 chain)

**Checkpoint:** `HEKI-CP-2026-11-22-02`
**Method:** same standard as Tranche 01 — front matter, status block, scope and relationship statements read in
context. Nothing edited. Machine extraction nominated the candidates; every judgement below is from reading.

**Coverage after this tranche:** 20 of 808 documents verified (2.5%).

**Correction to a prior statement:** Tranche 01 announced "the ten 2026-11-21 `research/MEMORY-09-*` records".
There are **six**. The count was wrong and is corrected here rather than in the old checkpoint.

---

## 1. `HEKI-D-0011` — Contract paper register `HF-CONTRACTS-2026-11-21-01` (contracts/README.md)

- **Purpose:** index and dependency register for the contract papers; records FC-7's deliberate absence.
- **Status (verbatim):** "**Review drafts, review-complete.** FC-1–FC-6 accepted as reviewed proposed drafts
  (2026-11-21); **none is approved**."
- **Authority chain:** policy packet `HF-CONTRACT-INTAKE-2026-11-20-01`; owner decisions D1–D8 as tabulated.
- **Dependency order (register §2):** FC-1 → (FC-2, FC-4) → FC-3 → FC-5 → FC-6; FC-7 waits on D8.
- **Material statement:** "FC-7 is deliberately absent … Filing FC-7 now would either pre-empt that decision
  or create a placeholder. The contract register records the block rather than manufacturing a paper."
- **Adoption:** review-accepted drafts; **approved: none**; implementation authorization: **none**.

## 2–7. `HEKI-D-0012`…`HEKI-D-0017` — FC-1 … FC-6

| Paper | Gate | As-written status | Declares |
|---|---|---|---|
| FC-1 Conversation State | M1 | PROPOSED — not approved | append-only event ledger is durable truth (D1 option C); state/branch graph is derived; six required identities; **no positional identity**; legacy identity **Unresolved** |
| FC-2 Composition and Context | M2 | PROPOSED — technology selection unresolved | PromptPlan preflight artifact; D4 is the whole pipeline, not a tokenizer choice |
| FC-3 Multi-Character Isolation | M3 | PROPOSED — not approved | hard flush on speaker change; isolation is a property of composed prompt, not display name; voice identity survives flush |
| FC-4 Content Ownership | M0B/M5 | PROPOSED — bypass inventory incomplete | canonical library owns authored objects incl. **Character**; memory record authority **Unestablished here** (`FC4-04a`); character/persona boundary **Unresolved** (`FC4-06`) |
| FC-5 Capability, Resource, Handoff | M4 | PROPOSED — reservation questions unresolved | selection≠admission≠allocation≠ready≠release; D6 option A temporary; `FC5-04`: **D6 is not T0.5**; config records FC5-27–FC5-31 present in the paper |
| FC-6 Evidence and Acceptance | All gates | PROPOSED — not approved | acceptance = evidence at a stated level; claim ladder Requested→Admitted→Attempted→Observed→Verified/Unknown/Failed; no promotion by adjacency |

**Uniform and decisive:** every paper says **"PROPOSED — not approved"** and **"Implementation authorization:
None"**. Review acceptance is not approval, and neither is implementation.

---

## 8. `HEKI-D-0025` — MEMORY-09 Capture Completeness audit (G1)

- **Purpose:** establish the capture envelope and the gaps, before any classifier is measured.
- **Finding:** "The app has the context; the memory service does not receive it."
- **Material rules:** `CCG-01` capture-incomplete ≠ classification-unresolved; `CCG-03` G1 is a prerequisite —
  classifying capture-incomplete examples "measures nothing except the classifier's ability to guess".
- **Connectivity:** **in = 0, out = 0, identifier mentions = 0.** Fully orphaned, despite being the origin of
  the G1 gate that later documents cite by name.

## 9. `HEKI-D-0026` — G1 Capture Fixture Contract

- **Purpose:** define what G1 had to demonstrate (cases F1–F10 + character-tell).
- **Status:** "fixture contract; the test-only fixture subsequently passed (G1R-22)"; production wiring
  unproven.
- **Material rules:** `GFC-01` no collapse of capture-incomplete into classification-unresolved; `GFC-02` a
  confident class cannot compensate for missing provenance; `GFC-04` G1 is a contract, not a classifier test.

## 10. `HEKI-D-0027` — G1 Capture Fixture Results

- **Result (Observed):** `go test ./internal/backend/memory -run 'TestG1' -v` → **PASS**, 11 subtests;
  package run `ok … 0.804s`.
- **Scope discipline:** `G1R-03` only new file is the harness; zero production files changed.
- **Critical limit stated by the document itself:** "The data-model contract passes; production capture wiring
  is still unproven." **E2 at best, on the fixture path only.**

## 11. `HEKI-D-0028` — Production Capture Boundary Proposal

- **Purpose:** trace the production path and propose the boundary change.
- **Two independent breaks:** frontend bindings absent (`ProcessMessageForEvents`, `ApproveCandidate` missing
  from `App.d.ts`, with a localStorage fallback that stores a *candidate from the user's own outgoing
  message*), and the Go append path has no stable source ID and never calls the memory service.
- **Source-ID defect (Observed):** `ID: fmt.Sprintf("%d", len(a.chatMessages)-1)` — positional, not stable
  across edits, deletions, reloads or snapshots (`PCB-05`, `PCB-06`).
- **Status:** proposal; **production change not made**; G2 **blocked**.

## 12. `HEKI-D-0029` — Formation and Classification relationship

- **Purpose:** investigate native classification (D), deferred formation (E), and the hybrid.
- **Owner determination recorded:** `CFD-01` — **owner rejects Option A (author-supplied class)**; Option A is
  **Rejected**, not merely unadopted.
- **Finding:** D classifies only explicit cases; E needs an L-2/Decision A amendment; hybrid is the only
  combination meeting both autonomy and no-fabrication — **recommendation, not adopted**.

## 13. `HEKI-D-0030` — Autonomous Classification Experiment Plan

- **Purpose:** fix the G0–G5 gate sequence before amending L-2.
- **State:** G0 accepted; **G1 fixture contract passed**; production wiring not established; **G3 not
  authorized** (no BGE fetch); corpus for class labels: none; L-2 and Decision A unchanged.
- **Separations that matter later:** capture = what happened; classification = which category; formation =
  how it becomes authoritative memory (`ACX-02`).

---

## 14. Findings

**T2-F1 — A live contradiction between three documents about FC-5 and amendment R6.**
- `docs/review/AMENDMENT-DRAFT-…-R6-2026-11-21.md` says: **"WORKING COPY — NOT APPLIED"**, and states the
  canonical paper still lacks every clause and "its registered hash still describes the unamended text."
- `docs/architecture/contracts/README.md` says: FC-5 "carries an owner-approved additive amendment (revision
  6, 2026-11-21) **applied** after that acceptance", and lists the added clauses.
- The FC-5 paper itself contains the identifier note and FC5-27–FC5-31 that only R6 introduces.

**Resolved on 2026-11-22 by measurement — not by preferring one statement over another.** FC-5 now measures
**23,056 bytes / 304 lines / `626ab4be…`**, against R6's recorded target of 12,326 bytes / 213 lines /
`c245096f…`. It contains `FC5-26a` twice and `FC5-27`–`FC5-31`. The document registry already records the
amended hash `626ab4be…`. Two independent witnesses agree, so **the amendment has been applied and R6's status
block is stale.**

No document was edited (R1). This is an **erratum candidate for owner action**, and it must not be smoothed
over in prose: a draft that says "NOT APPLIED" about an applied amendment will mislead every later reader who
trusts it.

What remains open is *not* whether R6 was applied, but whether anything followed: FC-5 still reads
`PROPOSED — not approved` with implementation authorization `None`. Applied ≠ approved ≠ implemented — three
separate milestones, and only the first is established.

**T2-F2 — "Review-accepted" is not "approved", and nothing here is implemented.** All six papers are
PROPOSED with zero implementation authorization. Any readiness statement built on the contract pack must say
**design-only**; nothing in it is evidence that code behaves as contracted.

**T2-F3 — The architectural blocker for MEMORY-09 is a deliberate hole.** `FC4-04a`: memory-record authority is
**Unestablished here**; `FC-7` is **blocked by D8**. The contract pack therefore cannot govern the memory
service, and MEMORY-09's production path is blocked by a delegated decision that has no record yet.

**T2-F4 — The G1 chain is connected only through the program direction.** Six research records: five have
`in = 1` (from MEMORY-09 program direction), **the audit has `in = 0`**, and all six have `out = 0` and **zero
program-identifier mentions**. The audit that created the gate cannot be reached from the results that cite it.

**T2-F5 — The identifier layer is a contract-internal mechanism.** 788 mentions across 281 identifiers exist in
only **45 of 808 documents**. It is dense inside the contract pack (FC-1 alone: 111 mentions) and absent from
the research corpus. Mentions therefore beat hyperlinks where clause-level lineage matters and add almost
nothing for raw evidence — the opposite of a general-purpose connectivity fix. It rescues only 15 of the 304
hyperlink-isolated documents (5%).

**T2-F6 — `FC5-04` already separates D6 from T0.5.** The paper states D6 proves neither process ownership nor
global reservation accounting. This is a useful guard against later misreading, and it is the kind of
distinction that would be lost in a summary document.

**T2-F7 — FC-4 assigns Character to the canonical library.** Direct text: `FC4-04` "Character | Canonical
library". This bears on ADR-001 §11.2's "standalone" wording, flagged earlier in the program. Not adjudicated
here — Phase 4 must compare the two texts rather than assume either controls.

---

## 15. Next tranche

**Tranche 03 — 2026-11-20:** the eight foundation-pack artifacts registered as rows 505–512
(`HOWL-FOUNDATION-CONTRACT-PACK-DEEPSEEK-INTAKE`, `HOWL-FOUNDATION-FRONT-TO-BACK-ARCHITECTURE-AND-DELIVERY-PLAN`,
the owner determinations D2/L2 and Decision A, `MEMORY-09-IMPLEMENTATION-PLAN`, and their neighbours), then
2026-11-19 and earlier.
