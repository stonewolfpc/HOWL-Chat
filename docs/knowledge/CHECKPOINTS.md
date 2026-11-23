# HEKI — checkpoint log (append-only)

**Rule:** an entry is never edited after issue. A correction is a new entry naming the entry it corrects.
**Resume token format:** `HEKI-CP-YYYY-MM-DD-NN` — search the repository for the token.

---

## `HEKI-CP-2026-11-22-01`

**Stage:** 1 — Discover (in progress)
**Corpus:** 808 Markdown documents under `docs/` (excluding the registry itself), 527 resolved doc-to-doc
hyperlinks, 304 documents with no link in either direction, 11 undated by filename and by head.

**Completed this checkpoint**

1. HEKI charter and rules (`README.md`) — R1–R8 in force.
2. Full census: `ledger.jsonl` (machine record) and `DISCOVER-LEDGER.md` (greppable table), with
   `HEKI-D-NNNN` identities assigned newest-first.
3. Verified tranche 01: 7 documents, 25 typed edges, 6 explicit non-relations
   (`PHASE-1-TRANCHE-01.md`, `RELATIONS.md`).
4. Registry-coverage audit of unregistered material: 322 unregistered `.md` (313 `docs/research/`, 7
   `docs/architecture/`, 1 review, 1 root). Classification and recommendations recorded in the conversation
   record; **nothing was registered.** Owner decision still open on whether to register the 9 recommended
   documents.

**Established (not assumed)**

- The registry has **no dangling rows** — 0 rows point at a missing file.
- Registry `Rows: 444` was a stale creation-time count; corrected to the measured 516 with the coverage
  reality stated. `423`/`424` each appear twice and were recorded, not repaired.
- Filesystem birth and mtime are **useless for provenance** here: all 808 files report 2026-11, while content
  dates reach back to 2026-11-11. Timestamps were not used to classify anything.

**Open findings carried forward**

| # | Finding | Status |
|---|---|---|
| F1 | Corpus cross-references by identifier string, not hyperlink; census misses those edges | **Blocks Tranche 02 quality** — add identifier-mention edges first |
| F2 | Scope authority is split: 27 original items vs A–E evidence categories (`HFS-06`) | Phase 4 |
| F3 | T0.5 contradiction: HFS-AUDIT preserves it; V2 dropped it; `asr-repair-design.md:10` calls it open | Phase 4; do not reconcile by assumption |
| F4 | "Complete" has a scope meaning and an E3 meaning | Must be labelled in every readiness statement |
| F5 | `DEF-FE-01` deliberately outside HF1-PLAN | Must not be lost |
| F6 | 9 unregistered documents are cited by the registered chain but locatable only by moving them | Owner decision; nothing registered |

**Next**

Tranche 02: `contracts/README.md`, FC-1…FC-6, and the ten 2026-11-21 `research/MEMORY-09-*` records.
Precondition: identifier-mention edges added to the census.

**Resume:** read this entry → open `DISCOVER-LEDGER.md` → continue from the first node whose ID appears in
neither `PHASE-1-TRANCHE-01.md` nor a later tranche file.

**Nothing outside `docs/knowledge/` was modified by this checkpoint.**

---

## `HEKI-CP-2026-11-22-02`

**Stage:** 1 — Discover (in progress) · **extends:** CP-01; corrects it only where stated below.

### Completed

1. **Identifier-mention layer** (safeguard 1). `mentions.jsonl`: one row per occurrence — document, identifier,
   line number, full source line. Mentions are **never** promoted to relationships.
2. **Census boundary settled** (safeguard 2). `README.md` §8: index artifacts are **excluded from the corpus by
   design**. Corpus = 808; index = 6; all `.md` under `docs/` = 814. No identity renumbered.
3. **Status separation** (added instruction). `README.md` §10, rule **R9**: as-written status is immutable
   against later references; current status is a separate judgement recorded only in `RELATIONS.md`.
4. **Tranche 02** (`PHASE-1-TRANCHE-02.md`): 13 documents — contracts register, FC-1…FC-6, six 2026-11-21
   MEMORY-09 research records. Cumulative verified: **20 of 808 (2.5%)**.

### Measured (not asserted)

| Metric | Value |
|---|---|
| Corpus nodes | 808 (unchanged) |
| Mention occurrences | 788 across 281 distinct program identifiers |
| Documents containing ≥1 identifier mention | **45 of 808** |
| Filename fragments rejected by the tightened regex | 72 occurrences / 13 distinct |
| Hyperlink-isolated documents | 304 |
| Of those, carrying identifier mentions | **15 (5%)** |
| Production documents modified | **0** |

**Regex correction:** the first mention pass captured filename suffixes (`CONTRACT-2026-11-11`,
`AMENDMENT-02-…-ROUTING-2026-11-11`) as identifiers. A document-noun blacklist plus a `(?!\.md)` guard removed
them; the rejected strings are counted above rather than silently dropped.

### Findings

| # | Finding | Status |
|---|---|---|
| T2-F1 | **R6 application status.** R6 says NOT APPLIED; the register says applied; FC-5 contains the R6-only clauses | **Resolved by measurement:** FC-5 = 23,056 B / 304 L / `626ab4be…` vs R6's recorded target 12,326 B / 213 L / `c245096f…`; registry already records the amended hash → **applied; R6's status block is stale** (erratum candidate, not edited under R1) |
| T2-F2 | Review-accepted ≠ approved: 6 papers, 0 approvals, 0 implementation authorizations | Recorded |
| T2-F3 | Memory-record authority is a deliberate hole (`FC4-04a`); FC-7 blocked by D8 | Blocks MEMORY-09 governance |
| T2-F4 | The G1 audit is orphaned (in=0, out=0, 0 mentions) while downstream documents cite its gate | Phase 2 candidate |
| T2-F5 | Identifier layer is contract-internal: 45 of 808 documents; 111 mentions in FC-1 alone | Adjusts the layer's expected value |
| T2-F6 | `FC5-04` already separates D6 from T0.5 and from global accounting | Guard retained |
| T2-F7 | `FC4-04` assigns Character to the canonical library; bears on ADR-001 §11.2 | Phase 4 comparison |

### Corrections to earlier statements

- Tranche 01 said "the ten 2026-11-21 `research/MEMORY-09-*` records". There are **six**. Corrected in
  Tranche 02; CP-01 left unedited per the append-only rule.
- CP-01 expected the identifier layer to materially reduce the 304 isolated documents. Measured: **15 (5%)**.
  The expectation was wrong and is recorded as wrong.
- Tranche 01 §3 said the R6 draft "cannot govern" because its target hash was unchanged. The hash *had*
  changed: the amendment was applied. Revised in both files with the measurement shown.

**Next:** Tranche 03 — the 2026-11-20 foundation pack (rows 505–512), then 2026-11-19 and earlier.
**Resume token:** `HEKI-CP-2026-11-22-02`

**Nothing outside `docs/knowledge/` was modified by this checkpoint.**

---

## `HEKI-CP-2026-11-22-03`

**Stage:** 1 — Discover (in progress) · **extends:** CP-02. Corrections to earlier entries are stated as
corrections to those entries, which are left unedited.

### Completed

1. **Tranche 03** (`PHASE-1-TRANCHE-03.md`): the 2026-11-20 foundation pack — `HEKI-D-0088` (front-to-back
   plan) and `HEKI-D-0089` (contract-pack intake). Cumulative verified: **22 of 808 (2.7%)**.
2. **Scope correction:** CP-02 announced Tranche 03 as "rows 505–512". Measured: rows 505–511 are 2026-11-21
   contract papers already verified in Tranche 02. The 2026-11-20 architecture population is **2 documents**.
3. **Third connectivity layer** (`filerefs.jsonl`): filename references, same shape and same non-promotion rule
   as `mentions.jsonl`. `README.md` §9 now carries three layers, not two.
4. **`ERRATA-CANDIDATES.md`** created: `E-01` R6 application, `E-02` filename-vs-content dates, `E-03` pointer to
   corrections already made. No original edited.

### Measured (not asserted)

| Metric | Value |
|---|---|
| Corpus nodes | 808 (unchanged) |
| Filename-reference occurrences | 2,642 |
| Distinct (citing document → target) pairs | 2,120 |
| Documents citing another by filename | 490 of 808 |
| Those references also hyperlinked | 42 of 2,642 (1.6%) |
| Hyperlink-isolated documents | 304 |
| …referenced by filename in one direction or both | 172 |
| …unreferenced in **all three** layers | **129 (16%)** |
| Documents whose head date post-dates their filename date | 15 of 622 |
| Production documents modified | **0** |

### R6 — now settled at provision level, not only by byte count

Baseline `0bc3b091…` holds the pre-amendment FC-5 at exactly R6's recorded `c245096f…` / 213 L / 12,326 B.
Diff to the current file: **+91 content lines, −0**, in 4 hunks. R6's 48 quoted content lines are **all 48
present verbatim**, and the added lines equal them **in order, exactly**. All eight provisions are in R6's
stated placement, `FC5-01`–`FC5-26` are untouched, and both revision-6 sentences in `FC5-29` are present.

Still **not** established, and recorded so they cannot be rounded up: FC-5 remains `PROPOSED — not approved`;
`FC5-27`–`FC5-31` and `FC5-26a` keep their `Proposed`/`Unresolved` labels; `FC5-CONFIG-RETENTION` is open;
M4 unchanged; no implementation. R6's "NOT APPLIED" is a **stale statement**, now with provision-level evidence
behind the finding (`ERRATA-CANDIDATES.md` E-01). No document records *who* applied it or *when*.

### Findings

| # | Finding | Status |
|---|---|---|
| T3-F1 | "304 disconnected" overstated isolation by 175; the honest figure is **129** | CP-01's number stands as written, was being read too strongly; recorded |
| T3-F2 | Filename references are a real third layer — 2,642 occurrences, ~98% unlinked | Materialised; occurrence index only |
| T3-F3 | Filename date ≠ content date; HEKI ordering is a **filename ordering** | Structural; stated in `README.md` §8 |
| T3-F4 | The parent plan that commissions the contract pack has 0/0 links and 0 mentions | Strongest Phase 2 candidate so far |
| T3-F5 | 23 `DQ-C/X/R/M` identifiers appear only in the plan; no crosswalk to `DQ-HRA-*` | Not an edge; do not read silence as closure |
| T3-F6 | The plan and `HRA-04` agree on modalities; HCDM was the outlier | Observation only, not adjudicated |
| T3-F7 | No implementation authorization anywhere in this tranche, in any tense | Recorded |

### Corrections to earlier statements

- CP-02's Tranche 03 scope ("rows 505–512") was wrong; the rows are 2026-11-21 papers except the intake.
- CP-01's expectation that identifier mentions would reconnect the isolated corpus was already recorded as
  wrong (15 of 304). The larger correction is that the 304 figure itself measured the wrong thing.

**Next:** Tranche 04 — the 62 `research/MEMORY-09-*` documents filename-dated 2026-11-20 (1,065,774 bytes):
the DQ-2/DQ-3/DQ-5 closures, P-1 and S-2/S-3 outcomes, and the implementation-authorization chain. Then
2026-11-19.
**Resume token:** `HEKI-CP-2026-11-22-03`

**Owner decisions still open, unchanged:** R6 erratum (now with provision-level evidence); registration of the
9 documents — which now includes `HEKI-D-0088`, the unregistered parent plan; Phase 2 link insertion.

**Nothing outside `docs/knowledge/` was modified by this checkpoint.**

---

## `HEKI-CP-2026-11-22-04` — strategy change: implementation becomes primary

**Stage:** 1 — Discover, **paused** · **extends:** CP-03. Owner direction, 2026-11-22.

### What changed

Chronological discovery stops at CP-03 (22 of 808 verified) and HEKI becomes the **supporting** workstream.
The primary workstream is foundational implementation, planned in
`docs/architecture/HOWL-FOUNDATIONAL-IMPLEMENTATION-EXECUTION-PLAN-2026-11-22.md` (`HFI-EP-2026-11-22-01`,
registry row 517).

- **R7 superseded as a gate.** It survives as a principle: a task names the evidence it relies on and may
  proceed when that evidence — not the whole corpus — is present.
- **Tranche 04 deferred, not cancelled.** The 62 `research/MEMORY-09-*` records dated 2026-11-20 remain
  unread and unindexed-as-verified. No identity renumbered; the ledger remains a valid census.
- Nothing in the corpus was edited to make this possible. The pivot is recorded, not enacted on originals.

### Why the index itself supports the pivot

Two failure modes it caught were not visible to a link graph: a document contradicting the state of the file it
describes (`ERRATA-CANDIDATES.md` E-01), and reachability that hyperlinks systematically miss (E-02, T3-F1).
Both are exactly what an implementation team needs when it hits an unresolved question — which is the role HEKI
now has.

### Census consequence, measured

| Metric | Before CP-04 | Now |
|---|---|---|
| `.md` under `docs/` | 816 | **817** |
| Index artifacts (`docs/knowledge/`) | 8 | 8 |
| Corpus | 808 | **809** |
| Ledger nodes | 808 | **808** |
| Registry rows | 516 | **517** |

One corpus document is therefore **unindexed**: the execution plan itself. It was created after the census and
before discovery resumed. It takes `HEKI-D-0809` when indexing continues; **no existing identity was renumbered**,
so "newest-first" now describes the original census rather than the live tree. Recorded so the gap is visible
rather than absorbed.

### The ownership findings the plan rests on

The architecture branch already contains the spine: ADR-001 §4 gives six services and two primitives with
explicit `Owns` / `Must not own` columns; ADR-002 §2 gives three state families with named owners; FC-1 separates
ledger, library and snapshot. The plan reuses these and adds nothing architectural.

Four conflicts and four missing owners were identified, none resolved:

- **C-1** Character: FC4-04 assigns it to the canonical library; ADR-001 §11.2 says preserve in place and do not
  merge into the library system.
- **C-2** Conversation state: five candidate implementations, no owner.
- **C-3** Reservation authority: two partial implementations, `DQ-HRA-01` open.
- **C-4** Configuration: four overlapping authorities with temporal precedence.
- **M-1** Context composition · **M-2** speaker isolation · **M-3** evidence/receipts · **M-4** memory records
  (deliberately absent, `FC4-04a`, FC-7 blocked by D8).

### The first task, and what blocks it

`TASK-01` — PathService substrate (ADR-001 contracts A, B, D half) — is bounded with eight observable acceptance
criteria. **It has zero blocking questions**, which is why it is first. It is **prepared, not authorized**: no
repository change occurs before the owner approves the boundary.

Phase 1 also splits into 1A (PathService + DurableStore), 1B (conversation identity and ledger), 1C (composition,
tokenizer boundary, speaker isolation) — because stable identity and recovery both require a trustworthy place to
write. That split **changes the stated first target** and is flagged for approval.

**Next owner decision:** approve or amend `TASK-01`.
**Resume token:** `HEKI-CP-2026-11-22-04`

---

## `HEKI-CP-2026-11-22-05` — the one corpus file the pivot modified

**Completes:** `HEKI-CP-2026-11-22-04`, which recorded the census consequence but not the resulting hash change.
CP-04 is left unedited (append-only).

Registering the execution plan required a registry row. The document registry's own maintenance rule —
"add a row for every new documentation file in the same change that adds the file" — makes that edit mandatory,
and the registry is a **live index**, not a frozen artifact. It is also a census node, so the census is now stale
for exactly one file.

| Node | Path | Census hash | Current hash | Bytes |
|---|---|---|---|---|
| `HEKI-D-0003` | `docs/DOCUMENT-REGISTRY.md` | `584dc186917b…` | `93be9664ddad…` | 269,201 → 270,007 |

**Measured: 1 of 808 corpus files differs from the census, and it is this one.** The other 807 are byte-identical.
No frozen artifact, contract, ADR or research record was touched. The registry's `Rows:` count moved 516 → 517 and
its file counts were corrected to the measured values (808 `.md` besides the registry, 860 files total, 23
sidecars).

`ledger.jsonl` was **not** regenerated: regenerating would change 800+ recorded hashes and dates to no purpose.
The stale `HEKI-D-0003` hash stays visible until discovery resumes, at which point the plan takes `HEKI-D-0809` and
this row is refreshed in one deliberate act.

**Resume token:** `HEKI-CP-2026-11-22-05`

---

## `HEKI-CP-2026-11-22-06` — TASK-01 implemented (first change to product code)

**Stage:** supporting (HEKI) · **implements:** `TASK-01` under `HFI-EP-2026-11-22-01` · **owner-authorized**
This is the first checkpoint in this program where the tree changed outside `docs/`.

### Authorized scope, as executed

PathService substrate only: root set, mode, category→location, class labelling, containment by identity,
external-root classification, typed failures. **Not** DurableStore, **not** migration, **not** call-site rewiring.
No existing `.go` file was modified and no storage location moved. The running app's behaviour is unchanged.

### Delivered

| File | Lines | Hash (first 12) |
|---|---:|---|
| `internal/paths/roots.go` | 385 | `afb97c13e07e` |
| `internal/paths/identity_windows.go` | 47 | `ec471a9116d9` |
| `internal/paths/identity_other.go` | 32 | `3f42e639f14e` |
| `internal/paths/roots_test.go` | 363 | `739f117ba7c6` |

Checks: `gofmt -l` clean · `go vet ./internal/paths/` clean · `go test ./internal/paths/` all pass (27 run
entries, 13 tests added) · `go build ./...` passes · `git diff --stat internal/paths/` empty.

### What the census changed

ADR-001 reads the defect as "three resolution rules". The call-site census shows something sharper: **21 of the 22
`ProjectRoot()` callers use it for read-only resource resolution** (`runtime/`, `models/`, kokoro roots) and are
correct. The storage error is one decision inside `StateRoot()` — it defaults to `ProjectRoot()`. That is a
single repair behind a migration, not 21 scattered bugs. `ProjectRoot` is retained and reclassified; it must
never be a storage authority.

### Findings recorded, not fixed

- `startup_log.go:24–27` writes `config/logs` **relative to the working directory** in normal mode, and
  `StateRoot/config/logs` under isolation — the same log lands in different places depending on launch directory.
- `audio/pipeline.go` `defaultTTSCacheDir` returns cwd-relative `"cache/tts"` when isolation is off.
- Whisper, `app-tts-catalog.go` and `voiceModelPath` build model paths from the working directory.
- `newIsolation` creates directories during resolution, although PathService is declared to have no writers.

All four belong to the call-site migration task. Fixing them here would have been the "quiet storage rebuild" the
owner ruled out.

### Stated limitations (not implied away)

No junction-escape test (needs privileges), no redirected-profile run, and the "three separate launches"
criterion was proven by `chdir` within one process rather than by three processes. All three are listed in the
task record §3.

### Next owner action

Audit `docs/architecture/TASK-01-PATHSERVICE-OWNERSHIP-AND-COMPLETION-2026-11-22.md`. On approval, `TASK-02` is
DurableStore — bounded to atomic replacement, the single-writer lock, pre-image retention and the four outcome
states, with ADR-001 §8.1's reader-sharing question (`UN`) resolved first.

**Resume token:** `HEKI-CP-2026-11-22-06`

---

## `HEKI-CP-2026-11-22-07` — TASK-01 audited, corrected, closed

**Stage:** supporting · **closes:** `TASK-01` · **does not open:** `TASK-02`

### The audit found a defect, which is the point of having one

The owner's first scrutiny area asked whether identity-based containment survives junctions. It did not.
The first implementation walked lexical ancestors looking for the root's identity; a junction inside the root makes
lexical and physical ancestry diverge, so a path physically outside the root reported **contained**.

Corrected to resolve physical paths (`GetFinalPathNameByHandle`) and compare those, failing closed when a path
cannot be resolved. A second trap was found while fixing it: **`filepath.EvalSymlinks` does not follow junctions on
go1.27.3/windows-amd64**, so the first correction still failed. Recorded as I-09.

| Test | Before | After |
|---|---|---|
| Junction escape (real `mklink /J`, no elevation needed) | **FAIL — defect** | **PASS** |
| Cross-volume candidate (8 volumes present) | not tested | **PASS** |
| Three separate processes, three working directories | `chdir` in one process | **PASS — byte-identical stdout** |
| Isolation creation (I-02) | asserted not to happen | **pinned by test** |

The third limitation — a redirected Documents profile — remains **unproven** and is accepted as such (L-1). It is
not reported as demonstrated.

### A defect in my own documentation

`Roots()` claimed "no filesystem creation". False under isolation, where the existing resolver creates its layout.
The comment now states the installed-mode guarantee and names the isolation exception. Pinning the behaviour by
test means changing it later is deliberate rather than accidental.

### Corrections applied

F-A1 containment algorithm · I-09 physical resolution · D-1 corrected doc claim · T-1 three tests added.
Artifacts as audited: `roots.go` 387 L `71f5ed0e…`, `identity_windows.go` 98 L `03fd0d44…`,
`identity_other.go` 42 L `be1dab43…`, `roots_test.go` 397 L `2409216d…`, `roots_windows_test.go` 99 L `ed1f0c8e…`.
`gofmt`/`go vet` clean, `go test ./internal/paths/` all pass (16 test functions), `go build ./...` passes,
no pre-existing file modified.

### Decision

**TASK-01 accepted and closed.** Conditions carried to `TASK-02`: C-1 `Writable` is classification, never
authorization; C-2 writes need handle-relative operations for the TOCTOU guarantee `CheckContained` does not give;
C-3 never use `EvalSymlinks` for containment here.

### Census bookkeeping, exactly as the owner required

`HEKI-D-0809`/`0810` are **reserved placeholders only**. Neither has a ledger node, hash, status extraction or
link count. The census remains **808**; the corpus is **811**; the gap is a **known, enumerated 3**, not a
reconciled total. Quoting 811 as a census figure would be wrong.

### Next

`TASK-02` — DurableStore — **not started**. Its first step is ADR-001 §8.1's reader-sharing question, which is
still `UN`. Awaiting the owner opening that task.

**Resume token:** `HEKI-CP-2026-11-22-07`

---

## `HEKI-CP-2026-11-22-08` — TASK-01 closure recorded

**Stage:** supporting · **closes:** `TASK-01` · authorized by the owner 2026-11-22

The audit's evidence is preserved as closure evidence in `TASK-01-AUDIT-2026-11-22.md` §8: the original
junction defect and why lexical ancestry produced a false containment result; the `GetFinalPathNameByHandle`
correction; I-09 (`EvalSymlinks` does not follow junctions on go1.27.3/windows-amd64); the passing junction,
cross-volume, traversal and legitimate-child tests; the literal three-process launch with byte-identical output;
limitation L-1 (redirected profile unproven); constraints C-1/C-2/C-3; clean `gofmt`/`go vet`/tests/`go build ./...`;
and no pre-existing Go file modified.

No second hash registry was created — the existing registry remains authoritative, and row 519 was updated in
place to the closure-edited hash after the edit, with the previous hash recorded in this checkpoint trail.

### Census preserved exactly

**808** historical census. **813** current corpus. Five reserved placeholders (**HEKI-D-0809 … 0813**), none of
which is a registered ledger node or has a recorded hash. **813 is not the census.**

**Resume token:** `HEKI-CP-2026-11-22-08`

---

## `HEKI-CP-2026-11-22-09` — TASK-02 DurableStore

**Stage:** foundational · **opens:** `TASK-02` · **status:** primitive complete, **not adopted**

### Decided (Step 1)

ADR-001 §8.1's `UN` reader-sharing question is resolved as **short-lived read + bounded writer retry**, on
measurement rather than preference: with a reader holding the document via `os.Open`, `os.Rename` returns
`Access is denied` while `os.WriteFile` **succeeds**. Go's `os.Open` shares read and write but not delete. So a held
reader blocks replacement, and a direct write is observably torn — the current `settings.Save` is not merely
non-atomic in theory.

Rejected: delete-sharing readers (needs `NtCreateFile`, capability `UN`, rewrites every call site); shared-lock
readers (turns reads into failures); snapshot pointers (a second "which document is live" authority).
**No ADR amendment is required**; the decision supplies a value without changing a normative statement.

### Implemented (Steps 2–4)

`internal/durable` — `store.go` 606 L, `lock_windows.go` 62 L, `lock_other.go` 45 L, `store_test.go` 430 L.
Order of a commit: lock → sweep dead temps → temp write + flush → pre-image rotation → manifest commit → rename
with 3 × 50 ms retry → trim → release. The live path is never opened for writing; its only mutation is a rename.

Two defects were found by the tests and fixed, not relaxed: manifest replacement needed the same bounded retry as
the document, and a reader's document+manifest read pair could straddle a commit and look like corruption (now
bounded re-read, with rollback distinguished from corruption by reason).

### Verified (Step 5)

`gofmt` clean · `go vet` clean · `go test ./internal/durable/ -race -count=2 -v` → **11 PASS, 2 helper SKIP** ·
`-count=3` pass · `go build ./...` passes. Evidence includes **cross-process** lock exclusion (`Busy`, then
success after release), a real child-process death between manifest commit and rename leaving the previous document
readable, continuous concurrent reads never observing a partial document, and rotation/recovery/schema refusals.

### Unproven

**Power-loss durability remains `UN`** — flush ordering exists, NTFS metadata durability does not. Network
filesystems unexercised. Reader long-holds are surfaced, not enforced. Crash injection covers one point only.

### Deferred

Adoption of the store by any call site — `gallery.writeFileAtomic`, `chat.saveJSON` and the direct `settings.Save`
remain, so **three replacement strategies still coexist** and this primitive is **not yet authoritative**. Also
deferred: the 20+ raw `os.WriteFile` sites, and `settings.Save`'s falsified rationale comment.

### Next step

Adoption starts at `internal/settings` — **not started, not authorized**.

**Resume token:** `HEKI-CP-2026-11-22-09`

---

## `HEKI-CP-2026-11-22-10` — TASK-02 self-audit

