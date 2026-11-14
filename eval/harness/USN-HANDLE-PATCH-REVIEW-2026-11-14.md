# `USN-HANDLE-PATCH-REVIEW`

**Identity:** `MEMORY-09-USN-PATCH-REVIEW-1` · **Date:** 2026-11-14 · **Track:** `WO-NEXT-03` Track 4
**Subject:** `eval/harness/USN-HANDLE-LIFETIME-PATCH-PROPOSAL-2026-11-14.md` (`5e88c9b0…`, 12,296 B, 292 lines)
**Status:** reviewed, **not applied**. `usn_reader.py` (`2b384069…`) and `usn_tests.py` (`f6152707…`) are unchanged; the suite still runs 85/85 on the unpatched code.

**Verdict: the architecture is sound and the documented defect is correctly addressed. Three defects remain in the proposal itself — two in exception handling, one non-deterministic test. Approval is recommended for a narrowly scoped correction once those three are fixed.**

---

## 1. The nine checks

| # | Check | Verdict | Evidence |
|---:|---|---|---|
| 1 | Does the lock truly prevent closure during an active ioctl? | **PASS** | `read()` acquires `self._lock` and calls `self._read_fn(...)` inside the `with` block; `close()` must acquire the same lock before touching the flag or the handle. A `close()` that arrives mid-ioctl blocks until the read returns. There is no window in which `close_fn` runs while `read_fn` is executing |
| 2 | Does exception handling close exactly once? | **PASS on count, DEFECT on fidelity** | `_closed` is set under the lock before `close_fn`, and `close_handle` re-checks its own flag, so exactly one `CloseHandle`. But **P-1**: if `CloseHandle` fails inside the `except BaseException` path, the resulting `OSError` replaces the original exception |
| 3 | Can `BaseException` handling suppress or alter cancellation semantics? | **DEFECT P-1** | It does not *swallow* anything — the bare `raise` re-raises. It does (a) delay a `KeyboardInterrupt` until `CloseHandle` returns, which is acceptable, and (b) **replace** it when `close_handle()` raises, which is not: a cancellation would surface as an unrelated `OSError` |
| 4 | Does the context manager close on all exit paths? | **PASS, with P-2** | `__exit__` unconditionally calls `close()` and returns `False`, so normal return, exception and `return` inside the block all close. **P-2**: if `close()` raises while an exception is already propagating, `__exit__`'s own exception replaces it |
| 5 | Is repeated closure safe? | **PASS** | Second `close()` returns `False` without touching the OS; `close_handle` also guards independently, which protects against a future caller reaching the raw closer |
| 6 | Do reads after closure fail deterministically? | **PASS** | `_require_open()` raises `RuntimeError("volume session is closed")` before any call; test `_h4` covers both `query()` and `read()` |
| 7 | Can a blocked read make `close()` wait indefinitely? | **Yes — LIMITATION, not a defect** | `DeviceIoControl` is synchronous and has no timeout for a local volume; the journal's `Timeout`/`BytesToWaitFor` parameters govern journal waiting, not I/O completion. If the device hangs, the read never returns and `close()` waits behind the lock forever. This is inherent to synchronous I/O and is **not** fixable by the session; the mitigation is a caller-side watchdog, not a longer patch |
| 8 | Is the two-second threaded test deterministic, or merely timing-dependent? | **DEFECT P-3 — timing-dependent** | `_h5` uses `time.sleep(0.05)` twice and `join(2.0)`. It can **false-pass**: if the closer thread has not been scheduled by the time of the assertion, `"close" not in order` holds trivially. It will not false-fail on the lock. It is a bounded smoke test mislabelled as proof |
| 9 | Does the API stay usable by the existing caller without silently leaking? | **PASS for this repository / residual risk P-4** | `scan()` takes `query_fn`/`read_fn` as parameters and is untouched; the only `volume_io` call sites in the suite (`_z1`, `_z2`) assert refusals and never reach a session. But a caller that ignores the returned object still leaks — the patch converts an *unavoidable* leak into an *avoidable* one, and that must be stated as a contract |

---

## 2. Defects in the proposal

### P-1 — exception masking during construction

```python
    except BaseException:
        close_handle()     # raises OSError if CloseHandle fails
        raise              # never reached
```

A failed `CloseHandle` replaces the original exception. Corrected:

