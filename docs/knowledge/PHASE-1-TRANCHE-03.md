# Phase 1 — Discover · Tranche 03 (2026-11-20 · the foundation pack)

**Checkpoint:** `HEKI-CP-2026-11-22-03`
**Method:** same standard as Tranches 01–02 — front matter, status block, scope, decision tables and
relationship statements read in full. Nothing edited.

**Coverage after this tranche:** 22 of 808 documents verified (2.7%).

**Scope correction, measured before reading.** Tranche 03 was announced as "the 2026-11-20 foundation pack
(registry rows 505–512)". That scope was wrong. Rows 505–511 are the **2026-11-21** contract papers and their
register, already verified in Tranche 02; row 512 is the intake. The 2026-11-20 architecture population is
**exactly two documents**, and both are read here.

The 62 `docs/research/MEMORY-09-*` documents carrying a 2026-11-20 filename date are a different population
(1,065,774 bytes) and are **not** covered by this tranche. They become Tranche 04, before 2026-11-19, so that
newest-first order is preserved and nothing is silently skipped.

---

## 1. `HEKI-D-0088` — Front-to-Back Foundation Architecture and Delivery Plan

`docs/architecture/HOWL-FOUNDATION-FRONT-TO-BACK-ARCHITECTURE-AND-DELIVERY-PLAN-2026-11-20.md`
32,384 B · 372 lines · identity `HF-ARCH-PLAN-2026-11-20-01` · **unregistered**.

- **Purpose:** the single reconciled system plan — conversation state, context construction, multi-character
  isolation, persistence and authoring, model/modality lifecycle, hardware admission, recovery, delivery gates
  and evidence. It states its own reason for existing: "deliberately front-to-back before work decomposition."
- **Status (verbatim):** "Planning and DeepSeek-paperwork preparation only. This addendum makes no
  implementation authorization, architecture adoption, schema/API decision, migration decision, test result,
  performance claim, or release claim."
- **Evidence labels defined here:** *Verified in current tree*, *Prior reported evidence*, *Owner-approved
  requirement*, *Proposal*, *Open decision*. This is the vocabulary the later papers inherit.
- **What it commissions:** §10 is the origin of the contract pack — "Produce a **Foundation Contract Pack**
  for conversation identity/state, composition/token budgeting, content ownership, and resource/capability
  lifecycle. Do not write implementation plans before unresolved decision rows are answered."
- **Gates:** M0A, M0B, M1–M8 plus chat and modality finalization, each with a dependency and an exit-evidence
  requirement (§5).
- **Evidence ladder:** §7.1 already states **requested → admitted → attempted → observed → verified**, or
  unknown/failed, with "UI state, process start, log output, and function return cannot advance a claim."
  The ladder in `FC5-22a` and FC-6 is therefore **not new** in the contract pack; it is inherited from here.
- **Decision register:** §9 mints 23 open decisions — `DQ-C01`–`C07`, `DQ-X01`–`X05`, `DQ-R01`–`R08`,
  `DQ-M01`–`M03`.
- **Stop conditions:** §8 — hardware claims from estimates or sleeps must report unknown, not ready; two
  instances racing for capacity must enforce a policy before concurrent heavy work; MEMORY-09 scope creep stops
  at the lifecycle-fact interface.
- **Connectivity:** inbound hyperlinks **0**, outbound **0**, identifier mentions **0**. It names ten controlling
  inputs in §1.1 and those are backticked filenames, not links. §5 and §7 name gates and receipts, also without
  links.
- **Adoption:** none. Planning artifact; explicitly no architecture adoption and no implementation.

## 2. `HEKI-D-0089` — Foundation Contract Pack and DeepSeek Intake

`docs/architecture/HOWL-FOUNDATION-CONTRACT-PACK-DEEPSEEK-INTAKE-2026-11-20.md`
35,184 B · 331 lines · identity `HF-CONTRACT-INTAKE-2026-11-20-01` · registry row **512**.

- **Purpose:** the bounded intake for the papers, the decision register, the delegation, and "the first
  implementation authorization that may be considered afterward."
- **Status (verbatim):** "**Owner decision register recorded; FC-1–FC-6 review-complete and accepted as
  reviewed proposed drafts (2026-11-21).** … Review acceptance is a record that the papers were reviewed; it
  does not approve FC-1–FC-7, authorize implementation, modify a schema/API, select a migration, close an
  audit, or report a test result."
- **Lineage declared:** "**Parent plan:** `HOWL-FOUNDATION-FRONT-TO-BACK-ARCHITECTURE-AND-DELIVERY-PLAN-2026-11-20.md`
  (`HF-ARCH-PLAN-2026-11-20-01`)." This is the only place in the corpus that names the parent plan.