**Stage:** foundational · **TASK-02:** `implemented → audited` · **accepted: no** (acceptance is the owner's)

### Mandate and outcome

Bounded self-audit against the owner's four points. **One blocking defect found, corrected and verified (F-A2).**
No other code changes; adoption remains deferred and unauthorized.

### F-A2 — a failed commit destroyed the newest retained generation

Rotation wrote `v1`…`vN` **before** the rename. When a reader denied replacement, `restoreManifest` restored a
manifest whose hashes no longer matched the rewritten files: `v1` held the current live document while the manifest
said it held the previous one. The newest real generation became unreachable, and recovery on a later corruption
returned `corrupt_no_valid_preimage` after a **routine** reader-contention event.

Reproduced before the fix: `audit_test.go:44: after failed commit, v1 = "d3", want the untouched d2`.

Corrected by staging (`.staged`) and promoting only after a successful replacement; a failed commit aborts staging
and touches nothing retained. `TestFailedCommitPreservesPreImageChain` (the reproduction) and
`TestFailedCommitLeavesNoStagedArtifacts` (cleanup) now pass. The task record's ordering description and defect count
were corrected; its claim that "a failed commit never rotates" was false as delivered.

### Point 1 — lock lifecycle: measured

| Step | Result |
|---|---|
| Probe while another live process holds the lock | `blocked (… another process has locked a portion of the file.)` |
| Probe after that process exits; lock file still present | `acquired` |
| Child dies mid-commit (`os.Exit(9)` while holding the lock) | exit `9` |
| Probe after that death | `acquired` |

A crashed writer cannot wedge the application, and a leftover lock **file** cannot admit a second writer. No code
path deletes the lock file, so the double-writer hole cannot be manufactured. Residual: no lease — a *hung* process
yields `Busy`, which is required D4 behaviour, not a bypass.

### Point 2 — commit and recovery protocol

Ordering: lock → sweep dead temps/staged → temp write + `Sync` → **stage** → manifest commit → [injection] →
rename → **promote** → trim → release.
**One interruption point is tested** (after manifest, before rename) with a real process death; five are reasoned,
not measured: after temp write, after staging, after rename before promote, during promote, during trim. Recovery
restores a complete hash-verified generation, or refuses — never partial bytes, never an empty document.

### Point 3 — reader consistency

Bounded re-read of both files (3 × 2 ms) on mismatch. Exhaustion yields, in order: rollback to the live bytes when
they equal a retained generation; a matching pre-image; or a refusal. Torn bytes are never returned. Stated
degradation: under commits faster than the retry window a reader can return the **previous** generation marked
`Recovered` — one behind, never mixed.

### Point 4 — acceptance evidence

Process-crash recovery is proven by a **real OS process death** (`exec.Command`, `os.Exit(9)` inside `Save`), not by
an in-process simulation; only the *timing* is injected. The two helper SKIPs are child entry points, verified to
exit `9` and `0` when driven directly — a skip means "runs in a child", not "does not run".
**Process crash ≠ power loss:** power-loss durability remains **`UN`** and is not claimed.

### Verification after the fix

`gofmt` clean · `go vet` clean · `go test ./internal/durable/ -race -count=2 -v` → **13 PASS, 2 helper SKIP** ·
`go build ./...` passes.

### Registry verification retained (actual results)

```
row 517 OK
row 518 OK
row 519 OK
row 520 OK
row 521 OK
row 522 OK
```

### Census

Historical census **808**. Current corpus **814**. Reserved placeholders `HEKI-D-0809`–`0814` are **not registered
ledger nodes and have no recorded hashes**. 814 is the corpus, not the census.

### Next step

Owner acceptance of the corrected TASK-02 code, then adoption at `internal/settings` — **not started, not
authorized**.

**Resume token:** `HEKI-CP-2026-11-22-10`

---

## `HEKI-CP-2026-11-22-11` — branch completion (TASK-01 + TASK-02)

**Stage:** foundational · **branch state: complete** · **next branch: not selected**

### Task states

| Task | State | Closure |
|---|---|---|
| `TASK-01` PathService | **Closed** | `TASK-01-AUDIT-2026-11-22.md` §8; `HEKI-CP-2026-11-22-08` |
| `TASK-02` DurableStore | **Closed — owner-accepted for its bounded objective** | `TASK-02-DURABLESTORE-2026-11-22.md` §8; audit `TASK-02-AUDIT-2026-11-22.md`; `HEKI-CP-2026-11-22-09`, `-10` |

### F-A2 preserved as a corrected blocking defect

Recorded in three places, none edited to remove it: the audit's full analysis (§2.2), this trail's pre-audit claims in
`HEKI-CP-2026-11-22-09` — **including the sentence F-A2 falsified, "a failed commit never rotates"** — and the task
record's original ordering sentence retained in place beneath a visible correction note.

- **Reproduction (against the delivered code):** `audit_test.go:44: after failed commit, v1 = "d3", want the untouched d2`
- **Regression:** `TestFailedCommitPreservesPreImageChain`, `TestFailedCommitLeavesNoStagedArtifacts`
- **Effect it had:** a routine reader-contention failure made recovery report `corrupt_no_valid_preimage` on a chain
  that was intact seconds earlier.

### Accepted limitations (carried forward, not resolved)

1. **Power-loss durability `UN`** — flush ordering exists; rename durability across a power cut is not demonstrated.
2. **Five interruption points have reasoned coverage, not crash-test coverage:** after temp write; after staging;
   after rename before promote; during promote; during trim. Only manifest-commit→rename is injected with a real
   process death.
3. **Reader contention may return a previous, complete generation marked `Recovered`** — one behind under sustained
   commits, never torn.
4. **Existing persistence implementations not migrated** — `settings.Save`, `chat.saveJSON`,
   `gallery.writeFileAtomic` remain; three replacement strategies coexist.
5. **Primitive, not the persistence authority** — no call site uses the store; no application behaviour changed.

### Acceptance condition on reader consistency — satisfied

`Recovered` is a distinct `Outcome`, explained by `Reason`/`From`, and `Result.Current()` is true only for a `Saved`
(current) read, so a caller that must not act on a superseded document has a predicate that cannot be confused with
`OK()`. Pinned by `TestRecoveredIsDistinguishableFromCurrentRead`. This addition post-dates the audit, which is why
the audit record says 13 tests and the closed record says 14; **the audit record was not updated to match.**

### Verification at branch completion

`gofmt` clean · `go vet ./internal/paths/ ./internal/durable/` clean ·
`go test ./internal/paths/ ./internal/durable/ -count=1` → both packages ok ·
`go test ./internal/durable/ -race -count=2 -v` → **14 PASS, 2 helper SKIP** · `go build ./...` passes ·
**no pre-existing Go file modified** (all new files untracked).

### Registry verification after the final edits (actual results)

```
row 517 OK  e2e80d2871af
row 518 OK  8bb5f657988c
row 519 OK  135f23716bb4
row 520 OK  d63f5d74f039
row 521 OK  8ab037fb2383
row 522 OK  8426fe571313
```

Six rows, six matches against the files on disk.

### Census

Historical census **808**. Current corpus **814**. Reserved placeholders `HEKI-D-0809`–`0814` are **not registered
ledger nodes and have no recorded hashes**. 814 is the corpus, not the census.

### Next

**Adoption is not authorized and was not started.** The next branch is selected from
`HOWL-FOUNDATIONAL-IMPLEMENTATION-EXECUTION-PLAN-2026-11-22.md` now that this branch is checkpointed — selection is
not made here.

**Resume token:** `HEKI-CP-2026-11-22-11`

---

## `HEKI-CP-2026-11-22-12` — TASK-03 prepared (implementation unauthorized)

**Stage:** foundational · **TASK-03:** `prepared` · **implementation: not authorized**

### Plan corrections applied (history preserved, nothing rewritten)

Six owner-approved corrections were **appended** to the execution plan, each naming the line it supersedes: the
ledger is `TASK-03` (not "`TASK-02` (ledger)"); §10's budget-proof task becomes `TASK-04`; `TASK-02` is closed; the
corpus is 814 at the time of that edit; §12's "0 of 808 files differ" described the tree before 1A; and 1A is
complete as `TASK-01` + `TASK-02`.

### What was prepared

`TASK-03-CONVERSATION-AUTHORITY-PREPARATION-2026-11-22.md` — objective, seven-question ownership contract, verified
inventory, four required decisions, three assessments, eleven acceptance criteria, handoff.

### The five implementations, verified

`a.chatMessages` (`app-chat.go:53,126,184`) · `chat.Service`'s own context (`service.go:52–75`) · `history.json`
(`persistence.go:211`, read `:515`) · `history.ndjson` (`O_APPEND` at `:678,696`, **plus** `RewriteHistory` at
`:631`) · frontend `localStorage`, including selection in `howl-chat-current-chat` (`app-wipe.js:16–24`).

Two facts carried forward: the append log already has **two** mutation paths, and `persistence.go:210` records that
DOM scraping has corrupted `history.json`. The branch must retire four of the five — the preparation states that
explicitly so the ledger cannot become a sixth.

### Assessment A — FC-1 triage

Compatible as written: FC1-01, 03, 05/06, 09, 10, 13a's invariant, 15, 17, 22/22a/22b, 24, 26.
Needs clarification: FC1-02 (graph now or later), FC1-08 (mandatory envelope fields), FC1-11 (Phase 1 event kinds),
FC1-22c/`FC1-STATE-RETENTION` (is the ledger deleted with its conversation), FC1-25.
Unresolved and blocking: FC1-12/`FC1-CONCURRENCY`, `FC1-SELECTION-AFTER-DELETION`, FC1-14b/14c. **FC-1 remains
PROPOSED; implementing a clause does not approve it.**

### Assessment B — DurableStore suitability: answer not assumed

Three options weighed (as-is / narrow segmented extension / separate journaling primitive). The preparation proposes
a small journaling primitive composed with DurableStore — matching ADR-001 §4's explicit separation of document
persistence from journaling — but **refuses to decide it** until one measurement exists: events per conversation
(p50/p95/max) and the interval between turns. That measurement is the only input that changes the answer. Four
properties must be decided regardless: record framing and checksum, flush policy per commit, the torn-tail
truncation rule (truncate to the last valid record and record a recovery event), and duplicate suppression when an
acknowledgement is lost. **Power-loss durability stays `UN`.**

### Assessment C — ContainmentChecker disposition

No second containment implementation. `internal/paths` already answers *where* and *is this inside* (TASK-01);
the gap versus ADR-001 §7 is component-level reparse inspection, pinned directory identity, post-issuance reparse
rejection and handle-relative creation (capability `UN`). `hasReparseComponent` already exists **four times**
(`audio/tts/helper/pathcheck_windows.go:16`, `pathcheck_other.go:14`, and the same pair in `fakehelper`). Proposal:
extend `internal/paths` additively when ExternalDataGuard work begins and retire the four copies then. **TASK-03 adds
no containment code**; the ledger calls `RootSet.Check`.

### Authorization state

Not authorized: implementation, migration of settings/chat/gallery, call-site changes, contract approval, branching
work, memory authority. Authorized and done: plan bookkeeping, preparation, three assessments.

### Registry verification after the final edits (actual results)

```
row 517 OK  e20d187df9d7
row 518 OK  8bb5f657988c
row 519 OK  135f23716bb4
row 520 OK  d63f5d74f039
row 521 OK  8ab037fb2383
row 522 OK  8426fe571313
row 523 OK  50c04965c167
```

### Census

Historical census **808**. Current corpus **815**. Reserved placeholders `HEKI-D-0809`–`0815` are **not registered
ledger nodes and have no recorded hashes**.

### Next

Resolve D-1…D-4 into the TASK-03 decision record, including the events-per-conversation measurement. **No code until
the owner reviews that record.**

**Resume token:** `HEKI-CP-2026-11-22-12`

---

## `HEKI-CP-2026-11-22-13` — TASK-03 ledger decision record

**Stage:** foundational · **TASK-03:** `prepared → decided for review` · **implementation: not authorized**

### Workload measurement: `Unestablished`, not zero

The resolved chat root (`~/Documents/Howl/Chats`) is **absent**; `config/chats` is absent; `%APPDATA%\Howl` and
`%LOCALAPPDATA%\Howl` are absent. **No conversation history exists on this machine**, so p50/p95/max cannot be
computed and are recorded as `Unestablished`. A measurement procedure is written into the record for when data
exists.

**Incidental finding:** `./Chats` exists **inside the installation root**, holding two chat directories with
`memories.json` only. That is the `ProjectRoot`-as-state-root defect TASK-01 characterised and deliberately did not
migrate (L-5), observed in the wild. Which code path created it is unverified.

### D-1 — representation: evolve, do not replace

Selected **E2**: `history.ndjson` stays the ledger file (one JSON object per line, per conversation), gains a
stable envelope, and is committed by **whole-document replacement through DurableStore**. Scored against the
guarantees, not against edit volume.

The decisive properties: it adds **no new file** (so no sixth authority), keeps **legacy lines readable** so
`FC1-07` stays unresolved without a migration, and **structurally eliminates torn appends** — nothing ever writes
into the live ledger, so the entire framing/checksum/truncation design leaves Phase 1 instead of being deferred.

**E1 (true append journal) is preserved as the measured upgrade path**, not discarded. **R (new document) rejected**
as a duplicate authority. Cost stated plainly: O(n) per commit and ~3× storage.

### D-2 — ordering and writer authority

Per-conversation `ledger.lock` reusing TASK-02's `LockFileEx` discipline · sequence assigned under the lock as
`last + 1` · gaps are errors that mark the ledger `unresolved`, never renumbered · second process refused with
contention, never merged.

### D-3 — selection after deletion

Nearest surviving ancestor, recorded as a `selection_resolved` event; explicit-empty fallback when no ancestor
survives. Silent fallback and resurrection prohibited. **`localStorage` selection is demoted to a cache** of the
ledger projection — removing that authority is part of the branch.

### D-4 — recovery and the journaling verdict

Duplicate submission → one effect; lost acknowledgement → retry by key returns the existing event as success;
damaged record → parsing stops, prefix stays valid, remainder reported `unresolved`. **Non-negotiable: recovery may
refuse or report a prefix unresolved, never silently discard valid history.**

**Journaling primitive: not required on current evidence, not built, not a service.** Upgrade trigger is explicit:
p95 records > 2,000, or max > 10,000, or p95 ledger > 8 MiB, or measured commit latency above the turn budget. The
calculated context for those numbers is labelled arithmetic on assumed record sizes, not measurement.

### Retention — recommended, but an owner decision

Recommended **A: ledger deleted with its conversation, plus a retained deletion receipt**, authored content
untouched. Options B (retain) and C (window) are preserved in the record. FC-1 leaves the chat-scoped case
unresolved (`FC1-STATE-RETENTION`), so this fills an `UN` field and **requires the owner**, not this checkpoint.

### Preserved open questions

`FC1-02` (graph now or later), `FC1-08` (mandatory envelope fields), `FC1-11` (Phase 1 event kinds), `FC1-14b/14c`,
`FC1-07`/`FC1-LEGACY`, `FC1-STATE-RETENTION`, `FC1-SNAPSHOT-CONTENT`, `FC1-PRIVACY-DELETE`. None was forced to a
false consensus.

### Registry verification after the final edits (actual results)

```
row 517 OK  e20d187df9d7
row 518 OK  8bb5f657988c
row 519 OK  135f23716bb4
row 520 OK  d63f5d74f039
row 521 OK  8ab037fb2383
row 522 OK  8426fe571313
row 523 OK  50c04965c167
row 524 OK  9c2c6c62db9c
```

### Census

Historical census **808**. Current corpus **816**. Reserved placeholders `HEKI-D-0809`–`0816` are **not registered
ledger nodes and have no recorded hashes**.

### Next

Owner review of the decision record as a whole. On approval, implementation may be authorized: ledger package first,
then retirement of the four competing implementations. **No code until then.** FC-1 remains `PROPOSED`.

**Resume token:** `HEKI-CP-2026-11-22-13`

---

## `HEKI-CP-2026-11-22-14` — TASK-03 amendments after owner review

**Stage:** foundational · **TASK-03:** `decided for review` · **implementation: still not authorized**

### D-1 condition recorded as an enforcement rule

Whole-document replacement preserves the **logical** append-only history; physical replacement grants no permission
to rewrite, reorder, re-encode, deduplicate or remove a committed event. A commit must carry prior records as
**bytes**, never re-serialise them from an in-memory projection. Observable form: after N commits, the first N
records are **byte-identical** to the previous file's prefix, with the pre-image as an independent check. **No
compaction in Phase 1**; any future compaction is a separate authorised act with its own receipt.

### The recovery guarantee is now an acceptance criterion

| # | Criterion |
|---|---|
| `A-12` | A damaged record yields a valid prefix, the **exact** boundary k, and the ledger reported `unresolved` from k — with the file's hash **unchanged after the read** (no repair) |
| `A-13` | A shortened ledger is never presented as complete; a "successful, complete" report is impossible when the boundary is not n |

### Installation-root directories separated as MF-1

`./Chats` inside the installation root (two directories, `memories.json` + `memories_backup.json`, no history file)
is recorded as **MF-1**, a migration finding, **outside** TASK-03. Provenance unverified. TASK-03 binds the
conversation authority to the PathService-resolved root and does **not** migrate, move, delete or clean up these
directories. **Preservation requirement:** both `memories.json` directories remain untouched — not "stale state",
not sweepable (FC1-16, FC1-23a).

### Retention: A, B and C defined exactly

| | **A — delete with conversation + receipt** | **B — retain** | **C — window then delete** |
|---|---|---|---|
| Deletes history? | Yes | No | Eventually |
| Restorable? | No | Yes | Only inside the window |
| Contract | **No amendment needed** | **Contradicts FC1-22c** (owner-approved D2: snapshots die with the conversation) | **Contradicts FC1-16** (no independent background deletion path) + needs a window nobody set |
| Receipt identifies the conversation? | No — salted hash of `conversation_id`, no title | Yes — the history itself is the identifier | Yes until expiry |
| Unattended deletion path? | No | No | **Yes** |
| Owner of the records after deletion | Receipt only | None — a deleted conversation still owns a writable ledger | Scheduler, then receipt |

Receipt contents proposed for A: operation id, timestamp, authority, record and snapshot counts, the last record's
hash as an end-of-history marker, and a **salted hash** of the conversation id. A human-readable title is not
retained; displaying "you deleted N conversations" would reintroduce identifying data and is a separate
`FC1-PRIVACY-DELETE` decision.

**Recommendation unchanged: A.** Decision remains the owner's.

### Still blocking the frozen ledger contract

`FC1-08` envelope fields · `FC1-11` Phase 1 event kinds · `FC1-02` graph scope · `FC1-07`/`FC1-LEGACY` · `FC1-14b`
· `FC1-14c` · the retention choice. Dispositions to prepare and review, not gaps for implementation to fill.

### Process note

The amendment was first written through a shell heredoc that truncated mid-table. It was repaired from the exact
partial line rather than rewritten, and the document was then verified structurally: 319 lines, all nine sections
present, no dangling partial rows.

### Registry verification after the final edits (actual results)

```
row 517 OK  e20d187df9d7
row 518 OK  8bb5f657988c
row 519 OK  135f23716bb4
row 520 OK  d63f5d74f039
row 521 OK  8ab037fb2383
row 522 OK  8426fe571313
row 523 OK  50c04965c167
row 524 OK  8bb216200953
```

### Census

Historical census **808**. Current corpus **816**. Placeholders `HEKI-D-0809`–`0816` are not ledger nodes and have
no recorded hashes.

### Next

Prepare the six dispositions listed above for review. No code, no migration, no contract approval.

**Resume token:** `HEKI-CP-2026-11-22-14`

---

## `HEKI-CP-2026-11-22-15` — TASK-03 contract preparation complete

**Stage:** foundational · **TASK-03:** `prepared + decided + dispositions ready for review` ·
**implementation: NOT authorized**

### Owner decision recorded — retention is no longer open

**Option A is selected. Options B and C are rejected**, with their definitions and rationale preserved in §9.4
because the reasons for rejection are part of the record. Recorded in the decision record §10.1 as an **owner
decision**, not a recommendation. No retained ledger after deletion; no window; no background expiry or sweep;
authored content outside the boundary; a receipt is preserved.

Observed justification for making deletion **enumerated and receipted**: `DeleteChat` is `os.RemoveAll(chatDir)`
(`persistence.go:471–479`) — no receipt, no enumeration, no partial-failure report, and it removes `memories.json`
inside the directory, which belongs to FC1-22e's lifecycle rather than this branch.

### The receipt's salted hash was withdrawn by its own privacy review

A salted hash is a pseudonym, and conversation directories are **timestamp-named**
(`Chats/2026-05-27T16-46-11Z`). With the salt stored beside the hash, anyone who can guess an ID can confirm it —
so the hash behaves like an identifier for every conversation that exists today. Answers recorded for all five
questions: salt ownership is **undecided** (nothing in the tree stores one); scoping is defeated by storing the salt
with the hash; linkage by ID+salt is **possible by construction**; retaining the salt is **not** necessary to prove
deletion; and the receipt **does** create a verification path if it carries any identifier.

**Candidate that survives: a receipt with no conversation identifier** — operation id, timestamp, authority, record
count, snapshot count, last-record hash. UI labels and counts remain prohibited pending `FC1-PRIVACY-DELETE`.

### Enforceable requirements restated

- **D-1** byte preservation: prior records carried as bytes, never re-serialised from a projection; no compaction;
  prefix byte-identity tested; pre-image as independent check. **Byte preservation, semantic validity and integrity
  validation are now distinguished explicitly** and reported separately.
- **A-12** valid prefix, exact damaged boundary, remainder `unresolved`, no repair during read, source hash unchanged
  after read.
- **A-13** a shortened ledger is never complete; **completeness is an explicit result field**, never inferred from a
  parse or a count.
- **MF-1** stays outside TASK-03; `./Chats` is not migrated, moved, cleaned or claimed; both `memories.json`
  directories untouched and not sweepable.

### Six dispositions prepared (proposals only, none adopted)

`TASK-03-FC1-DISPOSITIONS-2026-11-22.md`, each with clause, evidence, decision, alternatives, recommendation,
dependencies and freeze impact:

| Question | Recommendation | Blocks freeze |
|---|---|---|
| `FC1-08` | core of 8 fields + format version (`v`), state/branch fields deferred | Yes |
| `FC1-11` | 8 kinds including `turn_corrected`/`turn_retracted` — required so `RewriteHistory` can retire | Yes |
| `FC1-02` | linear single-parent chain now, graph-shaped API | Yes |
| `FC1-07` | legacy read-only `unresolved`, no invented IDs | No for Phase 1 |
| `FC1-14b` | branching excluded from Phase 1; fork-as-new-conversation noted, not decided | Partially |
| `FC1-14c` | one turn = one message, `turn_id` present | Yes |
| `FC1-PRIVACY-DELETE` | receipt with no conversation identifier; no UI labels or counts | Yes for receipt |

**Decisive evidence found while preparing them:** message identity is the array index — `app-prompt.go:303`
`ID: fmt.Sprintf("%d", len(a.chatMessages)-1)` and `:357` — which is the exact construction FC1-06 prohibits and the
reason legacy records cannot be promoted without inventing identity.

### Verification

Both documents verified structurally. Decision record: 10 sections, §10 with four subsections, ending complete.
Dispositions: summary + seven sections + scope + authorization, ending complete.

Registry rows verified against disk after the final edits:

```
row 517 OK  e20d187df9d7
row 518 OK  8bb5f657988c
row 519 OK  135f23716bb4
row 520 OK  d63f5d74f039
row 521 OK  8ab037fb2383
row 522 OK  8426fe571313
row 523 OK  50c04965c167
row 524 OK  faaf10d1227e
row 525 OK  3986fcbfdec8
```

Prior hashes remain visible: 524 moved `9c2c6c62db9c` → `8bb216200953` → `faaf10d1227e` across the amendment
sequence; 517 moved `e2e80d2871af` → `e20d187df9d7` when the plan corrections were appended. Nothing was rewritten
in place.

### Census and build

Historical census **808**. Current corpus **817**. Reserved placeholders `HEKI-D-0809`–`0817` are **not registered
ledger nodes and have no recorded hashes**.

**Zero Go files modified.** `go build ./...` passes — which validates only that the tree still compiles. It says
**nothing** about TASK-03 behaviour, because no TASK-03 code exists.

### Remaining owner decisions

1. Retention — **decided (A)**. Closed.
2. `FC1-08` core field set, including `v` and `actor`.
3. `FC1-11` eight event kinds.
4. `FC1-02` linear-not-graph for Phase 1.
5. `FC1-07` legacy read-only, migration deferred.
6. `FC1-14b` branching excluded from Phase 1.
7. `FC1-14c` one turn = one message.
8. `FC1-PRIVACY-DELETE` receipt without an identifier; UI disclosure none.
9. Whether per-conversation deletion proof is required (only then does a salt question return).
10. Adoption of DurableStore beyond the ledger — separate task, unauthorized.

### Remaining blockers to implementation

Every item 2–8 above; plus the workload measurement (`Unestablished`) that gates the upgrade threshold; plus the
`RewriteHistory` retirement ordering (ledger first, then the four competing implementations).

**Resume token:** `HEKI-CP-2026-11-22-15`

---

## `HEKI-CP-2026-11-22-16` — contract review materials complete

**Stage:** foundational · **TASK-03:** `prepared + decided + schema proposed` · **implementation: NOT authorized**

### Two counts in the earlier report were wrong, and are corrected in the record

| Earlier phrasing | Corrected |
|---|---|
| "eight core fields, plus `v` and `actor`" | **Nine required fields**: `v`, `event_id`, `conversation_id`, `sequence`, `event_type`, `timestamp`, `idempotency_key`, `actor`, `durability`; plus three conditional: `turn_id`, `predecessor_event_id`, `payload` |
| "eight event kinds in the ledger" | **Six ledger kinds + two evidence kinds.** `conversation_deleted` cannot be stored in the ledger it deletes; `recovery_recorded` may not be written during a read under A-12 |

Also recorded as deliberately absent: `state_id` (derived as the producing event's id), `affected_ids`, `operation`,
and `provenance` reduced to `actor`.

### Anti-erasure is now four testable properties

No in-place mutation · the target record is **byte-identical** after a correction · "erasure" exists only in the
projection · no compaction escape hatch. Documented cost: clearing 5,000 turns writes 5,000 retraction records.

### `FC1-14b`/`FC1-14c` reconciled

One message per turn; `message_id` equals `turn_id` by restriction; state derived, never stored; correction and
retraction each produce a new state while retaining `turn_id`; fork-as-new-conversation recorded as a **possibility**
only, expressible through optional fields so no required field changes if it is ever adopted.

### Receipt tightened after review

Last-record hash moved from required to **optional, proposed default absent** — it is a fingerprint that a surviving
ledger copy can match. Salted conversation-ID hash stays withdrawn. Timestamp and counts are retained with the
residual correlation risk stated rather than hidden.

### Counts kept distinct

**Seven disposition topics** vs **eight owner decisions** — the eighth being whether per-conversation deletion proof
is required at all, which is a different question from how a receipt is shaped.

### Verification

Registry rows verified against disk after the final edits:

```
row 517 OK  e20d187df9d7
row 518 OK  8bb5f657988c
row 519 OK  135f23716bb4
row 520 OK  d63f5d74f039
row 521 OK  8ab037fb2383
row 522 OK  8426fe571313
row 523 OK  50c04965c167
row 524 OK  0c50ab67c9f1
row 525 OK  4a740fb3ecbb
row 526 OK  22fd97675e0d
```

Hash transitions preserved: 524 `9c2c6c62db9c` → `8bb216200953` → `faaf10d1227e` → `0c50ab67c9f1`; 525
`3986fcbfdec8` → `4a740fb3ecbb`. A registry-update script failed mid-run on a wrong environment key; it wrote
nothing to the registry, the README had already been updated, and the rerun corrected both — recorded because a
half-applied registry update would have been a real defect.

### Census and build

Historical census **808**. Current corpus **818**. Placeholders `HEKI-D-0809`–`0818` are not ledger nodes and have no
recorded hashes.

**Zero Go files modified.** `go build ./...` passes; it validates compilation only and says nothing about TASK-03.

### Remaining owner decisions

`FC1-08` nine-field envelope · `FC1-11` six-plus-two event model · `FC1-02` linear-not-graph · `FC1-07` legacy
read-only · `FC1-14b` branching excluded · `FC1-14c` one message per turn · receipt without identifier · whether
per-conversation deletion proof is required.

### Remaining blockers

The decisions above; the workload measurement (`Unestablished`); and the retirement order (ledger first, then the
four competing implementations). The complete contract must be reviewed and explicitly accepted before any
implementation authorization.

**Resume token:** `HEKI-CP-2026-11-22-16`

---

## `HEKI-CP-2026-11-22-17` — final contract clarification round complete

**Stage:** foundational · **TASK-03:** `prepared + decided + schema proposed + contract final` · **implementation: NOT authorized**

### Three of my own earlier statements were wrong and are corrected

| Earlier | Now |
|---|---|
| `predecessor_event_id` "recommended on turn events; optional otherwise" | **Required for every `sequence ≥ 2`**, absent only at `sequence 1` |
| Left open whether every event produces a state | **Only `turn_committed`, `turn_corrected`, `turn_retracted` produce states.** Selection events never do |
| "Evidence store" named with no authority or lifecycle | Authority, two homes, two schemas and transitions defined |

The schema proposal's §2.2 row is struck through rather than rewritten, and points at the superseding clause.

### Failed commits cannot create a gap — structurally, not by policy

Sequence is **derived at commit from the committed file under the lock**, never pre-allocated. A failure before the rename leaves the file unchanged, so the retry reads the same last sequence and **reuses the same number**. Pre-reserved ranges, in-memory counters and "burn the number on failure" are named as contract violations. Crash-after-rename is recovered by the idempotency key, not by a second position.

### The evidence store had to become two homes

`conversation_deleted` must outlive the conversation, so it lives in app-state (`evidence/deletions.ndjson`). `recovery_recorded` is identified by location, so it lives in the conversation directory and dies with it under Option A. A pre-removal receipt is `intent`; completion requires a **second** record, and no reader may report deletion from a lone intent record — an intent with no outcome is `unknown outcome`, reported and never swept.

### Correction and retraction, settled

Corrections may target corrected turns and chain indefinitely; last by sequence wins per state. Retraction may target a corrected turn. **Retraction is irreversible in Phase 1** — restoring a message requires a new `turn_id`, which is a new fact, not an undo. Neither operation deletes or re-serialises the target; byte identity is a tested property.

### A linear event chain is not a chain of state changes

`predecessor_event_id` is **ordering only**. Consecutive events can share one visible state. Semantic ancestry is the derived **state** relation. This was the ambiguity the review exposed.

### Consistency

**Internally consistent within the stated Phase 1 restrictions; no contradiction found.** Residual risks retained rather than smoothed: power-loss `UN`; irreversible retraction; counts plus timestamp correlatable with a surviving ledger copy; legacy outside the ledger; workload `Unestablished`.

### Verification

```
row 517 OK  e20d187df9d7
row 518 OK  8bb5f657988c
row 519 OK  135f23716bb4
row 520 OK  d63f5d74f039
row 521 OK  8ab037fb2383
row 522 OK  8426fe571313
row 523 OK  50c04965c167
row 524 OK  0c50ab67c9f1
row 525 OK  4a740fb3ecbb
row 526 OK  c8ecbaf0f3b2   (was 22fd97675e0d; §2.2 superseded)
row 527 OK  50fcef6ffeb0
```

Hash transitions preserved: 526 `22fd97675e0d` → `c8ecbaf0f3b2`. Corpus **819**; placeholders `0809`–`0819`; registry 818 besides / 870 files.

**Zero Go files modified.** The verification command exits 1 because `grep -v '^??'` matches nothing — the zero-line result after filtering, not a failure.

### Remaining owner decisions

`FC1-08` nine fields with required predecessor · `FC1-11` six-plus-two · `FC1-02` linear · `FC1-07` legacy read-only · `FC1-14b` branching excluded and fork unadopted · `FC1-14c` one message per turn · receipt without the five excluded items · whether per-conversation deletion proof is required.

### Blockers unchanged

Those decisions; the `Unestablished` workload measurement; the retirement order. Contract must be explicitly accepted before any implementation authorization.

**Resume token:** `HEKI-CP-2026-11-22-17`

---

## `HEKI-CP-2026-11-22-18` — consolidated contract 1.0-rc; four review findings dispositioned

**Stage:** foundational · **TASK-03:** `contract 1.0-rc, awaiting acceptance` · **implementation: NOT authorized**

### Finding 1 — corrections target `turn_id`, and earlier states are immutable *by construction*

The guarantee no longer rests on discipline. A correction committed at `sequence m` is **not an input** to `P_k` for
any `k < m`, so `P_k` cannot change. That is C-24, and it is stronger than "we promise not to rewrite": there is no
code path whose input includes the past. Projection is defined as a pure fold over the event sequence with no clock,
locale, iteration-order or randomness dependency. Historical-state editing is excluded pending separate
authorization.

### Finding 2 — the crash window is named, and reconciliation is evidence-bounded

`W` = removal succeeded, terminal receipt not flushed. During `W` a lone `intent` is **`unknown outcome`** on every
path — never complete, never failed, never swept. `ReconcileDeletion(operation_id)` is explicit, bounded, idempotent,
deletes nothing, invents no identity. Its five-row table permits `complete` with **`attribution: unknown`**: the
effect is established, the actor is not. `aborted` and `partial` are the other outcomes; a receipt cannot be invented
for an operation that has none.

### Finding 3 — durability terminology pinned

`observed` = an observed successful rename after a flushed temp write under the lock. It does **not** prove
power-loss durability, post-restart presence, read-back verification, or metadata flushing. Power loss stays `UN`.
`verified` stays unwritten. Mapping `observed` to user-facing "saved"/"safe" is a **test failure** (T-21), not a style
choice.

### Finding 4 — recovery evidence is diagnostic, not a second authority

Not a replica, not a fallback, not reconstructable into a conversation. Losing `recovery.ndjson` with the directory is
**correct** under Option A, not data loss. The app-state receipt evidences the operation and never replaces the ledger;
its presence does not imply the conversation exists and its absence does not imply deletion succeeded.

### Supersession

`TASK-03-CONTRACT-FINAL` §5.3 was **not rewritten** — a supersession note extends its `state` enumeration to five
values and adds `attribution`, pointing at C-35. Hash `50fcef6ffeb0` → `a0f89160843f`.

### Verification

```
row 517 OK  e20d187df9d7     row 523 OK  50c04965c167
row 518 OK  8bb5f657988c     row 524 OK  0c50ab67c9f1
row 519 OK  135f23716bb4     row 525 OK  4a740fb3ecbb
row 520 OK  d63f5d74f039     row 526 OK  c8ecbaf0f3b2
row 521 OK  8ab037fb2383     row 527 OK  a0f89160843f
row 522 OK  8426fe571313     row 528 OK  9ff34f0bdbeb
```

Contract contains **53 clauses** and **24 acceptance tests**. Corpus **820**; placeholders `0809`–`0820`; registry 819
besides / 871 files. **Zero Go files modified**; `go build ./...` passes (compilation only).

### Remaining owner decisions

Nine: envelope with required predecessor · six-plus-two kinds · linear · legacy read-only · branching excluded and
fork unadopted · one message per turn · receipt without the five excluded items · whether per-conversation deletion
proof is required · **acceptance of contract 1.0-rc**.

### Blockers

The `Unestablished` workload measurement; the retirement order; and explicit acceptance. No implementation is
authorized by this checkpoint.

**Resume token:** `HEKI-CP-2026-11-22-18`

---

## `HEKI-CP-2026-11-22-19` — five consistency findings dispositioned

**Stage:** foundational · **TASK-03:** `contract 1.0-rc + clarification 1` · **implementation: NOT authorized**

### The one that changed the contract, not just its wording

`selection_resolved` had **no Phase 1 trigger**. C-22 keeps every state permanently addressable and C-03 excludes
state deletion, so nothing can delete a state and the event can never legally fire. Deferred: **five ledger kinds**,
not six. D-3's nearest-surviving-ancestor rule was misfiled there — it concerns a conversation that may no longer
exist, so it cannot be an event in that conversation's ledger. Re-homed to conversation-level selection in app-state,
carried forward as an open item rather than dropped.

### Reconciliation identity, resolved without putting an ID in the receipt

A transient `pending_deletions` journal in app-state maps `operation_id` → path, written **after** the intent and
**before** removal, deleted at the terminal record. Startup reports it read-only; the reconcile call consumes it. Two
consequences matter: a missing entry proves removal never started (`aborted`), and no directory is ever guessed by
scanning, counts or timestamp — that would rebuild the identity the receipt deliberately omits. Residual privacy
risk stated: a path can survive a crash in the window until reconciliation.

### Absence is not evidence of authorship

`absent_prior` is a terminal state for an already-absent directory — no intent, no removal, no claim. `attribution:
operation` requires **all three** of observed presence, a successful in-process removal, and observed absence; a crash
in `W` yields `complete` with `attribution: unknown`. Later absence alone never proves this installation deleted
anything.

### Selection precedence, made deterministic

Explicit selection persists until the next state-producing event, which re-selects the head. Nine-event worked example
recorded. Consequence accepted rather than hidden: viewing an earlier state and receiving a new turn moves you to the
head, because the alternative leaves new content invisible.

### Actor/role

Explicit matrix over eight rows with distinct errors: cross pairs `ROLE_ACTOR_MISMATCH`, `system` and `recovery`
`ACTOR_NOT_PERMITTED` on turns, per-kind actor allowances tabulated. `recovery` authors no ledger record in Phase 1.

### Verification

```
row 517 OK  e20d187df9d7     row 523 OK  50c04965c167
row 518 OK  8bb5f657988c     row 524 OK  0c50ab67c9f1
row 519 OK  135f23716bb4     row 525 OK  4a740fb3ecbb
row 520 OK  d63f5d74f039     row 526 OK  c8ecbaf0f3b2
row 521 OK  8ab037fb2383     row 527 OK  a0f89160843f
row 522 OK  8426fe571313     row 528 OK  9ff34f0bdbeb
row 529 OK  5fcbc8c7ee4f
```

**Row 528 is byte-identical** (`9ff34f0bdbeb`) — the 1.0-rc document was preserved as instructed; the clarification
supersedes named clauses without editing them. 23 `K-` clauses, 17 `KT-` tests, consistency matrix over 25 affected
clauses. Corpus **821**; placeholders `0809`–`0821`; registry 820 besides / 872 files.

**Zero Go files modified.** `go build ./...` passes (compilation only).

### Carried forward

1. `selection_resolved` and D-3's ancestor rule need a later contract.
2. The pending-deletions journal holds a path during the window.
3. Power-loss durability `UN`; workload `Unestablished`.
4. Eight owner decisions plus acceptance; decision 2 restated as five ledger + two evidence kinds.

### Blockers

Explicit acceptance; the measurement; the retirement order. No implementation is authorized by this checkpoint.

**Resume token:** `HEKI-CP-2026-11-22-19`

---

## `HEKI-CP-2026-11-22-20` — clarification 1a; four findings closed

**Stage:** foundational · **TASK-03:** `contract stack 1.0-rc + 1 + 1a` · **implementation: NOT authorized**

### The two contradictions were one defect each

**K-04 vs K-06.** Both inferred from the same observation — intent present, journal entry absent — and reached opposite
conclusions. The root cause was the write order: with `intent → journal`, absence could not be distinguished from
"never written". **The order is inverted**: the entry is now created *before* the intent, so an intent with no entry
can only mean the entry was created and later became unavailable. Absence therefore yields
`indeterminate`/`locator_unavailable`, never `aborted`.

**K-02 vs K-05.** Two readers were needed and only one was stated. Exactly two now: read-only startup discovery
reporting operation ids, and explicit reconciliation. Startup never reconciles.

### The safety property that makes the journal contractual

> Every journal failure — missing, corrupt, truncated, partially written, unknown version — degrades toward
> `indeterminate`. No journal failure may produce `aborted`, `absent_prior`, `complete`, or any claim in either
> direction.

Truncation is the dangerous case: a shortened file still parses, so an entry can vanish while the file looks valid.
That routes to `locator_unavailable` and claims nothing. Journal integrity, versioning, retention, permissions and
export exclusion are now clauses, not implementation choices.

### Actor matrix rebuilt per kind

K-21 was written for turns and generalised with the word "any". `system` authors `conversation_created` and
`selection_changed`; it authors no turn, correction or retraction. Corrections and retractions carry no `role` — their
authority derives from the target turn, with `TARGET_NOT_FOUND` evaluated before the actor check so errors are
deterministic.

### Terminal-then-cleanup crash rule

Reconciliation inspects the terminal receipt first. If one exists it is returned unchanged, the stale entry is
removed, and **no second receipt is written**. Reconciliation removes journal entries only; it never touches a
conversation directory, never re-verifies a `complete`, and never repeats a deletion. The journal's claimed lifetime
is thereby restored rather than merely asserted.

### Verification

```
row 517 OK  e20d187df9d7     row 524 OK  0c50ab67c9f1     row 529 OK  5fcbc8c7ee4f
row 518 OK  8bb5f657988c     row 525 OK  4a740fb3ecbb     row 530 OK  a929e63804bb
row 519 OK  135f23716bb4     row 526 OK  c8ecbaf0f3b2
row 520 OK  d63f5d74f039     row 527 OK  a0f89160843f
row 521 OK  8ab037fb2383     row 528 OK  9ff34f0bdbeb
row 522 OK  8426fe571313
row 523 OK  50c04965c167
```

**Rows 528 and 529 byte-identical** — neither earlier document was edited. 21 `L-` clauses, 20 `LT-` tests. Event-kind
count unchanged: **five ledger plus two evidence kinds**. Corpus **822**; placeholders `0809`–`0822`; registry 821
besides / 873 files.

**Zero Go files modified.** The "zero Go" check exits nonzero because the filter matches no lines; that is the zero
result, not a failure.

### Standing statement on evidence

`go build` passing and documents agreeing are **not** implementation acceptance. Not one `T-` or `LT-` test has been
executed, because no implementation exists. The contract is consistent as text; its tests are unrun by construction.

### Carried forward

1. `selection_resolved` and D-3's ancestor rule remain deferred.
2. The journal holds a path for the life of an unresolved deletion — bounded, versioned, sensitive.
3. Power-loss durability `UN`; workload `Unestablished`.
4. Eight owner decisions plus acceptance.

**Resume token:** `HEKI-CP-2026-11-22-20`

---

## `HEKI-CP-2026-11-22-21` — clarification 1b; three findings closed

**Stage:** foundational · **TASK-03:** `contract stack 1.0-rc + 1 + 1a + 1b` · **implementation: NOT authorized**

### The partial-attribution contradiction was mine, and it is resolved against me

L-03 let a partly emptied directory claim `attribution: operation`; L-21 forbade it. L-21 was right. **Partial now
always carries `attribution: unknown`**, and `attribution: operation` is narrowed to `state: complete` only. The error
was conflating a claim about **authorship** with a claim about **effect**: observing a directory proves it existed,
not that every later change was ours. A partial receipt still records what happened to our call (`removal_result`) —
a fact — without asserting we removed the files.

### Cleanup now has two classes, because a rewrite can destroy other operations

Entry-level invalid (journal parses, one entry fails, `operation_id` readable) may be removed only if the journal
parses, the entry is located by exact match, every other entry re-emits byte-equivalently, and the rewrite is atomic.
Journal-level unreadable (unparseable, truncated, unversioned, unreadable) is **preserved byte-identically**: no
rewrite, no prefix reconstruction, no quarantine rename. A terminal receipt never depends on the journal, so
corruption changes no outcome — it only leaves a stale entry that is reported, never forced away.

### Deletion now has an owner

The journal lock serialized single operations, not a lifecycle. A non-terminal entry keyed by path **is** the
ownership record, inserted under the journal lock in the same critical section as the presence observation, so no
process reaches `remove` without an entry it owns. A second request returns `busy_pending_deletion` with the owner's
id, writes no intent, removes nothing, writes no receipt. One physical deletion cannot produce two receipts.

### Verification

```
row 517 OK  e20d187df9d7   row 523 OK  50c04965c167   row 529 OK  5fcbc8c7ee4f
row 518 OK  8bb5f657988c   row 524 OK  0c50ab67c9f1   row 530 OK  a929e63804bb
row 519 OK  135f23716bb4   row 525 OK  4a740fb3ecbb   row 531 OK  e3131ad3eab2
row 520 OK  d63f5d74f039   row 526 OK  c8ecbaf0f3b2
row 521 OK  8ab037fb2383   row 527 OK  a0f89160843f
row 522 OK  8426fe571313   row 528 OK  9ff34f0bdbeb
```

**Rows 528, 529 and 530 byte-identical** — no earlier document was edited. 15 `M-` clauses, 14 `MT-` tests.
Event-kind count unchanged: **five ledger plus two evidence kinds**. New return values
(`busy_pending_deletion`, `stale_entry_reconcile_required`, `journal_cleanup_unavailable`,
`journal_cleanup_refused`) are **not** receipt states or events. Corpus **823**; placeholders `0809`–`0823`;
registry 822 besides / 874 files.

**Zero Go files modified.** The check exits nonzero because the filter matches no lines: the zero result.

### Standing statement on evidence

The contract is consistent **as text**. It is not implemented and not empirically validated. No `T-`, `KT-`, `LT-`
or `MT-` test has ever been executed. TASK-03 remains contract preparation; contract consistency is not
implementation acceptance.

### Carried forward

1. `selection_resolved` and D-3's ancestor rule remain deferred.
2. A permanently corrupt journal leaves an irremovable stale entry — reported, not repaired.
3. Partial deletions carry no authorship claim, so a summary may say "deletion did not complete", never "the app
   deleted some files".
4. Power-loss durability `UN`; workload `Unestablished`.
5. Eight owner decisions plus acceptance.

**Resume token:** `HEKI-CP-2026-11-22-21`

---

## `HEKI-CP-2026-11-22-22` — clarification 1c; deletion authority complete

**Stage:** foundational · **TASK-03:** `contract stack 1.0-rc + 1 + 1a + 1b + 1c` · **implementation: NOT authorized**

### Reader count corrected: three, by scope not by number

L-06 said two; M-11 required a third. Deletion acquisition is now a named reader with a bounded permission set: inspect
ownership for one candidate directory, insert one `prepared` entry, reject with a status, at most one `absent_prior`
receipt — no cleanup, no entry removal, no path disclosure.

### The abandoned-entry deadlock is closed by the state machine, not by a timer

A `prepared` entry had no release path. The fix uses what the durable ordering already proves: `prepared` and `pending`
mean the removal call **was never made**. Reconciliation may therefore remove a `prepared` entry with no receipt, but
only under the journal lock with three verifications — exact `operation_id`, state exactly `prepared`, no intent and no
terminal receipt. The intent *is* the request, so with no intent there is nothing to record.

**A correction beyond the finding, flagged for review:** because the barrier proves removal was never called,
`pending` + absent directory now yields **`absent_prior`/`prior_absence`**, not `complete`/`unknown`. L-03 was claiming
more than the ordering permits. This direction — claiming less — is why it was taken without asking, but it changes a
row the owner has already reviewed.

### Identity, not strings

Ownership is keyed by the canonical directory identity established through **PathService** — the existing TASK-01
authority that resolves final paths by handle and knows junctions, which `EvalSymlinks` does not follow. No second
resolver. Unresolvable identity fails closed (`identity_unresolved`), an absent target still yields `absent_prior`,
and identity is re-checked immediately before removal so a directory removed and re-created at the same path aborts
with `identity_changed` instead of being destroyed.

### Verification

```
row 517 OK  e20d187df9d7   row 523 OK  50c04965c167   row 529 OK  5fcbc8c7ee4f
row 518 OK  8bb5f657988c   row 524 OK  0c50ab67c9f1   row 530 OK  a929e63804bb
row 519 OK  135f23716bb4   row 525 OK  4a740fb3ecbb   row 531 OK  e3131ad3eab2
row 520 OK  d63f5d74f039   row 526 OK  c8ecbaf0f3b2   row 532 OK  6bf688cbef86
row 521 OK  8ab037fb2383   row 527 OK  a0f89160843f
row 522 OK  8426fe571313   row 528 OK  9ff34f0bdbeb
```

**Rows 528–531 byte-identical.** 17 `N-` clauses, 14 `NT-` tests. Event-kind count unchanged: **five ledger plus two
evidence kinds**. Corpus **824**; placeholders `0809`–`0824`; registry 823 besides / 875 files.

**Zero Go files modified.** `go vet ./internal/paths/` clean — vetting existing code, not evidence about this contract.

### Standing statement on evidence

The contract is consistent **as text**. Nothing is implemented or empirically validated; no `T-`/`KT-`/`LT-`/`MT-`/`NT-`
test has ever run. TASK-03 remains contract preparation.

### Carried forward

1. `selection_resolved` and D-3's ancestor rule remain deferred.
2. A permanently corrupt journal leaves an irremovable stale entry.
3. Identity is checked at acquisition and before removal; a rename between those points is detected, not prevented.
4. Partial deletions carry no authorship claim.
5. Power-loss durability `UN`; workload `Unestablished`.
6. Eight owner decisions plus acceptance.

**Resume token:** `HEKI-CP-2026-11-22-22`

---

## `HEKI-CP-2026-11-22-23` — clarification 1d; deletion claims now bounded by evidence

**Stage:** foundational · **TASK-03:** `contract stack 1.0-rc + 1 + 1a + 1b + 1c + 1d` · **implementation: NOT authorized**

### `absent_prior` was a false historical claim, and the fix claims less

Acquisition observes the directory **present** before writing the entry, so absence discovered later cannot be prior.
N-07 asserted a history that never happened — exactly the failure this contract exists to prevent. New outcome
`absent_after_acquisition` with `attribution: unknown` and `removal_attempted: false`, stating only: present at
acquisition, absent now, removal never called. `prior_absence` is now reserved for the acquisition observation itself.
A partly emptied directory under `prepared`/`pending` is likewise an external change and carries no authorship.

### The removal guarantee, stated as a limitation instead of a claim

Re-checking identity before `RemoveAll` narrows a race; it cannot close it, because the call re-resolves the path.
PathService is explicitly **not** a safe-deletion capability (TASK-01 separates classification from authorised writes).
Two bases are named: `handle_identity`, which is the only one supporting "the acquired directory was removed", and
`path_recheck`, which Phase 1 uses and which may claim only that the path is absent and the call returned success.
`identity_changed` proves an observed mismatch and nothing more.

**This is the first finding in the series that made the contract weaker on purpose.** A replaced directory can be
affected under `path_recheck`; the receipt is honest about its basis but honesty does not undo the effect. Raised to
the **tenth owner decision**: accept `path_recheck` with the limitation, or require `handle_identity` before Phase 1
deletion is enabled.

### Identity for something that does not exist

A composite identity — canonical parent identity plus a validated single-component child name, compared under the
filesystem's own case rule — used for classification only, never stored, never an ownership key. Five-way failure
taxonomy separates genuine absence from an inaccessible parent, an unresolvable parent, a target appearing
mid-classification, and a target whose own resolution fails. No string fallback; absence needs no owner, so two
concurrent requests may each write `absent_prior` without ambiguity.

### Verification

```
row 517 OK  e20d187df9d7   row 523 OK  50c04965c167   row 529 OK  5fcbc8c7ee4f
row 518 OK  8bb5f657988c   row 524 OK  0c50ab67c9f1   row 530 OK  a929e63804bb
row 519 OK  135f23716bb4   row 525 OK  4a740fb3ecbb   row 531 OK  e3131ad3eab2
row 520 OK  d63f5d74f039   row 526 OK  c8ecbaf0f3b2   row 532 OK  6bf688cbef86
row 521 OK  8ab037fb2383   row 527 OK  a0f89160843f   row 533 OK  6b3afd4f9b28
row 522 OK  8426fe571313   row 528 OK  9ff34f0bdbeb
```

**Rows 528–532 byte-identical.** 15 `P-` clauses, 13 `PT-` tests. Event-kind count unchanged: **five ledger plus two
evidence kinds**. Corpus **825**; placeholders `0809`–`0825`; registry 824 besides / 876 files.

**Zero Go files modified.** The zero-Go check exits nonzero because the filter matches no lines.

### Standing statement on evidence

Consistent as text; not implemented; not empirically validated; no test of any prefix has ever run. The deletion
section is now weaker by design, which is a gain in honesty and a real constraint on Phase 1.

### Carried forward

1. **Tenth owner decision:** `path_recheck` vs `handle_identity` for Phase 1 deletion.
2. `selection_resolved` and D-3's ancestor rule remain deferred.
3. A permanently corrupt journal leaves an irremovable stale entry.
4. Power-loss durability `UN`; workload `Unestablished`.
5. Nine prior owner decisions plus acceptance.

**Resume token:** `HEKI-CP-2026-11-22-23`

---

## `HEKI-CP-2026-11-22-24` — clarification 1e; deletion disabled pending a spike

**Stage:** foundational · **TASK-03:** `contract stack + 1e` · **implementation: NOT authorized · deletion: BLOCKED**

### D-10 recorded, and it turned out to be a schema decision

Owner decided identity-bound deletion is mandatory and `path_recheck` is rejected. The contract consequence I did not
anticipate: **a path cannot identify a directory**, so the journal could not keep storing one. Schema v2 stores
`(volume_serial, file_id_128, entries_at_acquisition)`. The path leaves durable state entirely — the K-07/L-20 privacy
risk (a path surviving a crash) disappears, and ownership comparison becomes identity comparison.

`entries_at_acquisition` exists because `partial` otherwise cannot be calculated honestly.

### Ordering corrected

`pending → removing → remove`. The `removing` marker is durable **before** the call, so it proves the boundary was
crossed and never that the call ran. Removal evidence is derived as `not_started` / `possible` / `verified` from
durable state alone. Ten crash windows tabulated; window 8 is the subtle one — a verification performed in memory
that never reached the `removed` write is **not evidence**, so it resolves to `absent_unverified`.

New outcome `absent_unverified` replaces `complete`/`unknown`: under D-10, path absence supports nothing.

### The eleven-pair table was audited and was wrong in three ways

Read from disk, not memory. It listed **11 of 12** matrix cells — `removed`+partly-emptied was missing. `removing`+absent
was `complete`/`unknown`, which D-10 forbids. `removed`+present was called inconsistent, but a replacement at the path
is *expected* and irrelevant once identity is the subject; the real inconsistency is the acquired identity still
resolving. Corrected to 16 rows including declared-impossible combinations.

### Feasibility: UNESTABLISHED, and that is the honest verdict

Documented (source-derived): `FILE_DISPOSITION_INFORMATION_EX` flags and DELETE-access requirement, `FILE_ID_INFO`
serial+128-bit ID identity, `FileIdExtdDirectoryInfo` enumeration, `OpenFileById` with reparse-point choice. Not
established: whether a `CreateFileW` handle is a valid `RootDirectory`, whether `FILE_DIRECTORY_FILE` +
`FILE_OPEN_REPARSE_POINT` is supported, whether withholding `FILE_SHARE_DELETE` actually prevents rename **and**
reparse conversion, and per-filesystem `OpenFileById`-failure semantics. Ten evidence items E-1…E-10 named; Q-01 keeps
deletion disabled until they are recorded.

### Absent target and case rules

`absent_prior` now requires a retained parent handle or a pinned parent; otherwise `absence_unverified` — absence is not
claimed. Case comparison reads per-directory class 71 on the parent and fails closed with `case_rule_unknown`; no
case-insensitive default, no application-level Unicode normalization.

### Verification

```
row 517 OK  e20d187df9d7   row 523 OK  50c04965c167   row 529 OK  5fcbc8c7ee4f
row 518 OK  8bb5f657988c   row 524 OK  0c50ab67c9f1   row 530 OK  a929e63804bb
row 519 OK  135f23716bb4   row 525 OK  4a740fb3ecbb   row 531 OK  e3131ad3eab2
row 520 OK  d63f5d74f039   row 526 OK  c8ecbaf0f3b2   row 532 OK  6bf688cbef86
row 521 OK  8ab037fb2383   row 527 OK  a0f89160843f   row 533 OK  6b3afd4f9b28
row 522 OK  8426fe571313   row 528 OK  9ff34f0bdbeb   row 534 OK  0529728ef74f
```

Rows 528–533 byte-identical. 32 `Q-` clauses, 28 `QT-` tests. Event kinds unchanged: **five ledger + two evidence**.
Corpus **826**; placeholders `0809`–`0826`; registry 825 besides / 877 files.

**Zero Go files modified.** No test of any kind executed. Feasibility is `UNESTABLISHED`, labelled as such in the
document rather than dressed as a conclusion.

### Carried forward

1. **Decision 11 (new):** enable deletion before the spike, or keep it disabled until E-1…E-10 are recorded.
2. `selection_resolved` and D-3's ancestor rule remain deferred.
3. A permanently corrupt journal leaves an irremovable stale entry.
4. Power-loss durability `UN`; workload `Unestablished`.
5. Nine earlier owner decisions plus acceptance.

**Resume token:** `HEKI-CP-2026-11-22-24`

---

## `HEKI-CP-2026-11-22-25` — clarification 1f; the recovery path now exists, deletion still gated

**Stage:** foundational · **TASK-03:** `contract stack + 1f` · **implementation: NOT authorized · deletion: UNAVAILABLE**

### Q-12 was not executable, and the fix added a locator

`OpenFileById` needs a volume handle; after a restart there was none and no stored locator. Schema v3 adds a **volume
GUID path** plus filesystem. The serial is demoted to a corroborator: it is written at format time and third-party
tools can change it, so it cannot locate anything. Reacquisition is now specified end to end, with six indeterminate
failure classes and no path fallback.

Two consequences worth the owner's attention:

- **FAT/exFAT are excluded from deletion entirely.** Their IDs are derived from directory-entry position and change on
  defrag or rename, so identity-bound deletion cannot hold. `FileId == 0` means "unsupported" and is refused. D-10
  leaves no alternative; deletion is simply unavailable for conversations on those filesystems.
- **A byte-cloned volume is indistinguishable** — same GUID, same serial, same IDs. This is the one remaining
  over-claim vector, and it is stated rather than denied. `complete` is therefore scoped to "the volume reached by the
  stored locator at verification time".

Reuse is safe in the direction that matters: NTFS increments a 16-bit sequence on reuse, so a stale ID fails to match
and the error is an under-claim, never a false `complete`.

### `absent_prior` was contradictory and is now gated

Table row 1 emitted it while §8 forbade it without a stable parent. The gate wins: `absent_prior` is **unreachable in
Phase 1**, `absence_unverified` is a result code with no receipt and no entry. `pinned_parent` is explicitly not
proven by holding a handle.

### The crash table now matches the state table

Four target conditions — identity absent, unchanged, partially changed, replacement at the path — across all ten
windows. An observed partial tree with no possible call is external change, and its attribution is `unknown`.

### Verified absence defined

`removed` requires **L3 only**: a *fresh* `OpenFileById` fails after every deletion handle is closed. A retained handle
proves L2, which is compatible with success. `complete` claims no erasure, no exclusivity, no absence of clones, and
no empty path.

### Decision 11 withdrawn

It was framed as a product choice and was not one. D-10 implies gate G-1: deletion stays unavailable until the
mechanism is established. Owner decisions are back to ten, one decided.

### Verification

```
rows 517-535 all OK against disk
row 534 (1e) unchanged: 0529728ef74f
row 535 (1f): d386317a316f
```

30 `R-` clauses, 22 `RT-` tests. Corpus **827**; placeholders `0809`–`0827`; registry 826 besides / 878 files.
**Zero Go files modified; no ledger created; no migration; no deletion enabled.** No test executed — 1f is source-derived
API analysis plus contract text, with E-11…E-18 and DT-1…DT-6 outstanding.

**Resume token:** `HEKI-CP-2026-11-22-25`

---

## `HEKI-CP-2026-11-22-26` — clarification 1g; the observation space is now honest

**Stage:** foundational · **TASK-03:** `contract stack + 1g` · **implementation: NOT authorized · deletion: UNAVAILABLE**

### "Intact" was a claim the record could not support

A child count is the only tree baseline stored (`entries_at_acquisition`). Count equality cannot mean "unchanged" —
one child out, one child in, same count. So five properties are now separated (count, entry set, child identity,
content, attribution) and **only the count is establishable**. The four-condition matrix became a 20-cell
state × observation matrix, `aborted`/`partial` were redefined to assert only what the evidence supports, and
`not_completed` was added for the removing-side case with no observed difference.

A manifest was assessed and declined: storage, acquisition cost, durable child names, and a second corruption path —
for a property that would still not attribute a later change to this operation.

### The reuse contradiction, resolved without reassuring anyone

1f said stale IDs fail because the sequence increments, and separately listed sequence wrap as an explanation for a
directory reappearing. Both true; together they overclaimed. NTFS is a 48-bit MFT index plus a 16-bit sequence, so a
stored ID can match a different object after 65,536 reuses of that record. "Stale IDs always fail" is wrong.

What survives is the asymmetry, and it is *proved* rather than asserted: `complete` requires a fresh open to **fail**,
so reuse can only make an ID match and produce a present-side outcome. Reuse and wrap can therefore only under-claim.
The clone case is openly admitted and kept separate — two different statements that 1f had blurred.

### Volume continuity: four layers established, one admitted missing

Locator, volume identity, filesystem identity and directory identity are established. **Physical device identity is
not**, so a different disk presenting matching GUID, serial and IDs is undetectable. `complete` is locator-scoped, and
that is stated rather than implied. Also corrected: `OpenFileById` needs *any* handle on the volume, not a special
volume handle — 1f implied otherwise.

### Post-terminal contradiction

Receipt returned unchanged, anomaly reported, nothing deleted, nothing re-verified, nothing rewritten. The anomaly is
**transient**: persisting it would put a mutable second history beside immutable receipts and add a third evidence
kind, which the existing two do not need.

### Verification

```
rows 528-536 all OK   536 (1g): 1db25754a88b   535 (1f) unchanged: d386317a316f
```

31 `S-` clauses, 18 `ST-` tests. Corpus **828**; placeholders `0809`–`0828`; registry 827 besides / 879 files.
**Zero Go files modified; no ledger created; no migration; no deletion enabled.** No test executed.

D-10 and G-1 unchanged. No new owner decision.

**Resume token:** `HEKI-CP-2026-11-22-26`

---

## `HEKI-CP-2026-11-22-27` — clarification 1h; absence must be proved, not inferred from a failure

**Stage:** foundational · **TASK-03:** `contract stack + 1h` · **implementation: NOT authorized · deletion: UNAVAILABLE**

### The sharpest evidence this round: absence and bad input may share an error code

A filesystem-filter engineering source records that an open by a **non-existent ID** fails with
`STATUS_INVALID_PARAMETER` — the same status a malformed request produces — and Win32's `OpenFileById` page documents no
codes at all. So 1g's O1 ("not reopenable") was unsound: a failed open cannot establish absence. Absence now needs five
preconditions, including a **control open** on a known-live identity with the same call shape, which is what separates
"the ID is gone" from "my request is wrong".

### ABA breaks the attribution proof, and the fix costs attribution

A removed → ID reused by B → B removed elsewhere → open fails. Absence survives; *authorship* does not. `complete` is
split from attribution, `reuse_excluded` requires a USN-window proof that Phase 1 lacks, so **no Phase 1 receipt
attributes removal to the operation**. That is a real loss of evidence value, taken because the alternative is a claim
the evidence cannot carry. Creation-time corroboration is also corrected: it can only downgrade a present-side
observation, never strengthen an absent-side one — the object is gone, so there is no time to read.

### Count deltas named honestly

`partial` → `count_delta_unattributed`, attribution always unknown, with an explicit list of what it does not establish
(which children changed, that the rest is intact, that we caused it). Equal counts remain never "unchanged".

### D-10 resolved without weakening it

Continuity is a **held handle**, not matching metadata. The destructive act happens on the handle obtained from the
acquired identity, in the process that acquired it, so D-10 is satisfied by construction; and because reconciliation
after a restart has no held handle, it may **never** write `complete`. That single rule bounds the clone exposure that
1g could only apologise for. Clone cases are now separated by detectability, and unsupported scenarios
(FAT/SMB/no capability flag/expected cloning/removable volume) are declared unavailable rather than degraded.

### Verification

```
rows 528-537 all OK   537 (1h): c247c9d34972   536 (1g) unchanged: 1db25754a88b
```

31 `U-` clauses, 20 `UT-` tests. Corpus **829**; placeholders `0809`–`0829`; registry 828 besides / 880 files.
**Zero Go files modified; no ledger created; no migration; no deletion enabled.** No test executed; E-24…E-27 name the
evidence required, and E-24 decides whether attribution can ever be restored.

No new owner decision. G-1 stands.

**Resume token:** `HEKI-CP-2026-11-22-27`

---

## `HEKI-CP-2026-11-22-28` — clarification 1i; a marker is not a handle

**Stage:** foundational · **TASK-03:** `contract stack + 1i` · **implementation: NOT authorized · deletion: UNAVAILABLE**

### The contradiction was real and it was the dangerous kind

1h asserted `removed` + K1 → `complete` unconditionally, and separately that reconciliation never writes `complete`. The
gap is a durable `removed` marker with no terminal receipt and a dead process. If the marker counted, a recovery path
would claim **stronger** continuity than the live operation ever held — exactly the failure this contract exists to
prevent. A marker records a past observation; it does not recreate a handle.

### Fixed by adding a dimension

Every observation is now qualified by execution context — `original` (acquisition handle still held, same process) or
`reconciliation`. The matrix is indexed by (state, context, observation), and **one cell in the whole table can produce
`complete`**: `removed` + K1 + `original`.

### Window 11 and `completion_unconfirmed`

New crash window: `removed` durable, terminal unwritten, process dies → restart with K1 yields
`completion_unconfirmed` — absence seen, removal call returned success, `continuity: none_after_restart`,
attribution `unknown`, immutable, never upgraded. Kept distinct from `indeterminate` because that one means the evidence
*conflicts*; here it is merely *incomplete*. No new ledger or evidence kind: five and two still hold.

The cost is stated rather than softened: a crash in window 11 permanently costs that operation's completion claim.

### Claims demoted to their real strength

U-31 now bounded to one continuously held volume object and the exact acquired identity — reopened locators,
replacement volumes and post-restart observations are out of scope, not merely unlikely. U-19/U-20 are design intent
gated on E-27, with the specific point that a held volume handle proves **nothing** if the directory is later resolved
by path; it must be opened and identified through that handle.

### Verification

```
rows 528-538 all OK   538 (1i): 3f00fc6cfa6d   537 (1h) unchanged: c247c9d34972
```

23 `V-` clauses, 12 `VT-` tests. Corpus **830**; placeholders `0809`–`0830`; registry 829 besides / 881 files.
**Zero Go files modified; no ledger created; no migration; no deletion enabled.** No test executed; E-24…E-27 unchanged.

No new owner decision. G-1 stands.

**Resume token:** `HEKI-CP-2026-11-22-28`

---

## `HEKI-CP-2026-11-22-29` — clarification 1j; 1i contradicted 1e, and the sequence was under-specified

**Stage:** foundational · **TASK-03:** `contract stack + 1j` · **implementation: NOT authorized · deletion: UNAVAILABLE**

### Check one found a defect, not a gap

The ordering was already fixed in **1e Q-09**: verify identity absent → `removed` → terminal, and `removed` only after
absence is shown. But 1i V-03 described the marker as evidence "that a removal was attempted and that its call
returned" — precisely what it would record if written *before* verification. Read literally, that would let a marker
plus a later observation be arranged into a completion claim. Corrected: the marker records the successful call **and**
the authorizing observation, in that order, and `complete` is a fresh in-process claim never reconstructed from a marker.

Fixing the ordering exposed a crash position the set had missed: removal returned, verification not done, no marker.
That is state `removing` + K1 → `absent_unverified` — already defined, so no new outcome was needed.

### Check two: V-17 stopped at step one

It required opening and identifying the directory through the held handle but said nothing about how children are
reached. Four required properties now stated — open by ID through the held handle and confirm via `FileIdInfo`,
enumerate through that directory object, remove through handles derived from it, verify the same identity absent by ID on
the same handle — with an explicit ban on path re-resolution at any step. A mixed sequence is not a weaker guarantee;
it is a different, unsupported design. E-27 splits into volume-object continuity and **E-28** (can the four steps be
performed at all), because it was carrying two questions. A partial demonstration enables nothing under G-1.

### Drafting stops here

No new outcome, state, kind or owner decision. The next evidence must be a bounded feasibility investigation of the
four-step sequence — not another document.

### Verification

```
rows 528-539 all OK   539 (1j): c02305359ab3   538 (1i) unchanged: 3f00fc6cfa6d
```

17 `W-` clauses, 9 `WT-` tests. Corpus **831**; placeholders `0809`–`0831`; registry 830 besides / 882 files.
**Zero Go files modified; no ledger created; no migration; no deletion enabled.** No test executed; nothing validated
against Windows.

No new owner decision. G-1 stands.

**Resume token:** `HEKI-CP-2026-11-22-29`

---

## `HEKI-CP-2026-11-22-30` — branch closed at document level; nothing else drafted

**Stage:** foundational · **TASK-03:** `contract review CLOSED` · **implementation: NOT authorized · deletion: UNAVAILABLE**

### Disposition

Clarification 1j is accepted at document level, subject to registry verification. That closes the contract review
series 1 → 1A → … → 1J. **No further contract text was written and none is warranted**: the remaining questions are
not about wording.

What the series established, in order:

| Round | Defect closed |
|---|---|
| 1f | volume reacquisition, absent gate, crash windows |
| 1g | count equality is not "intact"; reuse under-claims only |
| 1h | a failed open is not absence; ABA kills attribution, not absence |
| 1i | a durable marker is not a handle; execution context is part of the observation |
| 1j | `removed` means call **and** authorising observation; the sequence is four handle-bound steps |

### Verified state

```
rows 528-539 all verify against disk
539 (1j) c02305359ab3   538 (1i) 3f00fc6cfa6d   537 (1h) c247c9d34972
corpus 831 · placeholders 0809-0831 · registry 830 besides / 882 files
```

**Zero Go files modified. No ledger created. No migration. No persistence adoption. No deletion enabled.** No test
executed at any point in this series; every behavioural claim is labelled `Unestablished`.

### Next stage, recorded but NOT started

A bounded feasibility investigation, chartered as: E-27 and E-28 investigated **separately**; establish the exact
supported Windows APIs and required access rights; prove the **complete four-step sequence**, not isolated capabilities;
use a disposable test directory and a disposable test volume; exercise interruption and replacement scenarios; keep
disjoint from the failing `E:` drive and from all preservation data. Deletion stays unavailable unless the whole
acceptance gate passes.

Prerequisites this stage needs from the owner, not yet given: authorisation to run destructive spikes at all, and
elevation for disposable-volume creation. Neither is assumed here.

**TASK-03 remains in preparation. Not accepted, not closed.**

**Resume token:** `HEKI-CP-2026-11-22-30`

---

## `HEKI-CP-2026-11-22-31` — spike specification filed; investigation not started

**Stage:** foundational · **TASK-03:** `contract closed, investigation specified` · **implementation: NOT authorized · deletion: UNAVAILABLE**

### What was filed

`docs/architecture/TASK-03-SPIKE-SPEC-E27-E28-2026-11-22.md`, row 540 — 39 `SP-` clauses, six `EV-` criteria. The
charter turned into something executable-by-someone-else without further decisions.

### The parts that needed care

**Three permissions, kept apart.** Destructive spikes on disposable data, elevated VHD creation, and deletion
enablement are separate grants, all three currently **not given**. SP-03 says explicitly that permission to investigate
would not enable anything.

**E-27's criteria distinguish prevention from redirection.** If holding the volume handle causes the unmount to *fail*,
that is the platform enforcing continuity and it passes (EV-27c). Only a by-ID call returning the clone's object
falsifies the design. The two look like "something went wrong" and are opposite findings.

**E-28 detects path re-resolution per step**, with EV-28c capturing the sequencing hazard: a four-step sequence that
cannot be performed *in order* fails even if a reordered variant works, because the order is what makes the verification
meaningful.

**Neither result contaminates the other**, and only pass/pass makes `complete` reachable. Pass/fail leaves the design
unimplemented with no path fallback; fail/pass leaves `completion_unconfirmed` carrying the space permanently.

**The dangerous case stops the run.** If a call suggests the wrong volume was reached, the spike deletes nothing
further and reports first, so the evidence survives cleanup.

### Verification

```
rows 528-540 all OK   540 (spec): ffe7f28547a1   539 (1j) unchanged: c02305359ab3
```

Corpus **832**; placeholders `0809`–`0832`; registry 831 besides / 883 files.

**Nothing was executed: no VHD created, no volume attached, no file deleted, no Go file modified.** E-27, E-28 and
E-25 remain unestablished. No test run.

**Resume token:** `HEKI-CP-2026-11-22-31`

---

## `HEKI-CP-2026-11-22-32` — first executed evidence: the absence signal is 87, and reuse is immediate

**Stage:** foundational · **TASK-03:** `spike 01 executed` · **implementation: NOT authorized · deletion: DISABLED**

### Authorization used, and its limits

Destructive spikes on disposable data. Run 2 deleted three files it had just created inside
`build/tmp/spike/fixture/`, then removed the directory. Residue none, verified. No Howl data, no conversation directory,
no `E:` access. **Elevation was not used**, so E-27 has not been attempted.

### The finding that changes the design's cost

**A directory handle works as the `OpenFileById` hint, unelevated.** E-28 step 1 is feasible without administrator
rights, so 1j W-09 step 1 ("held volume handle") is satisfiable by something the app can open for itself.

### The finding that changes how absence must be read

A **deleted** identity returns `ERROR_INVALID_PARAMETER` (87), not "file not found". Out-of-range IDs return 2. A
malformed descriptor returns 87 as well. So one failure code covers absence, an unallocated ID and a broken request — and
the only thing that separates them is a successful control open of the exact same shape. 1h U-02 was written as caution;
it is the sole discriminator, and it works.

Corollary recorded, not applied: **U-02 clause 5 is vacuous**, because the control succeeds and therefore any failure
satisfies it. The load-bearing clauses are 2–4.

### Reuse, measured

Deleted `a.txt` (index 769116, seq 19) → the next create took **index 769116, seq 20**. Slot reuse is immediate; the
sequence is what keeps the old ID dead. Identical-ID collision therefore needs the 16-bit sequence to wrap — 65,536 frees
of one record. ABA is narrower than 1g implied and still unexcludable, so U-11 stands: no attribution in Phase 1.

### A bug caught before it became a finding

The first run said `FILE_SUPPORTS_OPEN_BY_FILE_ID = false`. Wrong constant: `0x00080000` is `FILE_READ_ONLY_VOLUME`; the
real flag is `0x01000000`. Re-run: set. Had it gone unchecked, it would have justified declaring NTFS unsupported.

Also recorded: `GetVolumeInformation` gives a 32-bit serial, `FILE_ID_INFO` a 64-bit one — comparing them directly would
mark every identity volume-mismatched.

### Verification

```
rows 528-541 all OK   541 (findings 01): 2bc92453203f   540 (spec) unchanged: ffe7f28547a1
go build ./... OK · zero product Go files modified · fixture residue: none
```

Corpus **833**; registry 832 besides / 884 files.

**Verdicts:** E-25 partially established (one build, one volume); E-28 step 1 established unelevated; steps 2–4, E-27 and
E-26 unestablished. No contract row edited. Deletion remains disabled.

**Resume token:** `HEKI-CP-2026-11-22-32`

---

## `HEKI-CP-2026-11-22-33` — the absence discriminator does not exist, and presence is not self-authenticating

**Stage:** foundational · **TASK-03:** `spike 01a executed` · **implementation: NOT authorized · deletion: DISABLED**

### F-3 was wrong, and the reviewer was right to block on it

Control success proves the mechanism works for a live identity. It does not prove a target failure means absence. My
own probe contained the counter-example — 87 for both a deleted identity and a malformed request — and I drew the wrong
conclusion from it.

Raw `NtCreateFile` statuses settle it at the lowest layer:

| Case | NTSTATUS |
|---|---|
| deleted identity | `0xC000000D` INVALID_PARAMETER |
| fabricated all-zero identity | `0xC000000D` |
| zero-length malformed name | `0xC000000D` |

A deleted identity is **indistinguishable** from a broken request. `0xC0000034` appears only for out-of-range IDs;
`0xC0000033` for a different malformed shape. So the codes do not partition "missing" from "broken" at all.

**Absence by open-by-ID is `unverifiable`, not weakly established.** 1j W-09 step 4 therefore cannot produce a
determinate result, and `complete` is unreachable by the specified mechanism.

### The finding that is larger than the absence problem

`D:\`'s root identity opened **successfully** through a `C:` hint, returning an unrelated `C:` object with the same
128-bit bytes and a different serial — both volumes number a record 5 with sequence 5. An ID is unique per volume, not
globally. **Presence is not self-authenticating**: opening "our" ID through the wrong volume returns success for
someone else's object, and nothing in the status says so.

This promotes F-6 (32-bit vs 64-bit serial) from implementation hazard to correctness requirement: every returned
handle must be checked for both ID and 64-bit serial. It also shows volume mismatch is ordinary, not clone-specific —
which weakens 1i's reassurance that clone cases are exotic.

### What survives

E-28 step 1 stands for this configuration, conditional on that check. E-25 is reclassified as unverifiable by this
mechanism. Three candidate positive mechanisms for absence are listed (parent enumeration, USN journal, held-handle
observation), all Unestablished; parent enumeration is cheapest and independent of open-by-ID.

### Verification

```
rows 540-542 all OK   542 (findings 01a): 249fdb58f2af   541 (findings 01) unchanged: 2bc92453203f
go build ./... OK · product Go files modified: 0 · residue: none
```

Corpus **834**; registry 833 besides / 885 files. No contract row edited. No `E:` access. `D:` opened read-attributes
only. Deletion remains disabled.

**Resume token:** `HEKI-CP-2026-11-22-33`

---

## `HEKI-CP-2026-11-22-34` — removal needs names, verification needs IDs

**Stage:** foundational · **TASK-03:** `spike 02 executed` · **implementation: NOT authorized · deletion: DISABLED**

### Three results, kept apart

**E-28 step 2 PASS** — `NtQueryDirectoryFile` on the directory handle listed all five children with file references, no
path involved. The Win32 wrapper returned `ERROR_NO_MORE_FILES` on the first call for a non-empty directory here; the
native call did not. Recorded as environment behaviour, not a general claim.

**E-28 step 3 PASS, but not the way 1j W-09 assumed.** Handles opened via `OpenFileById` **cannot be disposed** — both
`FileDispositionInfoEx` and `FileDispositionInfo` return 87 on by-ID handles while the same file disposes fine through a
path-opened handle. The working derivation is a **parent-relative name open** (`NtCreateFile` with
`RootDirectory = <directory handle>`): 5/5 removed, 5/5 identities verified by ID *and* 64-bit serial, zero path strings.

So the sequence is: open directory by ID, enumerate by handle, open child by relative name, **verify its ID**, dispose.
Removal needs names; verification needs IDs. The design does not get to use one mechanism for both.

**P1 FAILS, decisively.** A directory moved into a sibling is absent from the parent enumeration exactly as a deleted
one is, while its identity stays openable. Enumeration establishes "not a child of this parent", not absence — a moved
directory would be reported as removed. P1 is refused as an absence proof.

### Unchanged by any of this

Step 4 has no positive mechanism, `complete` is still unreachable, and E-27/E-26 remain unestablished. Removing a
directory object itself was not attempted — only files — so that part of step 3 is untested.

### One thing worth admitting

A path-opened **diagnostic** handle deleted a fixture file during what was labelled an enumeration experiment. The probe
reported the count mismatch instead of hiding it, but the lesson stands: in a deletion investigation, a diagnostic that
changes state is destructive. SP-35's discipline is what surfaced it.

### Verification

```
rows 540-543 all OK   543 (findings 02): 3a2ea810c872   542 (01a) unchanged: 249fdb58f2af
go build OK · product Go files: 0 modified · residue: none
```

Corpus **835**; registry 834 besides / 886 files. No contract row edited. No `E:` access, no elevation.

**Resume token:** `HEKI-CP-2026-11-22-34`

---

## `HEKI-CP-2026-11-22-35` — the directory object can be removed the same way, with a control that makes it mean something

**Stage:** foundational · **TASK-03:** `spike 03 executed` · **implementation: NOT authorized · deletion: DISABLED**

### One question, one verdict

**Directory-object removal: PASS** via parent-relative open + identity verification + disposition. The directory was
opened by name relative to its parent handle with `FILE_DIRECTORY_FILE`, its ID and 64-bit serial matched the values from
the parent enumeration, and then it was disposed through that same bound handle.

### The step that makes the verdict worth having

Disposing it **while non-empty failed with `ERROR_DIR_NOT_EMPTY` (145)**. Without that step, "dispose returned success"
would prove only that an API said yes. With it, the call rejects what it should reject and accepts what it should accept.

Step E adds the other half: the parent still listed the directory immediately after its child was removed. The
children-first ordering is enforced by the filesystem, not by our convention.

### by-ID is still not a route

A live **empty** control directory opened by ID failed to dispose with 87 — same as files (02 F-10). So the relative-name
open is the only route that works, for directories as well as children.

### Corrected step 3, complete

open directory by ID → enumerate by handle → per child: open relative, verify ID, dispose → dispose the directory
through its parent-relative handle, after it is empty. Every destructive act is preceded by an identity check on the very
handle being used.

### Still outside the verdict

Step 4 remains unresolved, so `complete` is unreachable. Nested depth, large directories, partial mid-pass failure,
reparse points, and foreign handles are untested. E-27/E-26 unestablished. P2/USN not investigated, by instruction.

### Procedure

Diagnostics were read-only by construction — `os.Stat` only. That rule came from 02 F-15, where a path-opened control
handle deleted a fixture file.

```
rows 540-544 all OK   544 (findings 03): bc7dcca3a22d   543 (findings 02) unchanged: 3a2ea810c872
go build OK · product Go files: 0 modified · residue: none
```

Corpus **836**; registry 835 besides / 887 files. No contract row edited. G-1 intact.

**Resume token:** `HEKI-CP-2026-11-22-35`

---

## `HEKI-CP-2026-11-22-36` — depth adds no mechanism; a truthful stop is observable

**Stage:** foundational · **TASK-03:** `spike 04 executed` · **implementation: NOT authorized · deletion: DISABLED**

### Both experiments passed, after two harness defects were found and named

The interrupted run had produced a log reading NESTED: FAIL and PARTIAL: FAIL. Both failures were mine, not the
platform's:

| Defect | What it produced | Fix |
|---|---|---|
| unwound a level while its handle stayed open | level 2 failed 145 → **false "nested FAIL"** | close before disposing the parent |
| obstruction was a handle with share mode `0x03` | it **did not obstruct**, so all three children were removed and nothing was tested | replaced with a non-empty child directory (145, already proven) |

The second is the worse class of error: an obstruction that doesn't obstruct turns a test into a tautology that passes
or fails for reasons unrelated to the question.

### The rule N produced

**An unlinked-but-still-open child directory remains visible to its parent.** Disposing level 3 and holding its handle
made level 2 non-empty. POSIX-semantics deletion alone is not enough — the handle must be released before ascending.
With that ordering, each parent reported 0 entries exactly when its child closed, and the top ended empty.

### What P produced

One child removed, the next blocked by 145, the pass stopped, and the state told the truth: exactly the unreached
remainder left, directory identity unchanged, directory non-removable, `os.Stat` true. Recovery then finished from a
fresh enumeration of the same handle — **no journal of prior progress needed**, which is what a crash path would want,
while saying nothing at all about absence.

### An anomaly, recorded not concluded

A held handle whose share mode excluded `FILE_SHARE_DELETE` did **not** prevent disposition of that file. That
contradicts a plausible reading of share-mode semantics. Testing "child open elsewhere" needs a different process, not
another handle in this one.

### Verification

```
rows 540-545 all OK   545 (findings 04): e10859b7b016   544 (findings 03) unchanged: bc7dcca3a22d
go build OK · product Go files: 0 modified · residue: none
```

Corpus **837**; registry 836 besides / 888 files. Step 4 unresolved, `complete` unreachable. No contract row edited.
No `E:` access, no elevation. G-1 intact.

**Resume token:** `HEKI-CP-2026-11-22-36`

---

## `HEKI-CP-2026-11-22-37` — P2 is blocked before the question can be asked

**Stage:** foundational · **TASK-03:** `spike 05 executed` · **implementation: NOT authorized · deletion: DISABLED**

### The measurement that ended the round early

USN journal FSCTLs are volume operations, so they need a volume handle. Three masks, three refusals:

| Mask on `\.\C:` | Result |
|---|---|
| `FILE_READ_ATTRIBUTES \| FILE_READ_DATA` | denied, 5 |
| `FILE_READ_DATA` | denied, 5 |
| `GENERIC_READ` | denied, 5 |

The handle itself was refused, so `FSCTL_QUERY_USN_JOURNAL` was never reached. **USN evidence is unavailable to an
unelevated process**, which means adopting P2 would require Howl to run elevated merely to confirm a deletion it had
already performed. That is a product-level cost, not a contract detail.

Forcing wrap or testing gap recovery needs `FSCTL_CREATE/DELETE_USN_JOURNAL` — volume metadata, not disposable data —
so that was not attempted.

### The convergence worth acting on

E-27 needs an elevated disposable volume for detach/clone tests. P2 needs an elevated volume to read the journal and a
disposable one to manipulate it. **One approved elevated session on a scratch VHD would serve both**, and neither can
progress without it.

### What was specified instead of measured

Five things USN evidence would have to show — journal ID stability, a captured start USN, `LowestValidUsn ≤ start`, a
`FILE_DELETE` record whose `FileReferenceNumber` equals the low 64 bits of the acquired ID, and restart continuity — all
Unestablished. Plus the two rules that must survive any later measurement: **a missing record proves nothing**
(disabled, wrapped, recreated or wrongly-started journals all say "no match" about a directory that really was deleted),
and **even a passing P2 would not restore attribution**.

### Deferred, not dropped

The share-mode anomaly is recorded with the five distinctions it must test (separate process, both flag sets,
disposition API, holder actually open, exact status). Foreign-handle interference stays unestablished.

```
rows 540-546 all OK   546 (findings 05): 6ddc91818f55   545 (findings 04) unchanged: e10859b7b016
go build OK · product Go files: 0 · read-only, no journal modified
```

Corpus **838**; registry 837 besides / 889 files. Step 4 unresolved, `complete` unreachable, G-1 intact.

**Resume token:** `HEKI-CP-2026-11-22-37`

---

## `HEKI-CP-2026-11-22-38` — deletion stops needing to be proved

**Stage:** foundational · **TASK-03:** `direction changed` · **implementation: NOT authorized**

### The pivot

Howl stops trying to prove a deletion happened. It deletes as it always did, writes a **reconciliation worklist**
entry (JSON), and a scanner observes whether the conversation's required data is still there. If it is gone, the
conversation leaves the UI and the index; the entry is retired.

The worklist is not a receipt. Nothing claims authorship, nothing claims the filesystem cooperated, nothing claims a
stronger guarantee than the app actually has.

### What this retires

The entire privileged-verification apparatus becomes unnecessary for normal deletion: no USN dependency, no elevated
volume handle, no identity-proof requirement, no E-25/E-27/E-28 gate on this workflow, no receipt asserting more than
is known. The reason is worth stating precisely, because it is not a retreat from rigour: the old question was
*"was this exact identity destroyed?"* — which needs privileged identity evidence — while the new question is
*"is this conversation's data still available?"* — which a path check answers, and which is the only question the UI
actually needs answered.

TASK-03 rows 528–546 stay as they are: a complete record of a harder problem that no longer needs solving. Not
abandoned — retired.

### Cadence recommendation (the detail asked for)

**Startup plus periodic, not startup-only.** Startup-only fails the ordinary case of an app left open for days: the
user deletes a conversation, the files go, and the UI keeps offering it until the next launch. Bounded work per pass,
fixed pass interval with jitter, retry with backoff for records that error, and a scan also on resume-from-sleep.

### Two gaps in the proposed table

1. **No row for partial absence.** A deletion that half-succeeded leaves some required paths gone and others present.
   The four-row table maps that to "required data exists → keep", which silently keeps a half-destroyed conversation.
   It needs its own row: data partially absent → mark as damaged, never silently keep, never auto-repair.
2. **A worklist-only scan cannot see orphans.** If the worklist write fails, or an entry is lost, an index entry whose
   data is gone is never revisited. A low-frequency **full sweep** (index vs disk) is needed as a backstop, and it is
   read-only like the rest.

Both are decisions for the owner, not assumptions to bury in a rewrite.

**Resume token:** `HEKI-CP-2026-11-22-38`

---

## `HEKI-CP-2026-11-22-39` — replacement contract drafted; the old series becomes history

**Stage:** foundational · **TASK-03:** `reconciliation contract drafted` · **implementation: NOT authorized**

### Filed

`docs/architecture/TASK-03-RECONCILIATION-CONTRACT-2026-11-22.md`, row 547 — 29 `RC-` clauses, 16 `RT-` tests, five open
questions. Rows 528–546 unedited: a complete record of a harder problem, retired rather than corrected.

### Grounded in the real layout, not in my assumptions

The manifest had to come from the storage contract, so I read it: `meta.json` always written, `history.ndjson` the
source of truth with legacy `history.json`, and `world`/`scenario`/`characters`/`lorebooks`/`memories` written
conditionally. That is why the manifest has three kinds of entry — required, alternatives, optional — instead of one
list. RC-07 is the rule that matters most from it: a manifest generated from disk would let a damaged conversation
redefine its own requirements and look healthy.

### The six states, and the transitions that are forbidden

Healthy, Pending reconciliation, Missing, Damaged, Unavailable, Orphaned — all **derived**, none stored. Six forbidden
transitions, of which two carry the old discipline forward intact: an error is never Missing, and Unavailable never
becomes Healthy. Damaged exists because a failed deletion produces exactly that shape, and both of the obvious wrong
answers (keep it looking fine, or discard the index entry) destroy something real.

### What got simpler

No receipt, no authorship, no elevation, no USN, no identity proof. Worklist ordering became a *speed* decision rather
than a correctness one (RC-19), because the full sweep is authority and observation is idempotent. That single property
is what removes the need for the journal the earlier contract was built around.

### What survived, deliberately

RC-23: deletion is confined to the conversation's own subtree — never a parent, never an uncontained path. That is the
one identity rule kept, and it constrains what may be *deleted* rather than what may be *claimed*.

### Verification

```
rows 540-547 all OK   547 (contract): 79d9de841191   546 (findings 05) unchanged: 6ddc91818f55
go build OK · product Go files: 0 modified
```

Corpus **839**; registry 838 besides / 890 files.

**Resume token:** `HEKI-CP-2026-11-22-39`

---

## `HEKI-CP-2026-11-22-40` — OQ-3 closed, and the placeholder file that would have hidden a deletion

**Stage:** foundational · **TASK-03:** `reconciliation contract amended` · **implementation: NOT authorized**

### The measurement

Scratch root, `NewPersistenceManagerWithDir`, one `NewChat`, one `AppendHistoryEntry`:

| Moment | Contents |
|---|---|
| after `NewChat` | `meta.json` (421 b), **`history.json` = `[]` (2 b)**, characters/lorebooks/memories — **no `history.ndjson`** |
| after first append | `history.ndjson` appears (132 b); **`history.json` still `[]`** |

`world.json`/`scenario.json` absent while nil, as the source said.

So the proposed rule holds and is now grounded: **empty conversations are Healthy; history is required only after content is committed.**

### The finding I did not expect

`history.json` is a lagging snapshot. That makes row 547's "at least one of ndjson or history.json" wrong **in both directions**:

- an empty conversation would be called Damaged for lacking a representation it never needed;
- a conversation whose `history.ndjson` was **deleted** would be called **Healthy**, because a 2-byte `[]` placeholder still sat there.

The second one is the whole point of the reconciliation system, defeated by an empty file. RC2-05 now scopes the accepted representation by contract version, and RC2-06 makes **non-empty a property of content, not existence**.

### Two more consequences

Expectation flips **after** the commit, so a crash between append and mark leaves Healthy rather than a false Damaged — the failure direction is chosen deliberately. And `UpdatedAt` is explicitly not the expectation, because `SaveChat` rewrites it for empty conversations too.

`meta.json` absent in a present directory is **Damaged, not Missing**: the messages may still be on disk, and calling it Missing would tell the UI to drop them.

### Verification

```
rows 546-548 all OK   548 (OQ-3): a40e5bfafd04   547 (contract) unchanged: 79d9de841191
go build OK · product Go files: 0 · scratch fixture removed, residue none
```

Corpus **840**; registry 839 besides / 891 files. Open: OQ-3a non-append writers, OQ-3b exact empty threshold, OQ-3c dual-state legacy conversations.

**Resume token:** `HEKI-CP-2026-11-22-40`

---

## `HEKI-CP-2026-11-22-41` — ten write paths, and the duplicate that writes no log

**Stage:** foundational · **TASK-03:** `writer-lifecycle inventory complete` · **implementation: NOT authorized**

### Inventory

P1 creation, P2 appends, P3 edit/delete, P4 clear, P5 final save, **P6 duplication**, P7 legacy migration, P8 wipe-restore, P9 import (does not exist), P10 delete. Each recorded with writer, representation, commit boundary, expectation, crash window.

### The finding that changes the contract

`SaveChatAs` copies content and calls `SaveChat`, which writes `history.json` and **never** the append log. Observed:

```
duplicate   history.json 503 b   history.ndjson ABSENT   LoadChat → 3 messages
```

A legitimate current-contract conversation with real content and no `history.ndjson`. Row 548's version-scoped representation rule would have called **every duplicate Damaged**. Superseded: acceptance is content-based now — any representation with at least one record, and the 2-byte `[]` is existence, not content.

That correction also fixes the row-548 worry properly: ndjson deleted but populated `history.json` present is **Healthy**, because the content is still loadable. Damage means content unavailable, not one file missing.

### Creation collides

`generateChatID` is one-second resolution, `NewChat` uses `MkdirAll` without exclusivity. Observed twice after two creations. The first probe run showed it for real: `SaveChatAs` inside the same second **reused the source directory**, overwrote its `meta.json` (422 → 425 b), and left the source's 375 b log in place — one directory, two conversations, content and identity disagreeing.

### The observation gap

Not hand-waved: the marker is **explicitly not claimed crash-safe**. Ordering stays append-then-mark (reversing it trades a silent failure for a loud false one), and the gap becomes representable instead of invisible through a three-valued marker written **before** the append:

| Marker | Content | State |
|---|---|---|
| `none` | none | Healthy |
| `none` / `pending` | present | Healthy, promoted to `committed` |
| `pending` | none | **Pending reconciliation** |
| `committed` | none | **Damaged** |

`pending` can never produce a false Damaged. Crash-and-vanish resolves to a visible question rather than silent reassurance. Promotion is monotonic, and the scan writes the marker only.

### Verification

```
rows 546-549 all OK   549 (OQ-3a): c2be874e51de   547 unchanged: 79d9de841191
go build OK · product Go files: 0 · both scratch fixtures removed, residue none
```

Corpus **841**; registry 840 besides / 892 files. Unestablished: collision frequency, `SaveChatAs` pointer intent, future import reuse, exact empty threshold.

**Resume token:** `HEKI-CP-2026-11-22-41`

---

## `HEKI-CP-2026-11-22-42` — three conclusions kept apart: content, marker, identity

**Stage:** foundational · **TASK-03:** `OQ-3b settled; identity defect separated` · **implementation: NOT authorized**

### Two rows, deliberately not one

Row **550** is the creation identity defect, on its own, because it is not a manifest adjustment: one-second IDs, non-exclusive `MkdirAll`, unconditional `SaveChat`, global `LastChatID`. Four normative invariants recorded (unique ID, exclusive creation, never adopt, fail loudly), and reconciliation's role limited to detection — never repair, rename, split, or delete either party's data.

Row **551** is OQ-3b.

### The rule, stated once

The question is **can Howl still load the committed history**, not whether a file exists. `history.json = []` is not content; a populated snapshot is valid without ndjson and vice versa; missing one is not damage; none-with-expectation is Damaged. Row 548's version-scoped rule stays explicitly superseded.

### What empty and malformed actually look like

`[]` = 2 bytes, `RewriteHistory([])` = 0 bytes, one newline = 1 byte and is neither.

The harder result: **malformed is indistinguishable from empty.** Truncated json, a non-array object and an array of strings all load as 0 messages with nil error; 3 records cut mid-record loads as 2; 2 good + corrupt + 2 good loads as 2. The loader is lossy by design — right for the app, useless as an observer. So `Healthy` means *the loadable prefix is present*, and promotion requires a successful parse yielding a record. Malformed never promotes.

One anomaly was my own fixture bug, checked in isolation before it reached a document.

### Cost, measured rather than assumed

100 → 0.50 ms · 1,000 → 1.54 ms · 10,000 → 10.2 ms · 100,000 → 104.7 ms, against `stat` at 36–39 µs and a final-byte read at 2–3 µs. Availability needs no parse; counting does. The scan is tiered, and the numbers are a measurement, not a budget.

`LoadChat` left the directory byte-identical over three loads: safe to call, still not authoritative.

### Marker

Persisted through `DurableStore`, which already has the lock, the manifest commit and `testHookBeforeReplace` for scheduling a death between manifest and replacement — with the marker's own crash behaviour required to be tested, or the ambiguity just relocates.

### Verification

```
rows 547-551 all OK   550 (identity): 9d1b51ff0c46   551 (OQ-3b): 2daf6165a652   547 unchanged
go build OK · product Go files: 0 · all scratch fixtures removed, residue none
```

Corpus **843**; registry 842 besides / 894 files. Next: OQ-3c, then the audit against all three findings.

**Resume token:** `HEKI-CP-2026-11-22-42`

---

## `HEKI-CP-2026-11-22-43` — the loader can prefer the lossier file

**Stage:** foundational · **TASK-03:** `OQ-3c settled; promotion tightened` · **implementation: NOT authorized**

### Four properties, kept apart

Loadability, structural validity, content presence, completeness. `LoadChat` answers only the first, and row 552 promotes that from a caution to a clause: it is a safe observation tool (it mutates nothing) and never a validity test.

### Coexistence measured on ten states

Precedence is `ndjson` if non-empty, else `history.json`. Three counterexamples to "benign":

| State | Result |
|---|---|
| log 1 record, snapshot 3 | loader returned **1**; two records on disk invisible |
| log torn, snapshot valid and fuller | loader returned **the torn prefix** |
| **log 0 bytes, snapshot populated** | returned **2** — deleted content came back |

The third is the `RewriteHistory` two-write window: the log is renamed into place first, the snapshot written second. So a crash there resurrects a clear.

Concluded: precedence is a display rule, not an evidence rule, and freshness cannot be read off size or timestamp — the snapshot lags during appends, the log lags after a tear, and after a clear the log is authoritative while the snapshot is stale.

### Promotion made impossible to satisfy by truncation

Row 551's "parses and yields at least one record" was satisfied by a log with two good records and a corrupt tail. Now: whole-representation validity (every line parses **and** a final newline; whole json document parses as an array), any invalid representation blocks promotion even if another is valid, disagreement uses the **larger** valid set and is reported, and a torn tail is real loss — **Damaged (reason `truncated`)**.

A zero-byte log beside a populated snapshot is `Inconsistent`, never content-present, and nothing is deleted or rewritten.

### Verification

```
rows 547-552 all OK   552 (OQ-3c): 909ac09d5713   547 unchanged: 79d9de841191
go build OK · product Go files: 0 · scratch fixture removed, residue none
```

Corpus **844**; registry 843 besides / 895 files. Next: the audit of row 547 against the three findings. Creation identity stays separate and unauthorized. Two product questions left open rather than folded in: whether the app should keep preferring a torn log, and whether the `RewriteHistory` window deserves its own defect record.

**Resume token:** `HEKI-CP-2026-11-22-43`

---

## `HEKI-CP-2026-11-22-44` — the audit says the contract cannot survive as written

**Stage:** foundational · **TASK-03:** `replacement contract audited` · **implementation: NOT authorized**

### Correction accepted

Count is not provenance. Disagreement fails closed: `Inconsistent`, no winner. Provenance is testable only as an exact **prefix on content in order** — and record IDs cannot carry that test, because they are positional indices (`len-1`, `i`), so deleting a message renumbers everything after it. Even a proven prefix only bounds the possibilities (stale-short vs resurrected-long); it resolves nothing.

### Six questions, one answered badly

Presence, structural validity, content presence, completeness, coexistence, promotion. Row 547 answered presence inadequately and did not address the other four.

### Contradictions, not gaps

`RC-05` calls `history.json` legacy while `SaveChat`, `RewriteHistory` and `SaveChatAs` all write it. `RC-09` demands a history file from an empty conversation. `RC-11` defines Healthy by **path existence**, which a 2-byte `[]` satisfies. `RC-11`'s Missing absorbs a present `meta.json`, contradicting `RC-12` in the same contract. `RC-14` forbids the marker write the design needs. And **`RT-05` is now wrong**: ndjson removed with a populated snapshot is Healthy, because the content still loads.

### The finding the audit produced

`RewriteHistory([])` leaves no trace. So `committed` + empty is indistinguishable from a deliberate clear, and the row-552 table would mark **every cleared conversation Damaged forever**. A clear is a retraction: the marker needs a fourth value.

### Boundaries held

Creation collision (row 550), the loader preferring a torn log (OQ-3c-a) and `RewriteHistory` resurrection (OQ-3c-b) stay separate. The scanner may report all three and must compensate for none — no choosing the good file, no rewriting a log.

### Verification

```
rows 547-553 all OK   553 (audit): 965eb92cc618   547 unchanged: 79d9de841191
go build OK · product Go files: 0
```

Corpus **845**; registry 844 besides / 896 files. Next: v2 as a consolidation, not another addendum.

**Resume token:** `HEKI-CP-2026-11-22-44`

---

## `HEKI-CP-2026-11-22-45` — v2, and the generation that four scalars could not carry

**Stage:** foundational · **TASK-03:** `reconciliation contract v2 drafted` · **implementation: NOT authorized**

### What was consolidated

Row **554**: 50 clauses, 20 tests, six transitions. A consolidation, not an addendum — each clause replaces its predecessor. Row 547 stays unedited as the superseded contract.

### The objection that changed the model

Four marker values were not enough. So the expectation record now holds **state + generation**: a clear closes generation G, the next commit begins G+1, and that is what stops a retracted snapshot being read as the new conversation's history. With one honest limit stated in the contract: the generation **bounds attribution, it does not identify records** — which bytes belong to which lineage is recorded nowhere today, and v2 does not pretend otherwise.

### One rule I had to derive rather than assume

Coexistence cannot be symmetric. `AppendHistoryEntry` writes only the log; `SaveChat` and `RewriteHistory` write both. So `snapshot ⊂ log` is **expected lag** and healthy, while `log ⊂ snapshot` is not expected and is `Inconsistent`. Direction is interpreted from the writer lifecycle, not guessed from count — and a prefix still never proves provenance.

### Six transitions, each with its crash points

Creation · first append · subsequent append · correction/retraction · clear · append-after-clear. The clear table is the one that matters most: intent recorded but rewrite not finished is **Pending reconciliation**, never Damaged and never Healthy-with-content. And append-after-clear is where the generation earns its keep — the retracted snapshot is `retracted-residue`: reported, never compared, never deleted.

### Safeguards written in, not implied

`cleared` grants no cleanup authority. No marker value authorises deleting residue. No state may derive from `LoadChat`'s count (VT-19 is an instrumented assertion). Row 550, OQ-3c-a and OQ-3c-b are named inside the contract with a clause forbidding the scanner from compensating for any of them.

### Verification

```
rows 547-554 all OK   554 (v2): d38124d29335   547 unchanged: 79d9de841191
V2 clauses=50 tests=20 transitions=6
go build OK · product Go files: 0
```

Corpus **846**; registry 845 besides / 897 files.

**Resume token:** `HEKI-CP-2026-11-22-45`

---

## `HEKI-CP-2026-11-22-46` — the gate was missing, and the mixture is real

**Stage:** foundational · **TASK-03:** `v2 amended: mutation gate` · **implementation: NOT authorized**

### The invariant did not hold

Gen G committed, a clear interrupted after its intent was recorded, restart, one ordinary append:

```
AppendHistoryEntry error: <nil>
log: 0=gen-G-0  1=gen-G-1  2=gen-G-2  3=gen-2-first-message
LoadChat -> 4 messages, err=<nil>
```

Two generations in one file, the marker never read. Source agrees: `AppendHistoryEntry` opens with `O_APPEND|O_CREATE` and checks nothing, the chat package has no `pending`/`unresolved`/recovery logic, and `App.startup` does nothing about chat persistence.

### Why it cannot be fixed afterwards

That log is byte-identical to a single-generation conversation of four messages. No boundary, no per-record stamp. So this is **prevention, not detection** — and it makes V2-OQ-1 a correctness question rather than an attribution refinement.

### The amendment

V2-51…V2-66. No mutation while a transition is unresolved; the writer completes or **refuses with a typed failure** (never queues — the base generation is unknown); "resolved" defined per operation (a clear is resolved only when the marker is cleared **and** both representations are valid **and** empty); G+1 accepts writes only after G closes; the app finishes or refuses, the scanner does neither. Promotion now needs generation-safe provenance; the prefix rule is qualified by log validity and transition settlement; the DurableStore test covers the `(state, generation)` pair, not a bare value.

All six audit areas answered. Eight tests, one of which — VT-28 — asserts the contract's limit instead of passing quietly.

### Left honest

The gate suffices for the current architecture **only if the five enumerated mutation paths are a closed inventory**, and a sixth path added later reopens the hole. Per-representation generation stamps are the only option that does not depend on that inventory staying complete.

### Verification

```
rows 547/554/555 all OK   555 (gate): c0d1e3d971e2   554 unchanged: d38124d29335
go build OK · product Go files: 0 · scratch fixture removed, residue none
```

Corpus **847**; registry 846 besides / 898 files.

**Resume token:** `HEKI-CP-2026-11-22-46`

---

## `HEKI-CP-2026-11-22-47` — the boundary the gate needed does not exist

**Stage:** foundational · **TASK-03:** `gate audited against the call graph` · **implementation: NOT authorized**

### Not enforceable as written

`PersistenceManager` has no mutex, no lock file, no atomic check-and-write; `Service.mu` covers in-memory service state only. So "mutations complete or refuse" is a convention with a race — exactly the caller-check-then-write sequence row 555 warned about.

### The inventory was an undercount

Four writers, **at least twelve entry points**: `SaveChat` from `NewChat`, `SaveChatAs`, `RenameChat`, `performSave`, `SaveChatData` and character edits; `AppendHistoryEntry` from user, assistant and audio turns; `RewriteHistory` from edit, clear and final save; `RemoveAll` from delete.

**Three writers sit outside `PersistenceManager` entirely** — memory service (`memories.json`, `memories_backup.json`), auto-capture (`auto_capture_config.json`), continuity (`continuity_config.json`) — and `SaveChat` also writes `memories.json`, so that file has two uncoordinated writers. Declared out of scope, not absorbed.

### Two bypasses

`App.LoadChat` saves the previous conversation before switching, so **a read path mutates**. `WipeAllData` renames the whole `chats` root, which no per-conversation gate can cover.

And row 547 RC-27 ("instances need no coordination") is superseded: with a marker and a gate, two instances interleave a check and a write.

### What the gate now requires

Check-and-write in **one critical section**; a **single** enforcement point; serialization visible to a **second process**. The three preserved requirements restated: idempotent clear retry, `cleared` plus both representations valid and empty, and an inventory rule where a new mutation path must add a gate-coverage test.

### Verification

```
rows 554-556 all OK   556 (call graph): 39fc707703ae   555 unchanged: c0d1e3d971e2
go build OK · product Go files: 0 · no fixture created (source reads only)
```

Corpus **848**; registry 847 besides / 899 files.

**Resume token:** `HEKI-CP-2026-11-22-47`

---

## `HEKI-CP-2026-11-22-48` — the pair is safe if the marker is one document

**Stage:** foundational · **TASK-03:** `DurableStore pair consistency closed` · **implementation: NOT authorized**

### Demonstrated, not assumed

`TestCrashBeforeReplaceLeavesPreviousDocument` kills a real child at the hook and asserts the previous document is intact with `Load -> Recovered(interrupted_commit_rolled_back)`. The live path's only mutation is a rename, so no reader sees a partial file. Both windows are safe: before the rename the previous document survives with a rollback reason; after it the new document matches its manifest and staged leftovers are swept.

**Therefore a mixed `(state, generation)` pair is impossible by construction — iff the marker is one document.** That is now a schema constraint (V2-85), not a test hope.

### Two requirements this produced

A `Recovered` marker read is internally consistent but **possibly stale** (`Current()` is false for it), so it never counts as evidence that a transition is settled — the gate refuses, reconciliation reports Pending.

And the boundary needs **compare-and-swap under one lock hold**: `Load` is lock-free, `Save` takes the lock internally, so read-then-save is precisely the caller-side check V2-74 forbids. `Save` cannot express conditional replace — a design requirement, not a store defect.

### Verification

```
rows 555-557 all OK   557 (durable pair): fd725e35599f   556 unchanged: 39fc707703ae
go build OK · product Go files: 0 · no test written or run (existing test source read)
```

Corpus **849**; registry 848 besides / 900 files.

**Resume token:** `HEKI-CP-2026-11-22-48`

---

## `HEKI-CP-2026-11-22-49` — the amendments contradicted the only executable surface

**Stage:** foundational · **TASK-03:** `v2 final audit` · **implementation: NOT authorized**

### One cause, eleven contradictions

The decision table is the contract's only executable part — "applied in order, first match wins". Every amendment since row 554 imposed preconditions on rows 5/6/8 from a different document. A scanner built from the table alone reproduces all of them.

### The sharpest

Content with marker `none` matches **no row** — and VT-18 asserts it is Healthy. That is precisely the probe-12 mixture being legitimised. Fail-closed resolution: Pending, `unattributed`, and duplication writes its own marker.

Also: a 0-byte log is invalid under V2-09 and empty under V2-10, and that file is the *normal* clear outcome. Row 8 can return Healthy on a `Recovered` marker V2-86 calls non-authoritative. Row 6 calls unattributable content `retracted-residue` — the provenance V2-16 denies. Deletion is exempt under V2-53 and gated under V2-75. V2-19 says no state is stored while the marker stores one. V2-13 overstates the store: a first marker write has no previous document, so the outcome is absent.

### Consequence

v2 is not closable as 554 + 555 + 556 + 557. The next artefact is a **v3 consolidation** with the conditions inline in the table, VT-18 replaced, and V2-91…V2-103 adopted or rejected one by one. CAS retained unchanged as an unimplemented design requirement.

### Verification

```
rows 554-558 all OK   558 (final audit): 166dac474051   557 unchanged: fd725e35599f
go build OK · product Go files: 0 · no fixture (document reads only)
```

Corpus **850**; registry 849 besides / 901 files.

**Resume token:** `HEKI-CP-2026-11-22-49`

---

## `HEKI-CP-2026-11-22-50` — v3: one ordered procedure, conditions inline

**Stage:** foundational · **TASK-03:** `v3 consolidation drafted` · **implementation: NOT authorized**

### What v3 is

Self-contained. 61 clauses, 20 decision rows, 8 lifecycle tables, 42 tests. It does not require rows 547–558 to be read to be executed; those stay byte-identical as evidence.

### The keystone

`gate_enforced` is a **declared writer property consulted by rows 18–19**. While the gate does not exist, content in generation 2+ is `Pending / promotion_requires_gate` — visible and unresolved — instead of Healthy on an assumption. Implement the gate, declare the property, and rows 16–17 become reachable with no other change.

That is the honest form of "the gate is specified, not implemented": the contract pays for the missing gate in its state table rather than in a footnote.

### Self-audit found two holes before registering

- **Usable content** had no definition, so a valid log beside an invalid companion satisfied no row. Now defined in V3-09 and used by rows 12–19.
- A void directory holding one lone-newline file would have been `Healthy empty`, swallowing the malformation. Row 7 is now "no usable content exists, and something present is invalid" — Damaged, with VT3-42.

### Other resolutions

Zero-byte NDJSON is vacuously valid, so a completed clear has a defined reading. `none`/absent marker with content is Pending `unattributed`; duplication writes its own marker so a conforming copy never reaches that row. `Recovered` reaches row 5 and stops. `retracted-residue` is vocabulary only — lineage cannot be established, so it is never emitted. Marker resolution is separated from content settlement. Deletion is exempt by statement and is not a resolution. First marker write is distinguished from replacement: a death yields `Absent` with no pre-image.

### Mechanical check

Every state-producing test names its row, and an assertion requires all 20 rows to be covered — so a rule change cannot leave a test preserving an obsolete expectation.

### Verification

```
rows 554/557/558/559 all OK   559 (v3): 475ca21dfe06
20 decision rows · 20 matrix rows · 42 unique VT3 ids
go build OK · product Go files: 0 · no fixture (document work only)
```

Corpus **851**; registry 850 besides / 902 files.

**Resume token:** `HEKI-CP-2026-11-22-50`

---

## `HEKI-CP-2026-11-22-51` — availability vs permission, and a derived gate flag

**Stage:** foundational · **TASK-03:** `v3.1 consolidation` · **implementation: NOT authorized**

### The fail-closed rule survives, on better grounds

Contamination and legitimate lag are byte-identical: settled clear then append, versus interrupted clear then append into a log still holding G. Same structure, size, ordering, content; the marker has no history. So availability itself is unestablished — the rule is not about write permission, and presenting such a conversation Healthy would show retracted content as confirmed.

### Two authorities

Availability is decided by the table; mutation permission by the gate. A state never grants permission: Healthy does not mean writable, Pending does not mean blocked.

### `gate_enforced` is derived

Build-bound self-test, default false, false on any doubt, never a marker field, config flag, setting or asserted boolean, never inferred from the absence of contamination. And enforcement does not certify the past: shipping the gate will not retroactively clear existing G≥2 conversations — OQ-7, with OQ-6 for recording clear settlement.

### Self-audit caught a self-contradiction

Rows 16–17 were unreachable while V3-39 called an unreachable row a defect. Now stated: an input that cannot exist is a defect; a capability not yet built is a state. Their tests are declared inactive with capability `gate_enforced`, enumerable and never silently skipped.

### Verification

```
rows 554/558/559/560 all OK   560 (v3.1): b7421bd051b3
75 clauses · 20 decision rows · matrix covers 1-20 · 47 tests
stale reason name: 0 · go build OK · product Go files: 0 · no fixture
```

Corpus **852**; registry 851 besides / 903 files.

**Resume token:** `HEKI-CP-2026-11-22-51`

---

## `HEKI-CP-2026-11-22-52` — baseline accepted; settlement evidence needs no new crash domain

**Stage:** foundational · **TASK-03:** `v3.1 baseline · OQ-6 investigated` · **implementation: NOT authorized**

### Accepted

V3.1 (row 560) is the contract baseline. OQ-6 and OQ-7 are open by decision, not defects. Acceptance is of the contract, not of any mechanism: the gate, CAS, `gate_enforced`, settlement and enforcement records are all specified and none implemented.

Two-record rule stated formally: settlement cannot certify enforcement, enforcement cannot certify that earlier generations were settled, neither is derivable from the other, neither from the absence of contamination.

### OQ-6's real question

Not "did the clear finish?" — V3-31 can check that while it is true. It is "how does a later reader know it finished *before* the next content was written?" The second half is a property of the writer, so **settlement evidence alone can never establish purity**.

### The design answer

Extend the existing marker with `settled_generation` rather than add a second document. One document means the same rename replaces pair and settlement — atomicity, CAS, pre-images and crash tests inherited, no new crash domain. The field matters only if the next generation's CAS compares `(state, generation, settled_generation)`; stored but uncompared it is a comment.

**Staleness fails closed:** a `Recovered` marker carries an older settlement, so it under-certifies. No recovery direction over-certifies — that is what makes the single-document option safe.

**No backfill.** Scanner, migration and repair never write it, so historical generation ≥2 stays Pending even after the field ships.

### Recommendation

Do not pursue OQ-6 before OQ-7. A settlement record without an enforced boundary proves what a writer asserted while appearing to prove more — worse than Pending.

### Verification

```
rows 560/561/562 all OK   561 (baseline): c078b00c91a3   562 (OQ-6): 7ba815b52716
21 S-clauses · 6 S-tests · go build OK · product Go files: 0 · no fixture
```

Corpus **854**; registry 853 besides / 905 files.

**Resume token:** `HEKI-CP-2026-11-22-52`

---

## `HEKI-CP-2026-11-22-53` — atomic is not truthful

**Stage:** foundational · **TASK-03:** `OQ-7 enforcement investigated` · **implementation: NOT authorized**

### The chain

A record is whole (rename exclusivity, present) → truthful (writer held exclusivity, absent) → complete (no bypass, unproven). Break a link and the settlement field is well-formed evidence of nothing.

### Four answers

**Authority:** the premise is false today — there is no single boundary to certify. Of three candidate proofs only an unforgeable capability type checked by the compiler is authority; a static build assertion is the second net; a startup self-test is regression evidence only, because its coverage is exactly the paths someone remembered.

**Exclusivity:** the primitive is adequate — OS locks, cross-process, fail-fast, and released on death, so no stale-lock problem. The composition is absent: no history write takes the lock, and the marker CAS does not exist. The lock is advisory, which is why the capability type is what makes the property provable.

**Lifecycle:** cannot be withdrawn while held; can be bypassed entirely while running. Unavailability is the dangerous case, and the product currently turns it into **silent non-persistence** — four call sites discard the error.

**Crash:** all seven death points are fail-closed, none reads as Healthy-but-wrong. Point 3 — both representations void and valid, settlement unrecorded — is why OQ-6's S-06 is superseded: **the precondition for settlement is the hold, not the observation.**

### Verdict

Conditional yes, under five conditions. Without authority the field records what an unverified writer asserted; without surfaced refusals the failure is silent, which is worse than Pending.

`settled_at` demoted to audit metadata no row may read.

### Verification

```
rows 560-563 all OK   563 (OQ-7): 371a27023eb3
25 O7-clauses · 8 O7-tests · go build OK · product Go files: 0 · no fixture
```

Corpus **855**; registry 854 besides / 906 files.

**Resume token:** `HEKI-CP-2026-11-22-53`

---

## `HEKI-CP-2026-11-22-54` — containment is an API shape, not a token

**Stage:** foundational · **TASK-03:** `OQ-7 clarification 1` · **implementation: NOT authorized**

### The correction

"Unforgeable capability" overstated Go. A token can be stored, copied, retained or reused after release, and nothing type-level proves the caller currently holds a lock — no goroutine-local state exists to check, and compiler/vet/race say nothing about ownership across a call. Split into four decidable properties — non-constructible, non-retainable, non-reusable, containment — a token satisfies one and struggles with three; an API shape satisfies all four.

### The shape

Raw writer unexported; one `boundary.Mutate(chatID, command)` acquires, validates, refuses when unresolved, mutates, records the CAS, releases. No capability value exists to outlive the call. Closed typed command set, no function-valued fields: a callback runs caller code inside the hold where it can block, re-enter, panic or retain what it was given; a caller-supplied expected pair moves the read-then-write race up one layer.

### Partial success is never success

Five content/marker orderings tabulated, each mapping to the fail-closed row the caller must see.

### The build assertion's limit

It must be negative — writing sites ⊆ boundary package — not a checklist of known writers. Completeness depends on conversation paths being centralised; where they are not, it is a heuristic and must be called one. **Unestablished today.**

### Verification

```
rows 560/562/563/564 all OK   564: 0328e436a57c
19 O7 clauses · discard sites in product: 4 (confirmed, release-blocking)
go build OK · product Go files: 0 · no fixture
```

Corpus **856**; registry 855 besides / 907 files.

**Resume token:** `HEKI-CP-2026-11-22-54`

---

## `HEKI-CP-2026-11-22-55` — confinement is false in three ways, so the assertion is a heuristic

**Stage:** foundational · **TASK-03:** `OQ-7a path centralisation` · **implementation: NOT authorized**

### The plain answer

History writes cannot currently be shown to pass through one boundary:

1. **a different package writes `history.json`** — `Service.SaveChatData` (`service.go:536-542`), non-atomic, from the in-memory inference context, not from either on-disk representation. Not another route to the same writer; a different writer with a different content source.
2. **no central path constructor** — `"chats"` literal repeated at nine sites in persistence.
3. **the identity basis differs** — persistence addresses by ID, `chat/service` by *display name*, under a capital-C root; on a case-insensitive volume both reach one directory under two identities, beside the known row 550 collision.

So the negative assertion would catch a new literal write and would miss a computed name — and it could not claim completeness. Stated as a heuristic, not proof.

### Correction recorded

A failure response is not a rollback. Refusal before any mutation leaves no consistency debt; content mutated with the marker unrecorded leaves an unresolved transition on disk that reconciliation must see. No compensating rollback — the resolution path is idempotent retry under the gate.

### What structural form needs

One directory constructor with basenames as constants; a `ConversationPath` type constructible only by it, so a runtime-built string cannot enter the API; and removal or routing of the second history writer. Then two checkable facts become proof within scope — and the non-history files still stay outside the claim.

### Verification

```
rows 560/563/564/565 all OK   565: 4d92ebc0b3e4
19 O7 clauses · 6 tests · go build OK · product Go files: 0 · no fixture
```

Corpus **857**; registry 856 besides / 908 files.

**Resume token:** `HEKI-CP-2026-11-22-55`

---

## `HEKI-CP-2026-11-22-56` — two trees, twelve files, one user-owned

**Stage:** foundational · **TASK-03:** `OQ-7d ownership inventory` · **implementation: NOT authorized**

### Three of my own claims corrected

1. **Withdrawn:** the case-collision claim from row 565. `chats` and `Chats` cannot collapse — persistence roots at `%APPDATA%\HOWLChat`, while `chat.Service` and the memory package root at the installation root. Two trees, not one directory under two identities. The ID-versus-name difference stands.
2. **Corrected:** `memories.json` is not one contested file. The two writers address different files in different trees with different schemas — the defect is two divergent memories stores, not a race.
3. **Corrected:** the second history writer is **dormant** — reachable only through the `NativeTestSaveChatData` seam. The finding stands; the severity I implied in row 565 does not.

### Ownership

Twelve non-history files, four writer packages. Identity (`meta`), projections (`world`/`scenario`/`characters`/`lorebooks`), **user-authored (memories — the only irreplaceable one)**, and application-owned settings and caches. **Seven of twelve are application-owned**, so adjacency establishes nothing.

Lifecycle: duplication copies seven files and drops five; wipe stages only the app-data tree; deletion removes only the app-data directory. **Every installation-tree file survives both.**

### Outside scope, deliberately

A shadow conversation tree at `ProjectRoot()/Chats/<id|name>` that V3.1 never scans. A `history.json` there would be reported `Missing` while present elsewhere. Recorded as OQ-7d-i, not absorbed.

### Verification

```
rows 560/565/566 all OK   566: 0cbcaff8d411
19 O7 clauses · 16 inventory rows · go build OK · product Go files: 0
```

Corpus **858**; registry 857 besides / 909 files.

**Resume token:** `HEKI-CP-2026-11-22-56`

---

## `HEKI-CP-2026-11-22-57` — a lock cannot live inside what it protects

**Stage:** foundational · **TASK-03:** `OQ-7b lock scope` · **implementation: NOT authorized**

### What the current lock covers

`LockPath()` is the document path plus `.lock`, held only during that document's write. So the marker would be serialised and **nothing else** — log, snapshot, companions and the active-chat pointer are unguarded, and `SaveChat` couples the pointer write to the transition on the hot path.

### Four scopes

Per-document: insufficient by construction. **Per-conversation: the minimum sufficient scope** — all four artefacts in one directory, contention O(1). Per-root: correct but serialises every conversation, turning a 105 ms rewrite into a global stall. Global: unnecessary, no operation spans the two trees.

### The constraint that decides it

A lock file *inside* the conversation directory breaks three paths: creation (no directory yet), deletion (an open handle blocks `RemoveAll` on Windows), and wipe (an open handle blocks a directory rename). So the lock belongs in a sibling namespace outside `chats/` — and creation then has somewhere to lock before `MkdirAll`, which is what row 550's second-resolution collision needs.

### Deadlock

Wipe needs a root-scoped lock, so the active-chat pointer must leave the conversation transition and be acquired only after the conversation lock is released — one upward edge, held by wipe alone.

### Verification

```
rows 560/566/567 all OK   567: f46e81e7aac2
20 O7 clauses · 8 tests · go build OK · product Go files: 0 · no fixture
```

Corpus **859**; registry 858 besides / 910 files.

**Resume token:** `HEKI-CP-2026-11-22-57`

---

## `HEKI-CP-2026-11-22-58` — the hold can be non-reentrant, but three paths must be flattened

**Stage:** foundational · **TASK-03:** `OQ-7c re-entrancy` · **implementation: NOT authorized**

### Seam exclusion is a build fact — with a caveat

`//go:build native_test_seam` with an `_off.go` counterpart, so the seam is absent from a default build. But `Service.SaveChatData` is untagged: the seam is build-excluded, the writer is present and uncalled. A static assertion must still list it.

### Root derivation is not stable

`getAppDataDir` uses the user home directory; isolation switches roots entirely. Different profiles, redirection or isolation resolve different locks for what a user calls one conversation. OQ-7c-i.

### The call graph does nest

`performFinalSave` = `performSave` + `RewriteHistory`. `RenameChat` and `SaveChatAs` route through `LoadChat`/`NewChat`. `App.LoadChat` mutates the conversation being left before reading the next.

**New finding:** every clean shutdown passes through log ⊂ snapshot — `SaveChat` writes the snapshot, `RewriteHistory` then rewrites both. Transient normally; persistent after a crash in that window.

### Answer

One non-reentrant per-conversation hold suffices, with three decompositions plus one merged final-save. No persistence call emits an event inside itself, so the boundary stays callback-free without changing behaviour; the synchronous `devEventSink` and the asynchronous frontend dispatch are excluded by rule, and a non-reentrant lock turns re-entry into a bounded `Busy` refusal.

### Verification

```
rows 560/567/568 all OK   568: e2961377d068
13 O7 clauses · 8 tests · 4 native_test_seam files · go build OK · product Go files: 0
```

Corpus **860**; registry 859 besides / 911 files.

**Resume token:** `HEKI-CP-2026-11-22-58`

---

## `HEKI-CP-2026-11-22-59` — the authoritative root is the history root, and `ProjectRoot()` is not stable

**Stage:** foundational · **TASK-03:** `OQ-7c-i root authority` · **implementation: NOT authorized**

### A hazard nobody had recorded

`ProjectRoot()` searches **upward from the executable** for `wails.json`/`go.mod` within six levels, else falls back to **the current working directory**. So a shipped build can adopt an enclosing repository as its data root, and the same binary launched from two directories resolves two different installation roots — splitting every non-history file while history stays put.

### The answer

The authoritative conversation root is the **history root**, `%APPDATA%\HOWLChat`: identical across production instances of one user on one machine, independent of launch location. Different users, isolated runs and roaming-redirected profiles are deliberately or necessarily different.

**The roaming case matters:** the tree lives under `AppData\Roaming`, so two machines can hold the same logical conversation with no shared filesystem and therefore no shared lock.

### Two rules

No lock may be derived from `ProjectRoot()` — launch dependence would give one conversation two lock identities. And the enforcement claim is **machine-scoped per user**; "two instances" unqualified is false.

### Defect recorded separately

The final-save window now has its own record (row 570): `SaveChat` writes the snapshot, `RewriteHistory` rewrites both, so every clean shutdown passes through row 6 — designed order, not a race. Reconciliation keeps reporting it and must not pick a winner.

### Verification

```
rows 560/568/569/570 all OK   569: 65d97ac26a6d   570: 79fb9497796d
14 O7 clauses · 11 FD clauses · go build OK · product Go files: 0 · no fixture
```

Corpus **862**; registry 861 besides / 913 files.

**Resume token:** `HEKI-CP-2026-11-22-59`

---

## `HEKI-CP-2026-11-22-60` — the pointer is a selection hint, and the app has no idea which conversation it holds

**Stage:** foundational · **TASK-03:** `OQ-7c-ii pointer consistency` · **implementation: NOT authorized**

### The structural finding

`App` holds `chatMessages`, `world`, `scenario`, `characters`, `lorebooks`, `triggeredLore` — and **no chat ID**. So the process knows which conversation it is holding only as `pm.config.LastChatID`: an in-memory field mirrored to a file every other process on the profile shares. Process-local state carried across a process boundary.

Five write sites, four of them side effects. `SaveChat` repoints on every save, so one field means both *what the user selected* and *what this process last persisted*. `saveConfig` is plain `WriteFile` where every other document in the package gets temp-and-rename — the field that coordinates selection has the weakest durability in the store.

### Authority, precisely

`LastChatID` is a **mutation target selector** — ten consumers resolve a content write target from it, including the shutdown rewrite — and **not a representation**: no content derives from it, so its loss cannot lose content, its corruption cannot fabricate content. It is outside the six-state reachability in both directions.

One narrow edge: `GetActiveChatMeta` backfills an empty `meta.ID` from the pointer. And recovery **rewrites** it — an unreadable `meta.json` makes startup fall back to the newest conversation and repoint the file. An automatic pointer repair reconciliation must neither imitate nor read as evidence.

### Verdict

**Selection hint:** one explicit writer, atomic replace, no CAS, no lock, no generation, loss tolerated — the worst outcome is a re-selection. Conditional on two unauthorized dependencies: an `App`-held active ID, and `SaveChat` no longer writing the pointer.

**The opening question, answered:** removing the pointer from the conversation lock removes an ordering edge. It does not make two instances agree on the active conversation, and it does not make the pointer-selected shutdown overwrite safe. That stays a writer-authority problem.

### Verification

```
row 560 OK b7421bd051b3   row 570 OK 79fb9497796d   row 571 OK c014f73f829e
go build OK · product Go files: 0 · no fixture
```

Corpus **863**; registry 862 besides / 914 files.

**Resume token:** `HEKI-CP-2026-11-22-60`

---

## `HEKI-CP-2026-11-22-61` — nothing needs re-entrancy; the autosave proves the identity rule instead

**Stage:** foundational · **TASK-03:** `OQ-7c-iii re-entrancy` · **implementation: NOT authorized**

### Five situations, one question

Synchronous nesting, async overlap, selection change, native callbacks, cross-conversation. Only the first is nominally re-entrancy, and it collapses by construction: one operation is one critical section over an unexported writer. **Re-entrancy is a symptom of leaving multi-write operations separately lockable.**

### What the investigation found instead

Two **unjoined autosave goroutines**, neither awaited at shutdown — so a final save can race an autosave of the same conversation. And `performSave` reads the pointer **at execution time**: an autosave captured in X that runs after selection moved to Y loads Y and overwrites it with X's in-memory content. OQ-7c-ii's identity defect, observed in shipped code.

New contract requirement (`O7-157`): every mutation uses the identity **bound to the content** it mutates; a deferred operation that finds its conversation no longer bound **refuses** rather than retargeting. No lock prevents retargeting — the second writer acquires the wrong conversation's lock legitimately.

### Verdict

No product requirement needs re-entrancy. The contract prohibits it and specifies bounded refusal: `acquired` / `busy` / `refused-reentrant` / `closed`. `refused-reentrant` must be diagnosed **above** the OS lock, or an internal composition bug surfaces as "another instance is writing". A refused write is not written, and must never become a content verdict.

### Verification

```
row 560 OK b7421bd051b3   row 571 OK c014f73f829e   row 572 OK f1aeed9dcc5f
go build OK · product Go files: 0 · no fixture
```

Corpus **864**; registry 863 besides / 915 files.

**Resume token:** `HEKI-CP-2026-11-22-61`

---

## `HEKI-CP-2026-11-22-62` — the shadow tree is a namespace, not a store; and "wipe all data" keeps the memories

**Stage:** foundational · **TASK-03:** `OQ-7d-i shadow-tree scope` · **implementation: NOT authorized**

### The answer is "a mixture", and that is the finding

`StateRoot()/Chats` holds four populations with four reachability classes: **live** user-authored memories and auto-capture settings (App bindings read and write them), a settings file with a live reader and no writer, **test-seam-only** `history.json`/`chat.json`/`memories_snapshot.json`, and a cache with an unestablished caller.

No conversation is ever opened from it — no writer there produces a `meta.json` — so it is neither a legacy conversation store nor a migration target. It is active product data with its own lifecycle, wrapped around test artifacts.

### Two corrections to row 566

The base is `StateRoot()`, not `ProjectRoot()` (equal in production, different under isolation). And `chat.json` has **no reader**: `ExportContentItem` refuses `library.Chats` before `contentItemFile` is reached.

### Two live lifecycle findings

`WipeAllData` removes only `<RootDir>/chats`, so **user-authored memories survive a "wipe all data"**. And the wipe's `SetChatID("")` is a no-op — empty is ignored — so the next memory write recreates the last conversation's shadow directory. Deletion leaves the namespace untouched too. Orphans accumulate by deletion, by wipe, and by launching from another directory.

### Boundary: unchanged

Reconciliation still scans `<rootDir>/chats` only. Candidacy is decided by `meta.json` — which makes the exclusion structural rather than a path rule. Shadow directories and orphans get **no verdict at all**. And a `Healthy` verdict must state that memories were not assessed, or it invites exactly the inference `O7-78` forbade.

Also recorded: `paths.Roots()` is a target, not an authority — its app-data root is `%AppData%\Howl` while persistence writes `%AppData%\Roaming\HOWLChat`, and it labels the installation root non-writable while the memory package writes beneath it.

### Verification

```
row 560 OK b7421bd051b3   row 572 OK f1aeed9dcc5f   row 573 OK 20bdd6ec0b08
go build OK · product Go files: 0 · no fixture
```

Corpus **865**; registry 864 besides / 916 files.

**Resume token:** `HEKI-CP-2026-11-22-62`

---

## `HEKI-CP-2026-11-22-63` — there is no retry identity, and a torn append is worse than a duplicate

**Stage:** foundational · **TASK-03:** `OQ-7d-ii retry identity` · **implementation: NOT authorized**

### The store proves non-application and cannot prove application

No operation identity exists anywhere in the conversation path: `ChatMeta`, `ChatHistory` and `AppConfig` carry no operation field, `AppendHistoryEntry` takes no key, and `internal/backend/chat` holds no journal, receipt or counter. Message identity is positional at three sites, so it is neither unique across time nor stable under a rewrite — it cannot key an idempotency check in either direction.

`AppendHistoryEntry` applies unconditionally. No marker read, no generation compare, no content inspection.

### The torn append is the sharpest finding

The writer discards the written count and never fsyncs. The reader stops at the first unmarshalable line and returns the prefix. Because the torn log is non-empty, the `history.json` fallback never engages. So a retry is appended after the tear, joins the invalid line, and every later successful append becomes invisible too. Retrying an ambiguous append does not risk a duplicate — it can make the conversation unreadable past the tear.

`RewriteHistory` is the only retry-safe writer, and only by substitution: it replaces the state rather than recognising an attempt.

### The existing heuristic is not a mechanism

`recordUserTurn` suppresses an identical adjacent turn in memory. It is adjacent-only, exact-content, field-blind, and sound only under an unstated caller protocol. Evidence that the product wanted idempotency; unusable as precedent.

### Settlement is not a retry ledger

`(state, generation, settled_generation)` has no operation identity and cannot answer "did this operation apply". Extending it would give one document two authorities — exactly what the row 573 domain separation removed.

### Verdict, and the two honest options

Only refusal before mutation is decidable, because "nothing was written" is a negative no missing evidence can hide. So: restrict retry to that class, or require a durable operation identity bound at capture, recorded before content, compared in one critical section, terminally settled, never content-derived, per-conversation, and re-execution-refusing. Retry of an ambiguous append is prohibited.

Also recorded: `_ = a.AppendChatMessage(...)` already leaves a failed message in memory and out of the authoritative log, restored only by the final rewrite that row 570 places in its own window; a refusal must return `not_attempted` / `not_applied` / `unknown` and may never discard captured content.

### Verification

```
row 560 OK b7421bd051b3   row 572 OK f1aeed9dcc5f   row 573 OK 20bdd6ec0b08   row 574 OK c0033256f2bf
go build OK · product Go files: 0 · no fixture
```

Corpus **866**; registry 865 besides / 917 files.

**Resume token:** `HEKI-CP-2026-11-22-63`

---

## `HEKI-CP-2026-11-22-64` — the caller is told nothing, and a switch destroys the only copy

**Stage:** foundational · **TASK-03:** `OQ-7d-iii caller-visible outcomes` · **implementation: NOT authorized**

### Severity, stated as such

The discarded append error is release-blocking silence. Three situations are indistinguishable to the user: displayed and stored; displayed but not stored; displayed with the storage outcome unknown. A missing log line is recoverable; a user told nothing cannot know that anything needs recovery.

### A sixth position row 574 did not name

`appendChatMessage` appends to `a.chatMessages` (`app-prompt.go:294`), then returns `nil` when `GetLastChatID()` is empty (`:299-301`). The message is captured and displayed, nothing is written, no conversation is identified, and the caller is told success. A second, different `nil` is the blank-content no-op (`:282-284`). So "no error" is compatible with applied, no-op, unbound and torn.

### The result is seven fields, over three domains

`Outcome` (`applied` / `no_op` / `not_attempted` / `not_applied` / `unknown`), `ConversationID`, `Stored`, `Captured`, `Displayed`, `Retry`, `Reason`. Only `Stored = observed` may be presented as saved, and `applied` means observed write success - never durability, since the count is discarded and no sync is called. `unknown` is terminal: never rendered, logged or summarised as saved.

### Switch and shutdown are disposition points

On a chat switch, `LoadChat` attempts a save, prints a warning and continues (`app-chat.go:96-99`), then `updateAppStateFromChat` replaces the in-memory history (`:106`). The outgoing turn is therefore unrecoverable in-process. A later rewrite *may persist* captured work only before the switch; after it, the possibility is gone. That is the accurate claim - not that the final rewrite restores it. On shutdown there is one attempt (`app-lifecycle.go:121`), reported only to stdout, and success is not authoritative until both row 570 writes complete.

### No conforming caller exists

`_ = a.AppendChatMessage(...)` at `app-prompt.go:138`; `.catch(() => {})` at `chat-input.js:247` and `chat-actions.js:77,172` - the three user-turn sites; four `_ = a.performSave()` sites; console-only autosave failures. The rule must be a negative assertion over discard sites, as with write containment.

### Binding is not operation identity

Conversation binding is required for refusal-only retry; durable operation identity is required only for ambiguous retry. Separating them is what lets honest failure reporting ship without a journal. Recommended baseline: refusal-only retry plus typed results plus switch/shutdown disposition. Unadopted; the owner selects.

### Verification

```
row 560 OK b7421bd051b3   row 574 OK c0033256f2bf   row 575 OK 64629096dcbd
go build OK · product Go files: 0
```

Corpus **868**; registry 867 besides / 919 files.

**Resume token:** `HEKI-CP-2026-11-22-64`

---

## `HEKI-CP-2026-11-22-65` — the register is closed to new branches except six baseline decisions

**Stage:** foundational · **TASK-03:** `open-question register audit` · **implementation: NOT authorized**

### There was no register

Every open question lived in a trailing table in one of ten documents. The state of the investigation therefore depended on ten footnotes - the condition the pack's governing principle forbids.

### Thirty-two entries, five classes

Six baseline, ten prerequisite, six evidence, six outside, four closed.

### Six can still change the baseline

- **OQ-1** - a disconnected UNC root stats as absent, so row 1 says `Missing` where row 2 is owed.
- **OQ-2** - path-string identity splits or merges case variants and junction aliases; the row 519 physical-identity finding decides it.
- **OQ-3b-b** - a non-empty log without a final newline is row 7 `malformed`; recommended strict, because the missing terminator is byte-indistinguishable from a tear.
- **Row 550** - a legitimate creation collision produces rows 6/12; `RC5-7` calls it an anomaly V3.1 has no state for.
- **`V3-44` wording** - row 572 retargeting can replace a conversation's content with another's without a generation transition, so Healthy must not be read as asserting authorship or binding.
- **Memories** - Healthy is a verdict about history only; memories were not assessed.

### No prerequisite can change a row

OQ-6, V2-OQ-2, OQ-7/7a/7b, OQ-3b-c, OQ-3a-c, OQ-7d-iii, operation identity and the gate self-test change schema, locks or result types. Their absence is already fail-closed: rows 16-19 unreachable, `gate_enforced` false.

### The consolidation point

A bounded six-clause V3.2 amendment - two decisions, one term, three limitations - plus the register table in the baseline. Not an extension. OQ-7d closes as investigation; the retry policy remains the owner's choice.

### Verification

```
row 560 OK b7421bd051b3   row 575 OK 64629096dcbd   row 576 OK 9d9ad823906b
go build OK · product Go files: 0
```

Corpus **868**; registry 867 besides / 919 files.

**Resume token:** `HEKI-CP-2026-11-22-65`


---

## `HEKI-CP-2026-11-22-66` — an unreachable share answers "not found", so row 1 would claim deletion

**Stage:** foundational · **TASK-03:** `OQ-1 disconnected-root classification` · **implementation: NOT authorized**

### Observed, not argued

On `go1.27.3 windows/amd64`, `os.Stat("\\nonexistent-host\share")` returns errno **53** (`ERROR_BAD_NETPATH`) and both `os.IsNotExist(err)` and `errors.Is(err, fs.ErrNotExist)` are **true** — at the root and at `\\host\share\chats\<id>`. A reachable host with a missing share returns errno **67** and both are **false**. Two unreachable roots, opposite verdicts from the same predicate, neither derived from evidence.

Local absence (errno 2/3) is correctly `ErrNotExist`. That is why every local test passes while the network case is inverted.

### The table defect

`V3-39` row 1 is "conversation directory absent" with no root precondition. Row 2's "I/O error" overlaps it, and first-match resolves the overlap toward row 1. So a disconnect produces `Missing` — the one-step form of the transition `V3-47` already forbids. `V3-60` does not cover it: the child `stat` is a complete observation of an error, not an incomplete enumeration.

### The amendment (proposed)

- `root_state`: `observed` / `absent` / `unobservable`.
- `root_class`: `local-fixed` (`DRIVE_FIXED`, 3) / `non-local` (UNC, including `\\?\UNC\`, or any other drive type) / `unknown`.
- Row 1 only when root observed, class local-fixed, child errno 2/3.
- Non-local or unknown: `Missing` unreachable; absence-shaped child is `Unavailable / absence_unverified`.
- Transport and access codes are Unavailable on every class; absence never from a sentinel alone.
- Double observation rejected: narrows but cannot close the window.

### Reachability and boundary

Production only conditionally, through redirected `AppData\Roaming` (source-derived, not reproduced) or an operator isolation root. A root unreachable at startup fails `MkdirAll`, so the realistic case is a mid-session disconnect.

Correction carried: prerequisites cannot change a row's rule, but they do change capability availability and row reachability — rows 16–17 are unreachable, not redefined. Rows 575–576 otherwise unchanged; `applied` / `stored` / `durable` stay distinct.

### Verification

```
row 560 OK b7421bd051b3   row 575 OK 64629096dcbd   row 576 OK 9d9ad823906b   row 577 OK 831641e401b2
go build OK · product Go files: 0
```

Corpus **869**; registry 868 besides / 920 files.

**Resume token:** `HEKI-CP-2026-11-22-66`


---

## `HEKI-CP-2026-11-22-67` — an alias is not a conversation, and a path string is never an identity

**Stage:** foundational · **TASK-03:** `OQ-2 alias paths and conversation identity` · **implementation: NOT authorized**

### Observed

A junction, case variant, trailing-dot form, 8.3 name and `\\?\` path all reach one directory with the same `(volume, file index)`; every path-string comparison is false. A copy has the same volume and a different index. `EvalSymlinks` returns the junction path unchanged; `GetFinalPathNameByHandle` returns the target. Physical identity survives rename; the junction then dangles with `IsNotExist` because it stores a path.

### PathService has the primitive, not the service

`FileIdentity` is exported and the Windows resolver/identity functions are correct, but both are unexported, `identityOf` has no caller, and no exported function returns or compares an identity. `CheckContained(base, alias)` returns nil through the junction — physical containment works; physical identity is not exposed.

### The live conflation

`GetChatList` enumerates by directory name and returns `ChatSummary.ID = meta.ID`; `LoadChat` treats its argument as the directory name. When name and `meta.ID` disagree, the list names a conversation that cannot be loaded. `RenameChat` changes only `Name`; nothing enforces agreement.

### Four identities

Locator (path; never a key), physical (volume+index; dedup and lock key; changes on copy/restore), logical (`meta.ID`; ownership, verdict attachment, retry binding), operation (row 574). `SaveChatAs` regenerates the logical ID; a manual copy does not.

### Rules

One candidate per physical directory. Dedup only with a nonzero, collision-free identity. Verdict attaches to logical identity. Two physical directories with one logical ID are reported, never merged. A dangling alias is not a conversation. Lock key is physical identity, verified before and after acquisition, refusing when uncomputable.

### Verification

```
row 560 OK b7421bd051b3   row 577 OK 831641e401b2   row 578 OK 341052761db4
go build OK · product Go files: 0
```

Corpus **870**; registry 869 besides / 921 files.

**Resume token:** `HEKI-CP-2026-11-22-67`


---

## `HEKI-CP-2026-11-22-68` — the newline is a writer promise, and the next append is the real hazard

**Stage:** foundational · **TASK-03:** `OQ-3b-b missing final newline` · **implementation: NOT authorized**

### Observed with the real loader

Three terminated records and three with the last unterminated both load **3**. A partial third loads 2. Two objects concatenated on one line load **1**. A lone `"\n"` loads 0. An unterminated two-record log beside a three-record snapshot loads 2 — the log wins.

Every production writer emits `\n`; the loader is a `bufio.Scanner` that returns the final token at EOF. So the terminator is a writer invariant, not a format requirement.

### Why strict rejection is wrong twice

`V3-08` would call a file whose records all load `malformed`, and row 7 would say Damaged when no snapshot holds content. And it protects nothing: the scanner is read-only, no writer consults a verdict, so a Damaged classification cannot stop the append that actually causes the loss.

### The amendment

- `V3-08`: every non-blank line parses; the final record may be unterminated.
- `tail = terminated | unterminated`: a diagnostic, not a state.
- Row 7 `truncated` only for an unparseable final line.
- Rows 14–15 keep Healthy and carry reason `unterminated_tail`.
- `V3-09` retained: a lone `"\n"` stays malformed.
- Gate: never concatenate onto an unterminated tail; terminate it in the same critical section or refuse.
- Rewrites repair the tail. Row 6 and the six states unchanged.

The classification and the writer rule are one change; without the rule, acceptance becomes silent loss.

### Verification

```
row 560 OK b7421bd051b3   row 578 OK 341052761db4   row 579 OK 17737609e32f
go build OK · product Go files: 0
```

Corpus **871**; registry 870 besides / 922 files.

**Resume token:** `HEKI-CP-2026-11-22-68`


---

## `HEKI-CP-2026-11-22-69` — creation can adopt a conversation and erase its memories

**Stage:** foundational · **TASK-03:** `row 550 creation collision` · **implementation: NOT authorized**

### Observed

Two child processes forced into the same UTC second both targeted one ID; one failed on a sharing violation on the fixed `meta.json.tmp` path, the other succeeded. A same-second `SaveChatAs` returned source ID == copy ID, renamed the source and replaced its snapshot.

Decisively: a conversation seeded at the predicted second with `memories.json = [{"text":"precious"}]` was silently adopted by `NewChat("intruder")` — same ID, nil error, meta renamed, memories replaced by `[]`, history.json `[]`, old log kept and served under the new name.

### Paths

`NewChat` uses `MkdirAll` and checks nothing. `SaveChatAs` inherits it. `SaveChat` and `RewriteHistory` also `MkdirAll` arbitrary IDs, so creation is not confined to `NewChat`. `AppendHistoryEntry` does not create.

### Sufficiency

Non-clock ID + exclusive claim (`Mkdir`, not `MkdirAll`) + no implicit creation elsewhere + safe commit ordering (meta last, or temp-dir move) are sufficient for "never silently adopt, merge, overwrite or redirect". A full transaction adds crash atomicity, not collision safety. A creation lock is not the primary guard; `LastChatID` is not ownership.

### Reconciliation

**Clarify, not Damaged.** The residue can be a valid, readable conversation. Rows 10/12/14–15 decide by evidence. V3.2 gains a limitation: Healthy/Pending assert content availability, not exclusive creation, ID uniqueness, companion ownership or uncompromised identity. No state, reason or repair. `RC5-7`'s anomaly is not implementable; `RC5-8`–`RC5-10` stand.

### Verification

```
row 560 OK b7421bd051b3   row 579 OK 17737609e32f   row 580 OK ec3cb718cee3
go build OK · product Go files: 0
```

Corpus **872**; registry 871 besides / 923 files.

**Resume token:** `HEKI-CP-2026-11-22-69`


---

## `HEKI-CP-2026-11-22-70` — `Healthy` means availability, not integrity

**Stage:** foundational · **TASK-03:** `V3-44` generation-1 wording · **implementation: NOT authorized**

### Finding

`V3-44`'s "nothing can contaminate it" is true only about **generation mixture** (`V3-42`/`V3-43`). Rows 10/11/14–17 read the two history representations and the marker; they read no companions, no ID uniqueness, no writer or operation identity. Row 572 (retargeted save), row 580 (adopted creation, memories erased) and row 578 (one logical ID, two locations) each leave a directory those rows call `Healthy`.

### Four properties

| Property | `Healthy` establishes it |
|---|---|
| Content structurally available | **yes**, within the assessed scope |
| Belongs to the intended logical conversation | no |
| Creation claimed an unused identity | no |
| User-authored companions preserved | no |

### Restated verdict

Observed content present, valid and mutually consistent within the permitted prefix relation, attributable to the marker's generation, or legitimately empty/closed; no row indicates loss. Non-claims: authorship, exclusive creation, identity uniqueness, correct write targeting, companion integrity. Assessed scope: directory + both history representations + marker; `meta.json` for candidacy and logical-ID attachment only; companions excluded.

### Not a seventh state

The wrongness is undetectable from the assessed bytes (no per-record provenance, `V3-18`). `Damaged` is false (the adopted directory loads cleanly); `Pending` is false (no unresolved transition). `V3-50` already separates ownership from availability. Amendment is `V3-44` wording + `assessed scope` + one verdict-meaning clause; `V3-73`/`OQR-10` close. No row, reason or transition changes.

### Carried from row 580 review

The creation-side requirement is a **foundational product correctness requirement**. **Collision safety** and **crash-safe creation** stay separate; staging and atomic publication may serve the second and must not be smuggled into the first. A replacement persistence system still may not certify what it has not observed.

### Verification

```
row 560 OK b7421bd051b3   row 580 OK ec3cb718cee3   row 581 OK 97b3363ca1da
go build OK · product Go files: 0
```

Corpus **873**; registry 872 besides / 924 files.

**Resume token:** `HEKI-CP-2026-11-22-70`


---

## `HEKI-CP-2026-11-22-71` — conversation health never implies memory health

**Stage:** foundational · **TASK-03:** memories limitation · **implementation: NOT authorized**

### Two stores

| Store | Path | Owner | Lifetime |
|---|---|---|---|
| conversation companion | `<RootDir>/chats/<id>/memories.json` (`[]interface{}`) | conversation | deleted with it; dormant writer |
| live records | `StateRoot()/Chats/<id>/memories.json` (`[]MemoryRecord`) | user | **outlives the conversation** |

`memories_backup.json` is literal `{}`, rewritten every save, no reader. `memories_snapshot.json` is read at init, no writer found. Settings and cache are app-owned; generated lore is global content. There is no global memory store.

### Lifecycle

Delete removes the companion and leaves the shadow store; the service keeps the deleted ID, so the next write recreates it. Wipe removes conversations and leaves every memory; `SetChatID("")` is a no-op. Rename keeps the ID. Copy/move and launch-root changes break association without deleting. `SaveChatAs` copies only the dormant companion; a same-second collision shares one shadow store.

### Disclosure

`unassessed`, `not_found_at_scope`, `unreadable`, `observed`, `orphaned` — findings, not states. No row, no lifecycle, no repair. Every verdict says memories were not assessed; no aggregation into "all data healthy".

### Verification

```
row 560 OK b7421bd051b3   row 581 OK 97b3363ca1da   row 582 OK 157b92142251
go build OK · product Go files: 0
```

Corpus **874**; registry 873 besides / 925 files.

**Resume token:** `HEKI-CP-2026-11-22-71`


---

## `HEKI-CP-2026-11-22-72` — v3.2 exists, and the audit passes with limitations

**Stage:** foundational · **TASK-03:** complete contract v3.2 + audit · **implementation: NOT authorized**

### What v3.2 contains

99 clauses, 21 decision rows (1, 2, 2b, 3–20), 58 tests. Six baseline outcomes inline, not appended:

- **Root** — `Missing` only on a local-fixed root with the root observed and errno 2/3, or absent with the index naming the conversation; non-local absence is row 2b `absence_unverified`.
- **Identity** — locator, physical, logical and operation; one candidate per physical directory; dedup by physical identity; lock key physical.
- **Tail** — an unterminated complete record is valid; `unterminated_tail` on rows 14–17; gate terminates-in-section or refuses.
- **Creation** — no state; V3-72a states the product requirement, collision safety separate from crash-safe creation.
- **`Healthy`** — availability only within the assessed scope; five non-claims; corrected `V3-44`.
- **Memories** — outside the verdict; separate five-finding disclosure; current implementation explicitly not the future architecture.

### Audit

99 unique IDs, no undefined references, all row reasons in V3-58, every row covered by a naming test, six states, no fail-closed violation, no hidden dependency. The added memory-language search finds only exclusion, prohibition and explicit implementation evidence; no memory state, lifecycle, retrieval or recovery model is defined.

**Result: pass with stated limitations.** Textual and structural; no test executed.

### Verification

```
row 560 OK b7421bd051b3   row 582 OK 157b92142251   row 583 OK 1994959107de
go build OK · product Go files: 0
```

Corpus **875**; registry 874 besides / 926 files.

**Resume token:** `HEKI-CP-2026-11-22-72`


---

## `HEKI-CP-2026-11-22-73` — v3.2 revised: the index is defined, and the remaining questions are bucketed

**Stage:** foundational · **TASK-03:** contract v3.2 revision + bounded audit · **implementation: NOT authorized**

### What the owner's review found

One hole. "The index" was an undefined authority in row 1, row 20 and the candidacy rule, so a clean
reconciliation could be read as proof that no unindexed conversation data exists.

### What changed in v3.2 (row 584)

- **`V3-13d`** — an index entry names one conversation by **logical identity only**; membership is byte-equality, no case folding or normalization; `index_state` is `present` / `absent` / `unreadable`, and absent is never treated as empty.
- **Rows 1 and 20** require `index_state = present`. An index-less scan reports neither `Missing` nor `Orphaned` — fail-closed silences, not defaults.
- **`V3-13a`** — coverage limit stated: a clean sweep means *every candidate is available*, never *no other conversation data exists*.
- **Row 4** extended to `meta.json` present but **unparseable** (`identity_lost`). That input previously satisfied candidacy by presence and could reach rows 14–15 as Healthy.
- **`V3-88`–`V3-92`** audit the revision: no duplicate clause ID, no undefined reference, no state, no row order changed.

**Finding, stated in the contract rather than assumed away:** no membership record is written by any code path. `GetChatList` enumerates `chats` with `os.ReadDir` and skips directories whose `meta.json` fails to load; `AppConfig` holds only `LastChatID`. So `Missing` and `Orphaned` are unreportable in today's product — opened as **OQ-8**.

### The four buckets (row 585)

| Bucket | Count | Blocking |
|---|---|---|
| 1 — changes the table's inputs (OQ-6, OQ-7, OQ-8) | 3 | none: all fail-closed already |
| 2 — gate implementability | 8 | 3 (OQ-7a, self-test, V2-OQ-2 for the marker half); OQ-5 conditional |
| 3 — caller behaviour | 4 | none feeds a row |
| 4 — separate defects and work | 8 | not absorbable |
| Residual — measurements | 2 | decide nothing |
| Closed | 10 | — |

**Result: no remaining question changes the rule of a reachable row.** Adoption conditions all hold; adoption remains the owner's decision and authorizes nothing.

### Verification

```
row 560 OK b7421bd051b3   row 584 OK 302519a3300b   row 585 OK 2746afc4ee4a
go build OK · product Go files: 0
```

Corpus **876**; registry 875 besides / 927 files.

**Resume token:** `HEKI-CP-2026-11-22-73`

---

## `HEKI-CP-2026-11-22-74` — v3.2 dispositioned: retained as reviewed draft, not adopted

**Stage:** foundational · **TASK-03:** disposition of v3.2 · **implementation: NOT authorized**

### The decision

The owner stopped expanding the contract and dispositioned it. Recorded as row 586, which edits nothing:

- **v3.2 is retained as a reviewed draft** — the row 584 text, 104 clauses, 21 rows, 62 tests.
- **Row 585 is its audit record** — not superseded, not folded into the contract.
- **Not adopted, not rejected.** Each outcome still needs an explicit owner act. Silence, the disposition,
  an implemented gate or a passing test is not adoption.
- **`V3-73a` still reads "subject to owner adoption"**, because adoption was not granted.

### The five conclusions kept, with their labels

| Conclusion | Standing |
|---|---|
| Reachable-state rules have no open-question dependency | register claim, not execution |
| Coverage incomplete without an index producer → `Missing`/`Orphaned` unavailable | observed in `GetChatList` (`os.ReadDir`, skips unloadable `meta.json`) |
| Gate not demonstrated implementable, three blockers | evidence claim, not proof of impossibility |
| Writer defects separately owned | unchanged ownership |
| Textual adoption conditions are not runtime evidence | `V3-88`–`V3-92` and `RQ-23` audit text only |

**On verification:** `go build ./...` exit 0 with 0 tracked or staged Go changes. That is compile evidence, not
behaviour. No test of the contract, of any proposed enforcement, or of any row's behaviour was executed.

### The next decision is two decisions, and neither is taken

| | Decision | What it settles |
|---|---|---|
| **D-2** | Resolve the gate's three blockers (OQ-7a, enforcement self-test, V2-OQ-2 marker half) | implementability |
| **D-3** | Formally accept the current coverage boundary | the operating envelope |

They are separable and unordered: D-2 without D-3 resolves implementability while leaving the boundary
unaccepted; D-3 without D-2 accepts a boundary whose enforcement is undemonstrated. Neither requires reopening
rows 577–582. D-3 is a coverage decision, not an adoption decision, and the fail-closed behaviour is not traded
away by either.

### Verification

```
row 584 OK 302519a3300b (unedited)   row 585 OK 2746afc4ee4a (unedited)
row 586 OK 59e0ca78edd4   go build exit 0
tracked/staged Go changes: 0 (12 untracked Go files from TASK-01/02 sit as they were)
```

Corpus **877**; registry 876 besides / 928 files.

**Resume token:** `HEKI-CP-2026-11-22-74`

---

## `HEKI-CP-2026-11-22-75` — the programmer takes D-2/D-3: two blockers designed, one excluded

**Stage:** foundational · **TASK-03:** D-2 gate blockers + D-3 technical input · **implementation: NOT authorized**

### The framing correction

D-2 was never a question analysis could close. Two of its three blockers are questions about what the code
should *become*, so the decision available is **"authorize a bounded slice, or leave D-2 open with its scope
precise"** — not "resolve or accept" (row 587, `PA-04`).

### The three blockers, from code

| Blocker | Evidence | Resolution |
|---|---|---|
| **B-1** path centralisation | the same `Join(pm.rootDir, "chats", …)` literal at **9** sites; **no** constructor anywhere in `internal/`; three root derivations | one `conversationDir(id)` containing against the history root, reusing the already-proven `CheckContained` |
| **B-2** enforcement self-test | **14** write calls (11 `persistence.go`, 3 `service.go`), **8** exported mutators, **13** App call sites, **7** probe programs | route writes behind one entry; replace the test with a **build-failing analysis** over the package (`V3-36`'s build-bound result) |
| **B-3** `V2-OQ-2` marker half | no marker, generation, settlement or retraction record exists in `internal/` at all | **excluded, not resolved** — outside the refusal-only slice |

**The hard part of B-2 is the exception list:** `service.go:531,541,605` write the shadow tree in the same
package. Naming them with reasons is acceptable for the slice; moving that writer out is row 573's ownership
change and is not part of it.

### Slice S-1, and what it costs

Constructor + containment (9 sites) · one entry with typed refusal · per-conversation lock in a sibling
namespace · tail rule at the append · the analysis test · App-layer refusal surfacing (13 sites). Explicitly
**not**: marker, CAS, settlement, index, repair, memories, `V3-72a`, the seven-field result, retry policy.

**One storage side-effect, named rather than denied:** transient lock files in a new sibling file class.
No document-format change; "no storage change" would be false.

### Recommendation

Accept **D-3** first — free, reversible, a product statement, and not a precondition of D-2 in either direction —
then authorize **S-1** as one bounded unit. If S-1 is not authorized, the honest status is "D-2 open, scope
precise", never "resolved". A trustworthy index is cheaper after S-1 than before it (row 573's two-authorities
precedent).

### Verification

```
row 584 OK 302519a3300b   row 585 OK 2746afc4ee4a   row 586 OK 59e0ca78edd4
row 587 OK 1a6e868fad74   go build exit 0 · tracked/staged Go changes: 0
```

No test was executed for the assessment; `PA-23`'s criteria are criteria, not results.
Corpus **878**; registry 877 besides / 929 files.

**Resume token:** `HEKI-CP-2026-11-22-75`

---

## `HEKI-CP-2026-11-22-76` — S-1 is implemented: the boundary exists, deletion is best-effort

**Stage:** foundational · **TASK-03:** S-1 code + audit · **authorization:** owner granted for S-1 · **uncommitted, undeployed**

### What now exists in the code

| | |
|---|---|
| **One path constructor** | `conversationDir(id)` / `chatsRoot()` — nine inline joins replaced; containment on physical paths, so a junction inside `chats/` is refused |
| **One mutation entry** | `withChatLock` per conversation, lock files in `<root>/.chat-locks`, **refusal not waiting**; same-process double entry is `refused_reentrant`, never reported as contention |
| **One write surface** | four helpers in `writers.go`; a package analysis test fails the build if a fifth write appears, and it proves it can fail |
| **Tail rule** | `appendLine` terminates an unterminated tail, never concatenates onto it |
| **Deletion** | intent appended to `deleted.ndjson` **before** removal; leftover directory and late autosave both refuse; the confirmation names the folder and links to it |

### Audit findings (the part that mattered)

The lock created an outcome the product never had: an append that can be **refused**. `app-prompt.go` was discarding that
result, so an unwritten assistant reply would have looked saved — now logged and surfaced as `chat:unsaved`, with four
`_ = a.performSave()` sites treated the same way. The first deletion-intent draft was a rewritten JSON document that could
lose a concurrent record; it is now an append-only log.

**Coverage statement kept honest:** seven of the eight mutators take the lock; `SetLastChatID` does **not** (root config
pointer, row 571 owns that race). Creation collision, the shadow-tree writer, wipe, repair, the final-save window and any
marker/generation work are explicitly outside.

### Verified

```
go build ./...                                            exit 0
go test . ./internal/backend/chat/ ./internal/durable/     ok
 go vet / gofmt                                            clean
node --check chat.js / delete-confirm.js                   OK
8 new tests, incl. an analysis self-test that rejects a synthetic bypass
```

### Not verified

No Wails build or app launch, so the dialog copy and the directory link are syntax-checked only. No real second process,
no crash between intent and removal, no performance measurement on save or append.

Corpus **879**; registry 878 besides / 930 files.

**Resume token:** `HEKI-CP-2026-11-22-76`

---

## `HEKI-CP-2026-11-22-77` — deletion recycles, creation claims, the final-save window closes

**Stage:** foundational · **TASK-03:** deletion architecture + S-1 corrections · **uncommitted, UI not launched**

### The three product changes

| Change | Before | After |
|---|---|---|
| Delete to Recycle Bin | `os.RemoveAll` | shell `SHFileOperationW` + `FOF_ALLOWUNDO`; **no permanent-delete path in this program** |
| Creation | `MkdirAll` adopted an existing same-second directory | exclusive `Mkdir` claim, suffix retry; row 580's data loss is unreachable |
| Shutdown | save snapshot → reload → rewrite both | one load, one `CommitChat`: companions → log → snapshot |

**No API guarantees recycling.** `FOF_ALLOWUNDO` is "if possible"; a volume without a bin still destroys and reports success. That is the platform deciding, not a fallback this code performs — and the dialog says where the files were.

**Tombstone kept, single-owner:** `deleted.ndjson`, append-only intent written before removal, owned by `PersistenceManager`, never proof of destruction, consumed (not duplicated) by the future index.

### Honest gaps

`SaveChat` / `RewriteHistory` / `AppendHistoryEntry` still `MkdirAll` their directory, so a write can still materialise a conversation — the creation side of `V3-72a`, deferred because its direct-save callers must be migrated. Membership still has **no authority**: it is derived from directory contents, so `Missing`/`Orphaned` remain unreportable. That is the next slice.

### Verified

```
go build ./... · go test . ./internal/backend/chat/ ./internal/durable/ · go vet · gofmt   all clean
new: real shell recycle on a disposable fixture · missing-path idempotence · 3 same-second creations stay distinct
```

**Not verified:** app not launched (dialog copy and folder link syntax-checked only); recycle-versus-destroved not distinguished; no second process; no path-cost measurement.

Corpus **880**; registry 879 besides / 931 files.

**Resume token:** `HEKI-CP-2026-11-22-77`

---

## `HEKI-CP-2026-11-22-78` — deleting from the front end now works, and the clicks prove it

**Stage:** foundational · **TASK-03:** deletion, user-visible path · **uncommitted, verified by clicking the packaged app**

### The defect the clicks found

Delete was unreachable. The action guarded on `window.currentChatId`, which nothing ever assigned; the load path stores the identifier on `window.dataSync.currentChatId`. The menu item highlighted, no dialog appeared, no intent was written, the folder stayed — and every backend test still passed, because the backend was never called. One owner now feeds the guard, `GetChatDirectory`, `OpenChatDirectory` and `DeleteChat`; the dead global is removed.

| Before | After |
|---|---|
| dead `window.currentChatId` | `activeChatId()` → `dataSync.currentChatId`, the only reader |
| folder link failed silently (console) | `Could not open the folder: <error>` in the window |

### What was clicked, in `build/bin/howl-chat-s1.exe` (`-tags "desktop,production"`), isolated store

| Step | Result |
|---|---|
| File ▸ Delete | dialog with copy, directory link, warning |
| dialog ▸ Delete | window cleared; `deleted.ndjson` written with the exact directory; folder gone from `chats/` |
| Recycle Bin | `$RA2V2NR` named `2026-11-22T00-00-00Z` present — recycled, not destroyed |
| folder restored by hand, app restarted, File ▸ Load | the deleted chat does **not** reappear |
| dialog ▸ directory link | `[OpenChatDirectory] opened …\chats6-11-22T19-19-37Z`, one Explorer window |
| dialog ▸ Cancel | dialog closes, chat and folder untouched |

Membership is decided by intent, not by filesystem success: a directory that survives a blocked removal is still not listed.

### Honest gaps

The folder link is only as large as its glyphs and its second line is the reliable hit area; UI Automation cannot see into WebView2, so all driving is coordinate-based with screenshot confirmation; only the packaged host was exercised, so browser-versus-executable agreement is **not** claimed; the dialog still says *Deletion is permanent and cannot be undone* while the code recycles — a copy decision left to the owner; the leftover case was simulated by restoring the directory because the shell refused nothing on this volume. Membership still has no index authority, so `Missing`/`Orphaned` remain unreportable.

### Verified

```
go build ./... · go test . ./internal/backend/chat/ ./internal/durable/ · go vet · gofmt · node --check   all clean
screenshots: build/verify/13, 14, 15, 20, 21
```

Corpus **889**; registry 888 besides / 932 files.

**Resume token:** `HEKI-CP-2026-11-22-78`

---

## `HEKI-CP-2026-11-22-79` — the gate is written down, and both hosts do the same thing

**Stage:** foundational · **TASK-03:** dialog truth, host parity, one integrity repair · **uncommitted, both hosts driven**

### Permanent gate

`AGENTS.md` now carries the rule and the two supported hosts: packaged WebView2 (`-tags "desktop,production"`) and the browser dev host (`-tags uiharness`, real `App` methods over `/__rpc`). Backend tests are supporting evidence. Any other runtime is not evidence.

### Corrections, verified in both hosts

| Correction | Result |
|---|---|
| Contradictory copy | dialog says the folder moves to the Windows Recycle Bin and can normally be restored; the irreversible sentence is now an overridable `warning` for the five dialogs that need it |
| Path hit area | `display:block`, measured 354x52 px; a click past every glyph opened the folder |
| Parity | identical text, identical link outcome, identical store effects and Recycle Bin entries in WebView2 and in the browser host, each on its own isolated root |

### One repair

`AppendHistoryEntry` and `RewriteHistory` refuse a missing conversation (`refused_unbound`) instead of creating it. `SaveChat` still creates; its caller migration is the next repair.

### D-3 status

Designed, not implemented, and no second index added. Listing is directories minus tombstones, so `Missing` requires an expectation record; the choice between a persisted manifest (a mirror that can disagree) and references only (a live reference pointing at an absent directory) is recorded as needing the owner's decision before code.

### Verified

```
go build ./... · go test ./internal/backend/chat/ (2 new tests) · go test . · go vet · gofmt · node --check · both host builds
```

**Not verified:** cancel in the browser host (the same dialog component, verified in the app); native pickers in the browser host; row 571 pointer race; autosave goroutine joining; D-3.

Corpus **890**; registry 889 besides / 933 files.

**Resume token:** `HEKI-CP-2026-11-22-79`

---

## `HEKI-CP-2026-11-22-80` — membership has an authority, and the app shows what it knows

**Stage:** foundational · **TASK-03:** D-3 · **uncommitted, demonstrated in both hosts**

### Decision

One append-only `lifecycle.ndjson` owns membership and deletion intent. Directories are payload. Deletion wins permanently. `deleted.ndjson` is absorbed once, never written again. Existing installs are bootstrapped once; a folder added later is never claimed.

### States, visible

| State | Meaning | Where it comes from |
|---|---|---|
| present | claimed and folder there | fold + stat |
| missing | claimed, folder gone | fold + stat |
| orphaned | folder nobody claims | directory scan |

The Load dialog shows all three. Only *present* is selectable.

### What clicking found

The browser host minted a conversation in a fresh root from its own cached session: the validator accepted any listed id, and a missing bridge skipped validation. Both are now fail-closed; the backend also refuses a payload whose id is not present, so the cache cannot write.

### Also closed

`SaveChat` no longer recreates a known conversation whose folder is gone; `CommitChat` requires presence; `config.lock` serialises `LastChatID`; startup background work is joined before the final save; Cancel is inert in the browser host (0 `DeleteChat` calls, hashes and tombstone unchanged).

### Verified

```
go build ./... · go test . ./internal/... · go vet · gofmt · node --check (4 scripts) · both host builds
screenshots: build/verify/25-membership-states.png ; harness: missing+orphaned cards, nothing written
```

**Not verified / open:** no recovery action for a missing conversation (product decision), frontend autosave timers not joined, row 571 beyond `LastChatID`, native surfaces differ between hosts.

Corpus **891**; registry 890 besides / 934 files.

**Resume token:** `HEKI-CP-2026-11-22-80`

---

## `HEKI-CP-2026-11-22-81` — Level 0 exists in code (HOC-H0)

**Stage:** foundational · **Direction:** HOC, accepted 2026-11-22 · **Work date:** 2026-11-22 18:17–18:25 CST (2026-11-23 00:17–00:25 UTC)
**Record:** `docs/architecture/HOC-H0-LEVEL0-2026-11-22.md`

### Delivered

`internal/host` — gate (admission + single-entry shutdown), background (bounded join),
streams (one generation token), processes (one child process, one `Wait`), stages (ordered
shutdown failures), and the `Authority` container. No domain value is stored in the package.

Removed from `App`: `shuttingDown`, `shutdownMu`, `shutdownCoord`, `backgroundWG`, `streamMu`,
`streamCancel`, `streamGen`, `serverMu`, `serverMonitor`, `serverProc`. `shutdownCoordinator`
deleted; its semantics live in `host.Gate.Once`. Call sites retargeted in `app-model.go`,
`app-prompt.go`, `app-image-model.go`, `app.go`, `app-lifecycle.go` and three test files.

### Verified — clicks, both hosts

Harness: startup, model-folder prompt dismissal, nav to World/Characters/Lorebooks/Image
Generator/Chat, typed send, File→Save, File→Exit (`[FinalSave] committed history with 1
entries`), restart restored the conversation and did not re-prompt.
Packaged exe: startup, nav World/Characters/Chat, typed send, window close (process exit, no
surviving child, no llama-server), relaunch restored both messages from disk.
`go build ./...`, `go vet`, `gofmt`, `go test ./...` (mic host probe skipped) all pass.

### Not done

No domain migrated; no competing implementation retired. `App` still holds conversation,
model, config, audio, image and content state — listed with current owner, intended
destination and status in the record's §3, all marked *retained*.

### Observed, not fixed

D-1 a failed send renders twice (two writers for one report); D-2 `never_prompt` is written but
the prompt path never reads it. Both pre-date HOC-H0 and belong to the domains that own them.

Corpus **892** (891 + this record's companion document).

**Resume token:** `HEKI-CP-2026-11-22-81`
**Next:** H1 — Media domain adopts Level 0 lifecycle (audio, image), verified by clicks in both hosts.

---

## `HEKI-CP-2026-11-22-82` — the Media domain owns audio and image (HOC-H1)

**Stage:** HOC — H1 Media domain, in working code
**Record:** `docs/architecture/HOC-H0-LEVEL0-2026-11-22.md` §8 (the H0 inventory extended in place)

### Done

`internal/media` now owns both media sessions. `media.Audio` holds the manager, desired settings,
last apply error and ASR activation pending, and drives `Apply` through three ports supplied by
`App` (Level 0 gate, microphone probe, error reporting); a refused capture returns
`media.CaptureUnavailableError`, which `App` maps to its existing UI error. `media.Image` holds
the engine, gallery, job manager, image settings store and the last model selection.
`App` holds `audio *media.Audio` and `images *media.Image` and no shadow copy. Eleven fields
were deleted from `App`; the four audio state helpers and the old `applyAudioSettings` body were
deleted, not wrapped, and every call site was retargeted to the owner. `TestA1SingleManagerOwnership`
now asserts App has no `audioManager` and that `media.Audio` has one manager and no pipeline.

### Verified — clicks, both hosts

Harness: startup with chat restored; Image Generator shows the gallery empty state and
`ScanImageModels FAILED` through the owner; audio panel reads the owner's desired settings
(volume 0.60); TTS enabled then **Save & Close** produced `SetSettings ok` and
`ApplyAudioSettings FAILED … voice model is not installed`, with the UI reporting "saved, but
activation failed"; File→Exit committed history and exited; restart restored both.
Packaged: startup, Image Generator and both generate gates ("Please enter a prompt." → "Please
select a model."), Audio Settings modal (blocked at hydration, D-4), Alt+F4 shutdown
(`wails.Run returned`), relaunch with restored history and a typed send.
`go build ./...`, `go vet`, `gofmt`, `go test ./...` (mic host probe skipped) all pass.

### Not done / blocked

Image generation has no SD checkpoint in the project, so it stops at model selection in both
hosts. Audio playback was deferred by instruction; settings, status and the apply refusal were
verified, playback/streaming/cancellation were not.

### Observed, not fixed

D-1 failed send still renders twice; D-2 `never_prompt` still unread; **D-4** packaged audio
settings stay at "Loading saved audio settings…" while backend calls work; **D-5** leaving Chat
can fail its save with `refused_reentrant`. All pre-date H1.

**Corpus:** 892 Markdown documents under `docs/` (unchanged; this record extends an existing
document).

**Resume token:** `HEKI-CP-2026-11-22-82`
**Next:** next domain adopts Level 0 against this seam; `modelRegistry` re-homes first.

---

## `HEKI-CP-2026-11-22-83` — packaged audio settings hydrate, and leaving Chat saves once

**Stage:** HOC — H1 follow-up: D-4 and D-5 closed
**Record:** `docs/architecture/HOC-H0-LEVEL0-2026-11-22.md` §9

### Done

D-4: `wails-bridge-iframe.js` deferred to a native `window.go` injected into the modal's
subframe, where `App.GetSettings()` can never settle (host IPC belongs to the top document).
A framed document now always uses the parent-host bridge; a standalone page may still prefer
native. D-5: the conversation boundary's `refused_reentrant` was two *concurrent* frontend
writers, not composition — `DataSync.cleanup()` on `beforeunload` and a raw `SaveChat()` in
`chat-state-bridge.js` on the same event. The duplicate raw save is deleted, `saveCurrentChatData`
is one serialized route with a coalesced trailing write, and `navigateTo` awaits it before
unloading so the write is not aborted mid-flight.

### Verified — clicks, both hosts

Packaged: Audio Settings hydrated with the owner's saved values (volume 0.40, speed 1.75, TTS
off, Push-to-Talk on, no loading label); toggled Push-to-Talk off and **Save & Close** →
`settings.json` `pushToTalk=false`; nav to Image Generator renders; File → Exit and Alt+F4 each
followed by relaunch restored chat history and the saved audio settings.
Harness: nav Chat → Image Generator logs `Chat data saved successfully via SaveChatData`, the
RPC log holds zero `refused` results and no save-failure console entries; audio volume changed
to 0.40 and saved; File → Exit committed history and exited.
`go build ./...`, `go vet`, `gofmt`, `go test ./...` (mic host probe skipped) pass; both
binaries rebuilt.

### Not done / limitations

Image generation remains blocked (no SD checkpoint). No audio playback (deferred). The browser
tooling cannot drive controls inside the audio-settings iframe, so harness save clicks used the
same page served top-level; packaged clicks covered the iframe path.

**Corpus:** 892 Markdown documents under `docs/` (unchanged; an existing record was extended).

**Resume token:** `HEKI-CP-2026-11-22-83`
**Next:** next domain adopts Level 0 against the Media seam; `modelRegistry` re-homes first.

---

## `HEKI-CP-2026-11-22-84` — D-4/D-5 verified by pointer clicks on the embedded panel (harness) and re-checked packaged

**Stage:** HOC — H1 follow-up, closing the harness verification gap
**Record:** `docs/architecture/HOC-H0-LEVEL0-2026-11-22.md` §10
**Work date:** 2026-11-22, 19:02–19:34 local (CST) — identifier date is local; the machine clock,
the harness logs and every `build/verify/h1-d5-*` artefact carry it. The same window is
2026-11-23 in UTC, which is why the H0 header states both.

### What was wrong with the previous proof

§9's harness evidence clicked `audio-settings.html` as a top-level page. The app embeds that
panel in an iframe, so the harness workflow itself had never been clicked.

### What was clicked (harness, real Chromium window + CDP `Input.dispatchMouseEvent`)

Settings menu (410,15) → Audio Settings item (432,197): the modal's iframe appeared at
`[262,0,900,600]`, and the panel inside it read **no** "Loading saved audio settings…", volume
**0.4**, speed **1.75**, Push-to-Talk off, Save & Close enabled. Scrolled the panel, clicked the
TTS Output Volume slider at (537,455) → **0.4 → 0.7**, clicked Save & Close (1044,565) → modal
closed with the toast "Audio settings saved, but not all supported audio settings are
effective." and `settings.json` `ttsVolume=0.7`. Then Image Generator nav (39.5,368) →
`image-generator.html` with console `Chat data saved successfully via SaveChatData`; Chat nav
back; File → Exit (220,15 → 255.6,372) → `[FinalSave] committed history with 4 entries`, process
gone. Harness stderr `grep -c refused` = **0**; only `ScanImageModels` FAILED (no checkpoints).

### Relaunch and packaged re-check

Harness relaunch on the same root: `#chat-messages` = 4 messages, and the embedded panel
hydrated volume **0.7**, speed 1.75, no loading label. Packaged: Settings → Audio Settings
hydrated volume 0.7 with no loading label, Push-to-Talk toggled on by click, Save & Close →
`settings.json` `pushToTalk=true`, File → Exit, relaunch → 4-message history and Push-to-Talk
**on** restored.

### Screenshots

Harness: `build/verify/h1-d5-harness-00-start.png` through `…-24-final-relaunch-audio.png`
(step numbers 00–24 as written, with 13 never produced), including `…-02-audio-open.png`
(hydrated embedded panel), `…-09-volume-changed.png` (slider at 0.7), `…-11-after-save.png`
(toast), `…-17-nav-away.png`, `…-19-back-to-chat.png`, `…-20/21` and `…-23/24` (relaunch).
Packaged: `build/verify/h1-d5-packaged-00-startup.png` … `…-06-relaunch-audio.png`.
Plans and console logs in `build/verify/plan-harness-*.json`, `build/verify/logs/`.

### Tooling note

The preview/browser tooling cannot address nodes inside frames and OS-level input does not reach
its offscreen panel, so the harness was driven in a real Edge window over CDP; the click is a
browser-level pointer event routed into the frame, not a DOM `click()`. No product code was
changed for this pass; `build/verify/cdp-drive.mjs` is the driver.

### Not done / remaining

Image generation still blocked (no SD checkpoint). No audio playback exercised. Non-chat pages ship
dead dropdowns (`.dropdown-item` with `onclick`: `image-generator.html` 19/5 with no File handlers,
`world.html`/`character.html`/`lorebooks.html` 14/0) — File → Exit is inert there; pre-existing,
recorded, not changed.

### Gates

`go build ./...` exit 0, `go vet ./...` exit 0, `go test ./... -skip TestProbeMicrophoneOnHost`
exit 0 with no `--- FAIL` lines, `gofmt -l` unchanged (27 pre-existing files, none in
`build/verify`). No Go file was touched in this pass.

Two environment flakes hit while running those gates, recorded so a later run is not read as a
regression: the Go 1.27.3 linker panicked in `cmd/link` on three attempts (`bad alignment value`,
`index out of range`) while linking different test binaries; `go clean -cache` and a re-run cleared
it. Separately `internal/audio`'s `TestConcurrentDifferentConfigurationsRemainIsolated` failed
twice on a loaded `go test ./...` run (`hits 2/5 misses 2, want 0/2`) and passed standalone and in
the final full run — a timing-sensitive single-flight assertion, not touched here.

**Resume token:** `HEKI-CP-2026-11-22-84`
**Next:** next domain adopts Level 0 against the Media seam; `modelRegistry` re-homes first.

---

## `HEKI-CP-2026-11-22-85` — the real image-model library is added through the UI and survives a restart

Work date: 2026-11-22 (local, `-06:00`). Continues `…-84`; no new document series.

### What was wrong

- `Image Generator` reported `no image model directories are available`: the only configured folder
  (`…\HOWL_Chat\models`) does not exist, and nothing had ever pointed the app at the real library
  (`H:\AI GGUF\Image models`, 83 entries).
- **Add Folder could not open a chooser.** `OpenImageModelsDirPicker` passed the stored
  `ModelsDir` as Wails' `DefaultDirectory`; the runtime refuses a default directory that does not
  exist, so the call failed with `directory dialog: default directory … does not exist` and no
  window appeared. (An earlier read of that status line as `unknown directory chooser` was a
  misreading of the small screenshot; the dialog was never reached.)
- **One malformed GGUF hung the scan.** `kitten-tts2-20261021-q4_0.gguf` makes `gguf-parser-go`
  panic (`makeslice: len out of range` in `ReadString`). The panic escaped `ProbeWeight`, so
  `AddModelDirectory` never returned: the page stayed on *Scanning …* / *Loading models…* with no
  CPU use, no sockets and no error (reproduced twice, five-plus minutes each).

### What changed

- `app-image-diagnostics.go`: new `dialogDefaultDir(preferred)` — the stored path only when it is
  an existing directory, otherwise the executable's directory.
- `internal/image/scanner/gguf_probe.go`: `inspectGGUF` recovers from the parser panic, resets the
  probe and logs `[Scanner] skipping malformed GGUF header <path>: <panic>`.
- Tests: `TestDialogDefaultDir{FallsBackWhenPreferredIsMissing,KeepsExistingPreference}`
  (`app-image-model-path_test.go`), `TestProbeWeightSurvivesMalformedGGUFHeader`
  (`internal/image/scanner/probe_test.go`, builds a `.gguf` with a `1<<62` string length).
- Harness tooling: `build/verify/pick-dir.ps1` (EnumWindows + UIA `FromHandle`),
  `build/verify/type.ps1` (SendKeys), `build/verify/plan-harness-image-models.json`.

### What was clicked

Packaged WebView2 (`howl-chat-s1.exe`, isolated root): Image Generator (619, 687) → **Add Folder**
(758, 645) → chooser opens at `…\build\bin` → Folder box + `H:\AI GGUF\Image models` + Enter →
breadcrumb `Crucial (H:) > AI GGUF > Image models` → **Select Folder** (1606, 983) →
`Added …/AI GGUF/Image models — 56 assets`, `2 folders`, `33 model(s) · 4 VAE · 8 LoRA discovered`
→ dropdown (818, 552) lists 31 models with family/format → `pony-v7-base (sdxl)` selected
(760, 815). Process stopped and relaunched: folder and `61 assets` restored, model field restored
to `pony-v7-base`, 1024×1024, same 33/4/8 status.

Harness (`howl-ui-harness.exe` on `127.0.0.1:8799`, Edge + CDP): Image Generator nav at 39.5, 368
→ `modelCount 32`, folders `…/AI GGUF/Image models  61 assets`, `2 folders`, console
`[IG] Scan: 33 models, 4 VAEs, 8 LoRAs, 3 ControlNets 1 T5, 2 CLIP-L, 0 CLIP-G`. The harness has
no native dialog host, so Add Folder cannot be clicked there — recorded, not worked around.

### Read back from disk (`%TEMP%\howl-hoc-h0`)

`model-scan-dirs.json` → `"H:\AI GGUF\Image models": "image"`; `model-registry.json` → 61
entries, all on `H:` (`diffusion 33, lora 8, embedding 5, llm 5, vae 4, controlnet 3, clip_l 2,
t5xxl 1`); `[Scanner] Scan complete: 33 models, 4 VAEs, 8 LoRAs, 3 ControlNets, 5 embeddings,
5 LLMs` about five seconds after Select Folder. No `image-settings.json` (nothing loaded or
generated); the restored model field is the renderer's `localStorage` `IG_STATE_KEY`, a second
authority beside the backend store — noted, unchanged.

### Screenshots and logs

`build/verify/h1-img-00-generator.png` (before), `…-02-dialog-screen.png` (chooser at the executable
directory), `…-04-window.png`, `…-09-window.png` (navigated to the library), `…-10-window.png`
(added, 56 assets), `…-11-list.png`, `…-13-zoom.png` (model list), `…-14-model.png` (pony selected),
`…-15-restart.png` (after relaunch), `h1-img-harness-01-generator.png` (harness);
`build/verify/logs/h1-img-harness2.log`.

### Not done / remaining

No generation or checkpoint load was exercised (VRAM and sd-server are a separate step). The
library is on a removable `H:` drive — when it is missing the row reads `missing` and the scan
skips it. `kitten-tts2-20261021-q4_0.gguf` is still listed as SD1.5 diffusion because its header
is unreadable and classification fell back to the filename. Folder-row `✕` and **Rescan** were not
clicked. The user's real profile was untouched, so the same clicks are needed there.

### Gates

`go build ./...` exit 0; `go vet ./...` exit 0; `go test ./... -skip TestProbeMicrophoneOnHost`
exit 0, no `--- FAIL`; `gofmt -l` unchanged (27 pre-existing files, none in `build/verify`).

**Resume token:** `HEKI-CP-2026-11-22-85`
**Next:** exercise a checkpoint load/generate against the added library, and decide whether
selection should move from `localStorage` to the image settings store it already has.
## HEKI-CP-2026-11-22-86 — Image Generator driven end to end on the real profile

**Finished.** The Image Generator workflow is verified by clicks on the owner's real profile in the
packaged WebView2 executable, and the same interaction surface is exercised in the browser harness.

**Clicked in the packaged app** (`build/bin/howl-chat-s1.exe`, 1400x900, `build/verify/clickwin.ps1`
window-relative): Image Generator nav → **Add Folder** → shell chooser → path → **OK** (row added,
`config/model-scan-dirs.json` updated) → folder-row **✕** → in-app **Forget** modal → **Forget**
(row and config entry gone, `3 folders`) → **Rescan** (`Image folders rescanned`, three
`[Scanner] Scan complete` passes, `66 model(s) · 8 VAE · 16 LoRA discovered`) → 768x768 → **Generate**
→ `Done!`. Then close → relaunch → dismiss "New model folder detected" with **Not now** → nav Image
Generator: `IMAGE FOLDERS 3 folders` (`…/AI GGUF/Image mod… 61 assets`, `…/InvokeAI/models 0 assets`,
`…/AI GGUF/Image mod… 61 assets`), `Illustrious-XL-v2.0` still selected, gallery intact → one further
generation completed (`Done!`).

Images produced on the real profile: `images/gallery/img_20261123_023950_85100.png`,
`img_20261123_025822_170700.png`, `img_20261123_030359_114300.png`, `img_20261123_032605_156600.png`
(768x768, 20 steps, euler_a, Illustrious-XL-v2.0, ~70-200 s each). Screenshots: `h1-real-52-done`,
`h1-real-54-gen2`, `h1-real-55-final` under `build/verify/`.

**Broke, and fixed:**
- Add Folder closed the app twice — Wails' owner-window dialog plus `go-webview2`'s
  `Chromium.Focus` `os.Exit(1)`. Replaced with `internal/host/folderdialog_windows.go`
  (`SHBrowseForFolderW`, no owner, foreground on init); `folderdialog_other.go` for other platforms.
- Folder-row **✕** did nothing in WebView2 (`window.confirm` is not shown there) — now the existing
  `DeleteConfirm` modal, included on `image-generator.html`, with an optional `confirmLabel`.
- Generation opened a Windows Terminal console (`sd-server.exe`) over the app — `proc.HideWindow` in
  `internal/image/engine/sdcpp/engine.go`.
- Harness showed "Lost contact with the image engine" beside a saved image (timer throttling in an
  occluded tab) — the job watchdog now probes `GetImageJobStatus` once before failing.

**Clicked in the harness** (`build/bin/howl-ui-harness.exe` on 127.0.0.1:8799, page
`image-generator.html`): real pointer clicks on the folder-row **✕** → in-app **Forget** → **Forget**
(row gone, backend back to three folders), **Rescan** (`Image folders rescanned`), **Generate**
(`Done!`, gallery 19 → 21 tiles, `img_20261123_031918_370200.png`, `img_20261123_032226_672900.png`).
Host difference: the native folder chooser is unsupported (`ChooseFolder` returns "not supported on
this platform"), so the temp folder was registered through the same bridge call the dialog feeds;
every other step was clicked.

**Owners.** Image engine state and wiring: `App` (`app-image*.go`) — destination Level 0 coordination
layer + image domain, migration **not started**. UI form state: renderer `localStorage`
(`howl-image-gen-state`); selected model name is written to both that and
`config/image-settings.json: last_model` — duplicate authority, destination backend store, migration
**not started** (they already disagree: backend `z-image-turbo-Q4_K_M` vs UI `Illustrious-XL-v2.0`).
Engine config/paths: backend `config/image-settings.json`, current owner, unchanged.

**Remaining.** Consolidate `IG_STATE_KEY` → backend store; fix the history column overlapping the
Generate button below ~900 px page width; SDXL GGUF still unloadable on this `sd-server`; duplicate
G:/H: model entries in the dropdown; `kitten-tts2…` misclassified as SD1.5.

**Stopped at:** all clicked verification complete for this slice; gates run
(`go build ./...`, `go vet ./...`, `go test ./... -skip TestProbeMicrophoneOnHost`, `gofmt -l`).

**Gates:** `go build ./...` exit 0; `go vet ./...` exit 0; `go test ./... -skip
TestProbeMicrophoneOnHost` exit 0, no `--- FAIL`; `gofmt -l` unchanged at 27 pre-existing files.

**Resume token:** `HEKI-CP-2026-11-22-86`
**Next:** migrate the renderer's image form state onto the backend settings record, then retire the
`IG_STATE_KEY` writer.
## HEKI-CP-2026-11-22-87 — Level 0 landed as a lifecycle authority; image engine moved to its domain

**Finished (code).** `internal/host` is the application lifecycle and coordination authority:
`Lifecycle` (lifecycle.go) runs an ordered `Stage` list through `Boot` (dependency order, stop at
first failure, per-stage reports, readiness only when all stages pass), `Shutdown` (gate closed once,
teardown stages through `Stages`, failures accumulated) and a reboot that reopens the gate
(`Gate.Reopen`), so restart runs the same composition. `Authority` exposes it beside the existing
gate/background/streams/processes primitives, all of which have real callers. `app-boot.go` holds the
seven boot stages (core services → audio owner → settings store → content libraries → model registry
→ image domain → inference capture) and `NewApp` is now a struct literal plus `Lifecycle.Boot`.
`app-lifecycle.go` supplies the shutdown stage list to `Lifecycle.Shutdown`.

**Image domain owns its engine.** `internal/media/image_boot.go` constructs the gallery, `sdcpp`
engine, job manager, job→gallery sink, last-model restore and version probe
(`Image.Boot`), and owns `Reconfigure`, `EngineConfig`, `SaveSettings`. Deleted from App:
`initImageEngine`, the job sink in `NewApp`, and `reinitEngine`'s construction body. sd-server version
probing moved to `sdcpp.BinaryVersion`/`IsKnownOldVersion`. Video, vision, memory, RAG, audio,
libraries and registry keep their existing owners.

**Clicked — packaged WebView2 (real profile).** Launch → dismiss "New model folder detected" (**Not
now**) → nav **Image Generator**: log `[lifecycle] ready: 7 stages in 54ms (boot 1)`,
`[lifecycle] webview attached: boot 1 ready true`, image domain `Initialized … commit: c6dfa9d`,
3 folders (`…/AI GGUF/Image mod… 61 assets`), 66 models, gallery intact. **Reload** → `Starting
sd-server … Illustrious-XL-v2.0`, badge `LOADED: ILLUSTRIOUS-XL-V2.0`. 768×768 **Generate** → VRAM
1656→7028 MiB, `[Image] Job … result saved to gallery`, screen `Done!`, image
`images/gallery/img_20261123_035218_966100.png`. **Unload** → log `UnloadImageModel called`,
sd-server 1→0, VRAM 1656→961 MiB. Window **close** → `[FinalSave] committed history with 15
entries`, `wails.Run returned`, clean exit. **Relaunch** → `[lifecycle] ready … boot 1`,
`Restored last model: "Illustrious-XL-v2.0"`, page shows the same 3 folders/61 assets/66 models with
Illustrious selected (screenshots `build/verify/h0-p1-start-view.png`, `h0-p1-loaded-view.png`,
`h0-p1-done-view.png`, `h0-p5-relaunch-view.png`).

**Clicked — browser harness.** Launch: `[lifecycle] ready: 7 stages in 23ms (boot 1)`, 3 folders,
Illustrious selected. **Reload** → `Starting sd-server`, status `Model ready`. **Unload** → log
`UnloadImageModel called`, sd-server 1→0, VRAM 1040→899 MiB. 768×768 **Generate** → `generate_image
768x768`, VRAM 6623 MiB, `[Image] Job … result saved to gallery`, page `Done!`, image
`img_20261123_040245_836600.png` (`build/verify/h0-h1-done.png`). Shutdown via the harness control
`GET /__shutdown` → `{"ok":true}`, `[FinalSave] committed history with 15 entries`, process exit.
Relaunch → `[lifecycle] ready … boot 1`, `Restored last model`, page shows 3 folders / 61 assets /
66 models. **Host differences:** no window or quit control in the harness, so shutdown is that HTTP
control (added in this slice) running the same sequence as Wails' `OnShutdown`; no native folder
chooser (`ChooseFolder` unsupported); no `GetGPUInfo` binding, so VRAM was read with `nvidia-smi`.
At 735 px page width a real Generate click hits the history column instead of the button (known
layout defect, untouched); generation was clicked at 1200×760 where `elementFromPoint` confirmed the
button is unobstructed.

**App state still in App** (owner · destination · status, full table in HOC §13.2): conversation
fields → `internal/backend/chat`, **not started**; inference/session/settings fields →
Runtime & Resources, **not started**; registries/content stores → Runtime & Resources/content,
**not started**; error reporter/logger/emitter → Level 0 diagnostics, **not started**;
`loaderDiagnosticsDir`/`discoveryDirOverride` → test seams, retained; Level 0 handle and `audio`/
`images` owner references are current owners with the state itself already in `internal/host` /
`internal/media`.

**Identifier date.** `2026-11-22` confirmed correct — `date "+%Y-%m-%d %H:%M:%S %z"` returned
`2026-11-22 21:43:31 -0600`; UTC-stamped gallery filenames (`img_20261123_040245…`) are 22:02 local.

**Remaining.** The App-state table above; `IG_STATE_KEY` vs backend settings duplicate authority
(owner/destination recorded in §12.5, migration not started); Generate-button overlap below ~900 px;
SDXL GGUF loading; duplicate G:/H: listings; `kitten-tts` classification.

**Gates.** `go build ./...` 0; `go vet ./...` 0; `go test ./... -skip TestProbeMicrophoneOnHost` 0
with no `--- FAIL`; `gofmt -l` unchanged at 27 pre-existing files.

**Stopped at:** Level 0 extraction complete in code and clicked through in both hosts; next slice
begins from the App-state table.

**Resume token:** `HEKI-CP-2026-11-22-87`
**Next:** move the conversation fields out of App into the chat domain as the first row of §13.2.
## HEKI-CP-2026-11-22-88 — the active conversation moved into the chat domain

**Finished (code).** `internal/backend/chat/session.go` adds `Session`, the one owner
of the active conversation: typed `SessionEntry` turns (role, content, sender,
character, characterId, hasImage, timestamp), the selected world/scenario/characters/
lorebooks, and the once-per-chat/scene/message trigger bookkeeping. It owns the
operations that were spread through App — `AppendUserTurn` (duplicate guard),
`AppendEntry`, `ReplaceEntries`, `ClearEntries`, `Reset`, `UpsertCharacter`,
`RemoveCharacter`, `SetCharacters`, `SetWorld`/`ClearWorld`,
`SetScenario`/`ClearScenario`, `SetLorebooks`, `ResolveLore` — plus the two
projections `History()` and `Project(chatData)`. Every method is nil-safe.

`App` lost `chatMessages`, `world`, `scenario`, `characters`, `lorebooks`,
`triggeredLore` and now holds `session *chat.Session`. `buildChatHistory`,
`updateChatDataFromAppState`, `updateAppStateFromChat`'s conversion and
`recordUserTurn`'s map handling are deleted; `performSave`/`performFinalSave` call
`session.Project`, `LoadChat` calls `session.Load`, New Chat / clear call
`ClearEntries`/`Reset`. `app-session.go` holds read-only projections for prompt
assembly and the `GetChatMessages` binding (`sessionEntriesToMaps`). No second
message slice or trigger map remains.

**Tests.** New `internal/backend/chat/session_test.go` (duplicate turns, Load
attribution + trigger reset, Project history preservation, membership recomputation).
Six app tests moved off `&App{chatMessages: …}` / `&App{lorebooks: …}` onto the
session; `newLorebookApp` now builds a session (a nil session had silently emptied
every lorebook assertion). Full suite: `go test ./... -skip TestProbeMicrophoneOnHost`
exit 0, no `--- FAIL`.

**Clicked — packaged WebView2.** Launch (built `-tags "desktop,production"`):
`[lifecycle] ready: 7 stages in 22ms`, restored conversation with 15 turns and status
`Chat session loaded: Chat 10/31/2026, 9:06:59 AM` (session Load on launch). **File →
New Chat** → toast `New chat created`, transcript and characters cleared
(`ClearEntries`). Typed via UIA and clicked **Send**: user turns persisted to
`%APPDATA%\HOWLChat\chats\2026-11-23T04-39-31Z\history.ndjson`; inference failed
because no chat model was loaded (`LLM-STREAM-ERR`, `dial tcp [::1]:8080`). Clicked
the message **Delete** twice: transcript dropped the turn and `history.ndjson` went
2 → 3 → 2 lines (ReplaceChatHistory → session.ReplaceEntries → RewriteHistory).
**Host limitation:** the packaged model selector did not open under synthetic clicks
(trigger at screen 1277,851 in a 1400×900 window; clicks landed on the composer), so
no model was loaded and no reply was generated in that host.

**Clicked — browser harness.** `[lifecycle] ready: 7 stages in 27ms`, page shows the
same collection. Selected `LFM2.5-8B-A1B-heretic.Q4_K_M`, clicked **Load**:
validation passed (4465 MB GPU estimate), llama-server health `{"status":"ok"}`, VRAM
874 → 5832 MiB. **File → New Chat**, typed `Reply with exactly: session ok`, pressed
**Send**: real assistant turn returned (`GetChatMessages` shows user + assistant;
`[inference] family="lfm2moe"`). (First choice `Gemma-4-E2B-it-Q8_0` + mmproj had
llama-server exit before healthy — recorded, not chased.) The packaged test chat was
then deleted through `DeleteChat`; `GetChatList` no longer lists it. Harness stopped
with `GET /__shutdown`, VRAM back to 714 MiB.

**Remaining for this row (14.2).** `chatPersistence` (collection/membership/append log)
still in App; send pipeline and prompt assembly still in App reading a projection;
`autoTemplate`/`currentModelProfile` untouched; `chat.Service`'s own `types.Context`
history is a second in-memory history nothing on the send path reads — retire, do not
sync.

**Gates.** build 0 · packaged 0 · harness 0 · vet 0 · test 0 (no FAIL) · `gofmt -l`
unchanged at 27.

**Stopped at:** conversation state has one domain owner and both hosts were clicked
through restore, new chat, send and message delete.

**Resume token:** `HEKI-CP-2026-11-22-88`
**Next:** bind the session to its conversation collection and move the send pipeline
behind it, so App keeps only the bindings.

## HEKI-CP-2026-11-22-89 — the collection and the send pipeline moved behind the conversation owner

**Owners now.** `chat.Session` owns the active conversation's turns (typed
`SessionEntry`, attribution, image flag), world, scenario, characters, lorebooks and
trigger bookkeeping. `chat.ActiveConversation` binds that session to one
`chat.PersistenceManager` and owns creation, loading, membership, the append log, the
snapshot commit, in-place history rewrites and deletion. `chat.Turn` owns the send
pipeline (stream and complete), and `Session.BuildMessages` owns prompt assembly,
reading the session directly. `App` holds one handle (`conversation`), the Wails
bindings, the sampler-option conversion, the display think splitter, the event sink and
the error log. `chat.Service` owns only world/scenario/character/lore reads and the
memory snapshot.

**What moved.** `chatPersistence` → `chat.ActiveConversation` (all routes, including
`SaveRecord`'s membership refusal, `FinalSave`, `ReplaceHistory`, `ClearHistory`).
`buildPromptMessages`, `truncateMessagesToFitContext`, `resolveLore*`,
`BuildCharacterLoreBlock`, `SendMessageStreamWithImage`'s prompt/stream/append body,
`GenerateCharacterResponse`'s stream/append body, `AppendChatMessage`'s log write and
`stripThinkBlocks` → `internal/backend/chat` (`prompt.go`, `turn.go`). `App` keeps
`buildInferenceOptions` (the settings authority's output), `newThinkSplitter`, events and
`logError`.

**What was retired.** `chat.Service`'s `types.Context`: the `context`, `inferenceOpts`
and `contextOpts` fields, the whole context-backed send/getter API, `buildPrompt`,
`retrieveRAG`, `convertToLorebookMessages`, `retryWithBackoff`, and the `history.json`
write in `SaveChatData` that serialised that context. `NewService(modelManager, client)`.
`currentSeed` now documents that the sampler settings are the seed authority. Exactly one
in-memory conversation history remains: the session's turns, which are what get sent,
persisted and rendered.

**Membership recomputation.** `Session.Project` now projects `Meta.WorldID`/`ScenarioID`
from the session when a world or scenario is present. Before this slice nothing ever
wrote those ids, so a saved conversation could hold content its membership record did not
name. An empty session still leaves the stored selection alone, so a conversation that
was never loaded cannot be emptied by a save.

**Tests.** `conversation_test.go` (append/log/restore, blank turn skipped, ReplaceHistory
rewrite, ClearHistory truncation, Save recomputes membership from the session, delete
active resets, foreign record refused, FinalSave commit, unbound no-collection);
`prompt_test.go` (session-read assembly, ordering, prefixes/override, lore, truncation,
StripThinkBlocks); `turn_test.go` (stream records user without a log line and assistant
with one, character attribution in log and snapshot, Complete writes no log line, stream
error leaves no assistant turn). App tests read the same route via `promptMessages(app)`
and `session().LoreBlockFor`.

**Clicked — browser harness (final tree).** Fresh isolated root;
`AddModelDirectory(.../lfm25-cpu/round2/models, "llm")` then reload; clicked the model
trigger, picked `qwen2.5-1.5b-instruct-q4_k_m`, clicked **Load** → llama-server health
`{"status":"ok"}`. **File → New Chat** (real clicks; new dir
`2026-11-23T05-11-35Z`), typed `Reply with exactly: conversation owner final ok`, clicked
**Send** → llama-server `Stream completed ... 28 bytes`; `history.ndjson` = user turn +
assistant `conversation owner final ok`. Harness stopped, relaunched, page reloaded:
`GetActiveChatMeta` → that chat and `GetChatMessages` → user + assistant. An earlier run
on the same tree produced `ok` for `Reply with exactly: conversation owner ok` and
restored it the same way.

**Clicked — packaged WebView2 (final tree).** Window moved fully on-screen (20,0,
1400×1020; its 580,270 placement had left the bottom controls past the screen edge).
Clicked the model trigger at window-relative 1266,819 → dropdown listed 8 chat models;
clicked `qwen2.5-1.5b-instruct-q4_k_m`; clicked **Load** (1248,857) → llama-server
`model loaded`, health ok. **File → New Chat** (dir `2026-11-23T05-08-24Z`), clicked the
composer, typed `Reply with exactly: packaged owner ok`, clicked **Send** → llama-server
`Stream completed ... 43 bytes`; `history.ndjson` = user turn + assistant
`package contains owner information correct`. **File → Exit** logged
`[FinalSave] committed history with 2 entries`; relaunch showed
`Chat session loaded: Chat 11/22/2026, 11:08:24 PM` with both bubbles. Selector fix:
`.model-select-trigger` gained `min-height: 34px`, flex centring and
`box-sizing: border-box` (its 16px text line plus a window whose content area clipped the
Load button at 862–892 is why earlier clicks missed); with the window fully visible a
real click selects and loads.

**Gates.** build 0 · packaged 0 · harness 0 · vet 0 · test 0 (no `--- FAIL`) ·
`gofmt -l` unchanged at 27 pre-existing files.

**Stopped at:** the conversation row of §13.2 is closed; both hosts were clicked through
model load, new chat, send, persisted history and restore.

**Resume token:** `HEKI-CP-2026-11-22-89`
**Next:** the Runtime & Resources row — `httpClient`, `modelManager`, `loadingProgress`,
`samplerSettings`, `modelSettings`, `appSettings`, `settingsStore`, `autoTemplate`,
`currentModelProfile`.

## HEKI-CP-2026-11-23-90 — the runtime session and the settings authority moved into internal/runtime

**What moved.** The Runtime & Resources row of HOC §13.2 is closed by a new domain package,
`internal/runtime`, with two owners:

- `runtime.Engine` owns the inference `Client` (interface moved out of `app.go`), the
  `*model.Manager` built on it, load progress and stage (atomics, progress×1000 as before),
  and the last load's `Profile` and injected `AutoTemplate` under one mutex. A nil engine is
  an absent runtime: reads return zero values, writes are dropped.
- `runtime.Settings` owns the unified settings document, the live sampler/model sections,
  the legacy `sampler-settings.json`/`model-settings.json` files and the `settings.Store`.
  `Load` restores the document and each section (unified section wins, legacy file when the
  section is absent) and returns the first failure without discarding the other section.
  `Replace` installs a whole document, replaces the sections it carries (removed keys stay
  removed; omitted sections are left alone) and persists it. Getters return snapshots;
  `Set*`/`Merge*` are the only writers.
- Deleted from `App`, not shadowed: `httpClient`, `samplerSettings`, `modelSettings`,
  `modelSettingsMu`, `appSettings`, `settingsStore`, `loadingProgress`, `loadingStage`,
  `autoTemplate`, `currentModelProfile`, the local `inferenceClient` interface and the four
  App-side settings file helpers. Boot constructs the owners; the capture wrapper wraps the
  engine's client; the discovery, reset and audio paths call the same bindings through the
  owner. Binding signatures and user-visible behaviour are unchanged.

**Tests.** New `internal/runtime/engine_test.go` (nil-safety, progress/stage, session
round-trips, concurrent reads/writes) and `internal/runtime/settings_test.go` (load and
legacy fallback, first-error continuation, replace-not-merge, snapshot isolation, persist
failure, legacy file round-trip/corruption, nil authority). Root tests were adapted to seed
the real owner rather than re-adding state fields.

**Gates.** build 0 · packaged 0 · harness 0 · vet 0 · test 0 (no `--- FAIL`) ·
`gofmt -l` unchanged at 27 pre-existing files.

**Clicked — browser harness (final tree).** Isolated root `%TEMP%\howl-13-2` reused; page
restored the previous chat. Viewport resized to 1500×1000 so the bottom bar was visible
(it sat below a 720px fold), clicked the model's **Load** → llama-server `model loaded`,
`server healthy`; page read progress 1 / stage `Ready` / loaded true through
`runtime.Engine`. **File → New Chat** (real clicks), typed
`Reply with exactly: runtime owner final ok`, clicked **Send** → assistant turn
`runtime owner final ok`; `history.ndjson` in `2026-11-23T05-27-07Z` holds both turns.
Harness killed and relaunched on the same root: `GetChatMessages` returned the same two
turns.

**Clicked — packaged WebView2 (final tree).** Window moved to (20,0, 1400×1020); clicked
**Load** (1276,858) → llama-server `model loaded` / `server healthy`. File menu: clicks at
window y=88 missed because the 32px menu row sits at y≈48–88; (230,68) opened it. Settings
in the menu opened (proving dropdown clicks land), and after closing it **New Chat**
(230,248) created `2026-11-23T05-33-07Z`. Clicked the composer, typed
`Reply with exactly: packaged runtime ok` as real keystrokes, clicked **Send** (1325,975) →
assistant turn `The package contains a functioning runtime environment.` in
`history.ndjson`. Killed and relaunched the binary: both bubbles are shown again
(`build/verify/13-2-relaunch.png`).

**Stopped at:** the Runtime & Resources row is closed; both hosts were clicked through model
load, new chat, send, append-log persistence and relaunch restore.

**Resume token:** `HEKI-CP-2026-11-23-90`
**Next:** the registries row of §13.2 — `modelRegistry`, `contentLibraries`,
`contentInventory`, `worldImport`, `scenarioImport`, `lorebookStore`, `modelIntel` (then
diagnostics).

## HEKI-CP-2026-11-23-91 — the content and model registries moved into internal/registries

**What moved.** The registries row of HOC §13.2 is closed by one composition owner,
`internal/registries`. It holds the existing authorities with no new state, index or cache:
the content `library.Registry`, the `library.InventoryStore`, the world and scenario
`contentimport.Service`, the canonical `lorebookstore.Store`, the `modelintel.Cache` and the
`image/registry.Registry`. `New(configDir)` builds them in the old boot's order and returns
every failure separately, so one corrupt store leaves only that part absent.
`Parts`/`Compose` wire pre-built parts for tests; every accessor (`Models`, `Libraries`,
`Inventory`, `Worlds`, `Scenarios`, `Lorebooks`, `Intel`) is nil-safe.

**What was retired.** The seven `App` fields `modelRegistry`, `contentLibraries`,
`contentInventory`, `worldImport`, `scenarioImport`, `lorebookStore`, `modelIntel`, plus
`bootContentLibraries` and `bootModelRegistry` (one `bootRegistries` stage now, boot stages
7 → 6). All 140 call sites read through `a.registries.X()`; no wrapper methods, no fallback
field. On-disk files and Wails binding signatures are unchanged, and the failure text for an
uninitialised store is the same.

**Tests.** New `internal/registries/registries_test.go` (empty root, nil set, Compose wiring,
corrupt libraries/model registry/intel, ignored corrupt inventory, duplicate ids, stale-plan
and cross-category import refusals, concurrent accessors) and `app-registries_test.go` (every
registry-backed binding refuses when the owner is absent; parts resolve independently).
`app-libraries`, `app-lorebooks-store`, `app-model-discovery_flow` and `app-wipe` tests build
the owner with `Compose` instead of assigning fields.

**Gates.** build 0 · packaged 0 · harness 0 · vet 0 · test 0 (no `--- FAIL`) ·
`gofmt -l` unchanged at 27 pre-existing files.

**Clicked — browser harness (final tree).** Isolated root `%TEMP%\howl-13-2`; log
`[lifecycle] ready: 6 stages`. Clicked the model selector: all 8 registered models listed;
picked `qwen2.5-1.5b-instruct-q4_k_m` and clicked **Load** → `GetLoadedModelName` reported it.
Fixture roots were registered through the real `AddContentLibrary` binding (the browser host
has no drivable native picker); **File → Settings → Files & Storage** showed Characters
`1 item discovered now`, and **Lorebooks → Import to Library** (its `confirm()` auto-accepted
for the click) wrote four `howl-lorebook-v1` files, after which `ListLorebooks` returned
count 4 / cutover true. `PreviewContentImport` for the scenario root planned
`destination_status: available, ready: 1`; `ExecuteContentImport` returned `completed`
(`imp_44224a287d956abe`) and wrote `fixtures/scenarios/scn-fixture-1/scenario.json`; replaying
the hash returned `plan_stale`; the journal listed the operation. **File → New Chat**, typed
`Reply with exactly: registry owner ok`, **Send** → assistant `Registry Owner Ok`;
`history.ndjson` in `2026-11-23T05-52-07Z` holds both turns. Relaunched: chat, four roots and
the four lorebooks restored.

**Clicked — packaged WebView2 (final tree).** Window moved to (20,0, 1400×1020). Model
selector (1280,819) listed the same 8 models; selected `qwen2.5-0.5b-instruct-q4_k_m` and
**Load** (1250,849) → llama-server `server healthy, model loaded`. **File → Settings →
Files & Storage** (210,66; 230,276; 196,128) showed all four fixture roots and Characters
`1 item discovered now`. Sidebar **Lorebooks** (60,348) read `Library files in 1 folder ·
4 entries`; **Scenarios** (60,295) and **Characters** (60,242) rendered. **File → New Chat**
(230,242), composer click (500,975), typed `Reply with exactly: packaged registries ok`,
**Send** (1335,975) → user + assistant turns in `2026-11-23T05-55-56Z/history.ndjson`.
Killed and relaunched: both bubbles shown again
(`build/verify/13-2-registries-packaged-relaunch.png`).

**Stopped at:** the registries row is closed in both hosts. The only remaining §13.2 row is
diagnostics (`errorLogger`, `errorReporter`, `errorEmitter`).

**Resume token:** `HEKI-CP-2026-11-23-91`
**Next:** the diagnostics row of §13.2 — `errorLogger`, `errorReporter`, `errorEmitter`, once
Level 0's reporting path owns escalation.
## HEKI-CP-2026-11-23-92 — the diagnostics row is blocked on Level 0 owning escalation

**Result: no code change.** The last §13.2 row is conditional, and the condition is not met.

**What the row says.** §13.2: `errorLogger`, `errorReporter`, `errorEmitter` → "Level 0
diagnostics, once its reporting path owns escalation". §3 keeps them in `App` with "retained —
H0 moved only a sink callback".

**Why it cannot be taken now.**
- HOC §2: "Escalation: the application supplies the sink (`a.logError`); the host supplies the
  ordering and the attribution." Level 0 "Must not own: any value a domain is responsible for."
- `internal/host/stages.go:11-15`: "The escalation sink is supplied by the caller, because
  error reporting is an application concern."
- `internal/host/lifecycle.go:29-34`: `ShutdownOptions` "carries the application's failure
  sink… reporting a failure is an application concern, so the sink is supplied rather than
  assumed."
- `app-lifecycle.go:33-45`: `App.shutdown` supplies `OnFailure` → `a.logError(…)`.
- `app-errors.go:17-49`: the three owners are built from the settings dir and the Wails
  context (`NewWailsEmitter(a.ctx)`), neither of which Level 0 holds.
- `internal/host` contains no logging/reporting/upload code, and no doc schedules it.

Moving the fields today would require a Level 0 contract/escalation change (out of scope), or
produce shadow copies in App (§4 forbids the duplicate authority). No shim or placeholder owner
was added.

**Untouched behaviour.** Same bindings (`GetErrorLog`, `SetErrorReportingEnabled`,
`IsErrorReportingEnabled`, `SendErrorReport`, `GetErrorReportPreview`), same log file
`<settings dir>/logs/howl-errors-YYYYMMDD.jsonl` with the same rotation, same
`error-reporting.json` toggle, same Wails UI emitter.

**Gates (tree unchanged).** build 0 · packaged 0 · harness 0 · vet 0 · test 0 (no `--- FAIL`) ·
`gofmt -l` unchanged at 27 pre-existing files. No host launched, nothing clicked: no behaviour
changed, and the §17 clicked evidence still applies to the identical tree.

**Stopped at:** the diagnostics row is open and blocked; it is the last row in §13.2, so there
is no other row to take before it.

**Resume token:** `HEKI-CP-2026-11-23-92`
**Next:** the diagnostics row, once Level 0's reporting path owns escalation (that contract
change comes before the move).

## HEKI-CP-2026-11-23-93 — diagnostics moved into internal/diagnostics; Level 0 unchanged

**Result: diagnostics has one owner; the diagnostics row of §3/§13.2 is closed (§19).**

**Contract correction.** The §18 block rested on the row naming Level 0 as the destination. The
destination was wrong: failure coordination (ordering, accumulation, invoking the sink) stays
with `host.Stages`/`host.Lifecycle`, while failure reporting (log + rotation, user-facing
reporting, event payloads) is application/domain infrastructure. Escalation policy and the sink
stay in `App`. Level 0 was not redesigned; §2 gained a "Diagnostics are not Level 0" bullet.

**What moved.** New `internal/diagnostics` (`logger.go`, `reporter.go`, `owner.go`) owns the
log (same `howl-errors-YYYYMMDD.jsonl`, keys `t,c,s,op,loc,msg,ctx`, 50 MB rotation, five-file
retention), reporting (same report shape/endpoint/`SaveReport`, mutex-guarded enabled flag) and
`app:error` payload construction. `internal/errors` keeps only the taxonomy and recovery
helpers (`sanitizeContext` → `errors.SanitizeContext`; `RecoverSilent` takes
`interface{ Log(*TrackedError) }`); `logger.go`, `reporter.go`, `wails.go` are deleted.

**Wiring.** `App` holds `diagnostics *diagnostics.Owner` instead of `errorLogger`,
`errorReporter`, `errorEmitter`. `NewApp` calls `diagnostics.Open(Config{LogDir,
AppVersion, UploadEnabled})` and attaches `appDiagnosticPort` → the guarded `emitEvent`
(Wails events or harness `devEventSink`), so no Wails type enters the package. This moved
initialisation out of Wails' `startup` into host-independent composition — which is why the
browser harness, which never ran `startup`, now has a working error system. `App.shutdown`
still supplies `OnFailure: a.logError(...)` and closes diagnostics last, so the
`SYS-SHUTDOWN` sink entry is written first. Bindings and on-disk formats are unchanged.

**Tests.** `internal/diagnostics/owner_test.go`: nil owner, unwritable/missing log dir, corrupt
line, duplicate reports, panicking port, write-after-close, 200 concurrent reports/events
(also `-race`). `app-diagnostics_test.go`: bindings' exact refusals without an owner;
owner-backed log/preview/toggle + one `app:error`; emitter binding with no `ctx` via
`devEventSink`; failing shutdown stage recorded at `app.shutdown`,
`app-lifecycle.go:shutdownStages` and reaching the application sink.

**Gates.** build 0 · packaged 0 · harness 0 · vet 0 · `go test ./... -skip
TestProbeMicrophoneOnHost` 0, no `--- FAIL` · `gofmt -l` = 26; `internal/errors/reporter.go`
(the pre-existing unformatted file this slice replaced) left the 27-file list.

**Clicked evidence.** Harness: speaker button with text on an isolated root → real
"audio manager has no active pipeline" → `TTS-SYNTH-FAIL` line in `howl-errors-20261123.jsonl`,
the same entry via `GetErrorLog`, and the visible `TTS-SYNTH-FAIL tts.synthesize failed` banner
from `app:error`; then model load → File→New Chat → send → assistant reply in
`chats/2026-11-23T06-38-31Z/history.ndjson` → `/__shutdown` → relaunch → restored. Packaged
WebView2 at (20,0) 1400×1020: same click produced a fresh log entry at `06:41:47Z` and the
banner; model load → new chat → send persisted at `chats/2026-11-23T06-42-52Z` → window closed
and relaunched → conversation restored. No host or `llama-server` left running; ports 8080 and
8799 closed.

**Stopped at:** every migration row in §13.2 is closed; §18's block is superseded by §19.
There is no next migration row in this HOC.

**Resume token:** `HEKI-CP-2026-11-23-93`
**Next:** none in this HOC. Candidate follow-ups (not started): move the two test seams
(`loaderDiagnosticsDir`, `discoveryDirOverride`) into the harness, and give
`frontend/dist/js/error-settings.js` a real Settings target (it still points at
`#settings-panel`, which the shipped page does not have).

## HEKI-CP-2026-11-23-94 — shutdown escalation proven live; error reporting reachable

**Result: both §19 loose ends closed with host evidence. Level 0 untouched. Nothing committed.**

**Live shutdown failure (harness).** A build-tagged seam can append one genuinely failing stage:
`app_shutdown_test_stage_harness.go` (`//go:build uiharness`) returns a `shutdown failure seam`
stage when `HOWL_TEST_FAIL_SHUTDOWN_STAGE=1`; `app_shutdown_test_stage_off.go` returns nothing,
so the packaged binary has no such switch. `shutdownStages` appends the seam after the real
stages; ordering, accumulation and the sink contract stay in `host.Stages`/`Lifecycle`
unchanged. With the switch on, `POST /__shutdown` answered
`{"ok":false,"error":"shutdown failure seam: injected shutdown failure (HOWL_TEST_FAIL_SHUTDOWN_STAGE=1)"}`;
the diagnostics log grew 8→9 with `SYS-SHUTDOWN` high `app.shutdown`
`app-lifecycle.go:shutdownStages` and that message; the live `/__events` stream delivered the
matching `app:error` payload. That is a runtime stage failure reaching the App sink, written by
the diagnostics owner and emitted on the UI path.

**Error settings were dead, now reachable.** `error-settings.js` was loaded but targeted
`#settings-panel`/`window.openSettings`, neither of which the shipped page has; the feature is
intended (N-05, bindings and `error-reporting.json` exist). Repaired: `#privacy-panel` gained an
Error Reporting card (toggle, status) and a Recent Errors card (Refresh, Preview Report, Send
Report, log list). The script now only drives that markup through the real bindings, with no
injection hook, no `window.openSettings` override and no localStorage mirror — the persisted
file is the single authority.

**Clicked evidence.** Packaged (20,0 1400×1020): Privacy & Data showed the card off with eight
log rows; toggling wrote `{"enabled": true}` and enabled Send; Preview rendered the report JSON;
after closing and reopening the app the toggle was still on, status "Error reporting enabled",
Send enabled, log list populated
(`build/verify/13-2-settings-privacy{,-on,-preview,-restarted}*.png`). Browser host: same panel
hydrated checked from the file with nine rows; toggling off wrote `{"enabled": false}` and
`IsErrorReportingEnabled` returned false. Send Report itself was deliberately not clicked (it
would upload); disabled-when-off was verified.

**Cleanup.** No compiled reference to `errorLogger`/`errorReporter`/`errorEmitter` remains (only
the `.freebuff/backup/app.go` tooling snapshot, outside the build); `internal/errors` has only
`types.go` and `recovery.go`; `grep` finds no `#settings-panel` in `frontend/dist` HTML/JS (only
legacy `.settings-panel` CSS selectors).

**Gates.** build 0 · packaged 0 · harness 0 · vet 0 · `go test ./... -skip
TestProbeMicrophoneOnHost` 0, no `--- FAIL` · `gofmt -l` = 26. No host or `llama-server` left
running; ports 8080/8799 closed.

**Stopped at:** nothing open in this HOC. The only unexercised surface is the actual report
upload, skipped by choice.

**Resume token:** `HEKI-CP-2026-11-23-94`
**Next:** none in this HOC.

## HEKI-CP-2026-11-23-95 — §13.2 diagnostics slice: implemented and tested, upload acceptance outstanding

**Result: status decision recorded; no code or product change.**

The diagnostics migration slice (§19, §20) is marked **implemented and tested** in both hosts:
one owner for logging/reporting/emission, a real shutdown-stage failure reaching the App sink at
runtime (log entry plus `app:error`), and a reachable Error Reporting panel whose toggle persists
across a packaged restart. It is **not** marked as acceptance of the diagnostics feature.

**Outstanding — upload acceptance.** Two things were never exercised or proven:
- **Transmission.** `Reporter.Upload` POSTs to `https://howl-error-reporter.vercel.app/api/report`
  asynchronously, ignores the response status and body, and drops network errors silently; Send
  Report was deliberately not clicked. No request leaving the machine, 2xx, non-2xx, timeout or
  offline behaviour has been observed.
- **Privacy boundary.** The payload after key-based `SanitizeContext` (exact keys `user`,
  `username`, `email`, `ip`, `host`, `password`, `token`, `key`; no value-pattern scan), the
  client IP/HTTP metadata the receiver sees, retention at the endpoint, and the UI's "error codes
  and line numbers" claim are all unverified. N-05 remains "Partially compliant".

HOC §20.5 records both; the HOC status field and the §13.2 diagnostics row point there, so §19/§20
cannot be read as acceptance of the whole feature.

**Gates.** Documentation only; no source changed, so the CP-94 run still describes the tree
(build 0 · packaged 0 · harness 0 · vet 0 · tests 0 with no `--- FAIL`). Nothing committed or pushed.

**Resume token:** `HEKI-CP-2026-11-23-95`
**Next:** if upload acceptance is wanted, verify transmission against a controlled endpoint
(request, 2xx, non-2xx, timeout, offline) and audit the built report/payload and privacy boundary;
until then the diagnostics feature is accepted only as "implemented and tested".

**Superseded in part:** `HEKI-CP-2026-11-23-96` executed the stub-endpoint transmission and
payload verification recorded here as outstanding; real-endpoint acceptance remains outstanding.

## HEKI-CP-2026-11-23-96 — report transmission and wire payload verified by executed tests

**Result: the two gaps CP-95 left open are closed at the wire boundary, with code fixes. The real
endpoint remains unobserved, so the diagnostics feature is still not fully accepted.**

**Defects the new tests exposed.**
- `Reporter.Upload` started a goroutine that discarded the response status, the reply body and any
  network/timeout error. A 503, an offline machine and a hung endpoint were indistinguishable from
  success, while the panel said "Report uploaded".
- The report carried the log's `context` map, redacted only when a key was exactly one of `user`,
  `username`, `email`, `ip`, `host`, `password`, `token`, `key`. A path, address or prompt under any
  other key would have been transmitted verbatim.

**Fixes (in `internal/diagnostics/reporter.go`).** `sendReport` now classifies the exchange: a
non-2xx reply is an error carrying the status and up to 512 bytes of reply text; the client is
injectable and bounded by its timeout; every failure is written to the diagnostics log as `NET-REQ`
/ `error_report.upload` at `internal/diagnostics/reporter.go:Upload`, so a lost report is never
silent, and an observer hook gives tests (and future callers) the outcome. `Upload` still returns
after the synchronous refusals only, so the fire-and-forget contract is unchanged. `ErrorPayload`
now carries sorted `context_keys` instead of `context` values, and `BuildReport` enforces the
documented 50-entry cap. `frontend/dist/js/error-settings.js` now says the report was queued, not
uploaded.

**Executed evidence** (`internal/diagnostics/reporter_test.go`, local `httptest` endpoints, no
network egress), all `--- PASS`:
`TestUploadPostsTheDocumentedRequestAndAccepts2xx` — actual method/path/headers/body captured, 2xx
→ nil; `TestUploadSurfacesANon2xxReply` — 503 with body → error + logged `NET-REQ`, later reports
still work; `TestUploadSurfacesANetworkErrorAndKeepsTheAppAlive` — refused connection → error +
logged, app alive; `TestUploadTimesOutOnAHangingEndpoint` — 120 ms client bound → `net.Error`
timeout, logged; `TestUploadDoesNotBlockWhileTheEndpointHangs` — returns in <250 ms while the
request is held open; `TestSerializedReportCarriesOnlyTheDocumentedFieldsAndNoContextValues` —
exact top-level (`report_id`,`generated`,`system`,`errors`,`error_count`), system and entry key
sets on the captured bytes, no `msg`, no `context` values, and no email/path/prompt/chat text even
under non-sensitive key names; `TestReportIsCappedAtFiftyEntries` — 60 logged → 50 sent.

**Gates.** build 0 · packaged 0 · harness 0 · vet 0 · full `go test ./...
-skip TestProbeMicrophoneOnHost` exit 0 with no `--- FAIL` and `ok howl-chat/internal/diagnostics`
· `gofmt -l` clean for both changed Go files · `node --check` on the panel script 0. A first full
run failed in `internal/audio` (untouched here); that package passed immediately in isolation and
in the complete rerun.

**Still outstanding.** No request has been sent to
`https://howl-error-reporter.vercel.app/api/report`: availability, TLS/DNS, reply, rate limits and
retention are unverified, and the receiver's view of client IP/HTTP metadata is unchanged. The
local log still keeps free-form context values for non-sensitive keys (never transmitted). The
"queued" copy was not clicked, because Send against a shipped binary would transmit.

**Resume token:** `HEKI-CP-2026-11-23-96`
**Next:** real-endpoint acceptance (controlled deploy or an agreed sandbox URL) and a decision on
redacting the local log's non-sensitive context values.

## HEKI-CP-2026-11-23-97 — endpoint deferred, code committed, local-state cutover owner decided

**Result: no behaviour change. Three owner decisions recorded.**

**Diagnostics endpoint deferred.** The owner will deploy and verify the Vercel endpoint after
the application is complete. Nothing sends automatically, so behaviour is unchanged; N-05 stays
partially compliant and the feature stays accepted only as wired and tested (HOC §20.5).

**Code committed.** `67c1a37` — "feat(hoc): establish Level 0 authority and domain state
ownership", 153 files, 13,200 insertions / 3,716 deletions. It contains the HOC implementation
(host, media, chat, runtime, registries, paths, durable, diagnostics, mic, harness, renderer
changes) and no documentation. Deliberately left uncommitted: `docs/` (300 paths), `eval/` (80),
`config/` runtime state (8), `build/` artifacts (4), `experiments/`, `AGENTS.md`,
`preview-frame.howl.json`. `go build ./...` is 0 on the committed tree. Nothing pushed.

**Local-state cutover owner: the domain stores** (HOC §21). The audit found lorebooks already
correct (`ImportLorebooks`, rollback, repository cache), worlds/scenarios owned by `worldstore`
but unreachable from the UI, characters split between an App-written directory and
`CharacterStorage` localStorage, settings written cache-first, and image/chat state duplicated.
Decision: each domain store owns an explicit, create-only, idempotent, journaled import and
rollback; `App` only composes the binding; the renderer gets one repository module per domain
(`lorebook-repository.js` is the template) and pages stop touching domain localStorage. Order:
settings inversion, character store, world/scenario bindings, image-state split, session pointer.

**Resume token:** `HEKI-CP-2026-11-23-97`
**Next:** the settings authority inversion (no migration), then the character store extraction.

## HEKI-CP-2026-11-23-98 — settings authority inverted, code committed

**Result: settings are backend-authoritative end to end; code committed as `5f32dd5`.**

`runtime.Settings` (GetSettings/SetSettings) is the sole authority.
`settings-state.js` is the only renderer seam: load asks the backend first and
treats an empty document as authoritative, the cache is a versioned rebuild
artifact read only when the backend is unreachable, and a mutation reaches
`SetSettings` before `current` or the cache changes. The save queue is
serialized, a failed write reverts and leaves the cache untouched, and
`flushSettings()` rejects so a "Saved" label cannot outrun disk.

Sampler, model, think-mode and image-path settings now go through that seam
instead of their own bindings. `GetSettings` returns the live sampler and model
sections with the document. UI-only state moved to `ui-prefs.js` (`howl-ui:`),
starting with the last settings panel; the legacy scale/animation keys are no
longer read.

**Verification.** `go build ./...` 0; `go vet ./...` 0; full
`go test ./... -skip TestProbeMicrophoneOnHost -count=1` 0;
`node --test tests/renderer/settings-authority.test.js` 7/7 pass. The new Go
tests prove `SetSettings` persists before returning and `GetSettings` includes
the live sections. The renderer tests execute the real module with a fake bridge
and prove conflicting/empty/legacy cache behaviour, backend-before-state order,
and failure revert without a cache write.

**Commit integrity.** The first form of the commit did not build from a clean
tree: `runtime/` in `.gitignore` matched `internal/runtime/`, so the package
`app-boot.go` imports was never committed (the earlier slice had the same gap).
The rule now matches only `/runtime/`, and `!internal/**/*.go` keeps the model
artifact globs (`*vae*`, `*encoder*`) from swallowing Go source. `internal/runtime`
and three model-named Go files are now tracked. A detached worktree at
`5f32dd5`, with no untracked files, passes `go build ./...`, `go vet ./...`,
the full Go suite and the 7 renderer tests.

**Commit contents (28 files, +1719/−246).** `app-chat.go`, `app-settings-authority_test.go`,
the settings JS modules and HTML entry points, `ui-prefs.js`,
`tests/renderer/settings-authority.test.js`, `.gitignore`, `internal/runtime`,
and the three model-named Go sources. No docs, eval, experiments, config or
build artifacts. Nothing pushed.

**Left open.** Settings backup/restore still copies localStorage wholesale and
does not restore settings (backend wins by design); `clearCache()` does not
target the settings cache key. Next HOC cutover: the character store
(§21.5 item 2).

**Resume token:** `HEKI-CP-2026-11-23-98`
**Next:** character store extraction, then world/scenario bindings.

## HEKI-CP-2026-11-23-99 — character store cutover landed

**Result: characters are backend-authoritative; code committed as `1d21b77`.**

`internal/characterstore` is the only writer of character documents: one folder
per identity, `howl-character-v1` envelope, revision-checked save, delete,
idempotent import with rollback. `CreateCharacter` now delegates to
`SaveCharacter`; `GetCharacters` lists the store and returns each revision;
delete/import/rollback bindings complete the route. `character-storage.js` is
deleted. `character-repository.js` is the only renderer seam: backend first, a
versioned derived cache written only after an acknowledged read/write, a
synchronous projection for display code, and a one-time import of the
pre-cutover raw localStorage array. All former direct readers and writers were
routed through it; no module writes `howl-chat-characters` outside the
repository.

**Verification (detached worktree at `1d21b77`, empty initial `git status`).**
`go build ./...` 0; `go vet ./...` 0; full
`go test ./... -skip TestProbeMicrophoneOnHost -count=1` 0;
`node --test tests/renderer/character-repository.test.js tests/renderer/settings-authority.test.js`
14/14 pass. New Go tests cover legacy listing, envelope updates preserving
unknown fields, stale revision refusal, delete isolation, import/rollback and
unsafe identities, plus the binding-level save/list/delete and import/rollback
paths. The renderer test covers cache-vs-backend, one-time raw import,
save/delete ordering, refusal revert, offline cache fallback, and a static check
that the cache key has no writer outside the repository.

**Commit contents (30 files, +1464/−424).** `app.go`, `app-characters.go`,
`app-characterstore_test.go`, `internal/characterstore`, the repository and
deleted storage module, the routed reader/writer modules, the six HTML entry
points, and `tests/renderer/character-repository.test.js`. No docs, eval,
experiments, config or build artifacts. Nothing pushed.

**Left open.** Legacy IndexedDB images are attached on read and persisted on the
next save rather than bulk-migrated. Next: world/scenario bindings (§21.5 item 3).

**Resume token:** `HEKI-CP-2026-11-23-99`
**Next:** world/scenario bindings over `internal/worldstore`.

## HEKI-CP-2026-11-23-100 — world/scenario cutover landed

**Result: worlds and scenarios are backend-authoritative; code committed as
`535439c`.**

`internal/worldstore` gained `List` and `Delete` and remains the only writer of
world/scenario documents. App now exposes `GetWorlds`, `SaveWorld`,
`DeleteWorld`, `ImportWorlds`, `GetScenarios`, `SaveScenario`, `DeleteScenario`
and `ImportScenarios`; saves cite a base revision, deletes surface store
refusals, and imports seed only absent identities. `world-repository.js` and
`scenario-repository.js` are the only renderer read/write seams: backend first,
versioned cache only after acknowledgement, one-time legacy-array import, and a
synchronous `all()` projection. Every former direct reader and writer was
routed, including `delete-service.js`; the grids await repository readiness.

**Verification (detached worktree at `535439c`, empty initial `git status`).**
`go build ./...` 0; `go vet ./...` 0; full
`go test ./... -skip TestProbeMicrophoneOnHost -count=1` 0;
`node --test tests/renderer/world-scenario-repository.test.js
tests/renderer/character-repository.test.js tests/renderer/settings-authority.test.js`
26/26 pass. The `internal/audio` concurrency test passed in the full suite.

**Browser dev-host evidence.** On `127.0.0.1:8799` with an isolated OS-temp root,
created, saved, edited, reloaded after cache deletion, and deleted one world and
one scenario. Canonical `world.json` and `scenario.json` were read back from
disk; edited titles survived cache removal. The pass found two defects that were
fixed and folded into the commit: `storeWorld` mutated the repository projection
before acknowledgement, and the grids rendered before repository readiness. The
packaged desktop host was not driven.

**Commit contents (31 files, +1521/−203).** `app-worlds.go`,
`app-worlds_test.go`, `internal/worldstore/list.go`, `list_test.go`, the two
repositories, the routed renderer modules, the six HTML entry points, and
`tests/renderer/world-scenario-repository.test.js`. No docs, eval, experiments,
config or build artifacts. Nothing pushed.

**Left open.** World/scenario image blobs remain in IndexedDB as references
inside canonical documents. Next: §21.5 item 4, image state split.

**Resume token:** `HEKI-CP-2026-11-23-100`
**Next:** image state split (§21.5 item 4).

## HEKI-CP-2026-11-23-101 — image state cutover landed

**Result: generation settings are backend-authoritative; code committed as the
amended §21.5 item 4 commit.**

`internal/image/settings` now owns a revisioned `generation` document in the
existing `image-settings.json`; `SaveEngineSettings` preserves it. App exposes
`GetImageGenerationSettings`, `SaveImageGenerationSettings` and
`ImportImageGenerationSettings`. `image-generation-repository.js` is the only
renderer seam: backend first, one-time legacy import, serialized saves, cache
only after acknowledgement. `IG_STATE_KEY` is version 3 UI-only
(`currentMode`, `seedLocked`, `aspectLocked`); no other renderer module touches
the generation keys through localStorage. `loadModels` restores the saved model
without replacing restored generation values with the family preset.

**Verification (clean checkout of the amended commit, empty initial `git status`).**
`go build ./...` 0; `go vet ./...` 0; full
`go test ./... -skip TestProbeMicrophoneOnHost -count=1` 0;
`node --test tests/renderer/*.test.js` 33/33 pass. New tests pin failed-save
revision retention and retry, save ordering, cleared-cache reload, one-time
legacy import, and the static no-direct-key rule. A mutation check that
reintroduces the pre-save backend reload makes the save-ordering test fail,
confirming the guard is effective rather than decorative.

**Browser dev-host evidence.** On `127.0.0.1:8799` with an isolated OS-temp root,
changed prompt and width, confirmed both plus a new `generation_revision` in
`data/config/image-settings.json`, confirmed `howl-image-gen-state` held only
UI fields, deleted `howl-image-generation-cache`, reloaded, and saw both values
render with repository source `backend`. The packaged desktop host was not
driven.

**Code contents (12 files, +937/−30), plus this checkpoint and the HOC §21.10
entry.** `app-image.go`,
`app-image-generation_test.go`, `internal/image/settings/settings.go` and its
test, `internal/media/image_boot.go`, the repository, the routed renderer
modules, `image-generator.html`, `app-wipe.js`, and
`tests/renderer/image-generation-repository.test.js`. No eval, experiments,
config or build artifacts. Nothing pushed.

**Left open.** Legacy import is VM-tested rather than browser-profile tested.
Next: §21.5 item 5, chat session pointer.

**Resume token:** `HEKI-CP-2026-11-23-101`
**Next:** chat session pointer cutover (§21.5 item 5).

## HEKI-CP-2026-11-23-102 — chat session pointer landed

**Result: `howl-chat-current-session` is a chat-id and UI-selection pointer;
full catalogs hydrate from their domain repositories.**

`chat-session-pointer.js` is the only writer and reader of the key. The pointer
is `{version: 2, chatId, selection: {worldId, scenarioId, characters,
lorebookIds, activeCharacterIds, removedCharacterIds, addedCharacterIds,
activeLorebookIds, removedLorebookIds, addedLorebookIds}}`. Character entries
keep `{id}` plus only the UI fields `status`, `active`, `stepOutReason` and
`stepOutMode`; no names or other catalog records are stored. `hydrate()` awaits
`WorldRepository`, `ScenarioRepository`, `CharacterStorage` and `LorebookRepo`,
resolves records by id, and merges the stored character UI fields. A legacy
session blob is converted to the id-only shape once, kept only in memory for
first hydration, and rewritten as a pointer. `LorebookRepo` gained `all()` and
is loaded on the chat page.

`context-data.js` delegates writes to `ChatSessionPointer.write` and loads via
`read()`/`hydrate()`. `sidebar-session.js` writes after session load and chat
data cache saves; its load/new-chat/delete paths use `clear()`. `data-sync.js`
uses `clear()` on reset paths. `app-wipe.js` still names the key only so a full
reset removes it. A static test rejects any other module naming the key.

**Verification (edited tree).** `node --check` on the five edited renderer
modules: 0 each. `node --test tests/renderer/*.test.js`: 39/39 pass, including
6 new pointer tests. `go build ./...`: 0. `go vet ./...`: 0. Full
`go test ./... -skip TestProbeMicrophoneOnHost -count=1`: 0.

**User-visible evidence.** Browser dev host `127.0.0.1:8799` with an isolated
OS-temp root: created a chat; selected the fixture world, scenario and
character through the UI. The raw pointer held only the chat id and ids
(`item5-world`, `item5-scenario`, `item5-character`) with no catalog text;
live hydration returned the canonical titles from the repositories; reload
restored all three in the UI. Packaged WebView2 host, rebuilt from the edited
tree with a separate isolated root: repeated the same selection workflow, read
the WebView2 Local Storage record back with only the chat id and ids, then
closed and reopened the app and saw all three panels restored. No packaged
devtools were needed; the record was read from the WebView2 profile after the
UI showed the final selection.

**Code contents (10 files):** `chat-session-pointer.js` (new), `context-data.js`,
`sidebar-session.js`, `data-sync.js`, `lorebook-repository.js`, `index.html`,
`story-tools.html`, `tests/renderer/chat-session-pointer.test.js` (new), HOC
§21.11 and this checkpoint. No Go source, eval, experiments, config or build
artifacts. Nothing pushed.

**Left open.** §21.5 item 6 remains dependent on per-domain cutover markers;
only `LorebookRepo` currently reports `cutover`. Other context-entry paths that
add lorebooks still read the pre-cutover lorebook cache directly. The pointer
itself no longer duplicates catalog records.

**Resume token:** `HEKI-CP-2026-11-23-102`
**Next:** per-domain cutover/rollback markers, then §21.5 item 6.
