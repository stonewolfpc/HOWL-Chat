# HOWL Engineering Knowledge Index (HEKI)

**Program:** documentation reconstruction — Discover → Connect → Organize → Reconcile → Implement.
**Status:** **Supporting workstream.** Phase 1 (Discover) is paused at `HEKI-CP-2026-11-22-03`; the primary workstream is implementation under `HFI-EP-2026-11-22-01`. Charter and rules remain in force from this file.
**Opened:** 2026-11-22 · **Checkpoint:** `HEKI-CP-2026-11-22-01`

---

## 0. What this is, and what it is not

**It is** a navigation structure over the Howl engineering corpus: stable identities, typed relationships
with stated reasons, group indexes, and an implementation-readiness view derived from evidence.

**It is not** a summary document, an authority, a decision record, or a replacement for anything it indexes.
A link in HEKI establishes **discoverability, not authority**. Where HEKI and an original disagree, the
original controls and HEKI is wrong.

**Naming collision, resolved.** `HEKI-D-NNNN` (this index, documents) is distinct from `HCDM-` (delivery map)
and `CP-`/`HPD-CP-` (program checkpoints). No HEKI identifier may be used to cite a decision.

---

## 1. Non-negotiable rules

| # | Rule | How it is enforced here |
|---|---|---|
| R1 | **Original documents are immutable during reconstruction.** | Nothing outside `docs/knowledge/` is written by this program. Corrections and additions go into HEKI files as additive amendments. |
| R2 | **Every document gets a stable identity.** | `HEKI-D-NNNN`, assigned newest-first, append-only. Never renumbered, never reused. |
| R3 | **Every relationship carries a reason.** | Every edge in `RELATIONS.md` names its type and quotes the evidence that establishes it. An unnamed hyperlink is recorded as a *citation*, not a relationship. |
| R4 | **Never infer adoption from a proposal or draft.** | Adoption status is a separate field. `Proposed`, `Owner-approved`, `Not applied`, `Unresolved` are carried verbatim, never rounded up. |
| R5 | **Never discard an experiment because a later one performed better.** | Superseded work is recorded as `superseded_by`, retained, and still indexed. |
| R6 | **Checkpoints must survive many conversations.** | Append-only `CHECKPOINTS.md` with a searchable resume token; no checkpoint is edited after issue. |
| R7 | **Do not begin implementation on an incomplete reconstruction.** — **superseded as a gate on 2026-11-22** (`HFI-EP-2026-11-22-01` §1). | Retained as a principle, not a gate: a task states which evidence it relies on and may proceed once *that* evidence exists. Readiness output still states its own coverage; gaps are still declared, never filled by assumption. |
| R8 | **Extracted ≠ verified.** | Every fact is labelled: machine-extracted from the ledger, or read and judged in a verified pass. The two are never mixed silently. |

**Consequence of R1:** this index can become stale. That is the correct failure mode. Staleness in an index
is visible and repairable; a rewritten original is not.

---

## 2. Stages

| Stage | Name | State |
|---|---|---|
| 1 | **Discover** — newest to oldest; identity, provenance, references | **Paused 2026-11-22** at `HEKI-CP-2026-11-22-03` — 22 of 808 verified, census complete (808 nodes). Tranche 04 deferred, not cancelled. Discovery is now on demand: it answers implementation's questions instead of preceding them. |
| 2 | **Connect** — add missing links without overwriting originals | Not started. Missing links are *recorded* in HEKI; insertion into originals needs separate authorization. |
| 3 | **Organize** — subsystem, experiment, artifact, verification indexes | Not started. Provisional multi-label groups exist in the ledger only. |
| 4 | **Reconcile** — stitch requirements, decisions, evidence, implementation | Not started. |
| 5 | **Implement** — traceable, evidence-backed roadmap | Out of scope; requires a completed reconstruction (R7). |

---

## 3. Evidence and verification levels

