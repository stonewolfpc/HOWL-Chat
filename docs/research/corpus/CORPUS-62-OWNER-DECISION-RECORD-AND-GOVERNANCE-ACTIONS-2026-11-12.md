# `CORPUS-62` — Owner Decision Record and Authorized Corpus Governance Actions

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identifier:** `CORPUS-62` · **Version:** 1 · **Date:** 2026-11-12.
**Record type:** owner decision record plus administrative correction and filed proposals.
**Assigned model identity (owner-supplied by the work order):** **DeepSeek**. Recorded as owner metadata. **Not** inferred from, and not substituted by, any self-description; **not** verifiable from the workspace.
**Executing agent:** Codebuff coding agent under the framework persona label "Buffy". The persona label is a framework identity, **not** model identity.
**Execution:** fresh conversation.
**Governing context:** `MEMORY-09` `G1`.
**Status:** Seven owner decisions **recorded**. One administrative correction **applied** (registry row 393). Five supporting proposals **prepared, filed, registered, and NOT applied**. No corpus item, standard, schema, runtime, model, or checker implementation was modified.

---

## 1. Purpose and authorization

The owner reviewed the seven recommendations arising from `CORPUS-61` and approved them. `CORPUS-61` closed with **no overall PASS** and recorded — explicitly — that `F3-01` diagnosticity, `R-2` consistency, the row-393 disposition, and the `F6-01` evidence gap were owner decisions. This record operationalizes those decisions and nothing else.

**Authorized:** recording the decisions; preparation of proposals; the specific administrative registry correction under Decision 1 (row linkage and role, with the original registration preserved); appraisal and registration of the resulting artifacts; running the existing authorized mechanical checker.

**Not authorized, not done:** general corpus cleanup; opportunistic improvements; applying any repair; amending `CORPUS-01`, `MEMORY-09`, or `CORPUS-07`; changing any item, tag, split, or acceptance condition; implementing or modifying any checker; model inference; beginning `G2`; sealing `G1`; annotating; altering historical reports; deleting or overwriting any historical artifact.

---

## 2. Governing boundaries maintained

| Boundary | State after this filing |
|---|---|
| `G1` | **Open** |
| Corpus | **Frozen**, except the registry administrative correction authorized below |
| Expansion | **Frozen** |
| `AM-01` | **Unissued** (verified: `CORPUS-42` records preparation only, not application) |
| `A-2a-R1` | **NOT IN FORCE** (verified: `CORPUS-50` close-out) |
| `D1` | **Unadopted** |
| Annotation | **Blocked** |
| Model inference | **None run** |
| Schema / runtime / model | **Unchanged** |
| Checker implementation | **Unchanged** (`check_corpus.py` hash recorded in §9.5) |
| Convergence | **0 of 2** |

The distinction the work order requires is maintained throughout: **owner-approved decisions** (§3) are separated from **documentation of them** (this record), from **proposals awaiting separate application** (`CORPUS-64`–`67`), from the **administrative correction actually performed** (`CORPUS-63`), and from **future implementation requiring separate authorization** (§10).

---

## 3. The seven owner decisions

Each entry states the approved decision, its rationale, its scope, its consequences, and the exact boundary of what it authorizes.

### Decision 1 — `CORPUS-58` registry disposition

**Approved:** the root-level `CORPUS-58` assignment brief is a **real historical artifact**, but registry row 393 **incorrectly represents it as the substantive review**.

**Rationale.** `CORPUS-61` §5 established that row 393's hash, byte count, and line count match a real root-level file exactly, that the linked corpus path has no target, and that the matched file's content is a work order with no findings. Both `CORPUS-59` ("dangling row, unbacked metadata") and `CORPUS-60` ("no file on disk") had scoped their searches to `docs/` and missed it. The row's numbers were never wrong; its link and its role were.

**Scope.** Registry row 393 and the registry's declared scope line. Nothing else.

**Required actions and disposition:**

| # | Required | Status |
|---|---|---|
| 1 | Preserve the root brief byte-for-byte | **Done** — no write; hash re-verified `292befdd…` |
| 2 | Preserve its original hash and metadata | **Done** — 4,311 bytes / 80 lines / `292befdd…` unchanged in the row |
| 3 | Prepare an auditable registry correction identifying artifact, location, true role | **Done and applied** — see `CORPUS-63` §4 |
| 4 | Preserve the fact that row 393 originally misidentified the artifact | **Done** — original role and original wrong path recorded inside the row and quoted verbatim in `CORPUS-63` §3 |
| 5 | Do not represent the brief as the completed review | **Done** — role now "Assignment brief (work order)" |
| 6 | Do not erase the original registration history | **Done** — row number, row existence, and original claim all preserved |

