# `USN-READER-FINAL-SOURCE-AUDIT`

**Identity:** `MEMORY-09-USN-FINAL-1` · **Date:** 2026-11-14 · **Track:** `WO-NEXT-01` Track 3
**Subject:** `eval/harness/usn_reader.py` at sha256 `2b384069…` (1,356 lines), `eval/harness/usn_tests.py` at sha256 `f6152707…` (1,040 lines)
**Status:** source-only audit. **No volume opened. No journal queried. No ioctl executed.**
**Verdict summary:** **8 PASS · 1 FAIL · 3 NOT TESTED · 3 LIMITATION.** The FAIL is a handle-lifetime defect, reported and **not** silently repaired.

A passing offline suite is not evidence of live Windows API correctness. Every PASS below is a statement about source, not about behaviour on this machine.

---

## 1. Audited items

| # | Item | Verdict |
|---:|---|---|
| 1 | `CreateFileW` / `DeviceIoControl` signatures, pointer widths, argument and return types | **PASS** |
| 2 | `INVALID_HANDLE_VALUE` handling | **PASS** |
| 3 | Lifetime of buffers passed to `DeviceIoControl` | **PASS** |
| 4 | `METHOD_NEITHER` output-byte accounting | **NOT TESTED** |
| 5 | Can a zero / unreliable returned-byte count produce a false clean result? | **PASS** |
| 6 | 1 MiB buffer vs record lengths vs partial trailing records | **PASS** + `L-1` |
| 7 | V4 grouping across multiple returned buffers | **PASS (source) · NOT TESTED (no fixture)** |
| 8 | 128 MiB requested-byte cap before every read | **PASS** |
| 9 | Report cap vs retained-record memory | **PASS against spec** + `L-2` |
| 10 | Every diagnostic outcome within Claim R′/T′ evidence limits | **PASS** |
| A | Handle lifetime and release | **FAIL** (`F-1`) |
| B | Journal-id fallback after a failed initial query | **LIMITATION** (`L-3`) |
| C | Volume opening gated only by a caller parameter | **LIMITATION** (`L-4`) |

---

## 2. Item-by-item evidence

### 1 — Signatures, widths, argument types · PASS

`volume_io` (L1229) declares both signatures before use:

| Function | `restype` | Declared arguments |
|---|---|---|
| `CreateFileW` | `wintypes.HANDLE` | `LPCWSTR, DWORD, DWORD, LPVOID, DWORD, DWORD, HANDLE` |
| `DeviceIoControl` | `wintypes.BOOL` | `HANDLE, DWORD, LPVOID, DWORD, LPVOID, DWORD, POINTER(DWORD), LPVOID` |

- Handle returned as `HANDLE` (pointer-width), not the ctypes default `int`. A default `int` return truncates a 64-bit handle — this was the previous defect and is closed.
- The input request is passed as `ctypes.byref(request)` with `ctypes.sizeof(request)`; `READ_USN_JOURNAL_DATA_V1` is 48 B, `USN_JOURNAL_DATA_V0` is 56 B, both asserted offline against field offsets.
- The output buffer is passed as `ctypes.cast(buffer, ctypes.c_void_p)` with an explicit `buffer_bytes` length — correct for an untyped `LPVOID`.
- `lpBytesReturned` is `POINTER(wintypes.DWORD)` over a live `wintypes.DWORD`.
- `lpOverlapped` is `None`; the handle was opened **without** `FILE_FLAG_OVERLAPPED`, so a synchronous call is consistent with the open flags. **PASS.**

### 2 — `INVALID_HANDLE_VALUE` · PASS

```python
invalid_handle = ctypes.c_void_p(-1).value      # 18446744073709551615
if handle is None or handle in (0, invalid_handle):
    raise OSError(ctypes.get_last_error(), ...)
```

Comparing against `c_void_p(-1).value` is the correct form: `CreateFileW` returns `INVALID_HANDLE_VALUE` (`(HANDLE)-1`), not `-1` as a signed `int`. `NULL` and `None` are also rejected. `use_last_error=True` is set on the `WinDLL`, so `get_last_error()` reports the call's own error rather than a stale one.

### 3 — Buffer lifetime · PASS

`read` creates `buffer = ctypes.create_string_buffer(buffer_bytes)` as a local, passes it to the synchronous `DeviceIoControl`, and only then does `return buffer.raw[: returned.value]` — which **copies** into a `bytes` object. No pointer outlives the call; there is no window in which the buffer could be freed while the kernel still writes to it. The `USN_JOURNAL_DATA_V0` output is likewise a local whose `byref` is valid for the call's duration. The request struct is a local too, so `byref(request)` cannot dangle.

