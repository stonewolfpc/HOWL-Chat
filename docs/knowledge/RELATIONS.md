# HEKI — RELATIONS (verified)

**Status:** verified edges only. Every row was established by reading the source document, not by string
extraction. Machine-extracted hyperlinks live in `ledger.jsonl` as `links_out`; they are **citations**, not
relationships, and are not promoted without a reason.

**Vocabulary:** see `README.md` §5. **Adoption is a field, not an edge** (rule R4).

---

## 1. The 2026-11-21 program cluster

Identity numbers are from `DISCOVER-LEDGER.md`.

| From | Type | To | Reason (quoted evidence) | Adoption |
|---|---|---|---|---|
| `HEKI-D-0019` MEMORY-09 program direction | `cites` | 25 records in its recovery index (§2.1–§2.3) | "the records to read before changing direction" | — |
| `HEKI-D-0022` HFS-AUDIT (original scope) | `derived_from` | `ISSUES.md` (V1) and `PROJECT_ISSUES_AND_EXPANSION.md` (V2) | "V2 is the later recovered list, but it silently narrows V1" — union 25 + 26 → **27** | — |
| `HEKI-D-0022` HFS-AUDIT | `contradicts` | V2 (`PROJECT_ISSUES_AND_EXPANSION.md`) **by omission** | "V2 drops T0.5 … without a recorded closure or rationale … The omission is therefore not evidence that T0.5 was fixed or deselected." | Unresolved — T0.5 preserved, not closed |
| `HEKI-D-0020` HFR-CP readiness audit | `validates` | A–E readiness of the Howl foundation | Establishes E0–E3 per subsystem; "no subsystem in A–E has end-to-end verification" | Audit complete |
| `HEKI-D-0020` HFR-CP | `contradicts` | any claim that Howl foundation is ready | `HFR-01`: "not ready for real-workload acceptance"; `HFR-04` forbids adjacency promotion | Verdict stands |
| `HEKI-D-0021` HF1-PLAN (Phase 1) | `derived_from` | `HEKI-D-0020` HFR-CP | "Readiness baseline `HFR-CP-2026-11-21-01` — remains the current baseline; not overwritten" | **Proposed only** |
| `HEKI-D-0021` HF1-PLAN | `governs` | nothing yet | "Separate implementation authorization required" | **Not authorized** |
| `HEKI-D-0018` HRA-DB (arbitration baseline) | `supersedes` | modality statements in `HEKI-D-0023` HCDM | `HRA-04`: "supersedes older statements in `HCDM-2026-11-21-01` that described image generation as deferred and audio generation as removed" | Owner-approved |
| `HEKI-D-0023` HCDM | `superseded_by` | `HEKI-D-0018` HRA-DB — **for modality scope only** | HCDM §0 scope-correction note reaches the same conclusion | Superseded in part |
| `HEKI-D-0023` HCDM | `derived_from` | `HEKI-D-0022` HFS-AUDIT | "Starting baseline `HFS-AUDIT-2026-11-21-01`"; `HCDM-02` names it the scope-recovery baseline | — |
| `HEKI-D-0023` HCDM | `depends_on` | `HEKI-D-0020` HFR-CP | "`HFR-CP-2026-11-21-01` remains the evidence snapshot. Its A–E categories are not the work list." | — |
| `HEKI-D-0023` HCDM | `reviews` | `HEKI-D-0021` HF1-PLAN | "`HF1-PLAN-2026-11-21-01` remains a **proposed** work package. This map reviews it; it does not authorize it." | **Not authorized** |
| `HEKI-D-0023` HCDM | `records` | corrected owner scope (2026-11-21) | `HCDM-04` table: audio capture in scope; audio generation complete; TTS deferred-required; image nearly complete; vision gaps; video in scope; 3D mesh and visual-character rendering removed | Owner-approved |
| `HEKI-D-0024` M0A review | `depends_on` | `HEKI-D-0023` HCDM | "Program map `HCDM-2026-11-21-01`" | — |
| `HEKI-D-0024` M0A review | `precedes` | `HEKI-D-0021` HF1-PLAN | `M0A-AR-03`: "Phase 1 will change chat state … If the send/render path is still untrusted, Phase 1 can 'fix' state while leaving a layer that hides or duplicates the user's message." | Recommendation only |
| `HEKI-D-0024` M0A review | `governs` | M0A gate (T0.1, T0.2, T0.5) | "This review defines a gate; it does not authorize implementation" | **No authorization** |
| `HEKI-D-0024` M0A review | `contradicts` | naive duplicate-guard reasoning | `M0A-AR-09/10`: guard is correct sequentially, not against the unawaited UI path; a synchronous test may pass while the UI still duplicates | Open risk |