```python
    try:
        try:
            journal["id"] = query().journal_id
        except OSError:
            journal["id"] = 0
        return VolumeSession(query, read, close_handle)
    except BaseException:
        try:
            close_handle()
        except OSError:
            pass           # a failed release must not mask the real failure
        raise              # bare raise preserves the original exception object
```

### P-2 — exception masking on exit

`close()` raising from `__exit__` replaces an in-flight exception. Corrected: `close()` records rather than raises, and the failure is surfaced explicitly.

```python
    def close(self) -> bool:
        """Close the handle once. Records a release error instead of raising it."""
        with self._lock:
            if self._closed:
                return False
            self._closed = True
        try:
            self._close_fn()
        except OSError as exc:
            self.close_error = exc      # caller checks this; nothing is hidden
        return True
```

with `self.close_error = None` in `__init__`. `__exit__` then stays as written and cannot mask anything. The cost is that a silent failure becomes silent again — so the contract must say: **a caller that cares about release asserts `session.close_error is None`**, and the gate script should do exactly that.

### P-3 — non-deterministic test

Replace `_h5` with two checks: one that is genuinely scheduling-independent, one that is explicitly bounded and labelled non-proof.

```python
@check("handle: racing closures close exactly once (deterministic)")
def _h5a():
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
    assert calls == [1], calls            # true under any interleaving


@check("handle: close does not complete under an in-flight read (bounded)")
def _h5b():
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
    # The lock is held by the reader, so close cannot have finished yet.
    # A false PASS is possible if the closer has not reached close(); a false
    # FAIL is not. This is a smoke test, not a proof.
    release.set()
    t.join(5.0); closer.join(5.0)
    assert order == ["read-start", "read-end", "close"], order
```

`_h5a` is deterministic by construction: whichever thread wins the flag, exactly one `close_fn` runs, under every interleaving. `_h5b` is honest about what it can and cannot show.

### P-4 — the residual leak path

A session that is created and never closed still leaks, and nothing in the patch prevents it. Two mitigations are compatible with the narrow scope:

1. **Contract**, stated in the docstring and in the gate script: `volume_io` must be used as `with ... as session:`; a bare call is a review defect.
2. **Optional** `__del__` that emits a `ResourceWarning` and **never** closes — a warning cannot double-close or run a syscall during interpreter shutdown in a way that matters.

No `weakref.finalize` is recommended: it would close during shutdown, which is the failure mode the proposal correctly avoided.

---

## 3. What the patch does not change (verified against the proposal)

| Preserved | Confirmed |
|---|---|
| `Caps` — 128 calls, 128 MiB requested, 4 MiB output, 1 MiB buffer | no hunk touches them |
| All 17 `INCONCLUSIVE_*` outcomes and their mapping | untouched |
| `claim_r_prime`, `_check_retention`, `absence_means`, `completeness` | untouched |
| `allow_open` gate and check ordering | untouched |
| `journal["id"] = 0` fallback (`L-3`) | deliberately untouched, and stated as such |
| Read, parse, V4 and `SourceInfo` behaviour | untouched |

---

## 4. Recommendation

**Approve for the narrowly scoped correction, conditional on P-1, P-2 and P-3.**

| Condition | Why it is a precondition |
|---|---|
| P-1 corrected | an exception path that masks cancellation is worse than the leak it fixes |
| P-2 corrected | `__exit__` must not replace the exception it was invoked to clean up after |
| P-3 replaced | a test that can false-pass is not evidence that the lock works |
| P-4 documented | without it the patch's own defect (a handle nobody closes) reappears one level up |

With those, the patch is: one new class, one declaration pair, one return-shape change, six checks (85 → 91). It remains source-only, offline-tested, and irrelevant to live behaviour until a gate opens.

**`L-3` should not be folded in.** Two behaviour changes in one patch is how a reviewable correction becomes an unreviewable one.

---

## 5. Claim ceiling

**Establishes:** that the lock provides the closure-vs-read exclusion it claims; that closure happens exactly once; that the context manager covers every exit path; that reads after closure fail deterministically; and that three specific defects remain in the proposal — two exception-masking paths and one non-deterministic test.

**Does not establish:** that `CloseHandle` succeeds on a real volume handle, that a live session behaves as the fakes do, or that a hung device would not block `close()` indefinitely — it would.
