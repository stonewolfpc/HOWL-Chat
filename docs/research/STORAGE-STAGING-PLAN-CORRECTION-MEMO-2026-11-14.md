# `STORAGE-STAGING-PLAN-CORRECTION-MEMO`

**Identity:** `HOWL-STAGE-MEMO-1` · **Date:** 2026-11-14 · **Applies to:** `STORAGE-DIAGNOSTIC-STAGING-PLAN-2026-11-14.md` (`967a6f2a…`)
**Status:** correction memorandum. **The staging plan is not edited.** Where this memo and the plan differ, this memo governs. Nothing here authorises a live step.

---

## 1. T-15: retire the overstated wording

Three phrasings in my own reports go beyond what the evidence supports and are **retired**:

| Location | Retired phrasing | Why it overstates |
|---|---|---|
| `HISTORICAL-H-READ-RECONCILIATION` §2.1 | "two read mechanisms returned different byte streams … **at the same time, from the same file object**", and "the file did not have a single well-defined content for readers at that instant" | the two reads were ~19 s and ~49 s apart; "same instant" is not established, and "inside the stack" prejudges the layer |
| `GGUF-DLL-COMPETING-HYPOTHESES` §2 | "a difference **at one moment, inside the stack**" | same |
| `WINDOWS-CACHE-HYPOTHESIS-REVIEW` §10 | the directional reading away from a writer, built partly on the DLL incident | the DLL's source was `H:` on a volume reporting `Full Repair Needed`; it is not evidence about `G:` (see `STORAGE-INVESTIGATION-CLAIM-CORRECTIONS` C-02) |

**Adopted wording, used verbatim from here on:**

> Two read paths returned different contents for the same file identity within one minute. This demonstrates read-path-dependent results, but the responsible layer and the correctness of either result remain undetermined.

### 1.1 The distinctions the retired wording blurred

| Distinction | Why it matters |
|---|---|
| Different read paths ≠ independent physical reads | both share the OS, driver stack, controller, device and host memory; a divergence narrows *where* but identifies no component |
| Unchanged size and mtime ≠ no write | a writer can alter content, and a copy tool can restore the timestamp afterwards (`/COPY:DAT`) |
| `in_use: false` = point-in-time | it says nothing about the seconds between observations, nor about kernel-mode actors without a user handle |
| Unbuffered ≠ uncached everywhere | it bypasses the Windows **system file cache** under the documented conditions; the device cache, mapped views and the filter stack are untouched |
| One valid run | the 19:15:48 attempt **aborted** (`RESULT: ABORTED`, missing intermediate file); T-15 rests on the 19:16:03 run alone |

---

## 2. Decision language: no single metadata field is conclusive

The plan's S-3 decision table converted a metadata reading into a conclusion. Replaced by a classification that permits — and expects — "inconclusive".

| Verdict | Meaning | What must be present |
|---|---|---|
| **supports X-1** (a re-copy occurred) | evidence that a second write happened | a copy-run record, a creation time that post-dates both readings **while no size/mtime-preserving tool is plausible**, or an ingestion manifest that disagrees with the file |
| **supports X-2** (no re-copy; two volumes varied) | the same file object read differently on two devices | identity evidence (file id + volume serial) for both readings, stable metadata across them, and no copy record |
| **supports X-3** (path or identity change) | the readings were of different objects | a file id or serial difference, or a path/rename record |
| **inconclusive** | none of the above is established | **the default and expected outcome for a robocopy-based copy chain** |

**Rule added to the plan:** a field may *contribute* to a verdict but may never *be* the verdict. Specifically prohibited: "creation time is unchanged, therefore no write occurred"; "mtime is unchanged, therefore the content is unchanged"; "exclusive open succeeded, therefore no writer exists".

This is not pedantry. `S3-METADATA-TEST-RECONCILIATION` shows that for robocopy with `/COPY:DAT`, a re-copy and no re-copy are **metadata-identical**, so the strongest available field proves nothing about the question S-3 was written to answer.

---

## 3. S-3 reclassified

| Original | Corrected |
|---|---|
| "`H:` metadata only — the X-1/X-2 discriminator" | "`H:` and `G:` metadata collection — **non-discriminating by itself**; expected verdict *inconclusive*" |
| Ranked first among storage-touching steps because "its result halves the hypothesis space" | Demoted. It no longer halves anything. It is still cheap and should still be collected — as evidence, not as a test |
| Compares `H:` timestamps | Also compares against `D:\H-MIGRATION-2026-11-14` staging metadata and, if it survives, `G:\AI GGUF\e-manifest.tsv` (written by `build/move-to-g.sh`, one row per expected file with size and mtime) |

---

## 4. S-4: handle enumeration, with a tool policy

### 4.1 Value asymmetry, stated plainly

| Result | Value |
|---|---|
| **Positive** — a process holds the exact file | high: names a candidate and, if the granted access mask is obtainable, distinguishes read-only from write-capable |
| **Negative** — nothing holds it | **near zero, historically.** It says nothing about a scheduled writer, an intermittent one, a write two hours ago, or a kernel-mode actor. It must never be reported as "no writer" |

