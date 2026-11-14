# `USN-HANDLE-LIFETIME-PATCH-PROPOSAL` — v2

**Identity:** `MEMORY-09-USN-PATCH-2` · **Supersedes:** v1 (`5e88c9b0…`, `eval/harness/USN-HANDLE-LIFETIME-PATCH-PROPOSAL-2026-11-14.md`)
**Status:** **PROPOSAL ONLY — not applied.** `usn_reader.py` `2b384069…` and `usn_tests.py` `f6152707…` are unchanged; the suite still reports 85/85 on unpatched code.
**Defect addressed:** `F-1` — `volume_io` opens a volume handle and never closes it.
**What changed from v1:** P-1 (exception masking during construction), P-2 (masking in `__exit__`), P-3 (non-deterministic test), P-4 (residual leak documented), a testable construction seam, an explicit `CloseHandle`-failure contract, and a ctypes declaration review.

---

## 1. Requirements, and where each is met

| # | Requirement | Met by |
|---:|---|---|
| 1 | Idempotent `close()` / context manager | `close()` returns `False` on repeat; `with` supported |
| 2 | Closes on success, exception, early return, cancellation | `__exit__` always closes; `_finish_open` closes before propagating any construction failure |
| 3 | Reads after closure fail deterministically | `_require_open()` raises `RuntimeError("volume session is closed")` |
| 4 | Repeated closure safe | `_closed` flag under the lock, plus a second guard inside `close_handle` |
| 5 | Not closed during an active call | `read()`/`query()` hold the lock for the whole ioctl; `close()` must take it to flip the flag |
| 6 | Preserve bounds, limits, outcomes, Claim R′/T′ | no change to `Caps`, `scan`, parsing, retention or any outcome string |
| 7 | Offline tests for closure, repeat, exception, use-after-close | injected callables; six of the ten checks are exactly that |
| 8 | **(new)** Preserve the original exception when cleanup fails | `_finish_open`'s `try/except OSError: pass` before a bare `raise` |
| 9 | **(new)** `CloseHandle` failure surfaced, never silent | recorded on `session.close_error`; documented in the docstring and asserted by callers |
| 10 | **(new)** Query-failure cleanup tested | `_h7`/`_h8` |

---

## 2. Corrected code

### 2.1 Module-level session (insert above `volume_io`)

```python
class VolumeSession:
    """One volume handle, owned for a bounded lifetime.

    The handle is closed exactly once. Reads and queries are refused after
    closure, and `close` cannot run underneath an in-flight call because the
    lock is held for the whole ioctl. If `CloseHandle` fails, the session still
    reports closed (a retry could release a recycled handle), the failure is
    stored on `close_error`, and a caller that requires a clean release must
    check it.
    """

    def __init__(self, query_fn, read_fn, close_fn):
        self._query_fn = query_fn
        self._read_fn = read_fn
        self._close_fn = close_fn
        self._lock = threading.Lock()
        self._closed = False
        self._reads = 0
        self.close_error = None

    @property
    def closed(self) -> bool:
        return self._closed

    @property
    def reads(self) -> int:
        return self._reads

    def _require_open(self) -> None:
        if self._closed:
            raise RuntimeError("volume session is closed")

    def query(self) -> JournalState:
        with self._lock:
            self._require_open()
            return self._query_fn()

    def read(self, start_usn: int, buffer_bytes: int) -> bytes:
        with self._lock:
            self._require_open()
            self._reads += 1
            return self._read_fn(start_usn, buffer_bytes)

    def close(self) -> bool:
        """Close the handle once. Returns True if this call closed it.

        A release error is recorded on `close_error`, never raised: raising from
        `__exit__` would replace the exception the context manager is already
        unwinding.
        """
        with self._lock:
            if self._closed:
                return False
            self._closed = True
        try:
            self._close_fn()
        except OSError as exc:
            self.close_error = exc
        return True

    def __enter__(self) -> "VolumeSession":
        return self

    def __exit__(self, exc_type, exc, tb) -> bool:
        self.close()
        return False
```

### 2.2 Construction seam (module level, testable without a volume)

