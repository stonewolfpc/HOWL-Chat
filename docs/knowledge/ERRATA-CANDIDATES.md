# HEKI — erratum candidates (additive; no original is edited)

**Rule:** this file records statements that a later measurement has shown to be stale or incomplete. It never
changes the document that contains them (R1) and never converts a proposal into an adoption (R4). An erratum is
applied only by an owner-authorized act on the original, which is outside this program's scope.

**Vocabulary for an erratum:** *stale statement* → *measured evidence* → *what was applied* → *what that does
**not** establish*.

---

## `E-01` — R6 says "NOT APPLIED"; FC-5 contains its text

### The stale statement

`docs/review/AMENDMENT-DRAFT-FC5-CONFIGURATION-AND-LAUNCH-RECORDS-R6-2026-11-21.md`:

> **Status:** **OWNER-APPROVED AMENDMENT DRAFT — REVISION 6 — WORKING COPY — NOT APPLIED.**

and, in the same document:

> It applies nothing. The canonical FC-5 paper still lacks every clause below, and its registered hash still
> describes the unamended text.

### Measured evidence (2026-11-22)

R6 names its target as 213 lines / 12,326 bytes / `c245096fbb58…`. That exact file is retrievable from the
baseline commit R6 records (`0bc3b091b5023ac867f6d060bd56c481d6f84a3a`), and its hash at that commit is
`c245096f…` — so R6's baseline was accurate when written.

| Check | Result |
|---|---|
| Baseline (commit `0bc3b091…`) FC-5 | 12,326 B · 213 L · `c245096f…` |
| FC-5 on disk now | 23,056 B · 304 L · `626ab4be…` |
| Diff baseline → now | **4 hunks: +91 content lines, −0** |
| R6 quoted content lines (§3–§6) | 48 |
| Of those present verbatim in FC-5 | **48 of 48** |
| Added lines equal R6's quoted lines, in order | **yes, exactly** |
| Contract register §6 | records `626ab4be…` · 23,056 B for FC-5 |
| Document registry row 509 | records `626ab4be…` · 23,056 B · 304 L |

### What was applied — provision by provision

Every one of R6's eight additions is present, in R6's stated placement, with no other change to the paper:

| R6 provision | Placement required | Present in FC-5 | Status as written |
|---|---|---|---|
| Identifier note (unnumbered) | front matter, after `Labels`, before `---` | yes, verbatim | unnumbered guidance |
| §4 lead-in (unnumbered) | under `## 4. Handoff phases`, above FC5-12 | yes, verbatim | unnumbered guidance |
| `FC5-27` configuration decision record | §4, after FC5-15 | yes, verbatim | **Proposed** |
| `FC5-28` launch construction record | §4, after FC5-27 | yes, verbatim | **Proposed** |
| `FC5-29` identities, references and retry | §4, after FC5-28 | yes, verbatim, **revision-6 wording** | **Proposed** |
| `FC5-30` records are not runtime evidence | §4, after FC5-29 | yes, verbatim | **Proposed** |
| `FC5-31` configuration-record retention | §4, after FC5-30, before §5 | yes, verbatim | **Unresolved** |
| `FC5-26a` boundary pointer | §10, after FC5-26, before closing line | yes, verbatim | **Proposed** |

R6's "Unchanged" list also holds: FC5-01 through FC5-26 appear in the paper with no deletion or rewrite — the
diff removes no content line anywhere.

Revision-6 specificity confirmed rather than assumed. R6 says "Revision 6 changes FC5-29 and nothing else". Both
revision-6 corrections are in FC-5: "The configuration decision identity remains the common reference across
those attempts; each construction record retains its own distinct identity. This identity relationship does not
establish any retention or persistence guarantee (FC5-31)." and "A failed construction record is subject to
whatever retention policy is eventually established (FC5-31)."

### What this does **not** establish

- **Not approval.** FC-5's own status line still reads `PROPOSED — reservation questions unresolved`, and
  `Implementation authorization: **None**`. Applying approved text to a proposed paper does not approve the paper.
