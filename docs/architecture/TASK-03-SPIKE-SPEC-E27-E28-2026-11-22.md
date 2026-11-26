# TASK-03 — spike specification: E-27 and E-28

| Field | Value |
|---|---|
| Purpose | Specify a bounded Windows feasibility investigation of the identity-bound deletion sequence |
| Status | **Specification only. Nothing executed. Not authorized to execute** |
| Authorizes | nothing |
| Requires before execution | explicit owner authorization to run destructive spikes **and** separate elevated approval for VHD creation/attachment |
| Relates to | 1j W-09 (four required properties), E-27, E-28; 1h U-02 (control open) |
| Recovery token | `TASK-03-SPIKE-SPEC-E27-E28-2026-11-22` |

Prefixes: `SP-` specification clause, `EV-` evidence criterion. Labels: **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished**.

---

## 1. Authorization boundary

**SP-01 — This document authorizes nothing.** It is a plan. Read as: *if* the owner authorizes destructive spikes, *then* these are the exact actions, in this order, with these effects and this cleanup.

**SP-02 — Three separate permissions.**

| # | Permission | Granted? |
|---|---|---|
| 1 | Investigate E-27 / E-28 with destructive operations on **disposable test data** | **no** |
| 2 | Create and attach a VHD (elevated, separate approval prompt) | **no** |
| 3 | Enable deletion inside Howl Chat | **no** — and permission 1 or 2 would not grant it |

**SP-03 — Permission to investigate is not permission to enable.** No result of this spike, however clean, enables deletion in the product. Enabling requires G-1's full acceptance gate, which is a separate owner decision and is not a decision this spike produces.

**SP-04 — Reporting rule.** Every run reports its **actual** effects — objects created, objects destroyed, volumes attached, residual state — whether or not the run succeeded. A spike that silently leaves an orphaned mount is a failed spike.

**SP-05 — No repository implementation.** No code is added to `internal/`, `cmd/` or any product package. The spike is a throwaway tool run from a temporary directory outside the repository, removed at cleanup. If a result suggests a product change, that is a *finding*, recorded — not code.

**SP-06 — Disjointness.** No experiment touches the failing `E:` drive, the HOWL Chat data directories, or any preservation data. Every fixture lives on a scratch volume created for the purpose.

---

## 2. Independence of the two experiments

**SP-07 — Two questions, two results.**

| | Question |
|---|---|
| **E-27** | Does a held volume handle keep subsequent by-ID calls bound to the volume object it was opened on? |
| **E-28** | Can the complete four-step sequence of 1j W-09 be performed with no path re-resolution at any step? |

**SP-08 — Results are independently reportable.** They run in separate sessions, produce separate logs, and are written up separately. Neither may be reported as a conclusion about the other.

**SP-09 — Contamination rule.** E-27 failing does not invalidate an E-28 result. E-28 is a statement about whether the sequence *can be performed*; E-27 is about whether continuity *survives adversarial volume manipulation*. A sequence that works is a real finding even if continuity cannot be proven, and continuity that holds is a real finding even if the sequence cannot be built.

**SP-10 — Neither result alone authorizes deletion.** Both must pass. Additionally:

| E-27 | E-28 | Consequence |
|---|---|---|
| pass | pass | `complete` becomes *reachable*; G-1 still needs the full acceptance run |
| fail | pass | `complete` stays unreachable; 1i's `completion_unconfirmed`/`unverifiable` carry the space permanently unless the design changes |
| pass | fail | the identity-bound design has no implementation; deletion stays unavailable; **no path-based fallback** |
| fail | fail | deletion is unavailable for this design |

---

## 3. Prerequisites

**SP-11 — Environment, recorded per run.**

| Item | Required | How recorded |
|---|---|---|
| Windows version and build | exact string | `[Environment]::OSVersion` / `winver` build number |
| Filesystem of scratch volume | NTFS (primary), ReFS if available (secondary) | `Get-Volume` |
| `FILE_SUPPORTS_OPEN_BY_FILE_ID` | asserted set | `GetVolumeInformation` via the tool |
| Elevation for VHD | available on request only | approval prompt, logged |
| Go toolchain | present, used only to build the throwaway tool | `go version` |
| Free drive letters | two reserved letters, verified unmapped immediately before use | `Get-PSDrive` + `Get-Volume` |

**SP-12 — Letter reservation.** Two letters are reserved for the spike (for example `X:` and `Y:`). Before each attach, both are verified absent from `Get-PSDrive` and from `Get-Volume`. **`E:` is never used and never probed.** A reserved letter that turns out to be occupied aborts the run — the spike does not "find another letter" at runtime, because that is how a fixture ends up on a real volume.

