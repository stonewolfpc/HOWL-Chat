# Registry Integrity Checker — Enhancement Specification

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identifier:** `CORPUS-67` · **Version:** 1 · **Date:** 2026-11-12.
**Authorizing decision:** `CORPUS-62` Decision 7 (owner-approved 2026-11-12).
**Assigned authoring model identity (owner-supplied metadata):** DeepSeek. Not inferred, not substituted.
**Execution:** fresh conversation.
**Status:** **SPECIFICATION ONLY. NO CHECKER IMPLEMENTATION OR MODIFICATION IS AUTHORIZED.** `check_corpus.py` is byte-identical to its pre-filing state.

---

## 1. Purpose and boundary

Decision 7: *future registry integrity validation should fail when a registered target is missing or incorrectly linked.*

This document specifies what such a validator must do, what it must refuse to do, and how each behaviour is demonstrated by fixture. It authorizes nothing.

**Explicitly out of scope for this document:** writing the validator; modifying `check_corpus.py`; repairing any registry row; regenerating the registry; changing the registry's schema, numbering, or family vocabulary; deciding whether the validator is a separate tool or a new mode of the existing checker (that is an implementation decision requiring separate authorization).

---

## 2. Verified starting point — what the current checker actually does

Read in full this session (`docs/research/corpus/check_corpus.py`, 10,098 bytes).

| Capability | Present? | Evidence |
|---|---|---|
| Parses corpus **item** files (`CORPUS-02`, `-03`, `-04`, or explicit paths) | **Yes** | `DEFAULT_FILES`, `parse()` |
| Required-field presence per scored unit | **Yes** | `PARTS`, presence loop |
| Structural rules `C-1`, `C-2`, `C-3`, `C-4`, `U-1` | **Yes** | corresponding loops |
| Recomputes split/totals and compares against `CORPUS-01`'s `Total` row | **Yes** | `declared` regex |
| **Opens `docs/DOCUMENT-REGISTRY.md` at all** | **No** | no occurrence of `REGISTRY` or `DOCUMENT-REGISTRY` in the source |
| **Resolves any registry target** | **No** | same |
| **Verifies any registry hash, byte count, or line count** | **No** | same |
| **Detects duplicate registry rows or conflicting registrations** | **No** | same |

**Consequence, verified by observation and not merely inferred:** on 2026-11-12, before row 393 was corrected, the registry carried a row whose linked target did not exist, and the checker returned:

```text
RESULT: PASS (0 fail, 0 warn, 0 unverifiable)
```

Both facts are true simultaneously and are not in tension: the checker validates item records, not registry linkage. **No claim in any prior report or in this document asserts that the existing checker validates registry target existence.** It does not, and this session's run does not establish that it does.

---

## 3. Required behaviour

Grouped by Decision 7's bullets. Each requirement is normative ("must").

### 3.1 Target resolution

| ID | Requirement |
|---|---|
| `R-1` | Every registry row that names a file must resolve to an existing file. |
| `R-2` | Resolution must be relative to a **declared, repository-root-based scope**, not to the directory the tool happens to run in, and not to a hard-coded subdirectory. |
| `R-3` | Resolution must support rows that point outside the main documentation directory — upward links (`../…`) and root-level paths — because the registry legitimately contains at least two such rows (the root requirements ledger; the root `CORPUS-58` assignment brief). |
| `R-4` | A row whose target exists but whose **identity does not match the registration** must fail as mislinked, not pass. Identity match means: the on-disk basename equals the registered basename, **compared case-sensitively**, and the resolved path is the one recorded. On case-insensitive filesystems this requires an explicit case comparison against the directory listing, because a naive `os.path.exists` will accept `FOO.md` for `foo.md`. |
| `R-5` | A row that resolves to a file already claimed by a different row must be reported as a **conflicting registration** (§3.3). |

### 3.2 Metadata verification

| ID | Requirement |
|---|---|
| `R-6` | **Byte count** must equal the on-disk size in bytes. |
| `R-7` | **Line count** must equal the number of LF (`0x0A`) characters in the file. This is the registry's current convention and must be stated in the report so that generator and validator cannot drift. |
| `R-8` | **SHA-256** must equal the lowercase hex digest of the raw bytes, unwrapped. |
| `R-9` | Binary rows (line count recorded as `—`) must skip the line check and still verify bytes and hash. |
| `R-10` | A row with any field that cannot be parsed must **fail closed** (§3.4 `UNPARSABLE_ROW`), never be skipped. |

### 3.3 Structural integrity

