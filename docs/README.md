# HOWL Chat — Documentation Authority and Navigation

**Status:** Navigation and authority map only. This file creates no decision, contract, freeze, acceptance, or implementation authorization.
**Last reconciled:** 2026-11-12
**Scope:** Every file under `docs/`, plus the requirements ledger at the repository root.

---

## 1. Start here

| If you need… | Read |
|---|---|
| The reconstructed history from lorebook ownership to the current governance state | `docs/GOVERNANCE-CHAIN.md` |
| The memory system's governing principle, four layers, capability dispositions, and acceptance criteria | `docs/architecture/ARCHITECTURE-DECISION-MEMORY-SYSTEM-2026-11-12.md` |
| The status and home of every document, including hashes | `docs/DOCUMENT-REGISTRY.md` |
| Every amendment, its status, and what supersedes it | `docs/AMENDMENT-REGISTER.md` |
| Every frozen artifact and its verification command | `docs/FROZEN-ARTIFACTS.md` |
| The current F3-C contract | `docs/review/F3C-MASTER-ARCHITECTURE-AND-IMPLEMENTATION-CONTRACT-2026-11-11.md` |
| The current F3-C ownership problem | `docs/research/ROOT-01-PRODUCTION-ROOT-OWNERSHIP-AND-CONVERSATION-IDENTITY-2026-11-11.md` |
| Memory systems: prior work, external evidence, small-model limits, failures, and gaps | `docs/research/MEMORY-01-RESEARCH-SYNTHESIS-2026-11-12.md`; sources in `docs/research/MEMORY-02-SOURCE-REGISTER-2026-11-12.md` |
| The memory evaluation protocol (design only; execution not authorized) | `docs/research/MEMORY-03-EVALUATION-PROTOCOL-2026-11-12.md` |
| The pre-registered evaluation package: corpus families, arms, rubrics, gates, and what evidence would justify freezing a choice | `docs/research/MEMORY-09-MEMORY-EVALUATION-PREREGISTRATION-2026-11-12.md` |
| Whether the five repaired corpus items survived independent review, what defects remain, and what the owner must decide next | `docs/research/corpus/CORPUS-14-RECONCILIATION-REGISTER-2026-11-12.md`; verdicts in `docs/research/corpus/CORPUS-13-INDEPENDENT-REVIEWER-REPORT-2026-11-12.md` |
| User personas: owner requirements for a first-class persona object, library/attachment/active-identity split, and the UI-only Role field | `docs/architecture/OWNER-DECISION-USER-PERSONAS-2026-11-12.md` |
| Automated lorebooks: organization/activation/maintenance/scope/deletion mechanics, local-model limits, orchestration, ranking, and pre-specified test fixtures | `docs/research/MEMORY-04-AUTOMATED-LORE-RESEARCH-ASSESSMENT-2026-11-12.md` |
| Why lorebook *generation* is not a separate subsystem, and what replaces it | `docs/architecture/OWNER-DECISION-MEMORY-FORMATION-2026-11-12.md` |
| Whether a context flush can be turned into safe memory maintenance | `docs/research/MEMORY-05-CONTEXT-FLUSH-AND-RECONSTRUCTION-RESEARCH-2026-11-12.md` |
| Whether Howl constructs each model family's prompt correctly | `docs/research/MEMORY-06-PROMPT-FORMAT-COMPATIBILITY-MATRIX-2026-11-12.md` |
| What the hardware can host, and how resource policy must be measured | `docs/research/MEMORY-07-HARDWARE-AND-INFERENCE-OPTIMIZATION-2026-11-12.md` |
| Which approaches rank where, with evidence classes and open questions | `docs/research/MEMORY-08-COMPARATIVE-RANKING-AND-UNRESOLVED-QUESTIONS-2026-11-12.md` |
| The T2.9.5B authorization control | `docs/t2.9.5b-authorization-decision-register-r14.md` and `docs/t2.9.5b-authorization-decision-register-r14.freeze.md` |
| The TTS governing architecture | `docs/review/DESIGN-howl-tts-helper-architecture-2026-11-06.md` |
| The product requirements ledger | `PROJECT_ISSUES_AND_EXPANSION.md` |
| The proposed corrected master plan (not adopted) | `docs/review/CORRECTED-MASTER-PLAN-HOWL-2026-11-11.md` |

The previous review index at `docs/review/README.md` remains the detailed T2.9.5B index. It is **not** a directory-wide authority map; its new header points back here.

---

## 2. Authority order