HEKI **inherits** the four readiness levels already defined in the corpus at
`HOWL-FOUNDATION-READINESS-AND-MEMORY-09-DECISION-REGISTER-2026-11-21.md` §2 (`HFR-03`), rather than
inventing a parallel scale:

| Level | Meaning |
|---|---|
| **E0** | Not established — no implementation, or implementation contradicts the claim |
| **E1** | Implemented — code exists on a live path, not exercised |
| **E2** | Working in a limited test — a narrow automated or controlled check passes |
| **E3** | End-to-end verified — a real user-relevant scenario passes across boundaries |

`HFR-04` also governs here: **no level promotion by adjacency.** A component at E2 does not raise its
subsystem. HEKI adds no fifth level and may not raise a level on its own authority.

Document-level labels are carried from the originals: **Observed · Source-derived · Calculated · Inferred ·
Unestablished · Proposed · Owner-approved · Unresolved · Superseded · Rejected**.

---

## 4. Identity scheme

`HEKI-D-NNNN` — 1-based, newest-first at first assignment, append-only thereafter.

- The identity is the **record**, not the file. A path rename does not change the ID; the ledger records both.
- Each node records the `sha256` observed at census. That hash is a **point-in-time observation**, not a
  freeze: an edited original keeps its ID and gets a new observed hash in the next pass.
- Later documents take the next free number even if they are older, because order of discovery is the only
  honest sequence available.

---

## 5. Relationship vocabulary

Typed, directional, always with a reason. Inverse forms are implied and may be stated.

| Type | Means | Typical evidence |
|---|---|---|
| `derived_from` | B was built out of A | "Starting baseline …", "sources and authority" |
| `depends_on` | B cannot be read correctly without A | declared baseline / prerequisite |
| `supersedes` | B replaces A for current scope | "supersedes older statements in …" |
| `superseded_by` | inverse of the above; A is retained, not deleted | |
| `implements` | B is the realisation of A | requires implementation evidence, not intent |
| `validates` | B is evidence that A holds | test result, audit, acceptance |
| `contradicts` | B and A cannot both be true as written | must be recorded, never silently reconciled |
| `supports` | B is consistent with and strengthens A, without proving it | |
| `cites` | untyped reference; B mentions A | weakest edge; never upgraded without a reason |
| `governs` | A sets a rule that binds B | contract → implementation |
| `precedes` | ordering constraint between work items | "M0A comes before Phase 1" |
| `proposes` | B recommends a change to A that is **not applied** | draft, amendment draft |
| `reviews` | B assesses A without authority over it | audit, acceptance review |
| `refines` | B narrows a decision in A without replacing it | ADR-002 refines D1 |
| `records` | B is the persistence of a decision made elsewhere | owner decision → file |

**Adoption is never a relationship type.** It is a separate field, because R4 forbids encoding it as one.

---

## 6. Files

| File | Role | Generated? |
|---|---|---|
| `README.md` | this charter | no |
| `ledger.jsonl` | full census, one JSON record per document | **generated** |
| `DISCOVER-LEDGER.md` | human/greppable census table | **generated** |
| `mentions.jsonl` | identifier-mention occurrences (document, identifier, line, source text) | **generated** |
| `filerefs.jsonl` | filename-reference occurrences (document, target, line, source text) | **generated** |
| `RELATIONS.md` | typed edges with reasons | no — verified only |
| `PHASE-1-*.md` | verified tranche records | no — verified only |
| `ERRATA-CANDIDATES.md` | statements a later measurement showed stale; originals never edited | no — measured only |
| `CHECKPOINTS.md` | append-only checkpoint log | no |

---

## 7. Checkpoint protocol

Every session ends with a checkpoint entry: what was verified, what remains extracted-only, the next tranche,
and a resume token of the form `HEKI-CP-YYYY-MM-DD-NN`. The token is searchable; a checkpoint entry is never
edited after issue — a correction is a new entry that names the one it corrects.

**Resume procedure:** read the latest checkpoint → open the ledger → continue newest-first from the first unverified node.

---

## 8. Census boundary (explicit)