### 4 — `METHOD_NEITHER` output accounting · NOT TESTED

`FSCTL_READ_USN_JOURNAL` is `ctl_code(46, METHOD_NEITHER)` = `0x000900BB`, so the kernel writes directly into the user buffer rather than through the system buffer, and `lpBytesReturned` semantics are the ones this code depends on but has never observed. No live call has been made; no fixture can exercise the kernel's behaviour.

**What is established:** the code uses the returned count only as a slice bound, and the slice is applied to a buffer of known length.

### 5 — Can a zero/unreliable count produce a false clean result? · PASS

| Returned value | Effect | Outcome |
|---|---|---|
| `0` | `buffer.raw[:0]` → `b""` → `parse_records` sees a buffer shorter than the 8-byte cursor | `I_MALFORMED_RECORD` → **INCONCLUSIVE** |
| Garbage larger than the buffer | slice clips at `len(buffer)`; parsing then fails at the first impossible record length | **INCONCLUSIVE** |
| A short but plausible buffer | parsed; a truncated record at the tail is rejected by `length > remaining` | **INCONCLUSIVE** |
| A correct short buffer holding only the cursor | cursor advanced, no records | `ADVANCED_WITHOUT_RECORDS` or `NO_MATCH_AT_CURSOR`, **never** a clean claim |

**There is no path from a bad byte count to `COMPLETED`.** That is the property that matters, and it holds structurally rather than by convention.

### 6 — Buffer size vs record size · PASS, with `L-1`

`parse_records` (L473) consumes the 8-byte cursor, then requires `length >= HEADER_BYTES[version]` and `length <= remaining`, and requires `remaining >= 8` before reading a header. A tail that cannot hold a header is reported, not skipped. This is the correct conservative behaviour for a trailing partial record.

**`L-1` — the buffer size assumption is not enforced.** The docstring reasons that "the buffer is sized to hold any single record, so a short tail is a defect". `Caps` rejects a buffer too small for *cursor plus header*, but nothing asserts that the buffer can hold the **largest possible record** — a V4 record carries an extent array whose size grows with the number of extents, and the documented maximum record size is not checked anywhere. A legitimate record larger than 1 MiB would therefore be classified `MALFORMED_RECORD`: inconclusive, but mislabelled, and it would look identical to a truncated tail. Fails safe; not precise.

### 7 — V4 grouping across buffers · PASS (source) · NOT TESTED (fixture)

`_collect_sequences` receives the accumulated record list, so a sequence whose records arrived in two different read buffers is still grouped by consecutive same-FRN V4 records inside `collect_v4_sequence` (L736). The recurrence `R[i+1] = R[i] − N[i+1]`, the non-terminal detection, the extent-array agreement and the companion rule (only after the final V4, only if V3 with `CLOSE`) are each covered by an offline check.

**Not tested:** no fixture places the boundary *between* two records of the same sequence, so the cross-buffer path is reasoned, not exercised. It is a one-fixture gap and it is stated rather than papered over.

### 8 — 128 MiB requested-byte cap · PASS

Both budget checks precede the read call:

| Check | Location | Counts |
|---|---|---|
| call cap (128) | L917, before `read_fn` | one call |
| byte cap (128 MiB) | L921, before `read_fn` | the **requested** buffer size, not the returned length |

Requested, not returned, is the correct quantity: a journal returning almost nothing per call must not be able to issue unbounded calls under a byte budget that never increments.

### 9 — Report cap vs retained memory · PASS against spec, with `L-2`

`_collect_matches` (L1183) bounds the **report**. **`L-2`:** `all_records` is not separately capped; it is bounded only by the 128 MiB requested-byte total, which at a 60-byte minimum record is on the order of 2.2 million records held in memory. v6.9 §4 constrains reporting, not retention, so this is conformant — and it is a real memory ceiling worth naming before the gate opens.

### 10 — Outcomes within Claim R′ / T′ · PASS

- Claim R′ (`claim_r_prime`, L998) is a three-way conjunction: traversal completed, records validated, retention stable across the scan.
- `NO_MATCH_AT_CURSOR` explicitly withholds R′, and `absence_means` states that an absent match is not evidence that no modification occurred.
- `claim_t_prime` is a fixed refusal string; no code path sets it.
- All 17 `INCONCLUSIVE_*` outcomes name a condition and none reports a negative finding about the file.
- `COMPLETED` carries a `completeness` string limiting the claim to the documented V2–V4 sequence under the selection semantics.

