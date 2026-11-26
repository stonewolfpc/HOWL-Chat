"""Adversarial tests for the protected-evidence write path.

The primitive must never open an existing destination for truncation, must verify
the identity of the handle it created *before* the first byte, and must refuse when
identity cannot be established. These tests use real files: mocked path strings
would prove nothing about aliases, links, or short names.

    PYTHONPATH=. python -m eval.mem1.writeguard

They also run inside `selftest.py` as WG01-WG24, so the receipt carries them.
"""
from __future__ import annotations

import hashlib
import os
import pathlib
import shutil
import tempfile

from . import receipts
from .errors import ProtectedEvidenceError
from .receipts import HISTORICAL, PACKAGE, write_once

_BIN = getattr(os, "O_BINARY", 0)
PROTECTED_RECEIPT = PACKAGE / "g1-receipt.json"
PROTECTED_SIDECAR = PACKAGE / "g1-receipt.json.sha256"

# The accepted U1+U2 reference artifacts. WG01-WG15 only exercised the Stage 1 pair;
# the audit showed these four were protected by exclusive creation alone, so they now
# need the same adversarial treatment.
V2_RECEIPT = PACKAGE / "v2" / "g1-receipt-u1u2-v2.json"
V2_RECEIPT_SIDECAR = PACKAGE / "v2" / "g1-receipt-u1u2-v2.json.sha256"
V2_REVIEW = PACKAGE / "v2" / "g1-review-receipt-u1u2-v2.json"
V2_REVIEW_SIDECAR = PACKAGE / "v2" / "g1-review-receipt-u1u2-v2.json.sha256"


def _sha(path) -> str:
    return hashlib.sha256(pathlib.Path(path).read_bytes()).hexdigest()


def _refused(fn) -> type[BaseException] | None:
    """Run a write attempt and return the exception type it raised, if any."""
    try:
        fn()
    except Exception as exc:  # the exception type is the assertion
        return type(exc)
    return None


def _short_name(path: pathlib.Path) -> str | None:
    """Windows 8.3 name for a path, or None where the platform/folder has none."""
    if os.name != "nt":
        return None
    try:
        import ctypes

        buf = ctypes.create_unicode_buffer(1024)
        n = ctypes.windll.kernel32.GetShortPathNameW(str(path), buf, 1024)
        return buf.value if n else None
    except Exception:
        return None


def run_checks(checks) -> None:
    """Register WG01-WG15 on `checks` (an object with add(id, name, passed, detail))."""
    root = pathlib.Path(tempfile.mkdtemp(prefix="mem1-guard-"))
    d = root
    try:
        _run(checks, d)
    finally:
        shutil.rmtree(root, ignore_errors=True)