**Consequences.** Row 393 now resolves. The registry's scope line records the second root-level exception. The substantive review is **still unfiled**. Any claim of the form "`CORPUS-58` found X" remains unverifiable at source.

**Authorizes:** the single row correction and the scope/count lines that must follow from it. **Does not authorize:** renaming the brief, filing a reconstruction, registering the brief a second time, or any change to `CORPUS-59`/`60`/`61`.

### Decision 2 — the missing substantive `CORPUS-58` review

**Approved:** do **not** reconstruct the original review and represent it as an original filed artifact. The conversation-supplied review may be retained as **conversation evidence**. If a reconstruction is useful, it must be a new identifier, explicitly labelled a reconstruction, sourced to the conversation-supplied text, with original file/hash/provenance recorded as unverifiable, with no byte-identity claim, and with no convergence credit.

**Rationale.** The original is unavailable; its metadata describes the brief, not the review; and reconstructing it from other reports' descriptions would produce a third-order paraphrase presented as a first-order artifact.

**Scope.** Provenance disposition only.

**Disposition taken.** **No reconstruction artifact was created.** The source text is not in this workspace and was not furnished in this execution, so a reconstruction would require inventing exactly the metadata Decision 2 forbids inventing. The full provenance-record schema, with every unverifiable field left unverifiable, is filed at `CORPUS-63` §§5–5.1. The two prohibitions (no byte-identity claim; no convergence credit) are recorded there so they bind any future reconstruction in advance.

**Consequences.** `CORPUS-58`'s substantive content remains unavailable as evidence. `CORPUS-60`'s citation checks and `CORPUS-61`'s items 60-4 and 60-12 keep their **UNVERIFIABLE** status. Convergence remains **0 of 2**.

### Decisions 3 and 4 — `F3-01` construct and adjudication

**Approved:** `F3-01` is intended to test **memory-dependent recognition of the persona's own historical experience**, not generic conversational appropriateness. The `CORPUS-07` adjudication is **historically valid evidence of an earlier proposed resolution**, but its **existing disjunction is insufficient to establish memory diagnosticity**.

**Rationale.** `CORPUS-07` §2 correctly identified a rule conflict and proposed a pass condition, but `CORPUS-61` §6.3 showed that its second disjunct ("offers rain-specific care unprompted") passes any warm hospitality without memory, and its first disjunct's letter admits the second-person generic ("The damp gets into your hands"). `CORPUS-60`'s "Type P if adopted" was therefore withdrawn. The adjudication's historical validity is fully preserved; only its sufficiency as a diagnosticity claim is denied.

**Scope.** `F3-01` only.

**Required work and disposition.** `CORPUS-64` is the revised repair proposal, satisfying all eight requirements: recognition of the persona's own historical information (construct); generic second-person hand references insufficient; generic rain advice or care insufficient; no pass via ordinary knowledge, generic empathy, scene context, or convention; the memory-absent counterpart cannot satisfy the same conditions; no verbatim quotation or unnatural recital required; natural expression and equivalent phrasings preserved; overconstraint avoided. It contains the pass-set comparison over identical non-memory context, and separates intended construct, required evidence, mandatory behaviour, acceptable variation, prohibited responses, memory-dependent diagnosticity, and remaining uncertainties.

**Consequences.** `F3-01` remains **UNRESOLVED** as a classification. It is **not** Type P. It is **not** repaired. `CORPUS-07` §2 remains an unapplied adjudication. The proposal **supersedes neither** the item nor the adjudication, and says so in its own §1.

**Authorizes:** preparation of the proposal. **Does not authorize:** applying it, adopting any pass condition, re-tagging the item, or beginning a scoring run.

### Decision 5 — `R-2` coverage discrepancy

**Approved:** treat the discrepancy as a **specification inconsistency**, not an assumed item-tag defect.

**Rationale.** `CORPUS-01` §3.1's `R-2` row lists three items; its own coverage line states `R-2 4`; four units carry the tag, including `F1-03`; and every downstream "four fixtures" declaration is consistent with the count and the tags, not with the row. The discrepancy is internal to `CORPUS-01` as well as between the specification and the items (`CORPUS-61` §7.1).

**Scope.** `CORPUS-01` §3.1, one row, plus a version notice. Not `MEMORY-09`.

**Required work and disposition.** `CORPUS-65` reconciles all five elements, names the complete affected-unit list (`F1-01`, `F1-02`, `F6-01`, `F1-03`), shows before/after wording, records the tag-side alternative without proposing it, and does **not** remove `F1-03`, change its tag, or infer which source was intended.

