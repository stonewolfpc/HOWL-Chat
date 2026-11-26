# TASK-03 — spike findings 01a: correction to F-2/F-3, and a volume-collision result

| Field | Value |
|---|---|
| Purpose | Correct the absence-discriminator interpretation in findings 01, and record two new measurements |
| Status | **Evidence record. Supersedes findings 01 §2 F-3 and the F-2/Run-1 interpretation. No contract row edited** |
| Relationship | Amends `TASK-03-SPIKE-FINDINGS-01-2026-11-22.md` (row 541) by addition only; row 541 unchanged |
| Artifacts | `build/tmp/spike2/probe2.go`, `build/tmp/spike2-ntstatus.log`, `build/tmp/spike2b-ntstatus.log` |
| Environment | Windows 11 build 26200 (25H2); hints on `C:` (NTFS, serial `0x885A56435A562DE6`) and `D:` (NTFS, serial `0x47661BC7661AEDC`) |
| Recovery token | `TASK-03-SPIKE-FINDINGS-01A-2026-11-22` |

Labels: **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished**.

---

## 1. What was run

| Run | Scope | Destructive |
|---|---|---|
| A | `NtCreateFile(FILE_OPEN_BY_FILE_ID)` raw NTSTATUS for seven identity classes | no (plus one fixture create/delete) |
| B | same, extended to report **which object was actually reached** on success | no (plus one fixture create/delete) |