```python
def _finish_open(query_fn, read_fn, close_fn, journal):
    """Resolve the journal id, build the session, and close on any failure.

    A failing initial query is not fatal: the id falls back to 0 and the scan's
    own retention check catches a changed journal (limitation L-3). Anything
    else closes the handle before re-raising the original exception.
    """
    try:
        try:
            journal["id"] = query_fn().journal_id
        except OSError:
            journal["id"] = 0
        return VolumeSession(query_fn, read_fn, close_fn)
    except BaseException:
        try:
            close_fn()
        except OSError:
            pass          # never mask the failure we are propagating
        raise             # bare raise preserves the original exception object
```

### 2.3 Inside `volume_io`

```diff
@@
+    kernel32.CloseHandle.restype = wintypes.BOOL
+    kernel32.CloseHandle.argtypes = [wintypes.HANDLE]
@@
     if handle is None or handle in (0, invalid_handle):
         raise OSError(ctypes.get_last_error(), f"CreateFileW failed for {volume_path}")
+
+    closed = {"done": False}
+
+    def close_handle() -> None:
+        # Idempotent at the Win32 layer too: a second CloseHandle on a stale
+        # numeric handle could release a recycled one.
+        if closed["done"]:
+            return
+        closed["done"] = True
+        if not kernel32.CloseHandle(handle):
+            raise OSError(ctypes.get_last_error(), "CloseHandle failed")
@@
-    try:
-        journal["id"] = query().journal_id
-    except OSError:
-        journal["id"] = 0
-    return query, read
+    return _finish_open(query, read, close_handle, journal)
```

The docstring changes from "Return (query_fn, read_fn) bound to a raw volume handle" to "Open a raw volume and return a `VolumeSession` bound to it. The caller owns the handle: use it as a context manager, or call `close()`. Reads after closure raise `RuntimeError`; a release failure is reported on `close_error`."

---

## 3. ctypes declaration review

| Declaration | Type used | Correct? | Note |
|---|---|---|---|
| `CreateFileW.restype` | `wintypes.HANDLE` (`c_void_p`) | **yes** | a default `int` return truncates a 64-bit handle — the original `F-1`-class defect |
| `CreateFileW.argtypes` | `LPCWSTR, DWORD, DWORD, LPVOID, DWORD, DWORD, HANDLE` | **yes** | `lpSecurityAttributes` and `hTemplateFile` are `LPVOID`/`HANDLE`; `None` is accepted for both |
| `DeviceIoControl.restype` | `wintypes.BOOL` (`c_long`) | **yes** | compares as truthiness; `0` is failure |
| `DeviceIoControl.argtypes` | `HANDLE, DWORD, LPVOID, DWORD, LPVOID, DWORD, POINTER(DWORD), LPVOID` | **yes** | `lpBytesReturned` must be `LPDWORD`; `lpOverlapped` is `None`, consistent with a non-overlapped handle |
| `CloseHandle.restype` | `wintypes.BOOL` | **yes** | 4-byte BOOL |
| `CloseHandle.argtypes` | `[wintypes.HANDLE]` | **yes** | passing the Python `int` returned by `CreateFileW` is accepted by `c_void_p`; no manual cast needed and none should be added |
| `INVALID_HANDLE_VALUE` | `ctypes.c_void_p(-1).value` | **yes** | equals `18446744073709551615`; comparing against `-1` would never match |
| Request struct | `ctypes.byref(request)` + `ctypes.sizeof(request)` | **yes** | 48 B, offsets asserted offline |
| Output buffer | `ctypes.cast(buffer, ctypes.c_void_p)` + explicit length | **yes** | untyped `LPVOID` with a truthful length |
| Returned bytes | `buffer.raw[: returned.value]` | **yes** | copies before the buffer dies; clips an impossible count |

**One hazard worth naming.** `CloseHandle` accepts any `int` as a `c_void_p`. After close, nothing in `ctypes` prevents the session from being made to pass a stale handle — which is precisely why the session refuses calls after `_closed` is set. That refusal is the safety property; the declaration is only its precondition.

---

## 4. Test plan — ten checks, 85 → 95

```python
@check("handle: a session closes on normal exit")                     # _h1
@check("handle: close is idempotent and reports the first call only") # _h2
@check("handle: an exception inside the context still closes, once")  # _h3
@check("handle: a closed session refuses reads and queries")          # _h4
@check("handle: racing closures close exactly once (deterministic)")  # _h5  barrier race
@check("handle: close does not complete under an in-flight read (bounded)")  # _h6
@check("handle: a failing initial query falls back to id 0")          # _h7
@check("handle: a fatal query error closes and preserves the original")  # _h8
@check("handle: a failing CloseHandle is recorded, not raised")       # _h9
@check("handle: CloseHandle is declared and confined to volume_io (static)")  # _h10
```