| ID | Requirement |
|---|---|
| `R-11` | **Duplicate row numbers** must fail. |
| `R-12` | **Duplicate identifiers** — two rows naming the same file — must fail. This is distinct from `R-11`: numbering can be unique while two rows still point at one artifact. |
| `R-13` | **Conflicting registrations** — two rows for the same artifact with differing hash, size, line count, or role — must fail and must name both rows. |
| `R-14` | **Unregistered files inside the declared scope** must be reported. `CORPUS-62` §9.4 records two such files at filing time. |
| `R-15` | The report must name the **exact row number** and the **registered path** for every failure. A failure that says only "1 missing" is non-conforming. |

### 3.4 Failure taxonomy — the classes must be distinguishable

Decision 7 requires the validator to distinguish missing from mislinked from metadata mismatch from hash mismatch. The taxonomy is therefore normative:

| Code | Meaning | Distinguishing test |
|---|---|---|
| `OK` | Row verified | Target exists, identity matches, all three metadata fields match |
| `MISSING_TARGET` | Nothing exists at the registered path | `stat` fails |
| `MISLINKED_TARGET` | A file exists at the path but is not the registered artifact (basename/case/identity mismatch), **or** the row's path points to a different artifact than the one the row's metadata describes | identity comparison (`R-4`) |
| `METADATA_MISMATCH_BYTES` | Size differs | `R-6` |
| `METADATA_MISMATCH_LINES` | Line count differs | `R-7` |
| `HASH_MISMATCH` | Bytes differ from the recorded digest | `R-8` |
| `DUPLICATE_ROW_ID` | Two rows share a number | `R-11` |
| `CONFLICTING_REGISTRATION` | Two rows claim one artifact, or rows disagree about its metadata | `R-12`, `R-13` |
| `UNREGISTERED_IN_SCOPE` | A file inside the declared scope has no row | `R-14` |
| `UNPARSABLE_ROW` | A table row could not be interpreted | `R-10` |

`MISSING_TARGET` and `MISLINKED_TARGET` must never be merged: the first means "the registration has no object", the second means "the registration has the wrong object". `CORPUS-58`/row 393 was the second, and every report that called it the first was corrected (`CORPUS-61` §5.2).

### 3.5 Failure, not silence

| ID | Requirement |
|---|---|
| `R-16` | Any non-`OK` finding must make the run **fail**. Exit status non-zero. A run that reports findings and exits 0 is non-conforming. |
| `R-17` | The report must state the scope it searched, so a narrow scope cannot masquerade as a clean result (`R-2`, `R-3`). |
| `R-18` | The summary line must separate counts by class, and must not print an aggregate "PASS" when any failure exists. |

### 3.6 Read-only guarantee

| ID | Requirement |
|---|---|
| `R-19` | The validator must not write, create, move, rename, delete, or normalise **any** file. |
| `R-20` | It must not repair, rewrite, or regenerate registry rows, and must not "suggest" a repair by writing one. Suggestions, if any, are report text only. |
| `R-21` | It must not touch version control, the network, or the corpus item files. |
| `R-22` | Its output must be deterministic for a given tree state. |

**Why `R-19`–`R-22` are not optional.** The corpus's own rule is that a tool which silently fixes records erases the difference between authored and corrected content (`CORPUS-01` §6.2). A validator that repaired a hash mismatch would destroy exactly the evidence a hash mismatch exists to surface.

---

## 4. Required test fixtures

Fixture harness: a temporary tree containing a synthetic registry and synthetic artifacts. Fixtures never touch the real registry, the real corpus, or `docs/`.

| # | Fixture | Construction | Expected result | Must-not |
|---|---|---|---|---|
| 1 | **Valid registration** | One artifact; row with correct path, hash, bytes, lines | `OK`; exit 0 | Report any finding |
| 2 | **Missing target** | Row points to a path with no file | `MISSING_TARGET`, row named; exit ≠ 0 | Pass; report `MISLINKED_TARGET` |
| 3 | **Mislinked target** | File exists at the path but is a different artifact (different basename), **or** same basename with different case on a case-insensitive filesystem | `MISLINKED_TARGET`, row named; exit ≠ 0 | Pass on filesystem leniency (`R-4`) |
| 4 | **Incorrect hash** | Correct path, bytes, lines; hash altered by one character | `HASH_MISMATCH`; exit ≠ 0 | Pass; report a byte/line mismatch |
| 5 | **Incorrect byte count** | Correct path and hash; byte field off by one | `METADATA_MISMATCH_BYTES`; exit ≠ 0 | Pass; report only lines |
| 6 | **Incorrect line count** | Correct path, hash, bytes; line field off by one | `METADATA_MISMATCH_LINES`; exit ≠ 0 | Pass; conflate with bytes |
| 7 | **Duplicate identifier** | Two rows, different numbers, same target; one row also differing in hash | `DUPLICATE_ROW_ID` **or** `CONFLICTING_REGISTRATION` (row numbering unique ⇒ conflicting registration), both rows named; exit ≠ 0 | Report a single clean row |
| 8a | **Artifact outside a previously searched directory scope — positive** | Artifact at repository root; registry row linking `../artifact.md`; validator run from a nested directory | `OK`; exit 0; report states the resolved root | Fail because the artifact is outside the documentation directory |
| 8b | **Artifact outside a previously searched directory scope — negative** | Same tree, but the validator is configured with a narrowed scope that excludes the root | `MISSING_TARGET` **or** `UNREGISTERED_IN_SCOPE`, explicitly naming the excluded path and the scope used; exit ≠ 0 | Pass silently because the search did not reach it |