- **Authority the contract pack depends on:** D1–D8 in §0.2, D1–D2 rationale in §3, and the §8 record
  "Contract papers approved: **None**. … M0A authorized: **No**."
- **Two rules that bind later documents:** the *history rule* — "A recorded owner decision supersedes a
  recommendation without deleting it. No later paper may cite a superseded recommendation as the current
  decision"; and the *authority rule* — recording decisions "does **not** approve FC-1–FC-7 and does **not**
  authorize M0A."
- **D8 is a delegation, not a decision:** "If the two agents disagree, the disagreement returns to the owner.
  It is not resolved by whichever agent writes the document last." FC-7 therefore has no authority to inherit.
- **Connectivity:** inbound hyperlink **1** (the contracts register), outbound **0**, identifier mentions **29**.
- **Adoption:** owner decisions D1, D3, D5, D6 selected; D2 selected with narrowed wording; D4 direction only;
  D7 requirement clarified; D8 delegated. Contracts unapproved; implementation unauthorized.

---

## 3. Findings

### T3-F1 — "304 disconnected documents" overstated isolation by 175

The census measured isolation from the hyperlink layer alone. Adding the filename-reference layer:

| Measure | Documents |
|---|---|
| No hyperlink in or out | **304** |
| …that cite another document by backticked filename and/or are named that way | **172** |
| …still unreferenced in the reference layer | **132** |
| Unreferenced in **all three** layers (hyperlink, filename reference, identifier mention) | **129** |

So the honest isolation figure is **129 of 808 (16%)**, not 304 (38%). The remainder are discoverable; they
simply do not use link syntax. Recorded as a correction to CP-01's "304 documents with no link in either
direction", which remains true as written but was being read as "disconnected".

### T3-F2 — Filename references are a third connectivity layer

Measured corpus-wide, excluding the index:

| Metric | Value |
|---|---|
| Filename-reference occurrences (backticked `*.md` naming an existing corpus file) | **2,642** |
| Distinct (citing document → target) pairs | **2,120** |
| Documents that cite at least one other document this way | **490 of 808** |
| Those references accompanied by a hyperlink to the same target | **42 (1.6%)** — 98.4% are unlinked |

Materialised as `filerefs.jsonl` (one row per occurrence: citing document, target basename, line, source text).
It is an **occurrence index**, on exactly the same footing as `mentions.jsonl`: it proves naming, never
dependency. The two layers are independent — `HF-ARCH-PLAN` appears zero times in `mentions.jsonl` because it
has no program-identifier form that the mention regex accepts.

### T3-F3 — A filename date is not a content date

Both 2026-11-20 documents contain 2026-11-21 material.

- The plan's §1.1 names six documents dated **2026-11-21** as controlling inputs (HFS-AUDIT, HFR, HF1, HCDM,
  M0A, HRA), and its §12 verification record claims to have inspected them.
- The intake's status line, §0 heading, §0.4 and §8 record all describe 2026-11-21 review and decision events.

This is not confined to these two. Corpus-wide, **15 of 622** filename-dated documents contain a head date later
than their filename date. The census already holds this evidence in `dates_head`; this tranche interprets it.

**Consequence for HEKI, and it is structural:** `HEKI-D` identity order, `DISCOVER-LEDGER.md` order and the
"newest-first" reading order are **filename orders**, not chronologies. They remain the correct stable ordering
(they never change, and R2 depends on that), but every statement of sequence must say "filename-dated" rather
than "dated". A document named for the 20th can post-date a document named for the 21st.

### T3-F4 — The document that commissions the contract pack is the least discoverable one in it

`HEKI-D-0088` has 0 inbound hyperlinks, 0 outbound hyperlinks and 0 identifier mentions. Its child — the intake —
declares it as "Parent plan" in a backtick, without a link. Nothing else in 808 documents names it.

So the chain *plan → intake → FC-1…FC-7* exists in the corpus but is traversable in only one direction and only
through one unlinked string. This is the strongest Phase 2 candidate found so far, and it is a different failure
from T2-F4 (the orphaned G1 audit): there, downstream documents cited a gate whose evidence was isolated; here,
the plan is the **upstream authority** and the isolation is in the inbound direction.

### T3-F5 — Two decision vocabularies, no crosswalk, and 23 orphaned identifiers

The plan mints 23 open decisions: `DQ-C01`–`C07` (conversation/authoring), `DQ-X01`–`X05`
(context/multi-character), `DQ-R01`–`R08` (resources), `DQ-M01`–`M03` (MEMORY-09 boundary).