---

## 2. Decision records (2026-11-22)

| From | Type | To | Reason | Adoption |
|---|---|---|---|---|
| `HEKI-D-0001` ADR-002 | `refines` | `HEKI-D-0002` ADR-001, decision **D1** only | ADR-002 header: refines exactly one ADR-001 decision and settles nothing else there | Adopted (architecture) |
| `HEKI-D-0002` ADR-001 | `records` | D1–D5 | two-population chat root, portable-mode deferral, retention N=2 provisional, lock contention → `Busy`, legacy stores read-only indefinitely | Adopted; **no implementation** |
| `HEKI-D-0002` ADR-001 | `depends_on` | `HEKI-D-0019` program direction + `HEKI-D-0023` HCDM | ADR-001 §5.1 defines L2 because no prior repo document defined it; discoverability links point at both | — |
| `HEKI-D-0019`, `HEKI-D-0023`, `HCDM-36a` | `cites` | `HEKI-D-0001`, `HEKI-D-0002` | Registered as decision records, deliberately **not** as contracts C13/C14 | Adopted |

---

## 3. FC-5 amendment draft (R6)

| From | Type | To | Reason | Adoption |
|---|---|---|---|---|
| `HEKI-D-0004` R6 draft | `proposes` | `HEKI-D-0013`…`HEKI-D-0016` FC-5 paper | "OWNER-APPROVED AMENDMENT DRAFT — REVISION 6 — WORKING COPY — **NOT APPLIED**" | **Text approved; applied to the paper — confirmed at provision level 2026-11-22: 48 of 48 quoted lines present verbatim, diff `+91 / −0`, order-exact. Not approved as a contract** |
| `HEKI-D-0004` R6 draft | `supersedes` | revisions 1–5 of the same amendment | Self-recorded revision table; "Revision 6 changes FC5-29 and nothing else" | Revision 5 superseded in FC5-29 only |
| `HEKI-D-0004` R6 draft | `contradicts` | FC-5's registered hash as a description of current intent | "The canonical FC-5 paper still lacks every clause below, and its registered hash still describes the unamended text." | Deliberate; not a defect |

**Four unsatisfied milestones, recorded so they cannot be mistaken for done:** application to the canonical
paper; register + document-registry reconciliation; any M4 change; implementation.

---

## 4. Explicit non-relations (recorded to prevent false inference)

| Non-relation | Evidence | Why it matters |
|---|---|---|
| HF1-PLAN does **not** authorize implementation | "No implementation … is authorized by this record" | A proposed work package must never be read as an authorization |
| HCDM does **not** create requirements | `HCDM-03`: "This map itself creates no new product requirement" | Prevents scope expansion by citation |
| R6's §1 "does not amend FC-5" is **not** the paper's state | R6 §1 and the paper disagree: FC-5 contains all eight additions (48/48 lines verbatim, `+91 / −0`); the baseline hash `c245096f…` is retrievable at `0bc3b091…`, current is `626ab4be…` | The boundary sentence describes the draft's authority, not the file. Recorded so a stale claim cannot govern a document that already contains the text |
| Applying an amendment is **not** approving the paper | FC-5 still reads `PROPOSED — reservation questions unresolved`; `Implementation authorization: **None**`; FC5-27–FC5-31 keep their own labels | Application, approval, adoption, M4 and implementation are five separate states |
| `HRA-02` and `HRA-09` are **Proposed**, not adopted | HRA labels | Adoption field, not an edge |
| `FC5-27`–`FC5-31` in R6 are **Proposed** | R6 §4 clause labels | Approved text ≠ adopted clause |
| M0A review does **not** authorize its own criteria | `M0A-AR-02` | Gate definition ≠ gate passage |
| A mention is **not** a relationship | `README.md` §9 | 788 mentions across 45 documents establish naming, not dependency or adoption |
| Contract review acceptance is **not** approval | All six papers: "PROPOSED — not approved" | 6 papers, 0 approvals, 0 implementation authorizations |

