# Phase 1 — Discover · Tranche 01 (newest first, 2026-11-21 → 11-22)

**Checkpoint:** `HEKI-CP-2026-11-22-01`
**Method:** read each document's front matter, status block, scope and relationship statements in full context;
record only what the document states about itself. Nothing was edited.

**Coverage:** 7 of 808 documents verified (0.9%). This tranche is deliberately small. It covers the documents
the entire current program hangs from, so that later tranches can be read against a correct skeleton.

---

## 1. `HEKI-D-0001` — ADR-002 Storage and Navigation Policy (2026-11-22)

- **Purpose:** decide where content lives (three state families), how the picker navigates, and how content
  is routed. Supplies N1–N19, criteria 1–11.
- **Subsystem:** storage, navigation, content routing.
- **Provenance:** authored 2026-11-22 as a decision record. Header states: refines exactly one ADR-001
  decision (D1) and settles nothing else there.
- **Relations:** `refines` ADR-001/D1; cited by program direction, HCDM §F, and the registry.
- **Adoption:** architecture decision; **no implementation implied.**
- **Never-linked check:** inbound = 3, outbound links = 0. It links to nothing, because it was written to
  state decisions rather than to index prior work — a navigation gap that Phase 2 must close from the
  side of the *index*, not by editing the ADR.

## 2. `HEKI-D-0002` — ADR-001 Foundational Backend Rebuild (2026-11-22)

- **Purpose:** five roots; six services + two primitives; contracts A–F; authority record as sole router;
  containment fail-closed; DurableStore; runtime evidence; Class A/B/C writes; decisions D1–D5.
- **Provenance:** authored 2026-11-22. §5.1 defines "L2" with an explicit provenance note, because no prior
  repository document defined the term.
- **Adoption:** adopted architecture; **no implementation.**
- **Never-linked check:** inbound = 3, outbound = 0. Same structural note as ADR-002.

## 3. `HEKI-D-0004` — R6 FC-5 amendment draft (2026-11-21)

- **Purpose:** persist owner-approved amendment text for FC-5 outside conversation history.
- **Status (verbatim):** "OWNER-APPROVED AMENDMENT DRAFT — REVISION 6 — WORKING COPY — **NOT APPLIED**".
- **Provenance:** revisions 1–5 existed only as conversation text; revision 6 text is the approved wording.
  Revision 6 changes FC5-29 and nothing else.
- **Relations:** `proposes` FC-5; `supersedes` revisions 1–5.
- **Critical reading — revised 2026-11-22, see Tranche 02 T2-F1.** The draft's claim that the target hash is
  unchanged is **stale**: FC-5 measures 23,056 B / 304 L and contains the R6-only clauses, and the registry
  records the amended hash. The amendment *was* applied; application and registry reconciliation are complete by
  measurement. Not established: M4, and implementation. The draft still cannot govern — an amendment record is
  not an authority — but it no longer describes the paper accurately.
- **Never-linked check:** inbound = 0, outbound = 0 — completely isolated. It names its targets as text, not
  links. This is the clearest single example of "important document that nothing links to."

## 4. `HEKI-D-0019` — MEMORY-09 Program Direction (2026-11-21)

- **Purpose:** program direction, recovery index and checkpoint `CP-2026-11-21-01` (token
  `HPD-CP-2026-11-21-01`). Declares itself the program entry point (`HPD-03`).
- **Provenance:** 2026-11-21; canon r33/v9/owner ledger/conformance map unchanged; map = 118 requirements.
- **Relations:** cites 25 records across §2.1–§2.3 (direction/state, MEMORY-09 canon and decisions, Howl
  foundation). References ADR-001/ADR-002 via §2.3.
- **Adoption:** direction and checkpoint only — "Not architecture approval, implementation authorization,
  a completion claim, or a MEMORY-09 rule change."
- **Never-linked check:** inbound = 3, outbound = 25. The densest node in the corpus and the correct entry
  point. Its outbound citations are the best available map of the MEMORY-09 corpus.

## 5. `HEKI-D-0020` — HFR-CP Consolidated Readiness Audit (2026-11-21)

- **Purpose:** establish per-subsystem readiness A–E with an evidence standard (E0–E3) and a tested-commands
  table; decision register for MEMORY-09.
- **Verdict (Observed):** "Howl has significant working code, but no subsystem in A–E has end-to-end
  verification, and the combined Howl + MEMORY-09 system is not ready for real-workload acceptance."