Fixture 8 exists because the real failure happened exactly this way: a hash sweep scoped to `docs/` reported "no file matches" for a file that existed one level up (`CORPUS-59` §3, corrected by `CORPUS-61` §5.1–5.2). 8a and 8b must both be present, because a validator that only passes 8a can still silently ignore out-of-scope artifacts.

**Additional fixtures recommended, not required by Decision 7:** unparsable row (`UNPARSABLE_ROW`), binary row with `—` lines (`OK`), and empty registry (explicit "0 rows" rather than PASS).

---

## 5. Report format

Minimum viable report, normative in structure only:

```text
scope: repository root C:\...\HOWL_Chat
targets checked: 402
line-count definition: count of LF characters

FAIL MISSING_TARGET      row 393  research/corpus/CORPUS-58-….md
FAIL HASH_MISMATCH       row 210  research/….md   registered 1a2b…  actual 9f8e…
WARN UNREGISTERED_IN_SCOPE  docs/research/corpus/CORPUS-26-….md
OK   verified: 400

RESULT: FAIL (2 fail, 1 warn)
No file was modified.
```

Rules: one line per finding, class first, row number second, path third, both values on any numeric or hash mismatch. No aggregate percentage. No "PASS" line when failures exist.

**Severity mapping (proposed):** `MISSING_TARGET`, `MISLINKED_TARGET`, `HASH_MISMATCH`, `METADATA_MISMATCH_*`, `DUPLICATE_ROW_ID`, `CONFLICTING_REGISTRATION`, `UNPARSABLE_ROW` = **fail**. `UNREGISTERED_IN_SCOPE` = **warn** by default, because a registry may legitimately lag a large file addition within a single change; the warn is a finding, not a pass, and a strict mode may promote it.

---

## 6. What this specification deliberately does not require

| Not required | Reason |
|---|---|
| Content inspection of any artifact | Registration is not inspection (`CORPUS-05` §7, rule `P-6`). A hash proves byte identity, not that anyone read the file, and not that a row's role description is accurate. A role/description checker is not proposed and should not be inferred from this document. |
| Registry self-consistency of the summary tables | The registry's §2 family and class breakdowns do not currently sum to the stated total; this was observed this session and is recorded in `CORPUS-62` §9.4 as a pre-existing finding. Adding a totals check here would be a scope extension; it is noted as a candidate for a separate decision, not part of Decision 7's required set. |
| Validation that `CORPUS-01`'s declared totals match the item files | Already implemented by the existing checker (§2). |
| Automatic row generation for new files | That would be a repair, forbidden by `R-19`/`R-20`. |

---

## 7. Acceptance criteria for the future implementation

The implementation is acceptable only when all of the following hold, verified by the owner or an independent reviewer:

1. All eight fixture families (§4) behave as specified, including the two halves of fixture 8.
2. Every failure class is distinguishable in output — specifically, fixtures 2 and 3 produce **different** codes.
3. Exit status is non-zero on any fail-class finding.
4. The tool is demonstrated to write nothing: a before/after hash of the fixture tree and the real registry are identical.
5. A deliberately introduced missing target in a **copy** of the real registry is caught; the real registry is never modified.
6. The report names the scope searched and the line-count convention.
7. The tool's own test run is recorded as a dated verification artifact with hashes.

---

## 8. Authorization status

| Item | Status |
|---|---|
| This specification | **Prepared and filed. Not implemented.** |
| `check_corpus.py` | **Unmodified.** No registry logic added. |
| Registry | **Not repaired by any tool.** Row 393 was corrected by hand under Decision 1 (`CORPUS-63`), not by a checker. |
| Implementation | **Requires separate authorization**, and would itself be a corpus-adjacent change requiring its own record, fixtures, and verification. |

---

## 9. Statement of no changes

**SPECIFICATION ONLY.** No checker, corpus item, standard, registry row (other than this document's own row), schema, runtime, or model was created or modified by this document. `AM-01` unissued; `A-2a-R1` NOT IN FORCE; `D1` unadopted; annotation blocked.

`G1` open · corpus frozen · expansion frozen.
