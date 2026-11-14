# `USN-HANDLE-LIFETIME-PATCH-PROPOSAL`

**Identity:** `MEMORY-09-USN-PATCH-1` · **Date:** 2026-11-14 · **Track:** `WO-NEXT-02` Track 4
**Status:** **PROPOSAL ONLY — not applied.** `eval/harness/usn_reader.py` is unchanged (`2b384069…`), `usn_tests.py` is unchanged (`f6152707…`), and the suite still reports 85/85.
**Defect:** `F-1` from `USN-READER-FINAL-AUDIT` — `volume_io` opens a volume handle and never closes it; the caller receives two closures and no way to release the handle.

---

## 1. What the patch must achieve

| # | Requirement | How this design meets it |
|---:|---|---|
| 1 | Idempotent `close()` or context manager | `close()` returns `False` on the second call and does not touch the handle again; `with` is supported |
| 2 | Closes on success, exception, early return, cancellation | `__exit__` always closes; `volume_io` closes before propagating any failure after `CreateFileW` |
| 3 | Prevent reads after closure | `read()` and `query()` raise `RuntimeError("volume session is closed")` |
| 4 | Repeated closure safe | guarded by the `_closed` flag under the same lock |
| 5 | Not closed while a read is active | `read()` holds the lock for the entire ioctl, so `close()` waits for it |
| 6 | Preserve bounds, limits, outcomes, Claim R′/T′ | no change to `Caps`, `scan`, `parse_records`, `_check_retention` or any outcome; the session wraps the same two closures |
| 7 | Offline tests for closure, repeat closure, exception cleanup, use-after-close | `VolumeSession` takes injected callables, so all four are testable with no volume present |

---

## 2. Design

Two changes, both inside `volume_io`'s file:

1. **A module-level `VolumeSession`** that owns a handle's lifetime. It takes three callables — `query_fn`, `read_fn`, `close_fn` — so its semantics are exercised offline. It holds one `threading.Lock`; `read` and `query` hold it across the call, `close` takes it only to flip the flag.
2. **`volume_io` returns the session** instead of a `(query, read)` tuple, declares `CloseHandle`, and closes the handle if anything fails after the handle was created.

**Deliberate non-goals.** No `__del__` (it can run during interpreter shutdown, where closing a handle is worse than leaking one), no `weakref.finalize` (same reason, plus it hides the leak the patch exists to remove), no change to the gate, no change to any read bound.

**Return-shape change is safe in this repository.** `scan` receives `query_fn`/`read_fn` as parameters and is unaffected; the only callers of `volume_io` in the suite are `_z1` and `_z2`, which assert the *refusal* paths and never reach a session. No in-repo code unpacks the success path.

---

## 3. Proposed diff

### 3.1 New module-level class (insert immediately above `def volume_io`)

```diff
@@
+class VolumeSession:
+    """One volume handle, owned for a bounded lifetime.
+
+    The handle is closed exactly once. Reads and queries are refused after
+    closure, and `close` cannot run underneath an in-flight call because the
+    lock is held for the whole ioctl. The three callables are injected so the
+    lifetime rules are testable without a volume.
+    """
+
+    def __init__(self, query_fn, read_fn, close_fn):
+        self._query_fn = query_fn
+        self._read_fn = read_fn
+        self._close_fn = close_fn
+        self._lock = threading.Lock()
+        self._closed = False
+        self._reads = 0
+
+    @property
+    def closed(self) -> bool:
+        return self._closed
+
+    @property
+    def reads(self) -> int:
+        return self._reads
+
+    def _require_open(self) -> None:
+        if self._closed:
+            raise RuntimeError("volume session is closed")
+
+    def query(self) -> JournalState:
+        with self._lock:
+            self._require_open()
+            return self._query_fn()
+
+    def read(self, start_usn: int, buffer_bytes: int) -> bytes:
+        with self._lock:
+            self._require_open()
+            self._reads += 1
+            return self._read_fn(start_usn, buffer_bytes)
+
+    def close(self) -> bool:
+        """Close the handle once. Returns True if this call closed it."""
+        with self._lock:
+            if self._closed:
+                return False
+            self._closed = True
+        try:
+            self._close_fn()
+        finally:
+            # A failed CloseHandle is not retried: the handle may have been
+            # released anyway, and a second close could release a recycled one.
+            self._closed = True
+        return True
+
+    def __enter__(self) -> "VolumeSession":
+        return self
+
+    def __exit__(self, exc_type, exc, tb) -> bool:
+        self.close()
+        return False
+
 def volume_io(volume_path: str, *, allow_open: bool = False):
```

`threading` joins the module imports; the file currently imports only `struct` and typing primitives, so the import is added once at the top rather than inside the function.

### 3.2 Handle creation, declaration and return (inside `volume_io`)

```diff
@@
     kernel32.DeviceIoControl.argtypes = [
         wintypes.HANDLE,
         wintypes.DWORD,
         wintypes.LPVOID,
         wintypes.DWORD,
         wintypes.LPVOID,
         wintypes.DWORD,
         ctypes.POINTER(wintypes.DWORD),
         wintypes.LPVOID,
     ]
+    kernel32.CloseHandle.restype = wintypes.BOOL
+    kernel32.CloseHandle.argtypes = [wintypes.HANDLE]
@@
     if handle is None or handle in (0, invalid_handle):
         raise OSError(ctypes.get_last_error(), f"CreateFileW failed for {volume_path}")
+
+    closed = {"done": False}
+
+    def close_handle() -> None:
+        # Idempotent at the Win32 layer as well as the session layer: a second
+        # CloseHandle on a stale numeric handle could release a recycled one.
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
+    try:
+        journal["id"] = query().journal_id
+    except OSError:
+        journal["id"] = 0
+    except BaseException:
+        # Anything unexpected still releases the handle before propagating.
+        close_handle()
+        raise
+    return VolumeSession(query, read, close_handle)
```