Both created and deleted a single probe-owned file under `build/tmp/spike2/fixture/`; residue none in both. Run B additionally opened `D:\` **read attributes only** to obtain a foreign-volume identity; no write, no delete, no enumeration of its contents.

---

## 2. Retraction: F-3 was wrong, and F-2 was overstated

**X-01 — What findings 01 claimed.** That a successful control open plus a valid descriptor plus a target failure establishes absence, and that U-02 clause 5 was merely vacuous.

**X-02 — Why that does not hold.** The control proves the mechanism and the request shape work for a live identity on the hinted volume. It says nothing about the *cause* of the target's failure. Findings 01 already contained the counter-example without drawing the conclusion from it: the same Win32 code 87 appeared for a deleted identity and for a malformed request.

**X-03 — The new measurement settles it at the lowest layer.** Win32 collapses statuses; `NtCreateFile` does not. The raw statuses:

| Case | NTSTATUS | Win32 |
|---|---|---|
| live identity (control) | `0x00000000` SUCCESS | 0 |
| **deleted identity** (was live, then removed) | **`0xC000000D` INVALID_PARAMETER** | **87** |
| fabricated all-zero identity | **`0xC000000D` INVALID_PARAMETER** | **87** |
| **live id, name length 0 (malformed)** | **`0xC000000D` INVALID_PARAMETER** | **87** |
| fabricated all-ones identity | `0xC0000034` OBJECT_NAME_NOT_FOUND | 2 |
| live id, name length 3 (malformed) | `0xC0000033` OBJECT_NAME_INVALID | 123 |
| live id, name length 8 (64-bit form) | `0x00000000` SUCCESS | 0 |

**X-04 — The conclusion, stated plainly.** A deleted identity is **not distinguishable** from a structurally invalid request at this layer: both return `STATUS_INVALID_PARAMETER`. A different malformed shape returns a different status, so the codes do not partition cleanly along "identity missing" versus "request broken" either. **No documented API result distinguishes them.**

Consequences:
- Control success is **necessary** and **insufficient**.
- Descriptor validation is **necessary** and **insufficient**.
- `0xC0000034`/`ERROR_FILE_NOT_FOUND` appears only for out-of-range identities and cannot be adopted as a general absence indicator.
- **Absence by open-by-ID is `unverifiable`.** Not "weakly established" — unverifiable.

**X-05 — What this does to `complete`.** 1j W-09 step 4 verifies absence "by ID open on the same held volume handle". The evidence says that step cannot produce a determinate result. So `complete` is currently **unreachable by the specified mechanism**, and 1h's fail-closed posture is the only correct one: `unverifiable` / `completion_unconfirmed` carry the space.

**X-06 — Recorded, not applied.** This is a finding against 1h U-02 and 1j W-09, not an edit to them. A future clarification is required, with its own review, before either is rewritten. Rows 528–541 are untouched.

---

## 3. New finding: the identity is volume-scoped, and the ID alone can resolve a foreign object

**X-07 — The measurement.** `D:\` root has identity `0500…0500` (record 5, sequence 5) on volume serial `0x47661BC7661AEDC`.

| Request | Hint volume | Result | Object reached |
|---|---|---|---|
| `D:` identity | `C:` | **SUCCESS** | id `0500…0500`, **serial `0x885A56435A562DE6` (`C:`)** |
| `D:` identity | `D:` | SUCCESS | id `0500…0500`, serial `0x47661BC7661AEDC` (`D:`) |
| `C:` identity | `D:` | `0xC000000D` | — |

**X-08 — Interpretation (Inferred, high confidence).** The hint scopes resolution to its own volume; it does not import the foreign identity. `C:`'s MFT record 5 also exists with sequence 5, so the **same 128-bit ID** denotes a *different object* there, and the open succeeded against that unrelated object.

**X-09 — Why this matters more than it first appears.**

- An ID is unique **per volume**, not globally. `(ID, volume serial)` identifies an object; the ID alone does not.
- **A present-side observation can therefore be a silent false positive.** Opening "our" ID through the wrong volume's hint returns success for a completely unrelated object. Nothing in the status says so.
- The only defence is comparing **`FILE_ID_INFO`'s 64-bit volume serial and ID from the returned handle** against the acquired pair. This elevates findings 01's F-6 from an implementation hazard to a correctness requirement.
- For `identity_not_found`-style reasoning it changes nothing directly, but it removes the last comforting assumption: not only is absence unverifiable, presence is not self-authenticating either.

**X-10 — Effect on 1h U-22 / E-27.** A "volume mismatch" is not hypothetical or clone-specific. It happens with two ordinary volumes whenever their record numbering coincides. Volume identity must be checked on **every** returned handle, not at acquisition only.

---

## 4. What survives, qualified

| Item | Status after this correction |
|---|---|
| E-28 step 1 (open by ID through a directory handle, unelevated) | **established for this configuration**, now conditional on verifying the returned object's ID **and** 64-bit serial |
| E-25 | **unverifiable by open-by-ID**; a different positive mechanism is required (below) |
| E-27, E-26 | unestablished, unchanged |
| 1h U-11 (no attribution) | stands, unchanged |
| `complete` | unreachable by the specified mechanism; deletion remains disabled |
| Immediate MFT slot reuse (findings 01 F-4) | stands; one observation, one filesystem, one build |
| 65,536-wrap calculation | stands as arithmetic, not as an established behaviour of every NTFS version |

**X-11 — Candidate positive mechanisms for absence (all Unestablished).**

| # | Mechanism | Why it might discriminate | Why it might not |
|---|---|---|---|
| P1 | Handle-relative enumeration of the **parent**: the acquired ID absent from the listing | enumeration is a different code path from open-by-ID | enumeration can be incomplete; the parent may itself be gone; a name can be recreated though the ID cannot |
| P2 | USN journal delete record for the reference (E-24) | positive evidence of a delete event | journal may be disabled or truncated |
| P3 | Holding an open handle across the removal and observing it close | direct causal observation | POSIX delete keeps the handle alive; contradicts the ordering |

None is established. P1 is the cheapest to test next and is the one the reviewer's recommendation points at, since it is independent of open-by-ID entirely.

---

## 5. Verification and authorization

Rows 541 and earlier remain byte-identical. This record is registered as row 542. Residue: none in both runs (`build/tmp/spike2/fixture` absent). No `E:` access. No Howl data touched. `D:` opened read-attributes only.

**No contract row edited. No deletion enabled. No attribution restored.** E-28 steps 2–4 remain the next investigative step, and they are independent of this correction — a handle-relative enumeration and removal sequence can be proven feasible even while the final absence method is unresolved.