When documents disagree, use this order. A later item never overrides an earlier one merely because it is newer or more detailed.

1. **Owner decisions and owner-characterized requirements.** The current requirements ledger is `PROJECT_ISSUES_AND_EXPANSION.md`. Owner decision records and resolved decision sections control within their stated scope.
2. **Frozen contracts.** F3-C data/ownership contracts, the F3-C master architecture, and the architecture freeze. T2.9.5B's frozen register and its freeze record. The TTS Stage 0 freeze inside the governing TTS design, as amended by adopted Clarification A2.
3. **Adopted protocols.** `GATE-F3C-IMPLEMENTATION-VERIFICATION-V2-2026-11-11.md` is the adopted F3-C implementation-verification protocol. It does not grant implementation authorization.
4. **Current amendment baselines.** Amendment 01 is a planning amendment, not architecture approval. Amendment 02 revision 9 is a closed requirements baseline, not architecture approval and not implementation authorization.
5. **Current proposals.** The corrected master plan revision 3 and the F3C-TARGET-2 R5 proposal are proposals only. They are not frozen and are not accepted.
6. **Current evidence and research.** F3-C implementation evidence, native verification, root inventory, D-ROOT closure, Gate B, and ROOT-01. These establish what was observed; they do not create authority.
7. **Historical evidence.** Reviews, audits, prior revisions, superseded proposals, and archived records. They remain traceable and must not be treated as current merely because they are present.

**Do not infer approval from a document's existence, length, filename, or hash.** A frozen contract is not an accepted implementation. An implementation report is not acceptance. A PASS review is not an owner decision.

---

## 3. Current state (2026-11-12)

### F3-C

- The architecture is **frozen** at `F3C-MASTER-ARCHITECTURE-AND-IMPLEMENTATION-CONTRACT-2026-11-11.md`; the freeze record is `FREEZE-F3C-ARCHITECTURE-2026-11-11.md`.
- The data and ownership contracts are **frozen** at `FREEZE-F3C-DATA-AND-OWNERSHIP-CONTRACTS-2026-11-10.md`.
- F3-C is **not accepted**. P1a implementation is **not authorized**. F3-D remains blocked.
- `ROOT-01` establishes, from the current dirty tree, that production chat/memory data does not implement `(rootID, chatID)`, that chat IDs can collide, that IDs are unvalidated path components, and that memory/delete/wipe boundaries are incomplete. It is source evidence, not runtime proof.
- Six F3-C defects remain open: `DEF-F3C-ROOT-01`, `WIPE-02`, `IDENTITY-03`, `REACH-04`, `STARTUP-05`, and `ALIAS-06`. The formal disposition request is `DISPOSITION-REQUEST-F3C-DEFECTS-2026-11-11.md`.

### Amendment 02

- Revision 9 is the **closed requirements baseline**; there is no revision 10.
- It is **not architecture approval** and does not authorize research execution.
- Muse's revision-9 `M6-01`–`M6-03` findings remain unresolved evidence obligations, recorded in `AMENDMENT-02-REVISION-CYCLE-CLOSURE-2026-11-11.md`.
- Revisions 6, 7, and 8 are archived; revision 3 is the independently reviewed identity. See `docs/AMENDMENT-REGISTER.md`.

### T2.9.5B / audio

- The controlling decision-control pair is `t2.9.5b-authorization-decision-register-r14.md` plus its checksum, frozen by `t2.9.5b-authorization-decision-register-r14.freeze.md`.
- The r14 register file's own header says **DRAFT — NOT FROZEN**. The freeze record controls that apparent conflict; the register was frozen as the exact r14 artifact.
- The T2.9.5B stage is **NOT FROZEN — NOT AUTHORIZED**. Implementation, measurement, build, runtime, benchmarking, MEAS1 closure, and T6 closure are not granted.
- The TTS Stage 0 protocol is frozen in the governing design and amended by adopted Clarification A2. That is design authority, not implementation authorization.

### Master plan and requirements

- The requirements recovery baseline is `REQUIREMENTS-RECOVERY-OWNER-ADDENDUM-HOWL-2026-11-11.md`.
- The reconciliation is `REQUIREMENTS-RECONCILIATION-HOWL-2026-11-11.md`.
- The complete-product master plan is a planning inventory, not authority.
- The corrected master plan revision 3 is **proposed, not adopted**. Its verification record verifies consistency, not adoption.

### Memory systems research

