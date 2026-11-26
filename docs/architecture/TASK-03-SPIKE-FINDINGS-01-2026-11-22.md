# TASK-03 — spike findings 01: open-by-ID, absence codes, and record reuse

| Field | Value |
|---|---|
| Purpose | Record the first executed evidence under spike spec row 540 |
| Status | **Evidence record. No contract row edited. No deletion enabled** |
| Scope executed | Non-elevated, read-only probing plus a self-isolated delete of probe-created fixture files |
| Not executed | E-27 entirely; E-28 steps 2–4; anything requiring an elevated VHD |
| Environment | Windows 11 build **26200** (25H2), NTFS volume `C:`, serial `0x5A562DE6` |
| Artifacts | `build/tmp/spike/probe.go` (throwaway), `build/tmp/spike-readonly.log`, `build/tmp/spike-e25-delete.log` |
| Recovery token | `TASK-03-SPIKE-FINDINGS-01-2026-11-22` |

Labels: **Observed** (measured this run) · **Source-derived** · **Calculated** · **Inferred** · **Unestablished**.

---

## 1. What was actually run

| Run | Command | Destructive? |
|---|---|---|
| 1 | `go run ./build/tmp/spike <project dir>` | no |
| 2 | `go run ./build/tmp/spike <project dir> --delete-test` | **yes** — deleted 3 files it created |

Run 2's deletions were confined to `build/tmp/spike/fixture/`, a directory created by the probe, containing three files written by the probe. Cleanup removed the two surviving files and the directory. **Residue: none** (verified). No Howl data, no conversation directory, no `E:` access.

`go build ./...` passes with the probe present; it is `package main` under `build/tmp` and is not product code.

---

## 2. Findings

### F-1 — The hint handle need not be a volume handle, and needs no elevation (**Observed**)

A directory handle opened with `FILE_READ_ATTRIBUTES | SYNCHRONIZE`, `FILE_SHARE_READ|WRITE|DELETE`, and `FILE_FLAG_BACKUP_SEMANTICS` works as `hVolumeHint`. The control open succeeded with both descriptor types.

**Consequence:** E-28 step 1 is feasible in an unelevated process. The identity-bound design does **not** require administrator rights merely to hold a usable hint. This materially improves feasibility: 1j W-09 step 1's "held volume handle" is satisfiable by a held *directory* handle, which is a narrower and cheaper requirement.

`CreateFile(\\.\C:, FILE_READ_ATTRIBUTES)` also succeeded unelevated (informational only; `FILE_READ_DATA` on a volume was not attempted and is expected to need elevation — **Unestablished**).

### F-2 — A deleted identity returns `ERROR_INVALID_PARAMETER`, not "file not found" (**Observed**)

| Identity opened | Descriptor | Result |
|---|---|---|
| live (control) | FileId(64) | success |
| live (control) | ExtendedFileId(128) | success |
| **deleted** (stored ID, after `DeleteFile`) | FileId(64) | **87 `ERROR_INVALID_PARAMETER`** |
| **deleted** | ExtendedFileId(128) | **87** |
| fabricated, record beyond MFT | both | 87 |
| fabricated, all-zero | both | 87 |
| fabricated, all-ones | ExtendedFileId(128) | **2 `ERROR_FILE_NOT_FOUND`** |
| live ID, malformed descriptor (`dwSize = 8`) | — | **87** |

Two consequences, both material:

1. **Absence has no single status.** 87 and 2 both appeared; which one depends on the ID's shape, not on the reason for failure. Any rule expecting one documented "not found" code is unimplementable.
2. **87 is also what a malformed request returns.** A failed open therefore cannot be interpreted at all without a successful control open of the same shape. 1h U-02's control requirement is not theoretical caution — it is the only discriminator, and it works.

### F-3 — U-02 clause 5 is satisfied vacuously, and should say so (**Inferred from F-2**)

1h U-02 clause 5 requires the target to "fail with a status that the control demonstrably cannot produce". The control succeeds, so *any* failure satisfies this — including a failure caused by a malformed descriptor. The clause has no discriminating power; the real preconditions are clauses 2–4 (well-formed descriptor, adequate access, working control).

**This is a contract observation, not an edit.** Recorded for a future clarification; rows 528–540 are untouched. The practical rule the evidence supports is: *control succeeds with the exact call shape, target fails → identity did not resolve*, with 87 and 2 both admissible.

### F-4 — Record reuse is immediate; identical-ID reuse is not (**Observed**)

