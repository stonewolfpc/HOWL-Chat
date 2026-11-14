# USN Reader — Implementation Review (v6.9 baseline)

**Identity:** `MEMORY-09-USN-IMPL-1` · **Version:** 3 · **Date:** 2026-11-14
**Status:** **IMPLEMENTED, OFFLINE-TESTED — no volume accessed, no journal queried.**
**Authority:** `MEMORY-09 WO-ARCH-06` §1; `WO-USN-REVIEW-02` (error mapping); `WO-USN-REVIEW-03` (final semantic audit). Controlling specification: USN Reader **v6.9**.

**Change history**

| v | Change |
|---|---|
| 1 | Initial implementation |
| 2 | 1178/1179/1181 separated into three outcomes with diagnostic text; unreachable validation branch removed; 75 → 80 checks |
| 3 | 1178/1179 wording states only what the status says (no replacement-journal inference, no "no journal exists"); Win32 structures hoisted so their layout is asserted offline; `CreateFileW`/`DeviceIoControl` signatures declared (a default-`int` return would truncate a 64-bit handle); 80 → 85 checks |

---

## 1. What exists

| File | Role |
|---|---|
| `eval/harness/usn_reader.py` | Bounded reader: control codes, Win32 structures, status mapping, cursor traversal, record parsing/validation, V4 sequences, `SourceInfo`, identity, named outcomes |
| `eval/harness/usn_tests.py` | 85 offline synthetic checks; exits non-zero on any failure |

Neither file opens a volume. `volume_io()` (L1229) is the only volume-facing code path and it refuses unless passed a raw volume path, on Windows, **and** `allow_open=True`. No test calls it.

---

## 2. Invariant → specification mapping

### 2.0 Journal status errors — three conditions, three outcomes

| Status | WinError.h name | Documented meaning | Outcome |
|---:|---|---|---|
| 1178 (0x49A) | `ERROR_JOURNAL_DELETE_IN_PROGRESS` | "The volume change journal is being deleted." | `INCONCLUSIVE_JOURNAL_DELETE_IN_PROGRESS` |
| 1179 (0x49B) | `ERROR_JOURNAL_NOT_ACTIVE` | "The volume change journal is not active." | `INCONCLUSIVE_JOURNAL_INACTIVE` |
| 1181 (0x49D) | `ERROR_JOURNAL_ENTRY_DELETED` | "The journal entry has been deleted from the journal." | `INCONCLUSIVE_JOURNAL_ENTRY_DELETED` |

**Source:** "System Error Codes (1000-1299) (WinError.h)", Microsoft Learn —
<https://learn.microsoft.com/en-us/windows/win32/debug/system-error-codes--1000-1299->

**Wording, and what it must not imply.** The three descriptions state the status and its immediate consequence, nothing more:

| Status | Reader text | Explicitly not claimed |
|---|---|---|
| 1178 | "the volume change journal is being deleted; the current operation cannot proceed" | that a replacement journal exists, or that a later `UsnJournalID` will differ |
| 1179 | "the volume change journal is not active; the current operation cannot proceed" | that no journal exists, was never created, or has a particular history |
| 1181 | "the requested journal entry has been deleted from the journal" | that retention was lost *during this scan*, or anything about the journal's activity |

A regression check (`text:` in the suite) asserts that none of the three contains "replacement", "no journal exists", "guarantee" or "discarded". Implementation: `OUTCOME_DESCRIPTIONS` (L241), `_exception_outcome` (L1080), `_detail_for` (L1117). Selection is numeric-first (`winerror`, then `errno`), with a narrow text fallback; unknown text stays `INCONCLUSIVE_INVALID_QUERY_RESULT`.

Why the distinction is kept, in the code's own words: 1178 says the journal **is being deleted** and nothing about what follows; 1179 says the journal **is not active** and nothing about its history; 1181 says the **requested entry has been deleted**, a statement about the range, not about the journal.

### 2.1 Query parameters and structures (v6.9 §1)

