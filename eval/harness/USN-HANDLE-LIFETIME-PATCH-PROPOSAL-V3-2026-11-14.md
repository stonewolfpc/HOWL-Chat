# `USN-HANDLE-LIFETIME-PATCH-PROPOSAL` — v3

**Identity:** `MEMORY-09-USN-PATCH-3` · **Supersedes:** v2 (`f1cf6c17…`), v1 (`5e88c9b0…`)
**Status:** **PROPOSAL ONLY — not applied.** `usn_reader.py` `2b384069…`, `usn_tests.py` `f6152707…`; suite still 85/85 on unpatched code.
**Defect:** `F-1` — `volume_io` opens a volume handle and never closes it.

---

## 0. What the review found, checked against the actual v2 text

Honest reconciliation before anything else, because two of the three findings attach to v1 rather than v2:

| Review finding | Status against v2 as filed |
|---|---|
| **P-1** — "`except BaseException: close_handle(); raise`; the claimed inner `except OSError: pass` is not present" | **v2 already has the fix.** v2 §2.2 lines 112–120 read `try: close_fn() / except OSError: pass / raise`. The quoted form is v1's. v2 §2.3's diff no longer contains that block at all |
| **P-2** — "there is no `close_error` attribute or implementation in the supplied class" | **v2 already has it.** v2 §2.1 line 50 `self.close_error = None`; lines 88–89 `except OSError as exc: self.close_error = exc`; `close()` does not raise. The description is accurate for v1, where `close()` re-raised |
| **P-3** — "`_h5` remains the sleep-based test; the eight-thread barrier race is not included" | **correct.** v2 *named* `_h5` as a barrier race in a checklist but never pasted its body. The only fully-specified concurrency test in the repository is still v1's sleep-based one. This is a real gap, and v3 closes it |
| **Missing**: body-exception + close-failure test | **correct**, not present in v2 |
| **Missing**: cleanup-failure during fatal query error | **correct**, not present in v2 |
| **Test count reconciliation** | **correct** — v2 claimed 10 checks → 95 without listing them all |

v3 therefore adds no new design; it supplies the missing implementations and tests, and reconciles the count.

---

## 1. Requirements

| # | Requirement | Where met |
|---:|---|---|
| 1 | Idempotent `close()` / context manager | `close()` returns `False` on repeat; `__enter__`/`__exit__` |
| 2 | Closes on success, exception, early return, cancellation | `__exit__` always closes; `_finish_open` closes before propagating |
| 3 | Reads after closure fail deterministically | `_require_open()` → `RuntimeError("volume session is closed")` |
| 4 | Repeated closure safe | `_closed` under the lock; second guard in `close_handle` |
| 5 | Not closed during an active call | `read()`/`query()` hold the lock across the ioctl |
| 6 | Bounds, limits, outcomes, Claim R′/T′ preserved | no change to `Caps`, `scan`, parsing, retention, outcomes |
| 7 | Offline tests for the four core behaviours | `_h1`–`_h4` |
| 8 | Original exception preserved when cleanup fails | `_finish_open`; tested by `_h8`, `_h11` |
| 9 | `CloseHandle` failure observable and retained | `close_error`, `released`; tested by `_h9` |
| 10 | Failure during initial journal query cleaned up | `_finish_open`; tested by `_h8`, `_h11` |
| 11 | Body exception not masked by a close failure | `close()` records; tested by `_h10` |
| 12 | Deterministic concurrency test | `_h5` (barrier race) |

---

## 2. Complete code

### 2.1 Session

```python
class VolumeSession:
    """One volume handle, owned for a bounded lifetime.

    The handle is closed exactly once. Reads and queries are refused after
    closure, and `close` cannot run underneath an in-flight call because the
    lock is held for the whole ioctl.

    A failed CloseHandle is never retried (the value may have been recycled) and
    never raised (that would replace the exception being unwound). It is stored
    on `close_error`; `released` is False in that case, and a caller that needs a
    clean release must check it.
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
    def released(self) -> bool:
        """True only if the handle was closed and CloseHandle succeeded."""
        return self._closed and self.close_error is None

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
        """Close the handle once. Returns True if this call closed it."""
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
        self.close()          # records, never raises: cannot mask the body exception
        return False
```

`closed` and `released` read state without the lock. That is deliberate: they exist for reporting, and taking the lock would deadlock a caller inspecting them from inside a fake `close_fn`. The flags are written under the lock and are monotonic.

### 2.2 Construction seam