**The artifacts of this index are excluded from the source-corpus census, by design.**

The census population is every `*.md` under `docs/` **except** `docs/knowledge/` — currently **808 documents**.
That population has not changed since `HEKI-CP-2026-11-22-01`, and no identity was renumerated when this
boundary was made explicit.

| Artifact | In the corpus census? | Why |
|---|---|---|
| `docs/knowledge/*` (this charter, ledger, relations, tranches, checkpoints) | **No** | The index is a view over the corpus. Counting it inside the corpus would inflate coverage, make every coverage figure self-referential, and let index prose masquerade as evidence. |
| `docs/DOCUMENT-REGISTRY.md` | **Yes** (`HEKI-D-0003`) | It is a pre-existing corpus document, not an artifact of this program. |

**Consequence, accepted deliberately:** HEKI asserts nothing about its own completeness, and its files carry no
`HEKI-D` identity. If an index artifact ever needs citing as *evidence* — for example this charter's rules —
it must be promoted into the corpus by a separate, explicit act, not by silent inclusion.

**Count arithmetic (so nobody has to guess):** 888 `.md` files exist under `docs/` at this checkpoint; 8 are
index artifacts; **880 are corpus**. Before this program, 808 existed. (`ledger.jsonl`, `mentions.jsonl` and
`filerefs.jsonl` are not `.md` and do not enter this arithmetic.)