- **Notable executed results:** Go suites PASS across chat/bench/modelintel/llama/gguf with explicit
  "what it does not prove" columns; `go test ./internal/resource` → **no test files**; frontend suite
  **15 pass, 1 fail** (`DEF-FE-01` territory).
- **Provenance:** 2026-11-21, current tree at that date. No model or GPU test.
- **Adoption:** audit only; authorizes nothing.
- **Why it matters to HEKI:** it defines the E0–E3 scale this index adopted rather than inventing another.

## 6. `HEKI-D-0021` — HF1-PLAN Phase 1 work package (2026-11-21)

- **Purpose:** define the chat-state/context-integrity package: stable conversation identity, branches,
  composition, real token accounting, flushing, model-change preservation, long/multi-character tests.
- **Carry-forward findings:** context budget not implemented (`TrimToMaxTokens` placeholder,
  `Message.Tokens` always 0); multi-character contamination untested; hardware management advisory only
  (`RequestAllocation` returns advice, no lease enforces it).
- **Adoption:** **proposed**; "Separate implementation authorization required."
- **Relations:** `derived_from` HFR-CP; reviewed (not authorized) by HCDM; preceded by M0A.
- **Never-linked check:** inbound = 1, outbound = 2.

## 7. `HEKI-D-0022` — HFS-AUDIT Original Scope and Evidence Chain (2026-11-21)

- **Purpose:** recover the original owner scope across two list generations and reconcile evidence.
- **Result (Calculated):** V1 = 25, V2 = 26, **union = 27** active items; Tier 4 holds 2 parked ideas;
  E.1–E.8 parked concepts; F2/F3/W are engineering ledgers outside the 27.
- **Material contradiction recorded:** V2 drops **T0.5** without a recorded closure; V1 has it at
  `ISSUES.md:69–86`; `docs/asr-repair-design.md:10` still calls T2.2 and T0.5 open. The audit preserves T0.5.
- **Also material:** `HFS-06` — "'foundation' must not be narrowed to A–E"; the readiness categories are an
  evidence organization, not the original list.
- **Adoption:** audit record, not a requirements baseline.
- **Never-linked check:** inbound = 0, outbound = 0 despite being the baseline for HCDM. Isolated.

---

## 8. Findings from this tranche

**F1 — The newest cluster is a citation island, not a network.** ADR-001, ADR-002, R6, HFS-AUDIT and HRA-DB
have **zero outbound doc links**; five of the seven have zero or near-zero inbound. The 2026-11-21 cluster
coordinates itself by **identifier strings** (`HCDM-2026-11-21-01`, `HFR-CP-2026-11-21-01`) rather than by
links. Identifiers survive renames; this is a deliberate house style, and it is why the census found 304
fully isolated documents. HEKI must therefore index **identifier mentions** as well as hyperlinks, or it will
report a corpus that is far less connected than it actually is.

**F2 — Scope authority is split and the split is honest.** `HFS-06` (A–E ≠ foundation) and `HFS-03` (V2 ≠
complete) both prevent the newest audit from being read as the final word on scope. Any readiness view built
later must union the 27 original items with the A–E evidence categories, never substitute one for the other.

**F3 — One genuine open contradiction between documents.**
T0.5 is preserved by HFS-AUDIT because V2 dropped it silently; `docs/asr-repair-design.md:10` says it remains
open. Recorded as `contradicts` in `RELATIONS.md` §1. Not reconciled here — reconciliation is Phase 4.

**F4 — "Complete" appears in the corpus with two different meanings.**
Owner scope statements ("audio generation complete", `HCDM-04`/`HRA-04`) are **not** E3 acceptance
(`HCDM-05`, `HRA-07` say so explicitly). Every readiness statement in HEKI must carry which meaning it uses.

**F5 — A defect is deliberately kept outside the work package.** `DEF-FE-01` is named in HF1-PLAN and
explicitly not absorbed. It must not be lost between packages.

---

## 9. Next tranche

**Tranche 02 — remaining 2026-11-21 nodes:** `contracts/README.md`, FC-1…FC-6, and the ten 2026-11-21
`research/MEMORY-09-*` records (capture boundary, G1 fixture contract/results, capture completeness,
formation/classification, autonomous classification plan).

**Then:** 2026-11-20 (ADR inputs, L2 and Decision A owner determinations, MEMORY-09 implementation plan),
working back toward the 2026-11-11 amendment package and the original T0–T3 ledger.

**Before Tranche 02:** add identifier-mention edges to the census (finding F1), so the graph reflects how this
corpus actually cross-references itself.