---

## 5. Contract pack (2026-11-21)

| From | Type | To | Reason (quoted evidence) | Adoption |
|---|---|---|---|---|
| `HEKI-D-0011` contracts register | `governs` | `HEKI-D-0012`…`HEKI-D-0017` FC-1…FC-6 | "Contract-paper index and dependency register"; §1 papers table assigns each paper its gate and dependency | Review-complete; **none approved** |
| `HEKI-D-0011` contracts register | `records` | FC-7 absence | "FC-7 is deliberately absent … records the block rather than manufacturing a paper" | Blocked by D8 |
| `HEKI-D-0012` FC-1 | `depends_on` | `HEKI-D-0015` FC-4 | Register §2.4: "FC-4 defines canonical authored-object authority and revision binding used by FC-1 and FC-2" | Proposed |
| `HEKI-D-0013` FC-2 | `depends_on` | FC-1, FC-4 | Register §2.2: composes "from FC-1 state and FC-4 canonical references" | Proposed |
| `HEKI-D-0014` FC-3 | `depends_on` | FC-1, FC-2 | Register §2: "FC-3 depends on both state and prompt composition" | Proposed |
| `HEKI-D-0016` FC-5 | `depends_on` | `HEKI-D-0018` HRA-DB | "Authority: D6 owner decision; `HRA-DB-2026-11-21-01`" | Proposed |
| `HEKI-D-0017` FC-6 | `governs` | acceptance surfaces of FC-1…FC-5 | "All gates"; `FC6-03`: "it does not replace each paper's specific criteria" | Proposed |
| `HEKI-D-0012` FC-1 | `refines` | D2 blanket-retention shorthand | Register §3 D2: "the earlier blanket shorthand is narrowed by FC-1's directed cleanup" | Owner-approved narrowing |
| `HEKI-D-0016` FC-5 | `contradicts` | any reading of D6 as process ownership or global accounting | `FC5-04`: "D6 is not T0.5"; `FC5-05`: "A per-instance ledger is not a global guarantee" | Calculated |
| `HEKI-D-0015` FC-4 | `contradicts` | any claim that memory-record authority is settled | `FC4-04a`: authority **Unestablished here**; FC-7 blocked by D8 | Unresolved |
| `HEKI-D-0004` R6 draft | `contradicts` | `HEKI-D-0011` contracts register **and** `HEKI-D-0016` FC-5 | R6 says "WORKING COPY — NOT APPLIED … its registered hash still describes the unamended text"; the register says the amendment was "**applied** after that acceptance" | **Resolved by measurement 2026-11-22: applied.** FC-5 = 23,056 B / 304 L / `626ab4be…` vs R6's recorded target 12,326 B / 213 L / `c245096f…`; `FC5-26a` and `FC5-27`–`FC5-31` are present; the registry already records the amended hash. **R6's status block is stale** — erratum candidate, R6 not edited (R1). Application ≠ approval: FC-5 still reads `PROPOSED — not approved` |

## 6. MEMORY-09 G1 chain (2026-11-21)