```python
def _finish_open(query_fn, read_fn, close_fn, journal):
    """Resolve the journal id, build the session, and close on any failure.

    An OSError from the initial query is the *soft* path (limitation L-3): the id
    falls back to 0, the handle stays open, and a session is returned. The scan's
    own retention check catches a changed journal.

    Any other failure — including a KeyboardInterrupt — closes the handle and
    re-raises the original exception. A failure during that cleanup is swallowed
    so it cannot replace the exception being propagated.
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
            pass
        raise
```

### 2.3 `volume_io`

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
+        if closed["done"]:
+            return                      # never a second CloseHandle on a stale value
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

---

## 3. ctypes declarations (unchanged from v2, restated)

| Declaration | Value | Verified |
|---|---|---|
| `CreateFileW.restype` | `wintypes.HANDLE` | pointer-width return; a default `int` truncates the handle |
| `CloseHandle.restype` | `wintypes.BOOL` | 4-byte BOOL |
| `CloseHandle.argtypes` | `[wintypes.HANDLE]` | accepts the int returned by `CreateFileW`; no manual cast |
| `INVALID_HANDLE_VALUE` | `ctypes.c_void_p(-1).value` | `18446744073709551615`; `-1` would never match |

Hazard retained for the record: `CloseHandle` accepts any `int` as `c_void_p`, so the safety property is `_require_open()` refusing calls after closure, not the declaration.

---

## 4. Twelve tests — every body written out

```python
@check("handle: a session closes on normal exit")
def _h1():
    calls = []
    with R.VolumeSession(lambda: "q", lambda a, b: b"r", lambda: calls.append(1)) as s:
        assert s.read(0, 1024) == b"r"
        assert s.query() == "q"
        assert s.closed is False and s.released is False
    assert calls == [1], calls
    assert s.closed is True and s.released is True


@check("handle: close is idempotent and reports the first call only")
def _h2():
    calls = []
    s = R.VolumeSession(lambda: None, lambda a, b: b"", lambda: calls.append(1))
    assert s.close() is True
    assert s.close() is False
    assert s.close() is False
    assert calls == [1], calls
    assert s.close_error is None


@check("handle: an exception inside the context still closes, once")
def _h3():
    calls = []
    try:
        with R.VolumeSession(lambda: None, lambda a, b: b"", lambda: calls.append(1)):
            raise ValueError("boom")
    except ValueError:
        pass
    else:
        raise AssertionError("the body exception was swallowed")
    assert calls == [1], calls


@check("handle: a closed session refuses reads and queries")
def _h4():
    s = R.VolumeSession(lambda: "q", lambda a, b: b"r", lambda: None)
    s.close()
    for name, call in (("query", s.query), ("read", lambda: s.read(0, 512))):
        try:
            call()
        except RuntimeError as exc:
            assert "closed" in str(exc), (name, exc)
        else:
            raise AssertionError(f"{name} was accepted after closure")


@check("handle: racing closures close exactly once (deterministic)")
def _h5():
    import threading
    calls = []
    s = R.VolumeSession(lambda: None, lambda a, b: b"", lambda: calls.append(1))
    barrier = threading.Barrier(8)

    def race():
        barrier.wait()
        s.close()

    threads = [threading.Thread(target=race) for _ in range(8)]
    for t in threads:
        t.start()
    for t in threads:
        t.join(5.0)
    assert calls == [1], calls        # true under every interleaving
    assert s.close_error is None


@check("handle: close does not complete under an in-flight read (bounded)")
def _h6():
    import threading
    order, started, release, attempted = (
        [], threading.Event(), threading.Event(), threading.Event()
    )

    def slow_read(start, size):
        order.append("read-start")
        started.set()
        release.wait(5.0)
        order.append("read-end")
        return b""

    s = R.VolumeSession(lambda: None, slow_read, lambda: order.append("close"))
    t = threading.Thread(target=lambda: s.read(0, 4096))
    t.start()
    assert started.wait(5.0), "read never started"
    closer = threading.Thread(target=lambda: (attempted.set(), s.close()))
    closer.start()
    assert attempted.wait(5.0), "closer never scheduled"
    # Bounded smoke test: can false-PASS if the closer has not yet reached
    # close(); can never false-FAIL, because the lock fixes the ordering.
    release.set()
    t.join(5.0); closer.join(5.0)
    assert order == ["read-start", "read-end", "close"], order


@check("handle: a failing initial query falls back to id 0 (soft path, no cleanup)")
def _h7():
    journal = {"id": 7}
    calls = []

    def soft_fail():
        raise OSError(1179, "the volume change journal is not active")

    s = R._finish_open(soft_fail, lambda a, b: b"", lambda: calls.append(1), journal)
    assert journal["id"] == 0, journal      # L-3 fallback preserved
    assert calls == [], calls               # handle deliberately NOT closed
    assert s.closed is False
    s.close()
    assert calls == [1]


@check("handle: a fatal query error closes once and preserves the original")
def _h8():
    journal = {"id": 7}
    calls = []

    def fatal():
        raise ValueError("programming error")

    try:
        R._finish_open(fatal, lambda a, b: b"", lambda: calls.append(1), journal)
    except ValueError as exc:
        assert "programming error" in str(exc), exc
    else:
        raise AssertionError("fatal query error was swallowed")
    assert calls == [1], calls


@check("handle: a failing CloseHandle is recorded, retained, and never retried")
def _h9():
    def bad_close():
        raise OSError(6, "invalid handle")

    s = R.VolumeSession(lambda: None, lambda a, b: b"", bad_close)
    assert s.close() is True
    assert s.closed is True             # closed, per policy
    assert s.released is False          # but the OS may still hold it
    assert isinstance(s.close_error, OSError)
    assert s.close() is False
    assert s.close_error is not None    # failure retained for the caller


@check("handle: a body exception survives a failing close")
def _h10():
    def bad_close():
        raise OSError(6, "invalid handle")

    s = R.VolumeSession(lambda: None, lambda a, b: b"", bad_close)
    try:
        with s:
            raise ValueError("body failure")
    except ValueError as exc:
        assert "body failure" in str(exc), exc     # NOT the close OSError
    else:
        raise AssertionError("body exception was masked")
    assert s.released is False
    assert isinstance(s.close_error, OSError)


@check("handle: cleanup failure during a fatal query keeps the original error")
def _h11():
    journal = {"id": 7}

    def fatal():
        raise ValueError("original failure")

    def bad_close():
        raise OSError(6, "cleanup failure")

    try:
        R._finish_open(fatal, lambda a, b: b"", bad_close, journal)
    except ValueError as exc:
        assert "original failure" in str(exc), exc
    else:
        raise AssertionError("original error was replaced or swallowed")


@check("handle: CloseHandle is declared and confined to volume_io (static)")
def _h12():
    src = open(R.__file__, encoding="utf-8").read()
    for needle in ("CloseHandle.restype", "CloseHandle.argtypes"):
        assert needle in src, needle
    index = src.find("CloseHandle(handle)")
    assert index > 0
    gate = src.rfind("def volume_io", 0, index)
    assert gate > 0 and gate < index, "CloseHandle called outside volume_io"
```

