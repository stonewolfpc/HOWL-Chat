# `USN-HANDLE-LIFETIME-PATCH-PROPOSAL` — v4

**Identity:** `MEMORY-09-USN-PATCH-4.1` · **Supersedes:** v4 (`e46dba2e…`), v3 (`daff292e…`), v2 (`f1cf6c17…`), v1 (`5e88c9b0…`)
**Status:** **PROPOSAL ONLY — not applied.** `usn_reader.py` `2b384069…`, `usn_tests.py` `f6152707…`; suite still 85/85 on unpatched code.
**Fixed in v4:** the `released` race the review found (a session could report `released == True` while `CloseHandle` was still running); `_h5` now checks the full API contract; a new `_h13` reproduces the race deterministically; the cleanup exception contract is stated for what it is.
**Fixed in 4.1:** the docstring claimed "the handle is closed exactly once", which the patch's own `_h9` contradicts — the contract is an **at-most-once release attempt**, success or reported failure; `_h6` and `_h13` now assert their worker threads terminated, as `_h5` already did.

---

## 1. The defect the review found, and why a boolean pair cannot fix it

v3's `released` was `self._closed and self.close_error is None`. Between `_closed = True` (under the lock) and the assignment of `close_error` (after `close_fn` returned), any observer sees `closed=True`, `close_error=None`, therefore `released=True` — **before the release was even attempted**.

Two booleans cannot express "in progress". v4 replaces them with **one state field** whose transitions are the contract:

| State | Means | `closed` | `released` | Reads/queries |
|---|---|---|---|---|
| `open` | handle held, usable | false | false | allowed |
| `closing` | release in progress; outcome unknown | true | **false** | refused |
| `released` | `CloseHandle` returned success | true | true | refused |
| `close_failed` | `CloseHandle` returned failure; not retried | true | false | refused |

`released` is true **only** in `released`. `closed` is true in all three terminal/in-progress states, which preserves v3's safety behaviour: no call may use the handle once closure has begun.

---

## 2. Session

```python
_STATE_OPEN = "open"
_STATE_CLOSING = "closing"
_STATE_RELEASED = "released"
_STATE_CLOSE_FAILED = "close_failed"


class VolumeSession:
    """One volume handle, owned for a bounded lifetime.

    The handle's release is attempted at most once. Reads and queries are
    refused once closure begins, and `close` cannot start underneath an
    in-flight call because the lock is held for the whole ioctl.

    At most once is the entire contract, and it is not the same as "closed":
    if CloseHandle returns success the session records `released`; if it
    returns failure the session records `close_failed` and does not retry. A
    failed CloseHandle call does not prove the OS released the handle, so a
    caller that needs a clean release must check `released`, never `closed`.

    State is explicit so callers can distinguish "closure started" from "the
    handle was released":

        open -> closing -> released
                        -> close_failed

    A failed release is never retried (the numeric value may have been
    recycled) and never raised from `close` (that would replace the exception
    being unwound). It is visible as `state == "close_failed"` and on
    `close_error`.
    """

    def __init__(self, query_fn, read_fn, close_fn):
        self._query_fn = query_fn
        self._read_fn = read_fn
        self._close_fn = close_fn
        self._lock = threading.Lock()
        self._state = _STATE_OPEN
        self._reads = 0
        self.close_error = None

    @property
    def state(self) -> str:
        return self._state

    @property
    def closed(self) -> bool:
        return self._state != _STATE_OPEN

    @property
    def released(self) -> bool:
        return self._state == _STATE_RELEASED

    @property
    def reads(self) -> int:
        return self._reads

    def _require_open(self) -> None:
        if self._state != _STATE_OPEN:
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
        """Begin closure. Returns True if this call performed it.

        A concurrent or repeated call returns False immediately: it did not
        close anything. The state, not a boolean, reports the outcome.
        """
        with self._lock:
            if self._state != _STATE_OPEN:
                return False
            self._state = _STATE_CLOSING
        try:
            self._close_fn()
        except Exception as exc:          # noqa: BLE001 - recorded, not raised
            with self._lock:
                self.close_error = exc
                self._state = _STATE_CLOSE_FAILED
            return True
        with self._lock:
            self._state = _STATE_RELEASED
        return True

    def __enter__(self) -> "VolumeSession":
        return self

    def __exit__(self, exc_type, exc, tb) -> bool:
        self.close()          # never raises: cannot mask the body exception
        return False
```

### 2.1 Interrupt during release — stated, not hidden

If `close_fn` raises a `BaseException` (a second `KeyboardInterrupt`, `SystemExit`), it is **not** caught. The session remains in `closing`: `closed` true, `released` false, further `close()` calls return `False`, and the outcome of the release is genuinely unknown. That is the honest representation, and it is preferable to claiming either success or failure. Recovery is process exit.

---

## 3. Construction seam and its exception contract

```python
def _finish_open(query_fn, read_fn, close_fn, journal):
    """Resolve the journal id, build the session, and close on any failure.

    OSError from the initial query is the soft path (limitation L-3): id falls
    back to 0, the handle stays open, a session is returned.

    Any other failure closes the handle and re-raises the original exception.
    A cleanup failure is swallowed when it is an Exception; a BaseException
    from cleanup still replaces the original, deliberately — an interrupt
    during cleanup outranks the failure being propagated.
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
        except Exception:
            pass
        raise
```

**Qualified, as the review asked:** the preservation guarantee covers `OSError` **and every other `Exception`** raised by cleanup. It does **not** cover a `BaseException` from cleanup, which replaces the original by design. v3 claimed the broader property without that qualifier.

### 3.1 `volume_io`

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