**Consequences.** `CORPUS-01` is **unchanged**; `F1-03`'s tag is **unchanged**; the amendment is **not applied**. The carrying form for future documents stays "3 specified / 4 tagged" until the amendment is authorized.

**Authorizes:** preparation of the dated amendment proposal. **Does not authorize:** applying it, re-tagging `F1-03`, or editing any prior report.

### Decision 6 — `F6-01` validity and evidence sufficiency

**Approved:** `F6-01` **cannot serve as a valid positive test in its current form**, because the scenario does not establish the cake's relevant contents.

**Rationale.** The behavioural threshold (`CORPUS-07` §2: observable accommodation; silence fails) resolves *form*, not *content*. With the cake's contents unstated, no reply can select the correct accommodation, and no rater can distinguish respect from ignorance. The two defects are independent and the second survives adoption of the first (`CORPUS-60` §4, confirmed by `CORPUS-61` §8).

**Scope.** `F6-01` only.

**Required work and disposition.** `CORPUS-66` supplies the restriction and the ingredient information in scenario evidence, preserves "silence does not satisfy", preserves observable accommodation, requires no invented ingredients, is not a recitation exercise, contains memory-present and memory-absent acceptance comparisons, and lists every change to `Scenario`, `Required`, `Variation`, `Forbidden`, `Notes` (and to the other affected fields). It keeps the behavioural threshold and evidence sufficiency as separate dimensions and states that adoption of either alone leaves the item unscorable.

**Consequences.** `F6-01` remains **unscorable as written** and **unrepaired**. The adjudicated threshold remains **unapplied**. The proposal also records a **pre-existing, unaddressed tension**: the item is tagged `R-2` ("the fact is never stated") while the threshold permits a warning, which states the fact. That tag is not changed here.

**Authorizes:** preparation of the repair proposal. **Does not authorize:** applying it, adopting the threshold, or changing the mode tag.

### Decision 7 — registry / checker integrity

**Approved:** future registry integrity validation **should fail when a registered target is missing or incorrectly linked**.

**Rationale.** The checkers validated item records and never opened the registry. A row pointed at a non-existent file and the check returned `PASS (0 fail, 0 warn, 0 unverifiable)` — verified from the source, which contains no registry logic, and observed in the same session in which the missing target existed.

**Scope.** A specification for future validation. Nothing implemented.

**Required work and disposition.** `CORPUS-67` specifies every required capability (target resolution, identity/path match, bytes, lines, hash, duplicate identifiers, conflicting registrations, root-level scope, explicit failure instead of passing, the four-way failure taxonomy, read-only behaviour) and the eight required fixture families, including both halves of the out-of-scope fixture. It states the current checker's verified limits.

**Consequences.** `check_corpus.py` is **unmodified**. No registry repair tool exists. The registry remains verified only by the manual sweep in §9.4.

**Authorizes:** a specification. **Does not authorize:** any implementation, modification, or automatic repair.

---

## 4. Artifact inventory

| # | Artifact | Identifier | Version | Role | Applied? |
|---|---|---|---|---|---|
| 1 | `CORPUS-62-OWNER-DECISION-RECORD-AND-GOVERNANCE-ACTIONS-2026-11-12.md` | `CORPUS-62` | 1 | This decision record | — |
| 2 | `CORPUS-63-CORPUS-58-REGISTRY-CORRECTION-AND-PROVENANCE-DISPOSITION-2026-11-12.md` | `CORPUS-63` | 1 | Registry correction record + provenance disposition | **Correction applied; no reconstruction** |
| 3 | `CORPUS-64-F3-01-MEMORY-DIAGNOSTICITY-REPAIR-PROPOSAL-2026-11-12.md` | `CORPUS-64` | 1 | `F3-01` repair proposal | **No** |
| 4 | `CORPUS-65-CORPUS-01-R2-COVERAGE-AMENDMENT-PROPOSAL-2026-11-12.md` | `CORPUS-65` | 1 | `CORPUS-01` `R-2` amendment proposal `A-R2-01` | **No** |
| 5 | `CORPUS-66-F6-01-SCENARIO-EVIDENCE-REPAIR-PROPOSAL-2026-11-12.md` | `CORPUS-66` | 1 | `F6-01` repair proposal | **No** |
| 6 | `CORPUS-67-REGISTRY-CHECKER-INTEGRITY-SPECIFICATION-2026-11-12.md` | `CORPUS-67` | 1 | Checker enhancement specification | **No implementation** |