Measured: **each of those 23 identifiers appears in exactly one document — the plan itself.** None is carried
into the intake, the contract papers, the resource baseline or the registry.

The resource questions were re-minted instead. `DQ-R01`–`R08` reappear as `DQ-HRA-01`–`DQ-HRA-12` in
`HRA-DB-2026-11-21-01`, which is what FC-5 actually cites. The overlap is close but not one-to-one — for example
the plan's single `DQ-R05` ("what survives crash and how are stale leases/orphaned processes reconciled") becomes
two questions, `DQ-HRA-04` and `DQ-HRA-05`, and HRA adds questions the plan did not ask (sub-budget sharing,
estimate overrun, another application taking memory).

**This is a proposed reading, not an established edge.** No document states the mapping; I am matching question
text. What *is* established: the older identifiers are dead in identifier space, and a reader following `DQ-R0*`
will find no successor. Recorded so that Phase 4 does not treat "no later citation" as "question closed".

The one explicit crosswalk that does exist is inside FC-5: `FC5-25` — "D6 option A is a temporary answer to one
part of `DQ-HRA-07`; it does not close the question."

### T3-F6 — The plan and HRA agree on modalities; HCDM was the outlier

The plan's §2.4 table (2026-11-20 filename) already says: audio generation owner-reported complete; TTS deferred
but required for chat finalization; image generation in scope and nearly complete; vision largely complete with
gaps; video in scope and incomplete; 2D/3D avatars out of scope.

`HRA-04` (2026-11-21) supersedes "older statements in `HCDM-2026-11-21-01` that described image generation as
deferred and audio generation as removed" — and then reaches the same position the earlier plan already held.

**Recorded as an observation, not a verdict.** It does not make the plan authoritative and it does not make HCDM
wrong in its own terms; HCDM's §0 already carries a scope-correction note. It does mean that a reader who trusts
filename order will read the *correction* before the *statement it agrees with*, which is precisely why T3-F3
matters.

### T3-F7 — No implementation authorization exists anywhere in this tranche, in any tense

Plan: "No stress, model, device, or destructive run is authorized by this planning artifact." Intake:
"Contracts unapproved; implementation unauthorized"; M0A "remains unapproved until the owner expressly adopts
it." Every decision recorded is *direction*, and the intake says so twice.

---

## 4. Edges established by reading (added to `RELATIONS.md` §7)

| From | Type | To | Reason | Adoption |
|---|---|---|---|---|
| `HEKI-D-0089` intake | `derived_from` | `HEKI-D-0088` plan | "**Parent plan:** `…FRONT-TO-BACK…` (`HF-ARCH-PLAN-2026-11-20-01`)" | Planning only |
| `HEKI-D-0011`…`HEKI-D-0017` FC-1…FC-7 | `derived_from` | `HEKI-D-0089` intake | Intake §2 assigns each paper its output, exclusions and gate; intake §10.2 is the commissioning instruction | **Proposed drafts** |
| `HEKI-D-0088` plan | `precedes` | every contract paper | "Do not write implementation plans before unresolved decision rows are answered" | Ordering only |
| `HEKI-D-0012`…`HEKI-D-0017` FC-1…FC-6 | `derived_from` | `HEKI-D-0088` §7.1 evidence ladder | Ladder identical in substance: requested → admitted → attempted → observed → verified, unknown/failed | Inherited |
| `HEKI-D-0088` plan | `contradicts` | any reading of the corpus as a filename chronology | §1.1 cites six 2026-11-21 documents; §12 claims they were inspected | **Open provenance finding** |

---

## 5. What this tranche changes about HEKI's purpose

Tranche 02 showed the index could catch a **status contradiction** between a document and the state of a file.
This tranche adds two failure modes a link graph cannot see at all:

1. **A connectivity figure that was an artifact of syntax.** The corpus was never 38% disconnected; it was 16%,
   and the difference is entirely backticks versus brackets.
2. **A chronology that was an artifact of filenames.** The reading order is stable and reproducible, and it is
   not chronological.

Both were found by measuring the index against itself rather than by reading more documents. Neither required
editing an original.

---

**Next:** Tranche 04 — the 62 `research/MEMORY-09-*` documents dated 2026-11-20 (the decision and evidence
spine: DQ-2/3/5 closures, P-1 and S-2/S-3 outcomes, the implementation authorization chain). Then 2026-11-19.
**Resume token:** `HEKI-CP-2026-11-22-03`

**Nothing outside `docs/knowledge/` was modified by this tranche.**