---

## 3. Defect — `F-1`, handle lifetime · **FAIL**

**Location:** `volume_io`, L1229–L1356 (whole function).

**Observed:** `CreateFileW` is called once and the resulting `handle` is captured by the two returned closures. **There is no `CloseHandle` call anywhere in `usn_reader.py`** — the only `CloseHandle` in the harness is in `unbuffered_digest.py` (L90–L92, L259, L312, L330), which is not modifiable under this work order and is not used by this module.

**Mechanism and consequence:**

| Aspect | Effect |
|---|---|
| Release | The volume handle stays open until the process exits. The caller has no way to close it — it holds only closures |
| Determinism | A caller cannot bound the window during which the volume is held open, which is exactly the property a bounded diagnostic is supposed to provide |
| Sharing | The handle uses `FILE_SHARE_READ | FILE_SHARE_WRITE`, so it will not block other readers or writers. It cannot cause a sharing violation or freeze a backup product |
| Data risk | **None.** No write, no flush, no lock, no dismount |
| Severity | **Low for safety, moderate for hygiene and for the "bounded window" claim** |

**Why it is a FAIL and not a limitation:** the module's own docstring promises a bounded read-only diagnostic whose live behaviour is a gated, discrete act. A handle that cannot be released contradicts that promise in the one function that would run at the gate, and the same file already guards against the *other* handle mistake (`INVALID_HANDLE_VALUE`), which makes the omission a gap rather than a design choice.

**Proposed narrow correction (not applied):** make `volume_io` return a third element — a `close()` callable that calls `CloseHandle(handle)` once and is idempotent — or return a context manager whose `__exit__` closes the handle in a `finally`. Either way the module stays read-only and the suites stay offline: a check can assert that `CloseHandle` is declared with `argtypes=[wintypes.HANDLE]` and that the returned closer is idempotent. This is a source change to a frozen-ish artifact and needs its own authorisation; it is **not** made here.

**`L-3` — journal-id fallback.** If the initial `query()` raises, `journal["id"]` is left at `0`. `FSCTL_READ_USN_JOURNAL` with `UsnJournalID = 0` addresses the current journal, so a read could in principle target a different journal instance than the one the scan intends. In practice the scan performs its own query and refuses with an inconclusive outcome if that fails, and `_check_retention` compares IDs across the scan — so the fallback is caught downstream rather than trusted. It is a deferred check, not a silent assumption, and it is recorded as a limitation rather than a defect.

**`L-4` — opening is gated by a caller parameter, not a capability boundary.** `allow_open=True` is keyword-only and defaults to `False`, and the raw-path and platform checks precede it, but any importer that passes the flag can open a volume. There is no import-time or scan-time handle and no ambient authority, so this is a design property to state plainly, not a hole to patch.

---

## 4. Offline verification

| Command | Result |
|---|---|
| `python -m eval.harness.usn_tests` | **85/85 passed**, exit 0 |
| `python -m eval.harness.validate` | **89/89 passed**, 0 failed, exit 0 |
| `python docs/research/corpus/check_corpus.py` | **PASS** (0 fail, 0 warn, 0 unverifiable) |

**What the suite can and cannot show.** It exercises parsing, status mapping, budget enforcement, claim ceilings, V4 arithmetic, `SourceInfo`, and the *declaration* of the Win32 signatures. It cannot show that a live ioctl returns what the documentation says, that `lpBytesReturned` is reliable under `METHOD_NEITHER`, that a sequence ever spans two real buffers, or that `CloseHandle` is ever reached — because no path in the module calls it.

---

## 5. Claim ceiling

**Establishes:** the reader's Win32 boundary is type-correct and its handles are compared correctly; buffer lifetimes are safe; a zero or unreliable byte count cannot yield a clean result; the budget is enforced before every call; report and retention caps are as specified; and every outcome stays inside Claims R′/T′.

**Does not establish:** that `FSCTL_READ_USN_JOURNAL` behaves as documented on this machine; that the return count is trustworthy; that a V4 sequence survives a real buffer boundary; or that the diagnostic will be conclusive when run. One source-level defect (`F-1`) remains open, and the three live gates remain closed.

**Code tested ≠ diagnostic validated ≠ investigation conclusive.**