**Why five separate artifacts rather than one.** Each has an independent authorization status: `CORPUS-63` contains a correction that *was applied* and a disposition that is *complete*; `CORPUS-64`, `CORPUS-66` and `CORPUS-67` each need a separate adoption decision and their own re-review; `CORPUS-65` is an amendment that would change a governing specification. Combining them would let one approval be read as covering another, which is precisely the conflation the governing rules forbid.

All six are filed in `docs/research/corpus/`, matching the existing `CORPUS-*` convention. No `.sha256` sidecar files were created, because no other `CORPUS-*` artifact has one; hashes live in the registry.

---

## 5. Cross-document reconciliation

Preserved without edit: `CORPUS-01`, `CORPUS-07`, `MEMORY-09`, `CORPUS-56`, `57`, `58` (root brief), `59`, `60`, `61`. Their hashes stand. No conclusion below is written back into them.

### 5.1 Confirmed

| Item | Source | Confirmed by |
|---|---|---|
| The checker verifies no registry targets, and returned PASS while a registered target was missing | `CORPUS-59` §3, `CORPUS-60` §1 | `CORPUS-61` §5.2; re-verified from source and by run (§9.5) |
| Row 393's hash/bytes/lines match the root brief; the linked path has no target; the brief's content is a work order | `CORPUS-61` §5.1–5.3 | Re-verified this session, byte for byte |
| `CORPUS-07` §2 contains a real adjudicated, unapplied `F3-01` pass condition; `CORPUS-59` failed to open it | `CORPUS-60` §2 | `CORPUS-61` §6.2; re-read |
| `F6-03` is structurally `N`; "weak P" is a category error | `CORPUS-59` §5 | `CORPUS-61` §9 |
| `F1-03` is `N` on scene-staging grounds, and its verdict is mode-independent | `CORPUS-59` §9 | `CORPUS-61` §7.2 |
| `F6-01` carries two independent defects | `CORPUS-60` §4 | `CORPUS-61` §8 |
| `R-2` set is 3 specified / 4 tagged | `CORPUS-60` §3 | `CORPUS-61` §7.1 |
| `R-2` root cause is conceptually decisive but not causally exclusive | `CORPUS-59` §6 | `CORPUS-61` §10 (59-7) |
| 45 units; 80-item floor; **35** remaining; 18 defect / 12 ambiguous / 15 clean | `CORPUS-59` §11, `CORPUS-07` §1 | Arithmetic and file counts re-verified |
| `AM-01` unissued; `A-2a-R1` NOT IN FORCE; `D1` unadopted; annotation blocked | `CORPUS-61` §12.8 | Re-verified against `CORPUS-42`, `CORPUS-50` |
| Convergence **0 of 2**; no reviewer lineage established; identities unverifiable | `CORPUS-59` §7, `CORPUS-61` §12.9–12.10 | Unchanged; this filing adds no convergence claim |

### 5.2 Corrected (by prior reports; preserved as corrections, not rewrites)

| Statement | Correction | Status |
|---|---|---|
| `CORPUS-58` "artifact absent; metadata unbacked; row dangling" (`CORPUS-59` §3) | Metadata was always backed by the root file; the **link and role** were wrong | **CORRECTED** (`CORPUS-61` §5.2) |
| `CORPUS-60` "no `CORPUS-58*` file exists on disk" | False repo-wide for the same scoping reason | **CORRECTED** (`CORPUS-61` §5.2) |
| `F3-01` attribution requirement was "reviewer interpretation" whose only source was `CORPUS-57`'s gloss | Real source: `CORPUS-07` §2 adjudication | **CORRECTED** (`CORPUS-60` §2, `CORPUS-61` §6.2) |
| `CORPUS-58`'s "pre-registered note" provenance | Substance right, provenance wrong — it is an adjudication, not a field annotation | **CORRECTED** (`CORPUS-60` §2.2) |
| "Four `R-2` fixtures" | Three specified; four tagged | **CORRECTED** (`CORPUS-60` §3) |
| "17 unique source documents" | 15 unique, from an **unfurnished** text — unverifiable here | **UNVERIFIABLE**, treated as reported only |
| "2 of 8" attributed to `CORPUS-57` | That wording is `CORPUS-56` §6; `CORPUS-57` §8 says 3 of 24 | **CORRECTED** (`CORPUS-60` §5) |
| "38 remaining of 80" | 35 | **CORRECTED** (`CORPUS-59` §11) |

### 5.3 Superseded by an owner decision