- `MEMORY-01-RESEARCH-SYNTHESIS-2026-11-12.md` reconciles the project's existing memory/continuity/RAG/automated-lore material and records external evidence on structured, episodic, semantic, temporal, graph, hybrid-retrieval, consolidation, provenance, isolation, deletion, and small/local-model approaches.
- It establishes evidence-supported constraints and failure modes. It does **not** select an architecture, approve a design, authorize implementation, or close an owner decision.
- `MEMORY-02-SOURCE-REGISTER-2026-11-12.md` is the traceable bibliography. No Howl runtime benchmark exists; the report says so explicitly.
- `MEMORY-03-EVALUATION-PROTOCOL-2026-11-12.md` defines the pre-registered conditions, model/hardware tiers, synthetic corpus, measures, thresholds, and authorization gates for a future evaluation. It is **design only**; no model was loaded and no experiment was run.
- The existing proposals remain proposals: `DA-M1`–`DA-M6`, `DA-R1`–`DA-R6`, `R-10`, `R-11`, and `REC-CH-09`–`REC-CH-11` are not decided by this research.

### Memory evaluation corpus (`G1`)

- `CORPUS-01`–`CORPUS-04` hold 45 scored units against an ≥80 floor; authoring is **paused** pending review (`P-5`).
- 20 authorised repairs were applied and mechanically verified (`CORPUS-09`, `CORPUS-10`). A checker `PASS` is **not** corpus validity.
- The five Group B units were independently re-reviewed (`CORPUS-13`) under a corrected packet (`CORPUS-12`): **1 accepted on item validity, 4 require rework, 0 retired.** Every targeted defect was corrected; `F6-04`'s repair is **incomplete** and left the item self-contradictory.
- Reconciliation and the new defect register `RD-1`–`RD-15` are in `CORPUS-14`. **Five owner decisions (`OD-7`–`OD-11`) are required before any item is edited.**

---

## 4. The governance chain in one paragraph

The trail starts with content ownership, not with F3-C. F3-A removed unsafe mirror writes but left lorebooks browser-owned without file identity. F3-B resolved the lorebook questions by making one canonical store the only writer, preserving identities, and requiring explicit, journaled, reversible import. F3-C applied that model to worlds and scenarios: stable identity, folder-per-item, pinned snapshots, no name fallback, and explicit migration. The resulting compatibility findings (especially the chat service as a name-keyed second writer) forced the data/ownership freeze, the root-authority decision, LF-1, single-writer audit, import-boundary audit, implementation evidence, the consolidated master contract, the architecture freeze, and Gate v2. The root-ownership inventory and D-ROOT closure then found that the implemented chat/memory stores do not actually carry `(rootID, chatID)`. Gate B and ROOT-01 carried that evidence into the current defect queue.

Read the full node-by-node chain in `docs/GOVERNANCE-CHAIN.md`.

---

## 5. Document families

| Family | Location | What it is | Current entry point |
|---|---|---|---|
| Master navigation | `docs/README.md` | Authority map and reading order | This file |
| Governance chain | `docs/GOVERNANCE-CHAIN.md` | Reconstructed historical chain | This file |
| Registry | `docs/DOCUMENT-REGISTRY.md` | Every file, hash, role, and current relevance | This file |
| Amendment register | `docs/AMENDMENT-REGISTER.md` | Amendment status overlay | This file |
| Frozen manifest | `docs/FROZEN-ARTIFACTS.md` | Frozen artifacts and verification | This file |
| Architecture | `docs/architecture/` | Integrated plan, designs, Amendment 01/02, verification design, owner persona requirements | `AMENDMENT-02-…-2026-11-11.md` for the closed requirements baseline; `INTEGRATION-AND-ACCEPTANCE-PLAN-2026-11-11.md` for the amended plan; `OWNER-DECISION-USER-PERSONAS-2026-11-12.md` for the user-persona requirement baseline |
| F3-C review/evidence | `docs/review/` | Lorebook/F3-C decisions, freezes, audits, implementation evidence, root defects, master-plan reconciliation | `FREEZE-F3C-ARCHITECTURE-2026-11-11.md` and `F3C-MASTER-…-2026-11-11.md` |
| Current research | `docs/research/` | Gate B, Q-06, R-07, ROOT-01, memory/retrieval evidence, evaluation protocol, automated-lorebook research package, context flush, prompt formats, hardware, comparative ranking | `MEMORY-01-…-2026-11-12.md` for findings; `MEMORY-03-…-2026-11-12.md` for the protocol; `MEMORY-08-…-2026-11-12.md` for the ranking and open-question register; `ROOT-01-…-2026-11-11.md` for identity |
| Audio architecture | `docs/audio-architecture/` | Audio subsystem implementation architecture | `README.md` |
| T2.9.5B authorization | `docs/` top level | Decision-control registers, proposals, briefings, measurements | `t2.9.5b-authorization-decision-register-r14.md` + freeze record |
| T2.9.5B review | `docs/review/` | Audio research, audits, Stage 0/1/2/3/C1/P17 evidence | `docs/review/README.md` (T2.9.5B sections) |
| Product/requirements | `docs/review/` + root ledger | Requirements recovery, reconciliation, master plan, adjudication | `PROJECT_ISSUES_AND_EXPANSION.md`; then recovery baseline |