| From | Type | To | Reason | Adoption |
|---|---|---|---|---|
| `HEKI-D-0025` capture-completeness audit | `derived_from` | the live capture path (`ProcessMessageForEvents`, append path) | Observed inventory of what the service actually receives | Audit only |
| `HEKI-D-0026` G1 fixture contract | `derived_from` | `HEKI-D-0025` audit | "It establishes the minimum capture envelope and the gaps that the later fixture tested" | Contract |
| `HEKI-D-0027` G1 results | `validates` | `HEKI-D-0026` G1 contract | `TestG1` PASS, 11 subtests, `ok … 0.804s`; `G1R-03` "no production type, function, signature, storage path or binding was changed" | **E2, fixture path only** |
| `HEKI-D-0028` production boundary proposal | `derived_from` | `HEKI-D-0027` G1 results | "The fixture contract passed; production capture is still unestablished" | Proposal; production change **not made** |
| `HEKI-D-0028` production boundary proposal | `contradicts` | any claim that capture works in production | `PCB-01`–`PCB-07`: missing bindings, localStorage fallback, no memory-service call, positional source IDs | G2 **blocked** |
| `HEKI-D-0029` formation/classification | `depends_on` | `HEKI-D-0030` experiment plan and owner direction | `CFD-01`: owner **rejects** Option A (author-supplied class) | **Option A Rejected**; hybrid recommended, **not adopted** |
| `HEKI-D-0030` experiment plan | `governs` | gates G0–G5 | "The gates G0–G5 are fixed"; G0 accepted, G3 **not authorized** | Plan; no implementation |
| `HEKI-D-0025`…`HEKI-D-0030` | `cites` | `HEKI-D-0019` program direction | Each has in-degree 1 from that record's recovery index; **the audit also has in-degree 0** | — |

## 7. 2026-11-20 foundation pack — the lineage of the contract pack

Read in `PHASE-1-TRANCHE-03.md`.

| From | Type | To | Reason (quoted evidence) | Adoption |
|---|---|---|---|---|
| `HEKI-D-0089` intake | `derived_from` | `HEKI-D-0088` plan | "**Parent plan:** `HOWL-FOUNDATION-FRONT-TO-BACK-ARCHITECTURE-AND-DELIVERY-PLAN-2026-11-20.md` (`HF-ARCH-PLAN-2026-11-20-01`)" — the only place in 808 documents that names the plan | Planning only |
| `HEKI-D-0011`…`HEKI-D-0017` FC-1…FC-7 | `derived_from` | `HEKI-D-0089` intake | Intake §2 assigns each paper its output, its exclusions and its gate; §10.2 is the commissioning instruction | **Proposed drafts; 0 approved** |
| `HEKI-D-0088` plan | `precedes` | every contract paper | "Do not write implementation plans before unresolved decision rows are answered" | Ordering only |
| `HEKI-D-0012`…`HEKI-D-0017` FC-1…FC-6 | `derived_from` | `HEKI-D-0088` §7.1 evidence ladder | Ladder identical in substance: requested → admitted → attempted → observed → verified, or unknown/failed | Inherited, not invented by the pack |
| `HEKI-D-0088` plan | `contradicts` | any reading of the corpus as a filename chronology | §1.1 names six **2026-11-21** documents as controlling inputs; §12 claims to have inspected them, under a 2026-11-20 filename | **Open provenance finding** (`ERRATA-CANDIDATES.md` E-02) |
| `HEKI-D-0088` plan | `cites` | 10 controlling inputs (HFS-AUDIT, HFR, HF1, HCDM, M0A, HRA, HSL, AMENDMENT-01/02, ISSUES pair) | §1.1 table, in backticks; **no hyperlink, no identifier mention** | — |

**Isolation of the parent plan, measured:** inbound hyperlinks 0, outbound 0, identifier mentions 0. It
commissions the entire contract pack and is reachable only through one unlinked string in its child. Different
failure from T2-F4: there the evidence was orphaned downstream, here the **authority** is orphaned upstream.

**Not an edge, recorded to prevent it becoming one:** the plan's 23 `DQ-C*`/`DQ-X*`/`DQ-R*`/`DQ-M*` identifiers
each appear in exactly one document — the plan. Their apparent successors (`DQ-HRA-01`–`12`) are a re-minting in
`HRA-DB`, not a rename stated anywhere. A question with no later citation is **not** a closed question.

## 8. Coverage of this file

**Verified so far:** 22 documents — the 2026-11-21 program cluster, the 2026-11-22 decision records, the
contract pack, the MEMORY-09 G1 chain, and the 2026-11-20 foundation pack. **Not yet verified:** 786 of 808 nodes.

Counts are deliberately not self-reported here; `CHECKPOINTS.md` measures them at issue time. Nothing in this
file is inferred from filename similarity or from mention frequency, and no edge exists that a reader cannot
re-derive by opening the two documents named.