| Object | MFT record index | Sequence |
|---|---|---|
| `a.txt` (created) | 769116 | 19 |
| `b.txt` (created) | 769118 | 16 |
| `a.txt` (deleted) | — | — |
| `c.txt` (next create) | **769116** | **20** |

The next file created after a deletion took the **same record index**, with the sequence advanced by one. Consequently:

- 1h's sequence analysis holds: the stored ID stops resolving after reuse, because the sequence differs.
- The ABA scenario's precondition — a *different* object answering to the **identical 128-bit ID** — requires the 16-bit sequence to wrap, i.e. **65,536 frees of that one record** (**Calculated**).
- Therefore slot reuse is the normal case, while identical-ID collision is rare. 1g's framing risked implying the opposite; 1h's ABA concession remains correct in structure but is quantitatively narrow.
- **`reuse_excluded` still needs positive evidence**, so 1h U-11 (no Phase 1 attribution) stands unchanged. This measurement does not restore attribution, and re-opening it would be a new owner decision.

### F-5 — The project directory's record has already been reused 14 times (**Observed**)

Its identity is `53D9 0900 0000 0F00 …`: index 645459, **sequence 15**. Sequence *n* means the record has been freed and reallocated *n−1* times. Reuse is ordinary, not exceptional.

### F-6 — Two different volume serials, and a naive comparison would fail (**Observed**)

| Source | Value |
|---|---|
| `GetVolumeInformation` | `0x5A562DE6` (32-bit) |
| `FILE_ID_INFO.VolumeSerialNumber` | `0x885A56435A562DE6` (64-bit) |

The low 32 bits match; the high 32 differ. **An implementation that compares the 32-bit `GetVolumeInformation` serial against the 64-bit `FILE_ID_INFO` serial will never match** and would report every identity as volume-mismatched. Identity comparison must use the 64-bit value from `FILE_ID_INFO` on both sides.

### F-7 — NTFS 128-bit layout confirmed (**Observed**)

Low 8 bytes = 48-bit record index + 16-bit sequence; high 8 bytes = zero. Matches the MS-FSA note already cited in 1h.

### F-8 — One probe bug caught before it became a finding (**Observed**)

The first run reported `FILE_SUPPORTS_OPEN_BY_FILE_ID = false`. That was wrong: the probe used `0x00080000`, which is `FILE_READ_ONLY_VOLUME`; the correct flag is **`0x01000000`**. Re-run: the flag is **set** (`flags 0x03E72EFF`). Recorded because a false "unsupported" verdict would have justified disabling deletion on NTFS entirely.

---

## 3. Verdicts

| Item | Verdict | Basis |
|---|---|---|
| **E-25** | **partially established** | Status set measured on one build and volume: 87 (absence), 2 (out-of-range ID), 87 (malformed). Not a general law; other filesystems unmeasured |
| **E-28 step 1** | **established, unelevated** | open by ID through a directory handle, identity confirmed by `FILE_ID_INFO` |
| **E-28 steps 2–4** | **unestablished** | enumeration and removal through derived handles not attempted |
| **E-27** | **unestablished** | requires detach/clone; needs an elevated VHD, approval not yet granted |
| **E-26** | **unestablished** | duplicate-GUID behaviour not exercised |

**Deletion in Howl Chat remains disabled.** No acceptance gate has passed; two of three experiments have not been attempted.

---

## 4. Environment facts for the next run

| Item | Value |
|---|---|
| OS | Windows 11, build 26200, 25H2 |
| `New-VHD` / `Mount-VHD` / `Dismount-VHD` | **absent** (Hyper-V module not installed) |
| `Mount-DiskImage`, `Initialize-Disk`, `New-Partition`, `Format-Volume` | present (Storage module) |
| Consequence | a scratch VHD needs `diskpart create vdisk` (elevated) — or Hyper-V cmdlets installed |
| Free letters | X, Y, V, W free; `E:` absent from the volume list and not probed |

---

## 5. What this does not license

- No contract row was edited. F-3 and F-4 are recorded for review, not applied.
- No attribution was restored (1h U-11 stands).
- No path-based fallback was adopted.
- Deleting anything in Howl Chat is still unauthorized and still disabled.

**Next, in order:** (1) E-28 steps 2–4 on a disposable fixture, unelevated, on a directory the probe creates; (2) E-27, which is the only part needing an elevated VHD — the exact command will be shown and approval requested before it runs.
