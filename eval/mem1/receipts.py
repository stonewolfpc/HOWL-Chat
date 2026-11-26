"""G1 receipt construction (user G1 §7).

A receipt that can be overwritten is not evidence, so writes are once-only and
carry a sidecar hash. The receipt separates three things the plan insists on
keeping apart: what was **observed**, what was **calculated**, and what was
**inferred**.
"""
from __future__ import annotations

import hashlib
import json
import os
import pathlib

from . import MEM1_VERSION
from .errors import ProtectedEvidenceError
from .identity import REPO_ROOT

PACKAGE = pathlib.Path(__file__).resolve().parent

# Immutable historical evidence. Never a write target, whatever an environment
# variable says: the only safe response to a request to regenerate one is to fail.
HISTORICAL = {
    "g1-receipt.json": "Stage 1 G1 receipt (96 checks, v1 semantics)",
    "g1-review-receipt.json": "historical G1 review, round 1",
    "g1-review-receipt-r2.json": "historical G1 review, round 2",
    "g4-preflight-receipt.json": "historical G4 preflight",
}

# The v2 (U1+U2) evidence set lives in its own directory so it can never be confused
# with, or overwrite, the Stage 1 artifacts beside it.
V2_DIR = PACKAGE / "v2"
V2_G1_RECEIPT = V2_DIR / "g1-receipt-u1u2-v2.json"
V2_REVIEW_RECEIPT = V2_DIR / "g1-review-receipt-u1u2-v2.json"

# This generation's evidence: the protected write path itself. New filenames, so the
# accepted v2 artifacts above remain byte-identical.
GUARD_DIR = PACKAGE / "write-guard-v1"
GUARD_G1_RECEIPT = GUARD_DIR / "g1-receipt-writeguard-v1.json"
GUARD_REVIEW_RECEIPT = GUARD_DIR / "g1-review-receipt-writeguard-v1.json"

# Coverage-correction generation: twelve protected artifacts, new directory.
GUARD2_DIR = PACKAGE / "write-guard-v2"
GUARD2_G1_RECEIPT = GUARD2_DIR / "g1-receipt-writeguard-v2.json"
GUARD2_REVIEW_RECEIPT = GUARD2_DIR / "g1-review-receipt-writeguard-v2.json"

# Accepted reference evidence: frozen, and now identity-protected exactly like the
# Stage 1 artifacts rather than merely refused by exclusive creation.
REFERENCE = {
    "v2/g1-receipt-u1u2-v2.json": "accepted U1+U2 v2 self-test receipt (118 checks)",
    "v2/g1-review-receipt-u1u2-v2.json": "accepted U1+U2 v2 review receipt (32 checks)",
}

# Recorded identity of all twelve protected artifacts: a generation states what it
# protects, and a review re-hashes exactly these files against these values.
REFERENCE_HASHES = {
    "eval/mem1/g1-receipt.json":
        "41cc95ff0ccee6ae573535b61d7ccf2a125974bc50164acaae6ae9e2d1e70d26",
    "eval/mem1/g1-receipt.json.sha256":
        "a930d039727819e15227317f6b8f13e9f7e7db57a11686a7596dce96ec9454aa",
    "eval/mem1/g1-review-receipt.json":
        "3a4aae7608703d3ec8c80095b1bfcf5da0fe6492c50d4722e8e25c7636e2b130",
    "eval/mem1/g1-review-receipt.json.sha256":
        "4d81510263d8db44f4899b353bd365929c549a89d7034e5c756ea3af5b640067",
    "eval/mem1/g1-review-receipt-r2.json":
        "3b69f1f9ab4ca03e52d05775244338404e4d6c6182fb22de0ad0d7f5f32b1e8d",
    "eval/mem1/g1-review-receipt-r2.json.sha256":
        "fae6e614b1e3afe1257ff474d0cfa9f17c3ce3931d2afe7c404ed13255fcacc6",
    "eval/mem1/g4-preflight-receipt.json":
        "1b6d8c0465fee21a9cacb4ff08f87169eaa69317447f6f781f8e3f7a5eb65d2e",
    "eval/mem1/g4-preflight-receipt.json.sha256":
        "19b6da2c9caa39880d8d01249bc277e244a1cfa3ce15db3e87934550c807f92e",
    "eval/mem1/v2/g1-receipt-u1u2-v2.json":
        "bc5766f975fb070421ebcb3fc587c4dde088ae94aa5c4019b6eed4aa08d53ad3",
    "eval/mem1/v2/g1-receipt-u1u2-v2.json.sha256":
        "64352aca15ba452f3c52f357e3284d14eaf4c2729b187ce2b792d3c6905f8138",
    "eval/mem1/v2/g1-review-receipt-u1u2-v2.json":
        "38cc18571dd50a12403d0a45c357630ba8bfbcd38fbb4a1509ed7e6b3fcb6df9",
    "eval/mem1/v2/g1-review-receipt-u1u2-v2.json.sha256":
        "46fc51afb0adbb61998bdbafecfa724159e1246a994d7600992826accf27803a",
}