| Invariant | Implementation |
|---|---|
| Complete documented-reason mask, reserved bits clear | `COMPLETE_REASON_MASK = 0x80FFFF77` (L50), `RESERVED_REASON_BITS = 0x7F000088` (L80) |
| Control codes are the documented CTL_CODE results | `FSCTL_QUERY_USN_JOURNAL` (L125), `FSCTL_READ_USN_JOURNAL` (L127), recomputable via `ctl_code()` (L130) |
| Structures match the documented layouts | `USN_JOURNAL_DATA_V0` (L135), `READ_USN_JOURNAL_DATA_V1` (L149) — sizes and offsets asserted offline |
| Version range 2–4 | `MIN/MAX_MAJOR_VERSION` |
| No completeness claim from a filtered read | `completeness` string in every report |

### 2.2 Records (v6.9 §1.3, §3)

Header sizes 60/76/64 (`HEADER_BYTES`, L107); `parse_records` (L473); `validate_record` (L595). Rejections: unsupported version, length below header or beyond buffer, unparseable extent array, record before the call's start USN, record at or above the returned cursor, no reason bit inside the mask, reason bits outside the documented set, non-zero `FILE_ID_128` high half on NTFS.

### 2.3 V4 range tracking (v6.9 §1, §2)

`collect_v4_sequence` (L736): consecutive same-FRN grouping over the **accumulated** record list, so a sequence spanning two read buffers is grouped correctly; exact recurrence `R[i+1] = R[i] − N[i+1]`; no underflow; extent-array and `ExtentSize` agreement; terminal zero; reconciliation as a parsing cross-check only; `BOUNDARY_TRUNCATED` when the last observed record is non-terminal; companion only after the final V4, and only if V3 with `CLOSE`; companion timestamp labelled a close-record timestamp.

### 2.4 SourceInfo (v6.9 §2)

`decode_source_info` (L677): raw value preserved, flags decoded, `DATA_MANAGEMENT` and `AUXILIARY_DATA` both reported as "a logged reason need not be an application-visible change", replication identified. `_collect_matches` (L1161) never drops a record for carrying a flag.

### 2.5 Identity (v6.9 §3)

`verify_volume_identity` (L646): V2 FRN direct, V3/V4 low 64 bits with zero high half required; 64-bit serial compared handle-to-handle; 32-bit serials a consistency check only.

### 2.6 Traversal and retention (v6.9 §1.2, §4)

Fixed boundary captured before the scan; cursor monotonic; over-read excluded (L~956); no-match withholds Claim R′; journal ID, `FirstUsn` and `NextUsn` compared across the scan (`_check_retention`, L1125); Claim R′ set at L998 only when traversal completed, records validated and retention stable.

### 2.7 Bounds (v6.9 §4)

| Limit | Value | Enforcement |
|---|---|---|
| Calls | 128 | before every read (L917) |
| Bytes (requested) | 128 MiB | before every read (L921) |
| Output | 4 MiB | `_collect_matches` summarises (L1183) |
| Buffer | 1 MiB | `Caps` (L270) rejects a buffer too small for cursor + header |
| Preflight | from `NextUsn − FirstUsn` | reads nothing |
| Completion guarantee | **never claimed** | `Preflight.completion_guaranteed` is always `False` |

### 2.8 Claim ceiling (v6.9 §5–§7)

Claim R′ is a three-way conjunction (L998). `claim_t_prime` is the fixed string "not established by this reader"; `absence_means` states that an absent match is not evidence of no modification; `completeness` limits the claim to the documented V2–V4 sequence under selection semantics; gaps carry "not interpreted" text.

---

## 3. Failure paths

17 distinct inconclusive outcomes, all named: journal delete-in-progress, not-active, entry-deleted, ID changed, retention advanced, `NextUsn` regressed, no-progress-with-records, advanced-without-records, cursor regressed, unsupported version, undocumented reason, malformed record, mask mismatch, three budget outcomes, invalid query result. None retries; none is reported as a negative result about the file.

---

## 4. Tests

**Command:** `python -m eval.harness.usn_tests` → **85/85 passed**, exit 0.