---

## 5. Test count, reconciled

| Section | Before | After | Delta |
|---|---:|---:|---:|
| Existing suite | 85 | 85 | 0 |
| New handle checks `_h1`…`_h12` | 0 | **12** | +12 |
| **Total** | **85** | **97** | **+12** |

Mapping, so nothing is claimed twice: `_h1` closure on exit · `_h2` idempotence · `_h3` exception cleanup · `_h4` use-after-close · `_h5` deterministic race · `_h6` bounded in-flight · `_h7` soft query failure (L-3, **not** a cleanup test) · `_h8` fatal query + cleanup · `_h9` CloseHandle failure recorded · `_h10` body exception + close failure · `_h11` cleanup failure during fatal query · `_h12` static declaration/confinement.

`validate.py` is untouched and stays at **89/89**.

---

## 6. The L-3 boundary, stated so it cannot be miscounted

| Initial query outcome | Handle | Session | Which test |
|---|---|---|---|
| succeeds | open | returned, id set | `_h1` |
| `OSError` | **stays open** | returned, id `0` | `_h7` — a **soft path**, not cleanup |
| any other exception | **closed** | not returned | `_h8`, `_h11` |

`L-3` is preserved unchanged: a failed initial query still falls back to id `0`, and `_check_retention` remains the thing that catches a changed journal. Describing that path as "handle-cleanup" would be wrong, and the tests are named to prevent it.

---

## 7. Risks

| Risk | Severity | Mitigation |
|---|---|---|
| Caller unpacking the old 2-tuple | low | no in-repo caller does; docstring states the contract |
| `close()` blocks behind a hung ioctl | inherent | documented; a hung device defeats any lock design |
| Caller ignores the session | medium | `with` required in the gate script; contract stated |
| `close_error` ignored | medium | gate script asserts `session.released` |
| `released`/`closed` read without the lock | low | monotonic flags, reporting only |

---

## 8. Claim ceiling

**Establishes:** complete code for all three fixes, twelve fully written checks with a reconciled total of 97, and an explicit boundary between the soft query path and the cleanup path.

**Does not establish:** that `CloseHandle` succeeds on a real volume handle, or that a live session behaves as the fakes do. Nothing is applied; the suite still runs 85 checks against the unpatched module, and this correction opens no storage gate.