The four new ones, in full, because they are the ones that answer P-1, P-2 and the query-failure case:

```python
@check("handle: a failing initial query falls back to id 0")
def _h7():
    journal = {"id": 7}
    calls = []

    def boom():
        raise OSError(1179, "journal not active")

    s = R._finish_open(boom, lambda a, b: b"", lambda: calls.append(1), journal)
    assert journal["id"] == 0, journal          # L-3 fallback preserved
    assert s.closed is False and calls == []    # handle NOT closed on the soft path
    s.close()
    assert calls == [1]


@check("handle: a fatal query error closes and preserves the original")
def _h8():
    journal = {"id": 7}
    calls = []

    def fatal():
        raise ValueError("programming error")

    try:
        R._finish_open(fatal, lambda a, b: b"", lambda: calls.append(1), journal)
    except ValueError as exc:
        assert "programming error" in str(exc)   # original, not a cleanup error
    else:
        raise AssertionError("fatal query error was swallowed")
    assert calls == [1], calls                   # handle released exactly once


@check("handle: a failing CloseHandle is recorded, not raised")
def _h9():
    def bad_close():
        raise OSError(6, "invalid handle")

    s = R.VolumeSession(lambda: None, lambda a, b: b"", bad_close)
    assert s.close() is True
    assert s.closed is True                      # documented: still reports closed
    assert isinstance(s.close_error, OSError)    # surfaced, not silent
    assert s.close() is False                    # no retry
```

`_h5` is the deterministic concurrency test from the review (eight threads behind a `threading.Barrier`; exactly one `close_fn` call under every interleaving). `_h6` is retained but labelled **bounded, may false-pass** — it can never false-fail, because the lock makes the ordering assertion one-directional.

**Expected outcome if applied:** 95/95, `validate` unchanged at 89/89, no live path exercised.

---

## 5. The `CloseHandle` failure contract (P-4 and the closed-flag question)

| Situation | Session state | `close_error` | Caller action |
|---|---|---|---|
| Release succeeds | closed | `None` | none |
| Release fails | **closed** (retry forbidden) | the `OSError` | record it; the OS may still hold the handle |
| Release fails during `__exit__` while the body raised | closed | the `OSError` | the body's exception propagates; check `close_error` afterwards |
| Caller never closes | **not closed** | `None` | a defect in the caller; the session cannot prevent it |

The last row is the honest limit: this patch makes the leak **avoidable**, not impossible. The docstring says so, and the gate script that eventually calls `volume_io` must use `with`, and must assert `session.close_error is None`.

---

## 6. Out of scope

`L-3` (the `journal["id"] = 0` fallback) is **preserved unchanged**. It is a separate limitation with its own review, and folding it in would make this patch unreviewable. No change to `Caps`, no change to any outcome, no change to `claim_r_prime`, `_check_retention` or the `allow_open` gate.

---

## 7. Risks

| Risk | Severity | Mitigation |
|---|---|---|
| A caller unpacking the old 2-tuple breaks | low | no in-repo caller does; docstring states the new contract |
| `close()` blocks behind a hung ioctl | inherent | documented; a hung device defeats any lock design. A caller-side watchdog is the only mitigation and is out of scope |
| Caller ignores the session | medium | documented contract + `with` in the gate script |
| `close_error` ignored | medium | the gate script must assert it; a test cannot enforce a caller |
| A future edit passes a stale handle | low | `_require_open()` refuses; `close_handle`'s own flag prevents a double release |

---

## 8. Claim ceiling

**Establishes:** a corrected, narrowly scoped design that closes the handle exactly once on every exit path, preserves the original exception when cleanup also fails, records a release failure instead of hiding it, refuses use after closure, serialises closure against in-flight calls, and is testable offline — with the ten checks that would demonstrate it.

**Does not establish:** that `CloseHandle` succeeds on a real volume handle, that a live session behaves as the fakes do, or that the leak is fixed in practice. Nothing is applied; the suite still runs the unpatched code, and this remains a maintenance correction that opens no storage gate.