**SP-13 — Access rights must be determined empirically, not assumed.** The tool tries the documented hint form (a handle to any file on the volume) and records, for each access mask, whether the by-ID open succeeds. This is E-25 territory: the required mask is `Unestablished`, and the spike records what actually worked rather than asserting it.

**SP-14 — Fixtures.**

| Fixture | Contents |
|---|---|
| `T1` | directory with 3 subdirectories, 20 files, nesting depth 3 |
| `T2` | directory with 0 children (empty-directory edge) |
| `T3` | directory with 500 files (enumeration volume) |
| `CONTROL` | a file guaranteed live for the duration, used for U-02's control open |
| `CLONE` | a second scratch volume, byte-identical copy of the first, for E-27's adversarial case |

**SP-15 — Fixture creation is non-destructive.** All fixtures are created by the spike on the scratch volume. Nothing is copied *from* user data — file contents are generated filler, not real conversations, so a mis-targeted delete destroys nothing of value.

---

## 4. E-27: does the held handle stay bound?

**SP-16 — What is being tested.** 1i V-17 and 1j W-09 step 1 assume by-ID calls through a held volume handle reach the object on *that* volume. This experiment tries to break that assumption.

**SP-17 — Method (adversarial, in order of increasing severity).**

| Step | Action | Observation |
|---|---|---|
| 1 | Open volume handle `Hv` on scratch volume `V1`, access mask recorded | success/failure + `GetLastError` |
| 2 | Open fixture `T1` by ID through `Hv`; read `FileIdInfo`, creation time, child count | baseline identity |
| 3 | Remove `V1`'s mount point while `Hv` is held | does the call succeed? is `Hv` still valid? |
| 4 | Attach `CLONE` (`V2`, identical GUID/serial) | what identity does the mount manager report? (`identity_ambiguous`, or a renamed GUID — E-26) |
| 5 | Open `T1`'s stored ID through `Hv` | which object answers: `V1`'s, `V2`'s, or an error? |
| 6 | Repeat 3–5 with the mount point reassigned from `V1` to `V2` | does the by-ID path follow the locator or the handle? |
| 7 | Repeat 1–2, then detach `V1` entirely (not just unmount) | is `Hv` invalidated? with what error? |

**SP-18 — EV-27a, pass.** In every adversarial step, calls through `Hv` either reach the object on `V1` or fail with an error — they never silently return `V2`'s object. Then a held handle is a continuity mechanism and 1j W-09 step 1 holds.

**SP-19 — EV-27b, fail.** Any step where a by-ID call through `Hv` returns an object from `V2`, or returns success with a mismatched `FileIdInfo`/creation time, falsifies the design assumption. Recorded as `E-27 = fail`; consequence per SP-10.

**SP-20 — EV-27c, invalidated handle.** If holding `Hv` causes the unmount/detach to fail because the handle prevents it, that is **not** a failure — it is the platform enforcing continuity, and it is recorded as such (with the exact error). The distinction between "prevented" and "redirected" is the whole point.

**SP-21 — Failure handling.** Any unexpected destructive effect outside the scratch volume stops the run immediately: detach, cleanup, report. An unexplained error is recorded with its `GetLastError`, not retried with a different access mask until something works.

---

## 5. E-28: can the sequence be built?

**SP-22 — What is being tested.** 1j W-09's four properties, end to end, with an explicit detector for path re-resolution at every step.

**SP-23 — Step observables.**

| Step | Required | Recorded |
|---|---|---|
| 1 open | `OpenFileById(Hv, T1_ID)` succeeds; `FileIdInfo` on the returned handle equals the stored identity | handle, `FileIdInfo`, creation time |
| 2 enumerate | children listed via a **handle-relative** directory query on that directory handle | enumeration class used, child count vs `entries_at_acquisition`, child IDs obtained |
| 3 remove | each child opened by ID through `Hv` and removed via `FILE_DISPOSITION_INFORMATION_EX` on that child handle | per-child result codes |
| 4 verify | `OpenFileById(Hv, T1_ID)` now fails, with the **control open** succeeding (1h U-02) | exact status, plus control result |

**SP-24 — EV-28a, pass.** All four steps complete on `T1`, `T2` and `T3`, with:
- no path string passed to any open, enumerate or remove call;
- the step-4 failure status identical across repeated runs, and distinguishable from the malformed-request status by the control open;
- child count falling to zero as observed by step 2's own mechanism.