| Group | Checks | Covers |
|---|---:|---|
| Parsing | 10 | valid V2, truncated header, over-long length, zero/under-header length, unsupported version, short buffer, record before start, journal-state mapping (dict, positional, rejected) |
| API structures / status text | 4 | struct sizes 56/48 with field offsets, control codes recomputed from CTL_CODE, description wording contains no unsupported inference |
| Mask/validation | 6 | documented reason accepted, mask mismatch (unit), undocumented bits, reserved bits, non-zero high half, mask arithmetic |
| Identity/volume | 7 | V2/V3 FRN extraction, other-file exclusion, 64-bit match, mismatch not rescued by 32-bit agreement, 32-bit-only difference, `from_file_id` |
| Cursor | 8 | completion, multi-call advance, no-match, no-progress with records, advance without records, backwards cursor, over-read exclusion, empty range |
| Retention | 6 | id change, `FirstUsn` advance, `NextUsn` regression, completed-cursor override, growth tolerated, failed closing query |
| Aborts/budget | 17 | 1178/1179/1181 separately, status distinctness, code-over-text precedence, `winerror`-only exception, message-only mapping, mask mismatch **through the scan**, malformed, unsupported version, undocumented reason, call cap, byte cap, preflight refusal, no completion promise, buffer rejection |
| Output | 1 | cap summarises rather than growing |
| V4 sequences | 11 | single, multi, underflow, recurrence skip, non-terminal, extent-array mismatch, extent-size mismatch, foreign-FRN interleave, missing companion, timestamp labelling, undated V4 |
| SourceInfo | 6 | data-management, auxiliary, replication, raw preservation, combined flags, flagged record not discarded |
| Ceiling/gaps | 5 | T′ refusal, absence wording, gap reporting, adjacency, outcome set |
| No live access | 5 | volume refusal, non-volume path, `CreateFileW`/`DeviceIoControl` confined to `volume_io`, **declared Win32 signatures**, scan holds no handle |

**Regression:** `python -m eval.harness.validate` → 89/89 passed; `python docs/research/corpus/check_corpus.py` → PASS (0 fail, 0 warn, 0 unverifiable).

---

## 5. Final source-level audit (`WO-USN-REVIEW-03`)

PASS = enforced and exercised. DEFECT = not as specified. NOT TESTED = implemented but unexercised, so behaviour on a real volume is unknown.

| # | Question | Verdict | Location and evidence |
|---:|---|---|---|
| 1 | Documented API structures and control-code parameters match | **PASS (offline)** | L125–L149: `ctl_code(61, METHOD_BUFFERED) == 0x000900F4`, `ctl_code(46, METHOD_NEITHER) == 0x000900BB`; `sizeof(USN_JOURNAL_DATA_V0) == 56`; `sizeof(READ_USN_JOURNAL_DATA_V1) == 48` (44 bytes of fields padded to 8-byte alignment); offsets asserted field-by-field |
| 2 | 1 MiB buffer and partial-record handling safe | **PASS (fails safe), NOT TESTED live** | `parse_records` L473 rejects a record whose declared length exceeds the remaining bytes as `MALFORMED_RECORD`, which is INCONCLUSIVE. A short tail can therefore never be silently dropped or misparsed. Whether a real journal ever returns one is unobserved |
| 3 | V4 grouping correct across separate buffers | **PASS (source), NOT TESTED** | `_collect_sequences` receives the accumulated list, so grouping spans calls; `collect_v4_sequence` L736 groups by consecutive same-FRN V4 records regardless of which buffer they arrived in. No fixture crosses a buffer boundary |
| 4 | Every `INCONCLUSIVE` outcome meaningful and reachable | **PASS** | 17 outcomes; a dedicated test now drives `MASK_MISMATCH` through `scan` (a record whose only reason bits are undocumented), closing the previous unit-level-only gap |
| 5 | Requested-byte cap applied before every IOCTL | **PASS** | L917 and L921 both run before `read_fn`; `requested_total` increments by the requested buffer size, not the returned length |
| 6 | Output cap bounds memory or only reporting | **PASS against spec / LIMITATION** | L1183 summarises the **report**. `all_records` is not separately capped: bounded only by the 128 MiB requested-byte cap, which at a 60-byte minimum record could hold ~2.2 M records in memory. Spec §4 constrains reporting, not retention; the risk is real but requires a journal far larger than a default configuration |
| 7 | Volume-opening reachable without caller authorization | **PASS by parameter / LIMITATION** | L1229: `allow_open` is keyword-only and defaults to `False`; the raw-path and Windows checks precede it. The gate is **explicit authorisation by the caller**, not a capability boundary — any importer that passes the flag can open a volume. No import-time or scan-time handle exists |
| 8 | Diagnostic text exceeds Claim R′/T′ | **DEFECT → fixed** | 1181's text previously added "the range was discarded under the reader", an inference about retention during *our* scan that the status does not support. Removed; the wording check now forbids "discarded", "replacement", "no journal exists" and "guarantee" in all three descriptions |
| 9 | Win32 handle path is type-correct | **DEFECT → fixed, NOT TESTED live** | `CreateFileW` and `DeviceIoControl` were called without `restype`/`argtypes`. ctypes then assumes an `int` return, which **truncates a 64-bit HANDLE**, and passes the handle back as 32-bit — so a valid handle would be corrupted on first use, and `INVALID_HANDLE_VALUE` (-1) would never have matched `c_void_p(-1)`. Signatures are now declared and the failure check compares against `c_void_p(-1).value` |
| 10 | `METHOD_NEITHER` output-byte accounting | **NOT TESTED** | `FSCTL_READ_USN_JOURNAL` is `METHOD_NEITHER` (derived from the control code), so `lpBytesReturned` reliability is unverified here. The code slices the buffer by that value; if it were 0 the result is `MALFORMED_RECORD` (INCONCLUSIVE), not a false clean scan |
| 11 | Structures behave as documented on a live volume | **NOT TESTED** | Layout is asserted only against the documented C structure; no live ioctl has run |
| 12 | Real record layouts | **NOT TESTED** | Fixtures are documentation-derived |