- **Not adoption of the clauses.** Every applied clause keeps its label — `Proposed` for FC5-27–FC5-30 and
  FC5-26a, `Unresolved` for FC5-31. Text being present is not a clause being adopted.
- **Not resolution of anything the text leaves open.** `FC5-CONFIG-RETENTION` is still an open question;
  `DQ-HRA-*` are untouched; the FC-6 criterion-ID question stays open.
- **Not an M4 change and not implementation.** R6's own §1 names four separate milestones — application,
  register/registry reconciliation, M4, implementation. Only the first is now established, and by measurement
  rather than by a record of it being performed.
- **Not a determination of who applied it or when.** No document records the application event. The commit that
  contains the amended bytes is not identified by this check; only that the working tree differs from `0bc3b091…`.

### Recommended erratum (proposed, not applied)

> **Erratum to R6 status (proposed).** The statement "NOT APPLIED" in the status block of
> `AMENDMENT-DRAFT-FC5-…-R6-2026-11-21.md`, and its §1 restatement, no longer describe the canonical paper.
> As measured on 2026-11-22, `FC-5-…-CONTRACT-2026-11-21.md` contains all eight of R6's additions verbatim, in
> R6's stated placement, with no other content change; the file's baseline form (`c245096f…`, 213 L, 12,326 B)
> is retrievable at commit `0bc3b091…`, and the amended form is 304 lines / 23,056 B / `626ab4be…`, as the
> contract register and document registry already record. This erratum records application only. It does not
> approve FC-5, adopt FC5-26a or FC5-27–FC5-31, resolve `FC5-CONFIG-RETENTION`, change M4, or authorize
> implementation.

**Owner decision required before any original is touched.**

---

## `E-02` — Two documents are named for a date their content post-dates

### The stale statement

Neither document contains a false statement about itself; the defect is in what the corpus lets a reader
conclude. Both are filename-dated 2026-11-20:

- `HOWL-FOUNDATION-FRONT-TO-BACK-…-2026-11-20.md` — §1.1 names six **2026-11-21** documents as controlling
  inputs, and §12 claims to have inspected them. **Unregistered**, so no registry field records a later
  revision.
- `HOWL-FOUNDATION-CONTRACT-PACK-DEEPSEEK-INTAKE-2026-11-20.md` — status line, §0, §0.4 and §8 describe
  2026-11-21 review and decision events. **Registered** (row 512), and that row *does* record "its D2 wording
  was reconciled and accepted 2026-11-21" — so for this document the registry already carries the correction.
- Neither document's own status block records the revision that added that material.

### Measured evidence

15 of 622 filename-dated corpus documents contain a head date later than their filename date (census field
`dates_head`). The two above are the only 2026-11-20 instances.

### What this does not establish

Not that either document is wrong, and not that its filename should change. A filename is an identity, and HEKI
depends on identities being stable (R2). The finding is about **reading order**, not about the files.

### Consequence recorded instead of an edit

HEKI's ordering is a filename ordering. `README.md` §8 and `PHASE-1-TRANCHE-03.md` now say so explicitly. The
second-order hazard is registration: an unregistered document revised in place leaves no record of the revision
anywhere, which is why `HEKI-D-0088` is listed among the nine documents awaiting an owner registration decision.

---

## `E-03` — Already recorded, not repeated as new

| Statement | Where it was stale | Where recorded |
|---|---|---|
| Registry header `Rows: 444` | creation-time count, actual 516 | CP-01; registry header corrected under earlier authorization |
| CP-01 "the ten 2026-11-21 `research/MEMORY-09-*` records" | six exist | CP-02 |
| CP-01 "the identifier layer will materially reduce the 304 isolated documents" | 15 rescued (5%) | CP-02 |
| CP-01 "304 documents with no link in either direction" read as "disconnected" | 129 are unreferenced in all three layers | CP-03 |
| Tranche 01 §3 "the R6 draft cannot govern because its target hash is unchanged" | the hash had changed | CP-02; `RELATIONS.md` §3 |

---

**No original document was modified in producing this file.**
