# TASK-03 — spike findings 02: enumeration works, by-ID handles cannot be disposed, P1 fails

| Field | Value |
|---|---|
| Purpose | Record the P1 and E-28 steps 2–3 investigation, with independent verdicts |
| Status | **Evidence record. No contract row edited. Deletion disabled** |
| Relationship | Adds to rows 540–542; none edited |
| Artifacts | `build/tmp/spike4/probe4.go`, `build/tmp/spike4-P1-E28.log`, `...4b`, `...4c`, `...4d` |
| Environment | Windows 11 build 26200 (25H2), NTFS `C:` serial `0x885A56435A562DE6` |
| Recovery token | `TASK-03-SPIKE-FINDINGS-02-2026-11-22` |

Labels: **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished**.

---

## 1. Results, independently reportable

| Investigation | Question | Verdict |
|---|---|---|
| **E-28 step 2** | Can children be enumerated entirely through the directory handle? | **PASS** |
| **E-28 step 3** | Can children be removed through handles derived from that object? | **PASS — with a corrected derivation** (see F-10, F-11) |
| **P1** | Can parent enumeration positively establish absence? | **FAILS** |
| **E-28 step 4** | Positive absence verification | still unresolved |
| **E-25** | Absence discriminator | remains `unverifiable` (row 542) |

Per the charter, step 3's success does **not** upgrade step 4, and P1's failure does **not** invalidate steps 2–3. Neither result authorizes deletion.

---

## 2. Findings

### F-9 — Handle-relative enumeration works (**Observed**)

`NtQueryDirectoryFile` on the directory handle opened by ID returned all 5 children with their 64-bit file references, **with no path involved anywhere**. `GetFileInformationByHandleEx` with the corresponding directory classes returned `ERROR_NO_MORE_FILES` on the first call for a non-empty directory in this environment — the Win32 wrapper did not work here, the native call did. That difference is recorded as environment-specific behaviour, not as a general claim.

### F-10 — Handles obtained by `OpenFileById` cannot be disposed (**Observed**)

Both disposition classes returned `ERROR_INVALID_PARAMETER` (87) on every by-ID handle:

| Handle | `FileDispositionInfoEx` (21) | `FileDispositionInfo` (4) |
|---|---|---|
| opened **by ID** | 87 | 87 |
| opened **by path** (same file) | **success** | **success** |

So the failure is a property of the handle, not of the class constant, the access mask, or the fixture. This matches previously observed behaviour for by-ID handles: operations that require a named file object are rejected with `STATUS_INVALID_PARAMETER`.

**Consequence:** 1j W-09 step 3's phrase "handles derived from that directory object" cannot mean by-ID opens. As literally written it was satisfiable, but the by-ID derivation of it is impossible.

### F-11 — Parent-relative name opens provide the missing derivation (**Observed**)

`NtCreateFile` with `RootDirectory = <directory handle>`, `ObjectName = <child name from enumeration>`, no path string anywhere, gave a **named** handle on which disposition worked:

```
by-id opens attempted 5 (none disposable); relative opens 5; removed 5/5;
ids verified 5/5; remaining 0
```

Every child was removed, and after opening each one its `FILE_ID_INFO` was compared against the ID returned by enumeration — **both ID and 64-bit volume serial matched for all five**.

**F-12 — The corrected sequence.** Steps 1–3, all handle-bound, no path strings:

| Step | Operation | Identity check |
|---|---|---|
| 1 | open the directory **by ID** through the held hint | `FILE_ID_INFO` ID + 64-bit serial match |
| 2 | enumerate **by handle** (`NtQueryDirectoryFile`) | child IDs collected from the listing |
| 3a | open each child **by name relative to the directory handle** with `DELETE` | `FILE_ID_INFO` must equal the enumerated ID **and** serial |
| 3b | dispose that named handle | — |
| 4 | verify absence | **unresolved** |

**F-13 — Why the name derivation is not a weakening.** The name comes from the enumeration, not from a path, and the ID check after opening is what makes it safe: if the name now refers to a different object — renamed, replaced, or recreated — its ID will not match and nothing is deleted. That closes the open-by-name race by *skipping*, not by assuming. The residual risk is that a matching ID is deleted while a second object with the same name appears later in the same pass, which only leaves residue for the next pass to observe.

Note the asymmetry this creates and record it plainly: **removal needs names; verification needs IDs.** The design does not get to use one mechanism for both.

### F-14 — P1 fails as an absence proof (**Observed**, decisive)

| Stage | Parent enumeration |
|---|---|
| before move | target **present**, matching ID |
| after moving the directory into a sibling directory | **absent** — 0 entries, ID absent |
| after that, is the identity alive? | **yes** — openable, ID and serial match |
| after a real delete | **absent**, same observation |

The observation "the ID is not among this parent's children" is produced identically by *moving the object elsewhere* and by *deleting it*. Enumeration therefore establishes **"not a child of this parent"**, not **"the identity is absent"**. A moved directory would be reported as removed.

That is the qualification the charter asked to test, and it does not survive. P1 is refused as an absence proof.

**Calculated consequence:** step 4 still has no positive mechanism. `complete` remains unreachable, and rows 542's conclusion stands unchanged.

### F-15 — A diagnostic deleted a file during an "enumeration" experiment (**Observed**)

To separate "wrong class" from "handle property", the probe disposed a file through a **path-opened** control handle. It succeeded, and that file was deleted. The probe reported the resulting count mismatch ("remaining 4") rather than hiding it, but the lesson is real: a diagnostic can be destructive, and in a deletion investigation every call that changes state must be assumed destructive even when labelled a control. SP-35's discipline — report effects before cleaning up — is what made this visible.

---

## 3. Verdicts and what they do not establish

| Item | Status |
|---|---|
| E-28 step 1 (open by ID through a directory handle) | established, conditional on ID + serial check |
| E-28 step 2 (enumerate by handle) | **established**, this build and volume |
| E-28 step 3 (remove via derived handles) | **established via parent-relative name opens**; by-ID derivation impossible |
| E-28 step 4 (positive absence) | **unresolved**; P1 refused; `complete` unreachable |
| Removing the directory object itself | **not attempted** — only files; the directory's own disposition is untested |
| E-27 (held-handle continuity) | unestablished; needs an elevated VHD |
| E-26 (duplicate GUID) | unestablished |
| E-25 | unverifiable by open-by-ID |

---

## 4. Authorization and state

Fixture created and removed under `build/tmp/spike4/fx`; residue none in the final run. One fixture file was deleted by the diagnostic described in F-15, then the whole fixture directory was removed. No Howl data, no conversation directory, no `E:` access, no elevation.

**No contract row edited. No deletion enabled. No attribution restored.** Rows 540–542 unchanged; this record is row 543.