_BINARY = getattr(os, "O_BINARY", 0)


def _protected_targets() -> list[pathlib.Path]:
    """Twelve protected files: four Stage 1 receipts + sidecars, two v2 + sidecars."""
    out: list[pathlib.Path] = []
    for name in HISTORICAL:
        base = PACKAGE / name
        out.extend([base, base.with_suffix(".json.sha256")])
    for rel in REFERENCE:
        base = PACKAGE / rel
        out.extend([base, base.with_suffix(".json.sha256")])
    return out


def _describe(path: pathlib.Path) -> str:
    """Human label for a protected path, used in refusal messages."""
    rel = path.relative_to(PACKAGE).as_posix()
    base = rel[:-7] if rel.endswith(".sha256") else rel
    if base in REFERENCE:
        return f"{REFERENCE[base]} (sidecar)" if rel.endswith(".sha256") else REFERENCE[base]
    if base in HISTORICAL:
        return f"{HISTORICAL[base]} (sidecar)" if rel.endswith(".sha256") else HISTORICAL[base]
    return rel


def _file_identity(fd: int) -> tuple[int, int]:
    """(device/volume, file id) of an open handle. Fails closed if unavailable.

    On NTFS `st_dev` is the volume serial number and `st_ino` the file id, so the
    pair distinguishes files that a path comparison cannot. Filesystems that report
    zero for either value cannot support this check, and refuse rather than guess.
    """
    st = os.fstat(fd)
    ident = (int(getattr(st, "st_dev", 0)), int(getattr(st, "st_ino", 0)))
    if not ident[0] or not ident[1]:
        raise ProtectedEvidenceError(
            "file identity is unavailable (st_dev/st_ino empty); refusing to write")
    return ident


def _open_protected() -> tuple[set[tuple[int, int]], list[int]]:
    """Identities of the protected evidence, read from open handles.

    The handles are held open for the whole write: on Windows an open handle without
    delete sharing also blocks replacement of the protected file.
    """
    ids: set[tuple[int, int]] = set()
    handles: list[int] = []
    try:
        for path in _protected_targets():
            try:
                fd = os.open(str(path), os.O_RDONLY | _BINARY)
            except OSError as exc:
                raise ProtectedEvidenceError(
                    f"protected evidence {path.name} cannot be opened: {exc}") from exc
            handles.append(fd)
            ids.add(_file_identity(fd))
    except BaseException:
        for fd in handles:
            os.close(fd)
        raise
    return ids, handles