ctypes declarations are unchanged from v2/v3 and remain correct: `CreateFileW.restype = wintypes.HANDLE` (a default `int` truncates the handle), `CloseHandle.argtypes = [wintypes.HANDLE]` (accepts the returned int; no manual cast), `INVALID_HANDLE_VALUE` compared as `c_void_p(-1).value`.

---

## 4. Tests — thirteen, all bodies written

```python
@check("handle: a session closes on normal exit")
def _h1():
    calls = []
    with R.VolumeSession(lambda: "q", lambda a, b: b"r", lambda: calls.append(1)) as s:
        assert s.read(0, 1024) == b"r"
        assert s.state == "open" and s.released is False
    assert calls == [1], calls
    assert s.state == "released" and s.released is True


@check("handle: close is idempotent and reports the first call only")
def _h2():
    calls = []
    s = R.VolumeSession(lambda: None, lambda a, b: b"", lambda: calls.append(1))
    assert s.close() is True
    assert s.close() is False
    assert s.close() is False
    assert calls == [1], calls
    assert s.state == "released"


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


@check("handle: racing closures initiate closure exactly once")
def _h5():
    import threading
    calls, results = [], []
    lock = threading.Lock()
    s = R.VolumeSession(lambda: None, lambda a, b: b"", lambda: calls.append(1))
    barrier = threading.Barrier(8)

    def race():
        barrier.wait()
        value = s.close()
        with lock:
            results.append(value)

    threads = [threading.Thread(target=race) for _ in range(8)]
    for t in threads:
        t.start()
    for t in threads:
        t.join(5.0)
    assert all(not t.is_alive() for t in threads), "a thread never terminated"
    assert len(results) == 8, results
    assert results.count(True) == 1, results
    assert results.count(False) == 7, results
    assert calls == [1], calls
    assert s.state == "released" and s.released is True


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
    # Bounded smoke test: can false-PASS if the closer has not yet contended
    # for the lock; can never false-FAIL. Not a proof of lock contention.
    release.set()
    t.join(5.0); closer.join(5.0)
    assert not t.is_alive() and not closer.is_alive(), "a worker never terminated"
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
    assert s.state == "open"
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
    assert s.state == "close_failed"
    assert s.closed is True
    assert s.released is False          # the OS may still hold the handle
    assert isinstance(s.close_error, OSError)
    assert s.close() is False
    assert s.close_error is not None    # retained for the caller


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
    assert s.state == "close_failed"
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


@check("handle: closed-but-not-yet-released is observable (regression)")
def _h13():
    import threading
    entered, release = threading.Event(), threading.Event()

    def slow_close():
        entered.set()
        release.wait(5.0)

    s = R.VolumeSession(lambda: None, lambda a, b: b"", slow_close)
    t = threading.Thread(target=s.close)
    t.start()
    assert entered.wait(5.0), "close_fn never entered"
    # The race the review found: before the state machine, released was True here.
    assert s.state == "closing"
    assert s.closed is True
    assert s.released is False, "released reported true before the release finished"
    assert s.close() is False           # a concurrent close reports it did nothing
    release.set()
    t.join(5.0)
    assert not t.is_alive(), "closer never terminated"
    assert s.state == "released" and s.released is True


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

| Section | Before | After |
|---|---:|---:|
| Existing suite | 85 | 85 |
| New handle checks `_h1`–`_h13` | 0 | **13** |
| **Total** | **85** | **98** |

`validate.py` untouched at **89/89**. `_h5` now asserts all eight threads terminate, exactly one `True` and seven `False`, the callback ran once, and the final state is `released`. `_h13` is the regression test for the newly found race: it is deterministic, because `entered` is set from inside `close_fn`, which runs only after the state has become `closing`.

---

## 6. Where each review finding landed

| Finding | Resolution |
|---|---|
| `released` true during `CloseHandle` | **fixed** — explicit four-state machine; `released` is true only in `released`; `_h13` reproduces the window |
| `_h5` did not check thread completion or return values | **fixed** — all eight joined and asserted not alive; results counted |
| `_h6`/`_h13` did not assert worker termination | **fixed in 4.1** — both assert `not is_alive()` after the bounded join |
| Docstring overclaimed a guaranteed close | **fixed in 4.1** — at-most-once release attempt, with `released` as the only success predicate |
| A concurrent `close()` returning `False` | **retained** — it means that call did not initiate closure, never that the handle was released; `state`/`released` is the outcome |
| `_h6` is not proof of lock contention | **retained and relabelled** as a bounded smoke test that can false-pass, never false-fail |
| `_finish_open` preservation claim too broad | **qualified** — covers every `Exception`; a `BaseException` from cleanup replaces the original, stated explicitly |
| L-3 soft path must not read as a cleanup test | **preserved** — `_h7` named "soft path, no cleanup"; §3 documents it |

---

## 7. Risks

| Risk | Severity | Mitigation |
|---|---|---|
| Caller unpacking the old 2-tuple | low | no in-repo caller does; docstring states the contract |
| `close()` blocks behind a hung ioctl | inherent | documented; a hung device defeats any lock design |
| Caller ignores the session | medium | `with` required; contract stated |
| Caller checks `closed` instead of `released` | medium | `released` documented as the release predicate; gate script asserts it |
| Interrupt during release leaves `closing` | low | documented; outcome genuinely unknown; process exit recovers |
| State fields read without the lock | low | monotonic transitions, reporting only |

---

## 8. Claim ceiling

**Establishes:** code that distinguishes "closure started" from "handle released", a deterministic regression test for the race that v3 could not detect, a concurrency test that checks the API contract rather than only the callback, and an exception contract stated for exactly what it covers.

**Does not establish:** that `CloseHandle` succeeds on a real volume handle, that a live session behaves as the fakes do, or that the leak is fixed in practice. Nothing is applied; the suite still runs 85 checks against the unpatched module, and this remains a maintenance correction that opens no storage gate.