| Prior conclusion | Superseded by | New state |
|---|---|---|
| "Row 393 requires an owner decision: complete, relocate, or withdraw" (`CORPUS-61` §13.1) | **Decision 1** | **Resolved** — link and role corrected in place; row preserved; no relocation, no withdrawal, no duplicate row |
| "Substantive review: file it or leave it as conversation material" (`CORPUS-61` §13.2) | **Decision 2** | **Resolved as a disposition** — no reconstruction; conversation evidence only; provenance record schema filed |
| "`F3-01`: adopt, amend, or reject the `CORPUS-07` condition" (`CORPUS-61` §13.3) | **Decisions 3–4** | **Resolved in direction, not in operative text** — construct fixed as memory-dependent recognition; the adjudicated condition is declared insufficient; a replacement proposal exists and is unapplied |
| "`R-2` set: stale row or wrong tag?" (`CORPUS-61` §13.5) | **Decision 5** | **Framed** — specification inconsistency; amendment proposal prepared; the historical-intent question remains an owner determination |
| "`F6-01`: adopt the threshold and/or repair the scenario — two separate decisions" (`CORPUS-61` §13.4) | **Decision 6** | **Half-addressed** — the evidence repair is proposed; the threshold remains unadopted and is explicitly kept separate |
| "Should the checker verify registry targets?" (`CORPUS-59` §15.5) | **Decision 7** | **Resolved as a specification requirement**; implementation still unauthorized |
| `CORPUS-60` §6 "`F3-01` **P conditionally**; 2 confirmed if adopted" | **Decisions 3–4** | **Withdrawn** — `CORPUS-61` §6.3 showed the adopted text is memory-permeable. No conditional-P status survives. |

### 5.4 Still unresolved

| Item | Why it stays open |
|---|---|
| `F3-01` classification | Unscorable on operative text (rule conflict); memory-permeable on the adjudicated text; the proposed replacement is unadopted and unproven (§9 of `CORPUS-64`) |
| `R-2` historical intent (stale row vs wrong tag) | Requires an owner determination; explicitly not inferred (`CORPUS-65` §8 `U-1`) |
| `F6-01` behavioural threshold | Adjudicated, unapplied; no new decision was taken here |
| `F6-01` mode tag | `R-2` vs a threshold that permits a warning — recorded, not decided (`CORPUS-66` §7 `U-3`) |
| `F1-03` repair direction | (i) preserves the construct, (ii) risks recitation drift; neither authorized |
| `F6-03` repair | `CORPUS-59` §5's consequence-based proposal is unapplied and unadopted |
| The other ten `CORPUS-07` §2 ambiguous units | Twelve proposals were unapplied and remain so; this work order touched two (`F3-01`, `F6-01`) |
| `CORPUS-57`'s own wording | Never edited; whether it "stands corrected" remains an owner decision (`CORPUS-61` §13.6) |
| Provenance sufficiency | Identities unverifiable; whether owner metadata suffices for any future credit is an owner decision |
| `CORPUS-58` substantive content | Still unfiled and unfurnished |
| Registry summary drift | §2.1/§2.2 breakdowns do not sum to the stated totals; pre-existing (§9.4) |
| Two unregistered `docs/` files | Pre-existing; recorded, not repaired (§9.4) |

### 5.5 Awaiting a separate implementation authorization

1. Applying `CORPUS-64` (with whole-record re-review).
2. Applying `A-R2-01` (`CORPUS-65`) and regenerating `CORPUS-01`'s registry row.
3. Applying `CORPUS-66` **and** deciding the threshold — both needed for a scorable item.
4. Implementing `CORPUS-67`'s validator, with fixtures and its own verification record.
5. Repairing `F1-03`, `F6-03`, and the remaining ten ambiguous units.
6. Any registry row for `CORPUS-26` / `CORPUS-32`; any reconciliation of the registry summary tables.
7. Any `MEMORY-09` amendment. **None was prepared** — Decision 5 targets `CORPUS-01`, and `MEMORY-09` was not to be modified outside an authorized amendment preparation, none of which exists.

---

## 6. `MEMORY-09` and the governing standards

`MEMORY-09` was **read** and **not modified** (hash `f77e47a9…`, §9.5). Nothing in this record amends its acceptance criteria, its six-mode reporting, `§3.5`'s perceptibility requirement, `R-11`–`R-14`, or `§10`'s refusal to invent sample sizes and thresholds.

Two points where a reader might expect an amendment and will not find one:

- **`R-2`'s criterion.** `CORPUS-57`/`CORPUS-59` identified a criterion gap: `§3.5` routes *perceptibility* to `R-2`/`R-3`, and perceptibility is not dependence. Decision 5 does **not** authorize a `MEMORY-09` amendment. The gap remains recorded and unaddressed at the standard level.
- **No threshold or sample size** was set, implied, or estimated anywhere in this filing. `MEMORY-09` §10's procedure is untouched.

---

## 7. What this record does not claim

- It does not claim `F3-01`, `F6-01`, `F1-03`, or `F6-03` is repaired, adoptable as written, or diagnostic.
- It does not claim any item is Type P. The only diagnostic unit remains `F2-02` on adopted text, as previously reported and not re-audited here.
- It does not claim the checker validates registry targets. It does not, in either the old or the current state.
- It does not claim the registry is now complete. Two `docs/` files remain unregistered and the summary breakdowns remain internally inconsistent (§9.4).
- It does not claim the registry correction repaired the review. The review is unfiled and no correction could file it.
- It does not claim convergence, independence, or verified model identity for any artifact in this chain.

---

## 8. Deviations from the work order

Recorded plainly, with reasons. No conflict was silently resolved and no authorization was broadened.

| # | Deviation / judgment call | Reason | Reversible? |
|---|---|---|---|
| 1 | **Decision 1's registry correction was applied, not merely described.** | Decision 1 requires "an auditable registry correction", §1 lists "administrative corrections explicitly authorized here", and §12 asks for "administrative correction completed". Describing a correction in a proposal would leave the registry knowingly false and would put the correction in the "proposals not applied" category, where it does not belong. The change is one row, quoted before and after in `CORPUS-63` §3–§4. | Yes — the original row text is preserved verbatim in `CORPUS-63` §3 |
| 2 | Decision 1 required actions 1 and 2 were **verified, not performed** (no write to the brief). | Preservation is an absence of action; any write would have destroyed the hash identity every prior report cites. | N/A |
| 3 | **No reconstruction artifact was created** under Decision 2's conditional clause. | The source text is absent; creating one would require inventing metadata Decision 2 forbids inventing. Recorded as a disposition with a ready schema (`CORPUS-63` §5.1). | A later authorized reconstruction can use the schema |
| 4 | **Registry summary tables were incremented, not re-derived.** | They were already internally inconsistent before this filing; re-deriving them would be a broad, unauthorized registry audit. Only the deltas from six new rows were applied. | Yes |
| 5 | **Two pre-existing unregistered `docs/` files were not given rows.** | Not among the six artifacts; adding rows for them is outside this work order's scope. Recorded as a finding for the Decision 7 work. | Yes, when authorized |
| 6 | **No `.sha256` sidecars created.** | Matches the existing `CORPUS-*` convention (hashes live in the registry). | Yes |
| 7 | The primary record does **not** contain its own hash. | A file cannot contain its own SHA-256. Its hash is in registry row 397, as the registry convention already does for every other row. | N/A |
| 8 | `CORPUS-66` records an **additional** finding (`F6-01`'s `R-2` tag vs a warning-permitting threshold) that the work order did not ask for. | It surfaced while specifying the required changes to `Variation`/`Prohibited`, and omitting it would have produced a proposal that silently deepens an existing contradiction. It is recorded as unresolved; **no tag was changed**. | N/A — finding only |

---

## 9. Verification results

### 9.1 Artifacts exist

| Artifact | Path | Exists |
|---|---|---|
| `CORPUS-62` | `docs/research/corpus/CORPUS-62-…-2026-11-12.md` | **Yes** |
| `CORPUS-63` | `docs/research/corpus/CORPUS-63-CORPUS-58-REGISTRY-CORRECTION-AND-PROVENANCE-DISPOSITION-2026-11-12.md` | **Yes** |
| `CORPUS-64` | `docs/research/corpus/CORPUS-64-F3-01-MEMORY-DIAGNOSTICITY-REPAIR-PROPOSAL-2026-11-12.md` | **Yes** |
| `CORPUS-65` | `docs/research/corpus/CORPUS-65-CORPUS-01-R2-COVERAGE-AMENDMENT-PROPOSAL-2026-11-12.md` | **Yes** |
| `CORPUS-66` | `docs/research/corpus/CORPUS-66-F6-01-SCENARIO-EVIDENCE-REPAIR-PROPOSAL-2026-11-12.md` | **Yes** |
| `CORPUS-67` | `docs/research/corpus/CORPUS-67-REGISTRY-CHECKER-INTEGRITY-SPECIFICATION-2026-11-12.md` | **Yes** |

### 9.2 Hashes, bytes, lines (supporting artifacts)

Computed from the working tree after finalisation. Each is also in its registry row, where it can be re-verified independently.

| Identifier | Bytes | Lines | SHA-256 |
|---|---:|---:|---|
| `CORPUS-63` | 14,632 | 198 | `8f620f4689cae45b7fea2b6a9fb9c46c6133649d8fa90fb7baa8712e483bfcca` |
| `CORPUS-64` | 22,781 | 257 | `f16a97dd09fa55e7fffb4cbc09649aeafe72e8a6db949f9e0265f2ea373a0504` |
| `CORPUS-65` | 12,436 | 200 | `a914cb73db5a6efd94a5ef6f74b948bb9f6619cff3c8aa133ecbac247ebcbea0` |
| `CORPUS-66` | 16,555 | 208 | `398f3e7db756e9f38752eae5a6912bfc025cdd5954523633147969185ea81460` |
| `CORPUS-67` | 14,726 | 208 | `3c842eb318ae82e512718e4638891b0c640f4072ffc671307893549222eabd65` |
| `CORPUS-62` (this file) | recorded in registry row 397 | recorded in row 397 | recorded in row 397 — a file cannot hash itself |

### 9.3 Registry entries

| Check | Result |
|---|---|
| Rows added | **6** (397–402), one per artifact |
| Row 393 | **Corrected** — link and role; number, hash, bytes, lines preserved; original claim recorded in the row |
| Total rows | **402** (400 `docs/` files + 2 root-level files) |
| Duplicate row numbers | **None** |
| Duplicate targets | **None** — no artifact is registered twice |
| Registered targets that fail to resolve | **None** after the row-393 correction |
| Hash / byte / line mismatches among all 402 rows | **None** (full sweep) |

### 9.4 Missing targets and unregistered files

Full manual sweep of all 402 rows: every target resolves; every hash, byte count, and line count matches. **No missing targets remain.**

Two `docs/` files remain **unregistered**: `CORPUS-26-INDEPENDENT-REVIEWER-CORRECTED-REPORT-2026-11-12.md` and `CORPUS-32-THIRD-REVIEW-PACKET-D1-ADDENDUM-2026-11-12.md`. **Pre-existing**, discovered during this verification, **not repaired** (outside scope). The registry's §2.1 family and §2.2 class breakdowns also do not sum to their own stated totals; that drift is **pre-existing**, was **not re-derived**, and only the six-row delta was applied.

Both findings are direct evidence for Decision 7's premise and are recorded here rather than fixed.

### 9.5 Existing authorized mechanical checker

Command: `python docs/research/corpus/check_corpus.py` — run after all edits.

```text
scored units: 45
split: tuning=19 heldout=26
declared in CORPUS-01: {'units': 45, 'tuning': 19, 'heldout': 26}
TOTALS MATCH
RESULT: PASS (0 fail, 0 warn, 0 unverifiable)
This check reports only; no file was modified.
```

Exit status **0**.

**Limitations of that result, stated so it is not over-read:**

- It validates **item records**, not the registry. It never opens `docs/DOCUMENT-REGISTRY.md` and verifies **no** registry target, hash, byte count, or line count. **The PASS is not evidence that the registry is sound.** The registry was verified separately by the manual sweep in §9.3–§9.4.
- `0 warn` is **not** a clean bill: `C-2` is explicitly WARN-only and "no warning is NOT a pass"; semantic `Required`↔`Variation` conflicts are declared undecidable, which is why the `F3-01` conflict cannot be caught mechanically; `C-5` evidence sufficiency is not machine-checkable, which is why `F6-01`'s defect could not be detected.
- `0 unverifiable` is the per-unit count; the whole-corpus `UNVERIFIABLE` block (semantic conflicts, `C-5`, `C-2` heuristic) prints unconditionally.

### 9.6 Nothing else changed

| Artifact | SHA-256 after this filing | Unchanged? |
|---|---|---|
| `CORPUS-01` | `be1f11b46383e98a58093bc498ab4ca7321c33bd002861daed5a621dccdfbea8` | **Yes** — matches `CORPUS-61` §4 |
| `CORPUS-02` | `75aa9be707cc0c4406f349d626f7f5c46914472e56dd362eac563b83bcf40b04` | **Yes** — matches `CORPUS-59` §2 |
| `CORPUS-03` | `d0dae1c0951e9ed2058bd7db80d548356fe6b399f851a3dab1d5e74ed3b82a2c` | **Yes** — matches `CORPUS-59` §2 |
| `CORPUS-07` | `72f6e2a22b86bb588ba477bc6377bc723320861e90aadc2a5b4d247d8bd6f168` | **Yes** — matches `CORPUS-61` §4 |
| `MEMORY-09` | `f77e47a98685e66028e5aa23e5be45d7cda70b7c0cf5f9c443589b76f252f8da` | **Yes** — matches `CORPUS-61` §4 |
| `check_corpus.py` | `be3e317c6975dbed4f4bbbd68b698165b4dc5f1e887e5c56850f814ec14096cc` | **Yes** — matches `CORPUS-61` §4 |
| Root `CORPUS-58` brief | `292befdd5a4529cef79aa7cb10b7e7ea5aa5fe99e15cdb21f5c1205d47cfad56` | **Yes** — 4,311 bytes, 80 lines |

No corpus item, label, fixture, acceptance condition, classification, split, tag, threshold, sample size, schema, runtime, model, or checker implementation was modified. No model was loaded or run. No experiment, scoring, or annotation occurred.

---

## 10. Actions requiring further authorization

| # | Action | Needs |
|---|---|---|
| 1 | Adopt / amend / reject the `F3-01` repair (`CORPUS-64`) | Owner decision, then application under a change manifest, then whole-record independent re-review |
| 2 | Apply `A-R2-01` (`CORPUS-65`) | Owner decision; then `CORPUS-01` edit, revision bump, registry regeneration with recorded hash change |
| 3 | Decide `F6-01`'s behavioural threshold; apply the evidence repair (`CORPUS-66`) | **Two** decisions; both required before the item is scorable |
| 4 | Reconcile `F6-01`'s `R-2` tag with a warning-permitting threshold | Separate decision; no proposal made |
| 5 | Implement the registry validator (`CORPUS-67`) | Separate implementation authorization, fixtures, verification record |
| 6 | Repair `F1-03` (direction i or ii) | Separate authorization |
| 7 | Repair `F6-03` (consequence-based acceptance) | Separate authorization |
| 8 | The remaining ten `CORPUS-07` §2 ambiguous units | Separate authorization |
| 9 | Any `MEMORY-09` amendment (including the `R-2` criterion gap) | Separate authorization; **not prepared** |
| 10 | Furnish and file the substantive `CORPUS-58` review; commit the root brief | Owner decisions |
| 11 | Register `CORPUS-26` / `CORPUS-32`; re-derive the registry summary tables | Administrative authorization |
| 12 | Resume authoring toward the 80-item floor; seal; pilot; `G2` | Each separately gated; none touched here |

---

## 11. Provenance summary

| Layer | What is known | Status |
|---|---|---|
| Assigned model for **this** work order | Owner-supplied: **DeepSeek** | Owner metadata; **not verifiable** from the workspace |
| Executing agent | Codebuff coding agent; persona label "Buffy" | Execution fact; persona label ≠ model identity |
| Conversation | Fresh; no prior substantive turns | Execution fact; independence from the `CORPUS-58`–`61` conversations is **asserted by context, not cryptographically established** |
| `CORPUS-58` review authorship | Owner-supplied: Nemotron Lightning | Unverifiable |
| `CORPUS-59`/`60` authorship | In-file declaration: DeepSeek (Buffy) | In-file only; unverifiable |
| `CORPUS-61` reviewer | Owner-supplied: Solar Pro 4 | Unverifiable |
| `CORPUS-58` substantive text | Absent | **Unfiled, unfurnished, unverifiable** |
| Root brief | Present, hashed, untracked | Verified this session |
| Registry correction | Applied, quoted before/after | Verified this session |
| Convergence | **0 of 2** | Unchanged; this filing makes no convergence claim |

---

## 12. Statement of no changes

**Changed by this filing:** six new files (§9.1); six new registry rows (397–402); **registry row 393** (authorized administrative correction, `CORPUS-63` §4); the registry's declared scope, row-count, and summary-total lines, incremented consistently with those six rows.

**Not changed:** `CORPUS-01`, `CORPUS-02`, `CORPUS-03`, `CORPUS-04`, `CORPUS-05`, `CORPUS-06`, `CORPUS-07`, `CORPUS-56`, `CORPUS-57`, the root `CORPUS-58` brief, `CORPUS-59`, `CORPUS-60`, `CORPUS-61`, `MEMORY-09`, `check_corpus.py`, and every other existing artifact — verified by hash in §9.6.

`AM-01` unissued · `A-2a-R1` NOT IN FORCE · `D1` unadopted · annotation blocked.
`G1` open · corpus frozen · expansion frozen · convergence 0 of 2.

**Completion of this work order means the records and proposals are filed, versioned, registered, and verified — not that any proposed change has become operative.**
