# TASK-03 — spike findings 04: nested depth and truthful mid-pass failure

| Field | Value |
|---|---|
| Purpose | Two independent questions: does the sequence work at depth, and does it fail truthfully partway |
| Status | **Evidence record. No contract row edited. Deletion disabled** |
| Scope excluded | P2/USN not investigated. E-27/E-26 not attempted. No re-litigation of steps 1–3 |
| Artifacts | `build/tmp/spike6/probe6.go`, `build/tmp/spike6-nested-partial.log` |
| Environment | Windows 11 build 26200 (25H2), NTFS `C:` serial `0x885A56435A562DE6` |
| Recovery token | `TASK-03-SPIKE-FINDINGS-04-2026-11-22` |

Labels: **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished**.

---

## 1. Verdicts

| Investigation | Verdict |
|---|---|
| **N) nested depth** (`a/b/c/f.txt`) | **PASS** — 3/3 levels identity-verified, removed deepest-first, top empty afterwards |
| **P) partial mid-pass failure** | **PASS** — one removal, truthful stop, identity intact, directory not removable until emptied, recovery completes |

Step 4 remains unresolved. Neither verdict authorizes anything.

---

## 2. Nested depth

| Step | Observation |
|---|---|
| open `a`, `b`, `c` by name relative to the current level's handle | identity match **at every level** (3/3) |
| enumerate deepest level, open `f.txt` relative, verify ID, dispose | ok |
| dispose level 3, close its handle, then check parent listing | ok; **parent now lists 0** |
| dispose level 2, close, check parent | ok; parent lists 0 |
| dispose level 1, close, check parent | ok; top lists 0 |
| top enumeration after unwind; `os.Stat(a)` | **0 entries**; exists=false |

**F-23 (Observed) — depth adds no new mechanism.** Each level is opened relative to the handle one level above, verified, and disposed in reverse order. No path is constructed at any depth, and the identity binding holds at every step.

**F-24 (Observed) — a child's handle must be closed before its parent is disposed.** This was measured the hard way: the first run disposed level 3 and *kept its handle open*, then failed to dispose level 2 with `ERROR_DIR_NOT_EMPTY` (145). **An unlinked-but-still-open child directory remains visible to its parent**, so the parent is not empty yet. The correct order is dispose → close → then dispose the parent, and the run above shows each parent reaching 0 entries at exactly that point.

That is an operational rule the implementation must carry: it is not enough to request POSIX-semantics deletion; the handle must also be released before moving up a level.

---

## 3. Partial mid-pass failure

Fixture: `part/` containing `a_ok.txt`, `m_nonempty/` (holding `inside.txt`), `z_ok.txt` — names chosen so a removal succeeds *before* the blocker is reached.

| Step | Observation |
|---|---|
| pass starts: 3 children | `[a_ok.txt m_nonempty z_ok.txt]` |
| `a_ok.txt` open relative, verify ID, dispose | **removed** |
| `m_nonempty` open relative, verify ID, dispose | **fails 145** (`ERROR_DIR_NOT_EMPTY`) → **pass stops** |
| remaining after stop | **2** — `[m_nonempty z_ok.txt]`, exactly the unreached remainder |
| directory identity | **unchanged** (ID equal to the pre-pass ID) |
| directory dispose attempt | **fails 145** — cannot be removed while content remains |
| read-only `os.Stat` | **exists=true** |
| recovery: empty blocker through its own handle, verify each ID, dispose | ok |
| recovery: dispose blocker, then `z_ok.txt` | ok |
| final | 0 entries, directory disposed, `os.Stat` false |

**F-25 (Observed) — the failure is truthful.** Partial progress is visible as a smaller remaining set, not as a completed removal. The directory survives, keeps its identity, and refuses disposition. Nothing in the reported state claims more than happened.

**F-26 (Observed) — the operation can resume from observed state alone.** Recovery needed no journal of what had already been removed: a fresh enumeration of the same handle showed what remained, and the pass continued. That matters because it is the property a crash-recovery path would need — **though it says nothing about absence verification, which remains step 4's problem.**

**F-27 — Two harness defects preceded this result, and both mattered.**

| Defect | Consequence | Fix |
|---|---|---|
| Unwind disposed a level while its handle stayed open | level 2 failed 145 → **false "nested FAIL"** | close before disposing the parent (F-24) |
| Blocker was a handle with share mode `0x03` (excluding `FILE_SHARE_DELETE`) | it **did not obstruct** — all three children were disposed, so nothing was tested | replaced with a non-empty child directory, a failure mode already proven in spike 5 |

Both were my errors, both produced results that would have been read as platform limitations, and both are recorded instead of being quietly re-run. The second is the more instructive: an obstruction that does not obstruct turns a test into a tautology.

---

## 4. What this still does not establish

| Item | Status |
|---|---|
| Step 4 (positive absence) | **unresolved**; `complete` unreachable |
| Reparse points, hardlinks, junctions | untested |
| Children held open by other processes at dispose time | untested — the share-mode attempt failed to produce this condition |
| Large directories (thousands of entries, buffer paging) | untested |
| Concurrent external mutation during a pass | untested |
| E-27 held-handle continuity | unestablished, needs elevated VHD |
| E-26 duplicate GUID | unestablished |

**F-28 — The share-mode result deserves its own note.** A held handle whose share mode excluded `FILE_SHARE_DELETE` did **not** prevent disposition of that file. That is unexpected, it contradicts a plausible reading of share-mode semantics, and it means a later test of "child open in another process" needs a different method (another process, not another handle in this one). Recorded as an anomaly to investigate, not as a conclusion.

---

## 5. Procedure and authorization

Fixtures removed; residue none. Diagnostics read-only (`os.Stat`, enumeration). No `E:` access, no elevation, no Howl data touched.

**No contract row edited. No deletion enabled. No attribution restored.** Rows 540–544 unchanged; this record is row 545. G-1 intact; TASK-03 remains in preparation.
