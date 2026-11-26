# TASK-03 — spike findings 05: USN evidence requires elevation

| Field | Value |
|---|---|
| Purpose | P2 feasibility: can USN journal evidence positively establish deletion of the exact acquired identity? |
| Status | **Evidence record. No contract row edited. Deletion disabled** |
| Result | **Blocked before the question could be asked**: no volume handle is obtainable without elevation on this build |
| Scope excluded | No journal was created, deleted or modified. E-27/E-26 not attempted. Share-mode anomaly deferred (§5) |
| Artifacts | `build/tmp/spike7/probe7.go`, `build/tmp/spike7-usn-access.log` |
| Environment | Windows 11 build 26200 (25H2), NTFS `C:` |
| Recovery token | `TASK-03-SPIKE-FINDINGS-05-2026-11-22` |

Labels: **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished**.

---

## 1. The blocking measurement

`FSCTL_QUERY_USN_JOURNAL` and `FSCTL_READ_USN_JOURNAL` are volume operations, so they need a volume handle. Three access masks were attempted on `\\.\C:`:

| Mask | Result |
|---|---|
| `FILE_READ_ATTRIBUTES \| FILE_READ_DATA` | **denied, 5 `ERROR_ACCESS_DENIED`** |
| `FILE_READ_DATA` | **denied, 5** |
| `GENERIC_READ` | **denied, 5** |

**F-29 (Observed) — no volume handle is obtainable without elevation.** Not "the journal query failed" — the handle itself was refused, so the FSCTL was never reached. This is consistent with the earlier observation that `FILE_READ_ATTRIBUTES` alone *was* granted: attributes are readable, volume data is not.

**F-30 — Consequence for P2 (Calculated).** USN journal evidence is unavailable to an unelevated process. If Howl were to verify deletion through USN records, it would have to **run elevated**. That is a product-level consequence, not a contract detail, and it is large: the application would require administrator rights merely to confirm that a conversation directory it already deleted is gone.

**F-31 — The wrap/gap question cannot be tested here either (Inferred).** Forcing journal wrap or testing recovery from a gap requires creating or deleting the journal (`FSCTL_CREATE_USN_JOURNAL` / `FSCTL_DELETE_USN_JOURNAL`). Those are **volume metadata** operations, not disposable test data, and the present authorization does not cover them. They would be safe only on a throwaway volume — which is the same VHD that E-27 needs.

**F-32 — Both remaining items converge on one prerequisite.** E-27 needs an elevated, disposable volume to test detach/clone behaviour. P2 needs an elevated volume to read the journal, and a disposable one to manipulate it safely. **Neither can progress in an unelevated session**, and one approved elevated session on a scratch VHD would serve both.

---

## 2. What USN evidence would have to show (specified, unestablished)

The question stays narrow: *can a USN record positively establish that the exact acquired directory identity was deleted?* On the documented shape of the API, a positive answer needs **all** of the following, and none of them has been measured:

| # | Requirement | Why it is needed | Status |
|---|---|---|---|
| 1 | Journal **ID** captured at acquisition and unchanged at verification | a recreated journal has a new ID; records from the old one cannot be trusted | Unestablished |
| 2 | A **start USN** captured at acquisition | the window must be bounded in both directions | Unestablished |
| 3 | `LowestValidUsn ≤ startUsn` at verification | otherwise records in the window were discarded and absence of a record means nothing | Unestablished |
| 4 | A `USN_REASON_FILE_DELETE` record whose `FileReferenceNumber` **equals the low 64 bits** of the acquired 128-bit ID | the record must match *this* identity, not merely the parent or the name | Unestablished |
| 5 | Continuity across a restart | the journal is on-disk and survives restarts, but this has not been observed here | Unestablished |

**F-33 — The asymmetry that must be written into any design. A missing record proves nothing.** A journal that is disabled, wrapped, recreated, or read with the wrong start USN all produce "no matching record" for a directory that really was deleted. Only a **present, matching record** is positive evidence, and only when requirements 1–3 hold. Absence of evidence must map to `unverifiable`, never to `complete`.

**F-34 — Even a passing P2 would not restore attribution.** A delete record shows that *something removed that identity* between two USNs. It does not show that this operation did it, which is 1h U-11's separate problem. P2 addresses absence, not authorship — the two must not be conflated when it is eventually measured.

---

## 3. Effect on the current state

| Item | Status after this round |
|---|---|
| Step 4 (positive absence) | **still unresolved** — P2 not measured, P1 refused |
| `complete` | **unreachable** |
| E-28 steps 1–3 | unchanged; PASS in tested configurations |
| E-25 | `unverifiable` by open-by-ID (row 542) |
| E-27 / E-26 | unestablished; now the shared prerequisite for P2 as well |
| G-1 | intact; deletion disabled |

No contract row was edited, and none should be until step 4 has a mechanism that either works or is shown to be impossible.

---

## 4. What would make P2 decidable

One approved elevated session on a disposable scratch VHD, in which:

1. the journal is queried for ID and USN ranges (requirement 1–3);
2. a fixture directory is deleted and the resulting records are read and matched by file reference (requirement 4);
3. the fixture is deleted **after** the volume is detached and re-attached, to test restart continuity (requirement 5);
4. the journal is deliberately recreated, to confirm that a changed ID is detectable and that an old start USN is then refused rather than silently honoured.

Step 4 is the one the contract actually depends on: whether a **stale** start USN is rejected, or silently yields "no records" — which would turn a gap into a false absence. Until that is observed, P2 cannot be recommended even as a fallback.

---

## 5. Deferred, not dropped: the share-mode anomaly

Findings 04 F-28 recorded that a handle whose share mode excluded `FILE_SHARE_DELETE` did not prevent disposition of that file. It is **not** dismissed, and it is **not** investigated here — the instruction was P2 first. When it is taken up, it must distinguish:

| Distinction | Why |
|---|---|
| same-process vs **separate-process** holder | a second handle in the same process is not a substitute for another holder |
| exact access and share flags on **both** opens | the result may hinge on the disposition open's own share mode |
| disposition API and flags (`FileDispositionInfoEx` vs `FileDispositionInfo`, POSIX vs not) | different paths may enforce sharing differently |
| whether the holder actually remained open through the attempt | a closed handle proves nothing |
| the **exact returned status**, not a boolean | 87 and 145 are different findings and were reported as "failed" in one place |

Until then, foreign-handle interference remains **unestablished**.

---

## 6. Procedure and authorization

Read-only. No journal created, deleted or modified. No `E:` access, no elevation requested or used, no Howl data touched, nothing deleted. No residue to clean.

**No contract row edited. No deletion enabled. No attribution restored.** Rows 540–545 unchanged; this record is row 546. TASK-03 remains in preparation.