def refuse_historical(path) -> pathlib.Path:
    """Supplemental name/path gate; file identity is the security boundary.

    Catches the cheap cases before a handle exists: case-folded names (Windows\
    aliases), traversal, absolute aliases and symlinked parents. It is deliberately
    not sufficient on its own -- a hardlink or a short-name alias passes here and is
    caught by the handle-identity check that follows exclusive creation.
    """
    target = pathlib.Path(path)
    resolved = pathlib.Path(os.path.realpath(target))
    package = pathlib.Path(os.path.realpath(PACKAGE))
    # The name test applies only inside the package directory: a file called
    # g1-receipt.json in a temp directory is not protected evidence, and refusing it
    # would break the fresh-run check. Resolution catches aliases wherever they live.
    names = {p.name.casefold() for p in _protected_targets()}
    for protected in _protected_targets():
        if resolved == pathlib.Path(os.path.realpath(protected)):
            raise ProtectedEvidenceError(
                f"{target} resolves to protected reference evidence "
                f"({_describe(protected)}); generate replacement evidence in a new "
                f"generation directory, never in {V2_DIR.name}/")
    if resolved.parent == package and resolved.name.casefold() in names:
        raise ProtectedEvidenceError(
            f"{target} names protected reference evidence in {PACKAGE.name}/; "
            "generate replacement evidence in a new generation directory")
    return target


def _write_all(fd: int, raw: bytes) -> None:
    view = memoryview(raw)
    while view:
        written = os.write(fd, view)
        if written <= 0:
            raise OSError("short write while writing receipt")
        view = view[written:]


def _cleanup_created(target: pathlib.Path, ident: tuple[int, int] | None) -> None:
    """Remove a file this call created, and only if its identity still matches.

    Never unlinks by name alone: a name match after a race could delete protected
    evidence. When identity is unknown the empty artifact is left in place.
    """
    if ident is None:
        return
    try:
        st = os.stat(target)
    except OSError:
        return
    if (int(st.st_dev), int(st.st_ino)) != ident:
        return
    try:
        os.unlink(target)
    except OSError:
        pass


def _write_bytes_exclusive(target: pathlib.Path, raw: bytes,
                           protected: set[tuple[int, int]]) -> None:
    """Exclusive create, then identity of the opened handle, then the first byte.

    Order matters: a path-only check leaves a check-then-truncate window, and an
    identity check after writing would already have corrupted the protected file.
    """
    flags = os.O_WRONLY | os.O_CREAT | os.O_EXCL | _BINARY
    fd = os.open(str(target), flags, 0o644)
    ident: tuple[int, int] | None = None
    try:
        ident = _file_identity(fd)
        if ident in protected:
            raise ProtectedEvidenceError(
                f"{target} resolves to protected evidence by file identity; refusing "
                "before writing any bytes")
        _write_all(fd, raw)
    except BaseException:
        try:
            os.close(fd)
        finally:
            _cleanup_created(target, ident)
        raise
    try:
        os.close(fd)
    except OSError as exc:
        _cleanup_created(target, ident)
        raise ProtectedEvidenceError(f"close failed for {target}: {exc}") from exc


def write_once(path: str | pathlib.Path, payload: dict) -> pathlib.Path:
    """Write a receipt and its sidecar exactly once; never touch protected evidence.

    An existing destination is never opened for truncation or replacement -- not
    even a previous generation's output. A new generation must use a new path.
    """
    target = refuse_historical(path)
    target.parent.mkdir(parents=True, exist_ok=True)
    body = json.dumps(payload, indent=2, sort_keys=True, ensure_ascii=False) + "\n"
    raw = body.encode("utf-8")
    protected, handles = _open_protected()
    try:
        _write_bytes_exclusive(target, raw, protected)
        digest = hashlib.sha256(raw).hexdigest()
        sidecar = target.with_suffix(target.suffix + ".sha256")
        _write_bytes_exclusive(
            sidecar, f"{digest}  {target.name}\n".encode("utf-8"), protected)
    finally:
        for fd in handles:
            os.close(fd)
    return target


def source_hashes() -> dict[str, str]:
    out: dict[str, str] = {}
    for path in sorted(PACKAGE.glob("*.py")):
        out[path.name] = hashlib.sha256(path.read_bytes()).hexdigest()
    return out


