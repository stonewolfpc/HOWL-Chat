# TASK-03 — spike findings 03: directory-object removal

| Field | Value |
|---|---|
| Purpose | One question only: can the acquired **directory object itself** be removed through a handle-bound operation, with identity binding? |
| Status | **Evidence record. No contract row edited. Deletion disabled** |
| Scope excluded | P2/USN explicitly **not** investigated. E-27 **not** attempted (needs elevated VHD). No child-file result re-litigated |
| Artifacts | `build/tmp/spike5/probe5.go`, `build/tmp/spike5-dir-removal.log`, `build/tmp/spike5b-dir-removal.log` |
| Environment | Windows 11 build 26200 (25H2), NTFS `C:` serial `0x885A56435A562DE6` |
| Recovery token | `TASK-03-SPIKE-FINDINGS-03-2026-11-22` |

Labels: **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished**.

---

## 1. Verdict

> **Directory-object removal: PASS** via parent-relative open + identity verification + disposition.

This verdict covers the directory object only. It does **not** upgrade step 4, and by itself it authorizes nothing.

---

## 2. The run, step by step

Fixture: `parent/subject/inner/` under `build/tmp/spike5/fx`, removed at the end; residue none.

| Step | Operation | Result |
|---|---|---|
| A | enumerate the **parent** by handle | 1 entry, `subject` found with its ID |
| B | open `subject` **by name relative to the parent handle**, with `DELETE \| FILE_LIST_DIRECTORY \| FILE_READ_ATTRIBUTES \| SYNCHRONIZE`, `FILE_DIRECTORY_FILE` | success; **identity bound** — ID and 64-bit serial match the enumerated values |
| C | dispose the bound handle **while non-empty** | **fails, Win32 145** (`ERROR_DIR_NOT_EMPTY`) |
| D | enumerate `subject` by handle; open `inner` relative; verify ID; dispose | **success**; ID verified before disposition |
| E | re-enumerate the **parent** after child removal | `subject` **still listed** — ordering confirmed |
| F | dispose the directory through the same bound handle | **success** |
| G | re-enumerate parent; open `subject` by ID; read-only `os.Stat` | `subject` not listed; by-ID **87**; `os.Stat` **exists=false** |

**F-16 (Observed) — the direct answer.** A directory can be opened by name relative to its parent's handle, bound to the enumerated identity, and disposed through that handle once empty. The guarantee that applies to child files applies to the directory object too.

**F-17 (Observed) — step C is the control that makes the rest meaningful.** Disposing a non-empty directory fails with `ERROR_DIR_NOT_EMPTY`, so the successful disposition in step F is not a no-op: the call rejects the case it should reject and accepts the one it should accept. Without C, F would prove only that the API returned success.

**F-18 (Observed) — the emptiness transition is observable through the same handle.** Step E shows the parent still listing `subject` after its child is gone, and step F then succeeds. The sequence's ordering (children first, directory last) is enforced by the filesystem, not merely by convention.

**F-19 (Observed) — by-ID handles remain non-disposable for directories too.** A live, **empty** control directory opened by ID and disposed failed with **87**, matching findings 02 F-10 for files. So the by-ID route is not an alternative for directory objects either; the parent-relative name open is the only route that worked.

**F-20 (Observed) — post-conditions confirmed twice.** Absence of the directory was confirmed by handle-relative parent enumeration *and* by a read-only `os.Stat`. The second is a path-based check used strictly as a read-only cross-check, never as an authorization input.

---

## 3. What this changes in the sequence

1j W-09 steps 1–3, corrected by findings 02 and now covering the directory object:

| Step | Operation | Identity binding |
|---|---|---|
| 1 | open the acquired directory **by ID** through the held hint | `FILE_ID_INFO` ID + 64-bit serial |
| 2 | enumerate **by handle** | child IDs from the listing |
| 3a | per child: open **by relative name**, with `DELETE` | `FILE_ID_INFO` equals the enumerated ID **and** serial |
| 3b | dispose that bound handle | — |
| 3c | dispose the **directory** through its parent-relative handle, after it is empty | same binding, established at open |
| 4 | verify absence | **still unresolved** |

**F-21 (Inferred) — the structure this reveals.** The parent handle anchors enumeration; relative-name opens acquire child handles; the post-open ID check binds each operation back to the enumeration evidence; and the directory is removed last through the handle that was bound at open. Every destructive act is preceded by an identity check on the very handle being used. That is a coherent design, and it is what the contract's W-09 was trying to say before it named the wrong API.

---

## 4. What remains outside this verdict

| Item | Status |
|---|---|
| E-28 step 4 (positive absence) | **unresolved** — P1 refused; no mechanism established. `complete` still unreachable |
| Nested directories (depth > 1) | **not attempted** — this run had one level |
| Large directories / partial failure mid-pass | **not attempted** |
| Reparse points, hardlinks, open-by-others children | **not attempted** |
| E-27 held-handle continuity | unestablished; needs elevated VHD, not requested |
| E-26 duplicate GUID | unestablished |
| P2 / USN journal | **not investigated**, by instruction |

**F-22 — Scope discipline.** Two of these — nested depth and partial mid-pass failure — are the obvious next questions for step 3 and are cheaper than P2. Neither was pulled into this run, because the instruction was one question with one verdict.

---

## 5. Procedure

Diagnostics were **read-only by construction** in this run: the only cross-check outside the handle path was `os.Stat`, which cannot change state. Findings 02 F-15 (a path-opened control handle deleting a fixture file) is what made that a rule rather than an intention. Residue: none; `build/tmp/spike5/fx` removed.

**No contract row edited. No deletion enabled. No attribution restored.** Rows 540–543 unchanged; this record is row 544. TASK-03 remains in preparation and G-1 is intact.