---

## 6. Defects found

1. **Status mapping conflated two failures** (first pass). 1178 and 1179 both routed to entry-deleted; a deletion-in-progress disappeared entirely, and the test that "covered" it passed only because of the conflation.
2. **Byte budget counted the wrong quantity** — returned rather than requested bytes, so a slow journal could issue unbounded reads.
3. **Sequence detail computed and dropped** — `collect_v4_sequence` never assigned its local `detail`.
4. **Termination test ordered before truncation** — a single unterminated V4 was called `INVALID` instead of `BOUNDARY_TRUNCATED`.
5. **Unreachable reserved-bit branch** read as coverage; removed.
6. **Undocumented ctypes signatures** on the volume path (audit row 9): a latent handle-truncation bug in the one function that would run at the next gate.
7. **Diagnostic over-claim in 1181's text** (audit row 8); removed.

---

## 7. Untested branches and known limitations

| Item | State |
|---|---|
| `volume_io()` ioctl path | Never executed; no volume opened |
| Live structure behaviour and `lpBytesReturned` | Unverified (`METHOD_NEITHER`) |
| Real record layouts | Documentation-derived fixtures |
| In-memory retention | Not separately capped; bounded by the 128 MiB requested-byte cap |
| V4 sequence across buffers | Source correct; no fixture crosses a boundary |
| Partial trailing record | Rejected as malformed; unobserved live |
| Output-cap accounting | Deterministic size proxy, not JSON serialisation |
| Message-only status matching | Narrow phrases; a localised message falls through to the generic outcome |
| Claim T′ | Out of scope; no event-log access |
| Concurrency | Single-scan, single-threaded |

---

## 8. What this establishes

**Establishes:** the reader implements v6.9 as written; structures and control codes match the documented layouts and arithmetic (offline); the three journal statuses are distinct in outcome and wording; all 17 outcomes are reachable; bounds are enforced by the component; and no code path opens a volume without explicit caller authorisation.

**Does not establish:** behaviour against a real journal, that `FSCTL_READ_USN_JOURNAL` behaves as assumed on this machine, that any claim about S1-B is true, or that the diagnostic will be conclusive.

Synthetic-test success is not live-volume validation. The three gates remain closed: `G:` disposable-file validation, event-log enumeration, and the bounded S1-B journal query.