def package_sha256() -> str:
    blob = "\n".join(f"{name}|{digest}" for name, digest in sorted(source_hashes().items()))
    return hashlib.sha256(blob.encode("utf-8")).hexdigest()


COMPONENTS = (
    ("model.py", "A' model: records, acquisition events, knowledge edges, supersession, revocation"),
    ("form.py", "deterministic formation from genuine spans; no model, RQ-MF-01 untouched"),
    ("retrieve.py", "context-explicit retrieval; eligibility before ranking; fail-closed ownership"),
    ("assemble.py", "wrapper, citation and shape signature (placebo structural fidelity)"),
    ("tokens.py", "measured per-turn accounting; memory competes with history"),
    ("placebo.py", "declared decoy pool; per-turn count/token/shape matching; irrelevance proof"),
    ("conversation.py", "long-conversation driver; arms A0/A1/A2/C/D; run plan for N=16"),
    ("probe.py", "probe quota enforcement, span-departure assertion, receipts"),
    ("score.py", "M1 and M2 definitions; A0 -> N/A; temporal veto; composite score refused"),
    ("stats.py", "C1/C2, exact sign test, alpha=.025, ties, exclusions, median, 4/39 MMD"),
    ("stopping.py", "fixed-sample stopping; outcome-blind instrument halts only"),
    ("seeds.py", "seed derivation and mechanical cross-arm matching"),
    ("distribution.py", "declared-vs-assigned scope report; outside the probe receipt"),
    ("identity.py", "frozen artifact and frozen preregistration verification"),
    ("receipts.py", "this receipt; exclusive create + handle identity + write-once"),
    ("writeguard.py", "adversarial tests for the protected write path"),
    ("selftest.py", "zero-model self-check and negative tests"),
    ("fixtures.py", "synthetic G1 fixtures (not experimental cases)"),
)


def build_g1_receipt(*, checks: list[dict], identity: dict, dependencies: list[dict],
                     acceptance: list[dict], confirmations: dict,
                     limitations: list[str], environment: dict,
                     observed: dict, calculated: dict,
                     evidence: dict | None = None) -> dict:
    return {
        "gate": "G1",
        "suite": "mem1.selftest",
        "instrument_version": MEM1_VERSION,
        "evidence": evidence or {},
        "status": "PASS" if all(c["passed"] for c in checks) and not any(
            a["result"] == "FAIL" for a in acceptance) else "FAIL",
        "components": [{"file": name, "responsibility": desc} for name, desc in COMPONENTS],
        "source_sha256": source_hashes(),
        "package_sha256": package_sha256(),
        "dependencies": dependencies,
        "frozen_identity": identity,
        "acceptance_criteria": acceptance,
        "checks_total": len(checks),
        "checks_passed": sum(1 for c in checks if c["passed"]),
        "checks": checks,
        "negative_checks": [c for c in checks if c.get("negative")],
        "observed": observed,
        "calculated": calculated,
        "inferred": {
            "note": "G1 makes no claim about capability; no result exists to interpret",
            "g4_executed": False,
        },
        "confirmations": confirmations,
        "limitations": limitations,
        "environment": environment,
    }


__all__ = ["write_once", "source_hashes", "package_sha256", "COMPONENTS", "build_g1_receipt",
           "REPO_ROOT", "HISTORICAL", "REFERENCE", "REFERENCE_HASHES", "V2_DIR",
           "V2_G1_RECEIPT", "V2_REVIEW_RECEIPT", "GUARD_DIR", "GUARD_G1_RECEIPT",
           "GUARD_REVIEW_RECEIPT", "GUARD2_DIR", "GUARD2_G1_RECEIPT", "GUARD2_REVIEW_RECEIPT",
           "refuse_historical", "_file_identity", "_open_protected", "_write_all",
           "_protected_targets", "_describe"]