**Seventy-two corpus documents are unindexed, deliberately.** They were created on 2026-11-22 after the census:
`architecture/HOWL-FOUNDATIONAL-IMPLEMENTATION-EXECUTION-PLAN-2026-11-22.md`,
`architecture/TASK-01-PATHSERVICE-OWNERSHIP-AND-COMPLETION-2026-11-22.md`,
`architecture/TASK-01-AUDIT-2026-11-22.md`, `architecture/TASK-02-READER-SHARING-DECISION-2026-11-22.md`,
`architecture/TASK-02-DURABLESTORE-2026-11-22.md`, `architecture/TASK-02-AUDIT-2026-11-22.md` and
`architecture/TASK-03-CONVERSATION-AUTHORITY-PREPARATION-2026-11-22.md` and
`architecture/TASK-03-LEDGER-DECISION-2026-11-22.md` and
`architecture/TASK-03-FC1-DISPOSITIONS-2026-11-22.md` and
`architecture/TASK-03-LEDGER-SCHEMA-PROPOSAL-2026-11-22.md` and
`architecture/TASK-03-CONTRACT-FINAL-2026-11-22.md` and
`architecture/TASK-03-CONTRACT-ACCEPTANCE-2026-11-22.md` and
`architecture/TASK-03-CONTRACT-CLARIFY-1-2026-11-22.md` and
`architecture/TASK-03-CONTRACT-CLARIFY-1A-2026-11-22.md` and
`architecture/TASK-03-CONTRACT-CLARIFY-1B-2026-11-22.md` and
`architecture/TASK-03-CONTRACT-CLARIFY-1C-2026-11-22.md` and
`architecture/TASK-03-CONTRACT-CLARIFY-1D-2026-11-22.md` and
`architecture/TASK-03-CONTRACT-CLARIFY-1E-2026-11-22.md` and
`architecture/TASK-03-CONTRACT-CLARIFY-1F-2026-11-22.md` and
`architecture/TASK-03-CONTRACT-CLARIFY-1G-2026-11-22.md` and
`architecture/TASK-03-CONTRACT-CLARIFY-1H-2026-11-22.md` and
`architecture/TASK-03-CONTRACT-CLARIFY-1I-2026-11-22.md` and
`architecture/TASK-03-CONTRACT-CLARIFY-1J-2026-11-22.md` and
`architecture/TASK-03-SPIKE-SPEC-E27-E28-2026-11-22.md` and
`architecture/TASK-03-SPIKE-FINDINGS-01-2026-11-22.md` and
`architecture/TASK-03-SPIKE-FINDINGS-01A-2026-11-22.md` and
`architecture/TASK-03-SPIKE-FINDINGS-02-2026-11-22.md` and
`architecture/TASK-03-SPIKE-FINDINGS-03-2026-11-22.md` and
`architecture/TASK-03-SPIKE-FINDINGS-04-2026-11-22.md` and
`architecture/TASK-03-SPIKE-FINDINGS-05-2026-11-22.md` and
`architecture/TASK-03-RECONCILIATION-CONTRACT-2026-11-22.md` and
`architecture/TASK-03-RECONCILIATION-OQ3-2026-11-22.md` and
`architecture/TASK-03-RECONCILIATION-OQ3A-2026-11-22.md`,
`architecture/TASK-03-CREATION-IDENTITY-DEFECT-2026-11-22.md` and
`architecture/TASK-03-RECONCILIATION-OQ3B-2026-11-22.md` and
`architecture/TASK-03-RECONCILIATION-OQ3C-2026-11-22.md` and
`architecture/TASK-03-CONTRACT-AUDIT-2026-11-22.md` and
`architecture/TASK-03-RECONCILIATION-CONTRACT-V2-2026-11-22.md` and
`architecture/TASK-03-RECONCILIATION-GATE-2026-11-22.md` and
`architecture/TASK-03-GATE-CALLGRAPH-2026-11-22.md` and
`architecture/TASK-03-DURABLE-PAIR-2026-11-22.md` and
`architecture/TASK-03-FINAL-AUDIT-2026-11-22.md` and
`architecture/TASK-03-RECONCILIATION-CONTRACT-V3-2026-11-22.md` and
`architecture/TASK-03-RECONCILIATION-CONTRACT-V3-1-2026-11-22.md`,
`architecture/TASK-03-BASELINE-ACCEPTANCE-2026-11-22.md` and
`architecture/TASK-03-OQ6-SETTLEMENT-2026-11-22.md` and
`architecture/TASK-03-OQ7-ENFORCEMENT-2026-11-22.md` and
`architecture/TASK-03-OQ7-LIFETIME-CONTAINMENT-2026-11-22.md` and
`architecture/TASK-03-OQ7A-PATH-CENTRALISATION-2026-11-22.md` and
`architecture/TASK-03-OQ7D-NONHISTORY-OWNERSHIP-2026-11-22.md` and
`architecture/TASK-03-OQ7B-LOCK-SCOPE-2026-11-22.md` and
`architecture/TASK-03-OQ7C-REENTRANCY-2026-11-22.md`,
`architecture/TASK-03-OQ7CI-ROOT-AUTHORITY-2026-11-22.md` and
`architecture/TASK-03-DEFECT-FINALSAVE-WINDOW-2026-11-22.md` and
`architecture/TASK-03-OQ7CII-POINTER-CONSISTENCY-2026-11-22.md` and
`architecture/TASK-03-OQ7CIII-REENTRANCY-2026-11-22.md` and
`architecture/TASK-03-OQ7DI-SHADOW-TREE-SCOPE-2026-11-22.md` and
`architecture/TASK-03-OQ7DII-RETRY-IDENTITY-2026-11-22.md`,
`architecture/TASK-03-OQ7DIII-CALLER-OUTCOMES-2026-11-22.md` and
`architecture/TASK-03-OPEN-QUESTION-REGISTER-AUDIT-2026-11-22.md` and
`architecture/TASK-03-OQ1-DISCONNECTED-ROOT-2026-11-22.md` and
`architecture/TASK-03-OQ2-ALIAS-IDENTITY-2026-11-22.md` and
`architecture/TASK-03-OQ3BB-TERMINATOR-2026-11-22.md` and
`architecture/TASK-03-ROW550-CREATION-COLLISION-2026-11-22.md` and
`architecture/TASK-03-V344-HEALTHY-MEANING-2026-11-22.md` and
`architecture/TASK-03-MEMORIES-LIMITATION-2026-11-22.md` and
`architecture/TASK-03-RECONCILIATION-CONTRACT-V3-2-2026-11-22.md` and
`architecture/TASK-03-REMAINING-QUESTIONS-AUDIT-2026-11-22.md` and
`architecture/TASK-03-V32-DISPOSITION-2026-11-22.md` and
`architecture/TASK-03-D2-D3-PROGRAMMER-ASSESSMENT-2026-11-22.md` and
`architecture/TASK-03-S1-IMPLEMENTATION-AND-AUDIT-2026-11-22.md` and
`architecture/TASK-03-DELETION-AND-S1-CORRECTION-2026-11-22.md`. The ledger still holds **808 nodes**.