---

## 6. How to trace a claim

1. Find the artifact in `docs/DOCUMENT-REGISTRY.md`.
2. Read its **role** and **current relevance** columns. Do not stop at its title.
3. Follow its **supersedes / superseded by** links in `docs/AMENDMENT-REGISTER.md` and `docs/GOVERNANCE-CHAIN.md`.
4. Verify its hash against the registry or `docs/FROZEN-ARTIFACTS.md`.
5. If it is frozen, do not edit it. A change requires a new dated owner decision or a new frozen revision with its own record.
6. If it is historical, use it as evidence of what was thought at that time, not as current authority.

---

## 7. Maintenance rules

- **Frozen artifacts are immutable.** Do not correct them in place. Add a dated decision/erratum and a new registry row.
- **Do not silently update a hash.** Every hash change needs a reason, a new artifact identity, and a registry entry.
- **Do not delete historical files.** Mark them archived/superseded and preserve their bytes.
- **Do not create a second current authority.** A new current document must name what it supersedes and link its predecessor.
- **Keep the chain bidirectional.** Every current document should be reachable from `docs/README.md`; every historical document should be reachable from the registry and its successor.
- **Update the registry when adding a file.** The registry is the completeness check, not a substitute for reading the file.

---

## 8. Known gaps deliberately not closed here

These are recorded so the documentation chain does not pretend to cover them:

- **Howl-specific memory-system evaluation** has not been performed. `MEMORY-01` records external evidence and explicit gaps; it does not produce a Howl benchmark, select an architecture, or close the identity/deletion defects. `MEMORY-03` designs the evaluation but authorizes no execution; `MEMORY-08` ranks approaches but selects none.
- **No memory, lore, flush, prompt-format, or hardware measurement exists.** `MEMORY-04`–`MEMORY-07` are research and compatibility assessments; every Howl-specific number in them is marked unmeasured.
- **No architecture has been chosen, and neither a framework nor a custom build is preferred.** The evidence does not yet justify either.
- **The memory architecture direction is adopted; its implementation is not.** The owner architecture decision of 2026-11-12 establishes one application-owned, AI-assisted, hybrid memory system with original evidence as authoritative, one formation pipeline, knowledge collections, and identity-aware retrieval. Schema, retrieval algorithm, chunking, reasoning budget, store, and model selection are deliberately unfrozen and must earn their place through testing.
- **Automated lore generation is not a planned subsystem.** The owner decision of 2026-11-12 collapses it into a single memory-formation pipeline with multiple destinations; lorebooks remain as an organization layer. The research is retained; the duplicate extractor is not.
- Runtime evidence for packaged root resolution, AppData redirection, filesystem aliases, Save As collision, cross-chat memory display, and ID reuse is not established.
- D-ROOT has no owner-selected canonical root. OD-18 remains unresolved.
- Candidate/continuity production reachability remains unestablished, not disproved.
- Amendment 02 `M6-01`–`M6-03` remain evidence obligations.
- T2.9.5B stage authorization remains not granted.
- The corrected master plan is not adopted.
- Some historical headers are stale (notably the r14 register's pre-freeze `DRAFT` header and the TTS §17.19 draft's `NOT ADOPTED` header). They are resolved administratively in `docs/AMENDMENT-REGISTER.md` and `docs/FROZEN-ARTIFACTS.md`; the underlying bytes are preserved.

---

## 9. Bottom line

There is now one entry point, one exhaustive registry, one amendment status overlay, one freeze manifest, and one reconstructed chain. No document is abandoned: every file has a row, a family, a role, a hash, and a current-relevance judgment. Historical files remain historical. Frozen files remain byte-identical. The current authority is still the owner decisions and frozen contracts, not this navigation layer.