**SP-25 — EV-28b, fail.** Any step that requires a path — including "open by ID, then enumerate by path", or "enumerate by handle, then remove by path" — fails E-28. This is 1j W-10: a mixed sequence is a different, unsupported design, not a weaker one.

**SP-26 — EV-28c, sequencing hazard.** If step 3 requires the parent directory to be removed before some children can be, or if a child cannot be opened by ID while its parent is open, record it: the four steps are then not performable in the specified order, which is a fail for the specified design even if a reordered variant works.

**SP-27 — The control open is part of the experiment, not a convenience.** Without it, step 4's failure cannot be distinguished from a malformed request (1h U-01). A run that skips the control is not evidence.

**SP-28 — Failure handling.** A child that cannot be removed leaves the fixture partially destroyed — acceptable, it is disposable. It is recorded, the run continues to the next fixture, and the residue is listed in the report. Nothing is retried in a loop.

---

## 6. Commands

**SP-29 — Elevated, separate approval, shown at the moment of use.**

```powershell
# create the VHD (elevated; approved separately, never assumed)
New-VHD -Path X:\spike\v1.vhdx -SizeBytes 256MB -Dynamic
Mount-VHD -Path X:\spike\v1.vhdx
Initialize-Disk -Number <N> -PartitionStyle MBR
New-Partition -DiskNumber <N> -UseMaximumSize -DriveLetter X
Format-Volume -DriveLetter X -FileSystem NTFS -NewFileSystemLabel SPIKE-V1 -Confirm:$false
```

**SP-30 — Non-elevated.**

```powershell
Get-Volume | Where-Object DriveLetter -in 'X','Y'      # must be empty before SP-29
Get-PSDrive -PSProvider FileSystem                      # confirm X: and Y: unmapped
```

**SP-31 — The tool**, built outside the repository and invoked per experiment:

```text
spike e27 --volume X: --fixture T1 --clone-volume Y: --log e27.log
spike e28 --volume X: --fixture T1 --log e28.log
```

**SP-32 — Cleanup, idempotent, run even after failure.**

```powershell
Dismount-VHD -Path X:\spike\v1.vhdx -ErrorAction SilentlyContinue
Dismount-VHD -Path X:\spike\v2.vhdx -ErrorAction SilentlyContinue
Get-Disk | Where-Object FriendlyName -like '*Virtual*' | Remove-Disk -Confirm:$false
Remove-Item X:\spike -Recurse -Force -ErrorAction SilentlyContinue
```

Cleanup reports what it removed and what it could not; a residue that survives cleanup is reported as a finding, not glossed over.

---

## 7. Cleanup and crash behaviour

**SP-33 — Residue inventory.** After every run, the report lists: attached VHDs, mounted volumes, drive letters still assigned, files remaining in the scratch directory, and whether any fixture was left partially deleted.

**SP-34 — Crash mid-run.** If the tool dies with the volume mounted, the recovery is: run SP-32, verify no `X:`/`Y:` in `Get-Volume`, verify no virtual disks in `Get-Disk`. An orphaned mount is reported, not silently cleaned in a way that hides it.

**SP-35 — No automatic cleanup on unexpected error.** If a call fails in a way that suggests the *wrong* volume was reached (E-27's dangerous case), the spike **stops without deleting anything further** and reports. Cleanup of the fixture waits until the report is written, so the evidence survives.

---

## 8. Deliverables of the spike

**SP-36 — Two independent write-ups**, one per experiment, each carrying: environment string, exact commands run, raw logs, the fixture identity baselines, per-step results, `EV-` verdict, residue inventory, and anything unexpected.

**SP-37 — Recorded as evidence, not as contract change.** Results are appended as evidence findings against E-27/E-28. They do not edit rows 528–539 and do not silently amend the contract. If a result *requires* a contract change, that is a new clarification with its own review — not an edit to this series.

**SP-38 — Bound.** Each experiment is one session, with a fixed fixture set. If it is not conclusive after the specified steps, it stops and reports "inconclusive" with the evidence it has. No exploratory extension into other volume types, other filesystems, or other API surfaces in the same run.

---

## 9. Explicitly not authorized by this document

- Executing any command in §6.
- Creating or attaching a VHD.
- Deleting anything, on the scratch volume or elsewhere.
- Adding the spike tool to the repository.
- Any Go change to `internal/`, `cmd/`, or product packages.
- Ledger creation, persistence migration, or deletion enablement in Howl Chat.
- Treating a passing spike as a contract amendment or an acceptance.

**SP-39 — Current state on filing.** E-27 unestablished. E-28 unestablished. E-25 unestablished. Deletion unavailable. TASK-03 in preparation, not accepted, not closed.