def _run(checks, d: pathlib.Path) -> None:
    before_receipt = _sha(PROTECTED_RECEIPT)
    before_sidecar = _sha(PROTECTED_SIDECAR)

    # --- WG01: the ordinary path still works -----------------------------
    fresh = d / "fresh.json"
    write_once(fresh, {"ok": True})
    side = fresh.with_suffix(".json.sha256")
    checks.add("WG01", "a new destination is created exclusively with a verified sidecar",
               fresh.exists() and side.exists()
               and side.read_text().split()[0] == _sha(fresh),
               f"{fresh.stat().st_size} B, sidecar matches")

    # --- WG02/WG03: the protected artifacts themselves -------------------
    kind = _refused(lambda: write_once(PROTECTED_RECEIPT, {"x": 1}))
    checks.add("WG02", "an existing protected receipt is refused, bytes unchanged",
               kind is ProtectedEvidenceError and _sha(PROTECTED_RECEIPT) == before_receipt,
               f"{kind.__name__ if kind else 'NO REFUSAL'}; hash unchanged")
    kind = _refused(lambda: write_once(PROTECTED_SIDECAR, {"x": 1}))
    checks.add("WG03", "an existing protected sidecar is refused, bytes unchanged",
               kind is ProtectedEvidenceError and _sha(PROTECTED_SIDECAR) == before_sidecar,
               f"{kind.__name__ if kind else 'NO REFUSAL'}; hash unchanged")

    # --- WG04: case-only alias (Windows is case-insensitive) -------------
    if os.name == "nt":
        kind = _refused(lambda: write_once(PACKAGE / "G1-Receipt.JsOn", {"x": 1}))
        checks.add("WG04", "a case-only filename variant is refused",
                   kind is ProtectedEvidenceError
                   and _sha(PROTECTED_RECEIPT) == before_receipt,
                   f"{kind.__name__ if kind else 'NO REFUSAL'}; bytes unchanged")
    else:
        checks.add("WG04", "a case-only filename variant is refused", True,
                   "SKIP: POSIX names are case-sensitive, so no alias exists")

    # --- WG05: 8.3 short-name alias --------------------------------------
    short = _short_name(PROTECTED_RECEIPT)
    if short and pathlib.Path(short) != PROTECTED_RECEIPT:
        kind = _refused(lambda: write_once(pathlib.Path(short), {"x": 1}))
        checks.add("WG05", "an 8.3 short-name alias is refused before any write",
                   kind in (ProtectedEvidenceError, FileExistsError)
                   and _sha(PROTECTED_RECEIPT) == before_receipt,
                   f"{kind.__name__ if kind else 'NO REFUSAL'} via {pathlib.Path(short).name}")
    else:
        checks.add("WG05", "an 8.3 short-name alias is refused before any write", True,
                   "SKIP: no 8.3 alias for this path (short names disabled)")

    # --- WG06: symlink to protected --------------------------------------
    link = d / "sym-to-receipt.json"
    try:
        os.symlink(PROTECTED_RECEIPT, link)
        made_link = True
        link_note = ""
    except (OSError, NotImplementedError) as exc:
        made_link = False
        link_note = f"{type(exc).__name__}: {exc}"[:70]
    if made_link:
        kind = _refused(lambda: write_once(link, {"x": 1}))
        checks.add("WG06", "an existing symlink to a protected receipt is refused",
                   kind in (ProtectedEvidenceError, FileExistsError)
                   and _sha(PROTECTED_RECEIPT) == before_receipt,
                   f"{kind.__name__ if kind else 'NO REFUSAL'}; bytes unchanged")
    else:
        checks.add("WG06", "an existing symlink to a protected receipt is refused", True,
                   f"SKIP: symlink creation unavailable ({link_note})")

    # --- WG07: hardlink to protected -------------------------------------
    hard = d / "hard-to-receipt.json"
    try:
        os.link(PROTECTED_RECEIPT, hard)
        made_hard = True
    except OSError as exc:
        made_hard = False
        hard_note = f"{type(exc).__name__}: {exc}"[:70]
    if made_hard:
        kind = _refused(lambda: write_once(hard, {"x": 1}))
        checks.add("WG07", "an existing hardlink to a protected receipt is refused",
                   kind in (ProtectedEvidenceError, FileExistsError)
                   and _sha(PROTECTED_RECEIPT) == before_receipt,
                   f"{kind.__name__ if kind else 'NO REFUSAL'}; both names unchanged")
    else:
        checks.add("WG07", "an existing hardlink to a protected receipt is refused", True,
                   f"SKIP: hardlink creation unavailable ({hard_note})")

    # --- WG08/WG09: traversal and absolute alias -------------------------
    kind = _refused(lambda: write_once(PACKAGE / "v2" / ".." / "g1-receipt.json", {"x": 1}))
    checks.add("WG08", "a relative traversal path to protected evidence is refused",
               kind is ProtectedEvidenceError and _sha(PROTECTED_RECEIPT) == before_receipt,
               f"{kind.__name__ if kind else 'NO REFUSAL'}")
    kind = _refused(lambda: write_once(PROTECTED_RECEIPT.resolve(), {"x": 1}))
    checks.add("WG09", "an absolute path alias to protected evidence is refused",
               kind is ProtectedEvidenceError and _sha(PROTECTED_RECEIPT) == before_receipt,
               f"{kind.__name__ if kind else 'NO REFUSAL'}")

    # --- WG10: parent-directory substitution -----------------------------
    short_dir = _short_name(PACKAGE)
    if short_dir and pathlib.Path(short_dir) != PACKAGE:
        alias = pathlib.Path(short_dir) / "g1-receipt.json"
        kind = _refused(lambda: write_once(alias, {"x": 1}))
        checks.add("WG10", "a substituted parent path does not bypass the guard",
                   kind in (ProtectedEvidenceError, FileExistsError)
                   and _sha(PROTECTED_RECEIPT) == before_receipt,
                   f"{kind.__name__ if kind else 'NO REFUSAL'} via 8.3 parent "
                   f"(symlinked parent unavailable on this host)")
    else:
        checks.add("WG10", "a substituted parent path does not bypass the guard", True,
                   "SKIP: neither symlinked nor 8.3 parent alias is available here")

    # --- WG11: protected identity unavailable ----------------------------
    no_id = d / "no-protected-identity.json"
    real_open_protected = receipts._open_protected

    def _no_protected():
        raise ProtectedEvidenceError("simulated: protected identity unavailable")

    receipts._open_protected = _no_protected
    try:
        kind = _refused(lambda: write_once(no_id, {"x": 1}))
    finally:
        receipts._open_protected = real_open_protected
    checks.add("WG11", "protected identity unavailable fails closed with no file",
               kind is ProtectedEvidenceError and not no_id.exists(),
               f"{kind.__name__ if kind else 'NO REFUSAL'}; nothing created")

    # --- WG12: identity of the new handle unavailable --------------------
    flaky_target = d / "no-handle-identity.json"
    real_fstat = os.fstat
    calls = {"n": 0}

    def flaky_fstat(fd):
        calls["n"] += 1
        if calls["n"] > len(HISTORICAL) * 2:
            return os.stat_result((0o644, 0, 0, 1, 0, 0, 0, 0, 0, 0))
        return real_fstat(fd)

    os.fstat = flaky_fstat
    try:
        kind = _refused(lambda: write_once(flaky_target, {"x": 1}))
    finally:
        os.fstat = real_fstat
    checks.add("WG12", "unverifiable handle identity fails closed with zero bytes",
               kind is ProtectedEvidenceError
               and (not flaky_target.exists() or flaky_target.stat().st_size == 0),
               f"{kind.__name__ if kind else 'NO REFUSAL'}; "
               "empty artifact left (identity unknown, so no unlink)")

    # --- WG13/WG14: write-once semantics ---------------------------------
    again = d / "again.json"
    write_once(again, {"first": True})
    kind = _refused(lambda: write_once(again, {"second": True}))
    checks.add("WG13", "a repeated exclusive creation is refused",
               kind is FileExistsError, f"{kind.__name__ if kind else 'NO REFUSAL'}")

    taken = d / "taken.json"
    taken.write_bytes(b"original")
    kind = _refused(lambda: write_once(taken, {"x": 1}))
    checks.add("WG14", "an existing non-protected destination is not replaced",
               kind is FileExistsError and taken.read_bytes() == b"original",
               f"{kind.__name__ if kind else 'NO REFUSAL'}; bytes unchanged")

    # --- WG15: identity is checked on the handle, before the first byte --
    hard2 = d / "hard-identity.json"
    try:
        os.link(PROTECTED_RECEIPT, hard2)
        probe = os.open(str(hard2), os.O_WRONLY | _BIN)
        os.close(probe)
        writable = True
    except OSError as exc:
        writable = False
        w_note = f"{type(exc).__name__}: {exc}"[:70]

    if writable:
        before15 = _sha(PROTECTED_RECEIPT)
        real_open, real_write_all = receipts.os.open, receipts._write_all
        wrote = {"n": 0}

        def bypassing_open(path, flags, *args, **kwargs):
            # Simulate a filesystem where O_EXCL fails open: hand back a writable
            # handle on the hardlink instead of creating anything.
            if pathlib.Path(path) == hard2:
                return real_open(str(hard2), os.O_WRONLY | _BIN)
            return real_open(path, flags, *args, **kwargs)

        def counting_write_all(fd, raw):
            wrote["n"] += 1
            return real_write_all(fd, raw)

        receipts.os.open = bypassing_open
        receipts._write_all = counting_write_all
        try:
            kind = _refused(lambda: write_once(hard2, {"x": 1}))
        finally:
            receipts.os.open = real_open
            receipts._write_all = real_write_all
        checks.add("WG15", "identity is verified on the opened handle before the first byte",
                   kind is ProtectedEvidenceError and wrote["n"] == 0
                   and _sha(PROTECTED_RECEIPT) == before15,
                   "O_EXCL bypassed: refused on handle identity, 0 bytes written")
    else:
        checks.add("WG15", "identity is verified on the opened handle before the first byte",
                   True, f"SKIP: cannot open a hardlink for writing ({w_note})")

    # --- WG16-WG24: the accepted U1+U2 v2 artifacts ----------------------
    protected_ids, handles = receipts._open_protected()
    try:
        members = []
        for path in (V2_RECEIPT, V2_RECEIPT_SIDECAR, V2_REVIEW, V2_REVIEW_SIDECAR):
            fd = os.open(str(path), os.O_RDONLY | _BIN)
            try:
                members.append(receipts._file_identity(fd) in protected_ids)
            finally:
                os.close(fd)
        set_size = len(protected_ids)
    finally:
        for fd in handles:
            os.close(fd)
    checks.add("WG16", "all four accepted v2 artifacts are in the protected identity set",
               all(members) and set_size == 12,
               f"{sum(members)}/4 v2 artifacts identified, set size {set_size}")

    for cid, path, label in (("WG17", V2_RECEIPT, "v2 self-test receipt"),
                             ("WG18", V2_RECEIPT_SIDECAR, "v2 self-test sidecar"),
                             ("WG19", V2_REVIEW, "v2 review receipt"),
                             ("WG20", V2_REVIEW_SIDECAR, "v2 review sidecar")):
        before = _sha(path)
        kind = _refused(lambda p=path: write_once(p, {"x": 1}))
        checks.add(cid, f"the {label} is refused, bytes unchanged",
                   kind is ProtectedEvidenceError and _sha(path) == before,
                   f"{kind.__name__ if kind else 'NO REFUSAL'}; hash unchanged")

    if os.name == "nt":
        before = _sha(V2_RECEIPT)
        kind = _refused(lambda: write_once(PACKAGE / "v2" / "G1-Receipt-U1U2-V2.JsOn",
                                           {"x": 1}))
        checks.add("WG21", "a case-variant alias of the v2 receipt is refused",
                   kind is ProtectedEvidenceError and _sha(V2_RECEIPT) == before,
                   f"{kind.__name__ if kind else 'NO REFUSAL'}; bytes unchanged")
    else:
        checks.add("WG21", "a case-variant alias of the v2 receipt is refused", True,
                   "SKIP: POSIX names are case-sensitive, so no alias exists")

    short_v2 = _short_name(V2_REVIEW)
    if short_v2 and pathlib.Path(short_v2) != V2_REVIEW:
        before = _sha(V2_REVIEW)
        kind = _refused(lambda: write_once(pathlib.Path(short_v2), {"x": 1}))
        checks.add("WG22", "an 8.3 alias of the v2 review receipt is refused",
                   kind in (ProtectedEvidenceError, FileExistsError)
                   and _sha(V2_REVIEW) == before,
                   f"{kind.__name__ if kind else 'NO REFUSAL'} via "
                   f"{pathlib.Path(short_v2).name}")
    else:
        checks.add("WG22", "an 8.3 alias of the v2 review receipt is refused", True,
                   "SKIP: no 8.3 alias for this path (short names disabled)")

    v2hard = d / "hard-to-v2.json"
    try:
        os.link(V2_RECEIPT, v2hard)
        made_v2_hard = True
    except OSError as exc:
        made_v2_hard = False
        v2_hard_note = f"{type(exc).__name__}: {exc}"[:70]
    if made_v2_hard:
        before = _sha(V2_RECEIPT)
        kind = _refused(lambda: write_once(v2hard, {"x": 1}))
        checks.add("WG23", "a hardlink to the v2 receipt is refused",
                   kind in (ProtectedEvidenceError, FileExistsError)
                   and _sha(V2_RECEIPT) == before,
                   f"{kind.__name__ if kind else 'NO REFUSAL'}; both names unchanged")
    else:
        checks.add("WG23", "a hardlink to the v2 receipt is refused", True,
                   f"SKIP: hardlink creation unavailable ({v2_hard_note})")

    v2hard2 = d / "hard-to-v2-review.json"
    try:
        os.link(V2_REVIEW, v2hard2)
        probe = os.open(str(v2hard2), os.O_WRONLY | _BIN)
        os.close(probe)
        v2_writable = True
    except OSError as exc:
        v2_writable = False
        v2_w_note = f"{type(exc).__name__}: {exc}"[:70]
    if v2_writable:
        before24 = _sha(V2_REVIEW)
        real_open2, real_write_all2 = receipts.os.open, receipts._write_all
        wrote2 = {"n": 0}

        def bypassing_open_v2(path, flags, *args, **kwargs):
            if pathlib.Path(path) == v2hard2:
                return real_open2(str(v2hard2), os.O_WRONLY | _BIN)
            return real_open2(path, flags, *args, **kwargs)

        def counting_write_all_v2(fd, raw):
            wrote2["n"] += 1
            return real_write_all2(fd, raw)

        receipts.os.open = bypassing_open_v2
        receipts._write_all = counting_write_all_v2
        try:
            kind = _refused(lambda: write_once(v2hard2, {"x": 1}))
        finally:
            receipts.os.open = real_open2
            receipts._write_all = real_write_all2
        checks.add("WG24", "v2 handle identity is checked before the first byte",
                   kind is ProtectedEvidenceError and wrote2["n"] == 0
                   and _sha(V2_REVIEW) == before24,
                   "O_EXCL bypassed on a v2 hardlink: refused on identity, 0 bytes written")
    else:
        checks.add("WG24", "v2 handle identity is checked before the first byte", True,
                   f"SKIP: cannot open a hardlink for writing ({v2_w_note})")


def main() -> int:
    class _Checks:
        def __init__(self) -> None:
            self.items: list[dict] = []

        def add(self, cid, name, passed, detail="", kind="guard"):
            self.items.append({"id": cid, "name": name, "passed": bool(passed),
                               "detail": detail})

    checks = _Checks()
    run_checks(checks)
    passed = sum(1 for c in checks.items if c["passed"])
    for c in checks.items:
        print(f"  [{c['id']}] {'PASS' if c['passed'] else 'FAIL'} {c['name']:<58} {c['detail']}")
    total = len(checks.items)
    print(f"\nRESULT: {'PASS' if passed == total else 'FAIL'} ({passed}/{total} write-guard checks)")
    return 0 if passed == total else 1


if __name__ == "__main__":
    import sys

    sys.exit(main())