**A reserved identifier is not a reconciled count.** 0809/0810/0811/0812/0813/0814/0815/0816/0817/0818/0819/0820/0821/0822/0823/0824/0825/0826/0827/0828/0829/0830/0831/0832/0833/0834/0835/0836/0837/0838/0839/0840/0841/0842/0843/0844/0845/0846/0847/0848/0849/0850/0851/0852/0853/0854/0855/0856/0857/0858/0859/0860/0861/0862/0863/0864/0865/0866/0867/0868/0869/0870/0871/0872/0873/0874/0875/0876/0877/0878/0879/0880
are **placeholders** for those seventy-two files when indexing continues. None has a ledger node, a recorded sha256, an extracted status or link counts, so none
may be cited as census evidence. Until they do, the census is **808** and the corpus is **811**, and the difference
is a known enumerated gap — not a total. No existing identity is renumbered, so "newest-first" describes the
original census, not the live tree.

**This number now moves whenever implementation adds a document.** That is expected: the corpus is no longer
frozen by the index, by the owner's 2026-11-22 decision. The census remains a dated measurement, not a live count.

**Ordering caveat (measured, `HEKI-CP-2026-11-22-03`).** Identity order, ledger order and "newest-first" reading
order are **filename orders, not chronologies**. 15 of 622 filename-dated documents contain a head date later
than their filename date; both 2026-11-20 architecture documents are among them. The ordering is still the
correct stable ordering — it never changes — but a sequence statement must say "filename-dated".

---

## 9. Mentions are not relationships

Three connectivity layers are recorded separately and must never be conflated:

| Layer | File | Evidence value |
|---|---|---|
| **Hyperlink** — a document links to another document | `ledger.jsonl` `links_out` / `links_in` | Navigational intent |
| **Mention** — a document names an identifier | `mentions.jsonl`, one record per occurrence with file, line number and the full source line | A named connection of unknown type |
| **Filename reference** — a document names another document's file, usually in backticks, without linking it | `filerefs.jsonl`, same shape as `mentions.jsonl` | A named connection of unknown type |

The third layer exists because the first two miss most real citation: 2,642 filename references stand in 490 of
808 documents, and 98.4% of them (2,600 of 2,642) are not accompanied by a hyperlink to the same target. Reading `links_in = 0` as
"unreferenced" is therefore wrong — 172 of the 304 hyperlink-isolated documents are referenced this way, leaving
**129** unreferenced in all three layers.

A **mention** proves only that a document names an identifier. It is *not* evidence of dependency, authority,
adoption, agreement, or sequence. `RELATIONS.md` is the only place a semantic type (`implements`, `supports`,
`supersedes`, `contradicts`, `governs`, `derived_from`, …) may be asserted, and only with a quoted reason from
a verified read.

**Never do this:** promote a mention to an edge because it recurs. Repetition is not adoption (R4). Ten
mentions of a proposal still leave it a proposal.

---

## 10. Status: as-written vs current (R9)

Every node carries `status` exactly as the document writes it. That field is **immutable with respect to
later documents**: no citation count, no later reference and no later acceptance review may overwrite it.

A **current status** is a separate judgement recorded only in `RELATIONS.md`, and it may differ from the
as-written status (`Proposed` → `Owner-approved` → `Applied` → `Superseded`). When the two differ, both are
shown. This is the mechanical form of R4: adoption is something a *document* declares, never something a
*reference* implies.