### 4.2 Tool policy

No third-party handle utility is to be installed to satisfy this step. In order:

| Option | Status |
|---|---|
| **Exclusive-open sharing probe** (already used at 13:03:38: `GENERIC_READ`, `dwShareMode = 0`, six paths, 0 blocked) | **available, native, already exercised** — proves absence-by-exclusion at one instant; reports no access mask; briefly holds the file |
| **Restart Manager** via a small in-repo PowerShell/C# shim | read-first choice for a *named* file: returns the processes using it, native API, no install. Requires the shim to be reviewed before use |
| `openfiles /query` | **unavailable as configured** — `fltmc` output states "The system global flag 'maintain objects list' needs to be enabled"; enabling it is a system change and out of scope |
| `handle.exe` (Sysinternals) | **not authorised.** Provenance, signature and privilege review required first; not installed and must not be installed merely to check a box |

Elevation: the probe needs none; Restart Manager needs none for same-session processes; `handle.exe` would need `SeDebugPrivilege` and therefore elevation.

---

## 5. S-7 / S-8: clean counters are never exoneration

The plan says so in prose; the results template must enforce it. The following conversions are **prohibited**, and the template ships with the replacement sentence:

| Collected | Prohibited conclusion | Permitted sentence |
|---|---|---|
| SMART attributes normal | "the drive is healthy" | "no SMART attribute is currently flagged; SMART is firmware-reported and blind to many localised read faults" |
| No WHEA events | "memory is fine" | "no WHEA record exists for the period; non-ECC faults of this kind are typically invisible" |
| Memory test passed | "RAM is ruled out" | "the test found no fault in the tested region during the tested interval" |
| No filter of interest | "no agent can act" | "no *registered* minifilter matches the hypothesis; kernel callbacks and unregistered drivers are not covered" |
| Control file reads consistently | "the volume is healthy" | "one other file read consistently; this certifies the tool and the path resolution, not the volume, device, cache or memory" |

---

## 6. S-9a: conditional, and redefined

Retained as **proposed, not authorised**, pending review of the specification (`S9A-RANGE-PROBE-SPECIFICATION-2026-11-14.md`). Three corrections to the plan's version:

1. **"Drop the file's cache" is withdrawn.** Documented Win32 has no per-file cache invalidation available to an unprivileged caller; `FlushFileBuffers` needs `GENERIC_WRITE`, and global cache control needs privileged tooling. The specification must not describe an experiment it cannot perform.
2. **Tail handling.** T-15's unbuffered digest covered the body unbuffered and the last 4,064 bytes buffered. Any repeat must hash body and tail **separately** and report both, so attribution is exact rather than inferred.
3. **A ground-truth variant exists** and should be preferred: a known-pattern scratch file written unbuffered, flushed, verified, then read buffered. With known content, a divergence finally has an oracle. That variant writes to `G:` and therefore belongs to gate 1, not to a read-only step.

---

## 7. Ordering, revised

| Rank | Step | Change from the plan |
|---:|---|---|
| 1 | Record and artefact reconciliation (S-0) | unchanged — and materially expanded by `S3-METADATA-TEST-RECONCILIATION` |
| 2 | External reference hashes for the three candidates | **new, and first among live-adjacent steps** — public metadata, no storage, and the only thing that can say whether the file was ever right |
| 3 | Build / CU / boot history (S-1) | unchanged, still gated on the event log for the shutdown part |
| 4 | Filter *configuration* and driver versions (S-2 remainder) | membership is already known; this is now a narrower step |
| 5 | Volume metadata incl. `e-manifest.tsv` (S-3) | demoted, verdict expected *inconclusive* |
| 6 | Restart Manager / sharing probe on the exact path (S-4) | native route first; positive-only value |
| 7 | ACLs (S-5), scheduled writers (S-6) | unchanged |
| 8 | SMART and link counters (S-7), WHEA (S-8) | unchanged, with §5's template rules |
| 9 | S-9a range probe | conditional on the specification being reviewed |
| 10 | Three cold reads of a candidate (S-9b) | last; a resumption gate, not a diagnostic |

**Standing rule added:** the investigation may not run *more* reads because the read-path signal is interesting. Interest is not authorisation. Every read beyond the plan needs its own decision.

---

## 8. What this memo does not do

It does not edit the staging plan, open a gate, schedule a step, or reinterpret an owner decision. It restates T-15 conservatively, removes a step's false discriminator status, replaces a categorical decision table with one that can return *inconclusive*, sets a tool policy that prefers native APIs, and forbids the four conversions from "no evidence of fault" to "no fault".

**Establishes:** the corrected decision language, the corrected T-15 statement, the handle-tool policy, the counter-interpretation rules, and the revised order.

**Does not establish:** that the revised order will produce a verdict. It removes one step that could not, and adds one — external hashes — that can.