A failure *inside* `query()` while resolving the journal id is already caught (`OSError` → id 0), so the new `except BaseException` only covers a programming error or an interrupt. It closes and re-raises rather than swallowing.

### 3.3 Docstring update on `volume_io`

```diff
-    """Return (query_fn, read_fn) bound to a raw volume handle.
+    """Open a raw volume and return a `VolumeSession` bound to it.
+
+    The caller owns the handle: `close()` it, or use the session as a context
+    manager. Reads and queries after closure raise RuntimeError.
 
     This is the only code path that can open a volume. ...
     """
```

---

## 4. Offline test plan

Six new checks, appended to `eval/harness/usn_tests.py` in the existing "no live access" section, naming style matching `_z*`. No volume, no ioctl: every check constructs `VolumeSession` with fakes.

```python
@check("handle: a session closes on normal exit")
def _h1():
    calls = []
    with R.VolumeSession(lambda: "q", lambda a, b: b"r", lambda: calls.append("close")) as s:
        assert s.read(0, 1024) == b"r"
        assert s.query() == "q"
        assert not s.closed
    assert calls == ["close"], calls
    assert s.closed


@check("handle: close is idempotent and reports the first call only")
def _h2():
    calls = []
    s = R.VolumeSession(lambda: None, lambda a, b: b"", lambda: calls.append(1))
    assert s.close() is True
    assert s.close() is False
    assert s.close() is False
    assert calls == [1], calls


@check("handle: an exception inside the context still closes, once")
def _h3():
    calls = []
    try:
        with R.VolumeSession(lambda: None, lambda a, b: b"", lambda: calls.append(1)):
            raise ValueError("boom")
    except ValueError:
        pass
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


@check("handle: close waits for an in-flight read")
def _h5():
    import threading, time
    order = []
    release = threading.Event()

    def slow_read(start, size):
        order.append("read-start")
        release.wait(2.0)
        order.append("read-end")
        return b""

    s = R.VolumeSession(lambda: None, slow_read, lambda: order.append("close"))
    t = threading.Thread(target=lambda: s.read(0, 4096))
    t.start()
    time.sleep(0.05)
    closer = threading.Thread(target=s.close)
    closer.start()
    time.sleep(0.05)
    assert "close" not in order, "close ran while a read was in flight"
    release.set()
    t.join(2.0); closer.join(2.0)
    assert order == ["read-start", "read-end", "close"], order


@check("handle: CloseHandle is declared and confined to volume_io (static)")
def _h6():
    src = open(R.__file__, encoding="utf-8").read()
    for needle in ("CloseHandle.restype", "CloseHandle.argtypes"):
        assert needle in src, needle
    index = src.find("CloseHandle(handle)")
    assert index > 0
    gate = src.rfind("def volume_io", 0, index)
    assert gate > 0 and gate < index, "CloseHandle called outside volume_io"
```

**Expected outcome if applied:** 85 → **91 checks**, all passing; `validate` unchanged at 89/89; no live path exercised. `_h5` is the only threaded check and is bounded by a 2-second join, so a deadlock fails the suite rather than hanging it.

---

## 5. What the patch deliberately leaves alone

| Preserved | Why it matters |
|---|---|
| `Caps` (128 calls, 128 MiB requested, 4 MiB output, 1 MiB buffer) | the budget is a separate concern; touching it would make this patch unreviewable |
| Every `INCONCLUSIVE_*` outcome and `_exception_outcome` mapping | outcome semantics are frozen by the review record |
| `claim_r_prime`, `_check_retention`, `absence_means`, `completeness` | Claim R′/T′ boundaries |
| The `journal["id"] = 0` fallback (`L-3`) | separate limitation; changing two behaviours at once is how a bounded patch becomes an unbounded one |
| The `allow_open` gate and its ordering | `F-1` is about lifetime, not authorisation |
| Read, parse and V4 behaviour | untouched entirely |

**`L-3` is not fixed here.** It is recorded as a limitation in the final audit and would need its own proposal.

---

## 6. Risks of the patch

| Risk | Severity | Mitigation |
|---|---|---|
| A caller that unpacked the old 2-tuple would break | low — no in-repo caller does | stated in §2; the docstring change makes the new contract explicit |
| `close()` holding a lock means a hung ioctl blocks closure | inherent | the ioctl is synchronous and bounded by the call itself; a hung driver is not something a lock can fix |
| `CloseHandle` failure is swallowed by the session (flag already set) | low | `close_handle` raises `OSError`, which propagates from `close()`; the session still reports closed, because a retry could double-close |
| A future caller may simply never call `close()` | medium — the original defect, at one remove | context-manager form is the documented usage; a lint on the gate script can require it |

---

## 7. Claim ceiling

**Establishes:** a narrowly scoped, offline-testable design that gives the handle a deterministic lifetime, refuses use after closure, serialises closure against in-flight calls, and closes on every exit path — plus the six checks that would demonstrate it without a volume.

**Does not establish:** that `CloseHandle` succeeds on a real volume handle, that a live session behaves as the fakes do, or that the defect is fixed in practice. Nothing is applied; the suite still runs the unpatched code.
