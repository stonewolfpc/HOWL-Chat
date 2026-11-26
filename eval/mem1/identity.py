"""Frozen-identity enforcement.

G1 may not run on an apparatus that has drifted. Two registers are checked:

  * the eight frozen artifacts recorded in plan v3 §1.1 (prefix as recorded there);
  * the four frozen preregistration documents, by full hash, from the freeze record.

A mismatch raises `FrozenIdentityError`. `report()` returns the same facts as data
so the G1 receipt can quote them instead of restating them by hand.
"""
from __future__ import annotations

import hashlib
import pathlib

from .errors import FrozenIdentityError

REPO_ROOT = pathlib.Path(__file__).resolve().parents[2]

# plan v3 §1.1 — prefixes exactly as recorded (the plan records prefixes, so the
# check is a prefix check; full hashes are printed in the report).
FROZEN_ARTIFACTS: dict[str, str] = {
    "eval/s3/rerank-spec.json": "bbf400c225e6c90a",
    "eval/s3/rerank-spec-v1.2.json": "5414d0a7ced42386",
    "eval/s3/rerank_eval.py": "f59959f30ec336ae",
    "eval/s3/rerank_set.py": "23badd5e7825718d",
    "eval/s3/formation.py": "d8aac5003c66229d",
    "eval/s3/records.py": "1432fb63da5ae561",
    "eval/s3/evidence.py": "19db697a0f82d29f",
    "eval/s3/temporal-spec.json": "43d81e07e236df8c",
}

# Freeze record §1 — full hashes, because these four *are* the design.
FROZEN_PREREGISTRATION: dict[str, str] = {
    "docs/research/MEMORY-09-PREREGISTRATION-2026-11-15.md":
        "8c788a96ea6a866d4fd06e13496ddbfa46b649da5a5e0c8ef9bace960844ff12",
    "docs/research/MEMORY-09-PREREGISTRATION-AMENDMENT-1-2026-11-15.md":
        "975033401ae1beaa9caaa5d5ae55aab0ca3ddb4f5159eb41e9feb8e01a4c9b1a",
    "docs/research/MEMORY-09-PREREGISTRATION-AMENDMENT-2-2026-11-15.md":
        "c3ffd7ed735ed4cf68f8b455555af2a1d10eb38273b07a756f269076cd917f8f",
    "docs/research/MEMORY-09-PREREGISTRATION-AMENDMENT-3-2026-11-15.md":
        "c5af036265b03662b26feaf40d358ad08c3063fa0055c72d2fb0c2f41c236281",
}

# Plan v3 and the G0 receipt: preserved, and quoted in the G1 receipt.
PRESERVED_DOCUMENTS: dict[str, str] = {
    "docs/research/MEMORY-09-EXPERIMENT-PLAN-2026-11-15.md":
        "6c4078bd8c1fd29f1c74691afd0e25273000ac8bbddf6cc5c7f162bcfbe6bbb9",
    "docs/research/MEMORY-09-G0-PREFLIGHT-2026-11-15.md":
        "5c245283c1712ed59ade4ab51f6310b35b4ad428d6b446061c50b6f577f162d0",
}


def sha256_file(path: str | pathlib.Path) -> str:
    return hashlib.sha256(pathlib.Path(path).read_bytes()).hexdigest()


def _resolve(rel: str) -> pathlib.Path:
    return REPO_ROOT / rel


def check_artifacts(strict: bool = True) -> list[dict]:
    """Verify the eight frozen S3 artifacts by the prefixes the plan records."""
    results: list[dict] = []
    for rel, prefix in FROZEN_ARTIFACTS.items():
        path = _resolve(rel)
        if not path.exists():
            results.append({"path": rel, "expected_prefix": prefix, "actual": None,
                            "match": False, "reason": "missing"})
            continue
        actual = sha256_file(path)
        ok = actual.startswith(prefix)
        results.append({"path": rel, "expected_prefix": prefix, "actual": actual,
                        "match": ok, "reason": "" if ok else "sha256 prefix mismatch"})
    _raise_if_bad(results, strict, "frozen artifact")
    return results


def check_preregistration(strict: bool = True) -> list[dict]:
    """Verify the four frozen preregistration documents by full hash."""
    results: list[dict] = []
    for rel, expected in FROZEN_PREREGISTRATION.items():
        path = _resolve(rel)
        if not path.exists():
            results.append({"path": rel, "expected": expected, "actual": None,
                            "match": False, "reason": "missing"})
            continue
        actual = sha256_file(path)
        ok = actual == expected
        results.append({"path": rel, "expected": expected, "actual": actual,
                        "match": ok, "reason": "" if ok else "document edited after freeze"})
    _raise_if_bad(results, strict, "frozen preregistration")
    return results


def check_preserved(strict: bool = True) -> list[dict]:
    results: list[dict] = []
    for rel, expected in PRESERVED_DOCUMENTS.items():
        path = _resolve(rel)
        actual = sha256_file(path) if path.exists() else None
        ok = actual == expected
        results.append({"path": rel, "expected": expected, "actual": actual, "match": ok,
                        "reason": "" if ok else ("missing" if actual is None else "changed")})
    _raise_if_bad(results, strict, "preserved document")
    return results


def verify_all(strict: bool = True) -> dict:
    """All three registers, as one report."""
    return {
        "frozen_artifacts": check_artifacts(strict),
        "frozen_preregistration": check_preregistration(strict),
        "preserved_documents": check_preserved(strict),
    }


def _raise_if_bad(results: list[dict], strict: bool, what: str) -> None:
    bad = [r for r in results if not r["match"]]
    if bad and strict:
        detail = "; ".join(f"{r['path']}: {r['reason']}" for r in bad)
        raise FrozenIdentityError(f"{what} identity check failed: {detail}")


__all__ = [
    "REPO_ROOT", "FROZEN_ARTIFACTS", "FROZEN_PREREGISTRATION", "PRESERVED_DOCUMENTS",
    "sha256_file", "check_artifacts", "check_preregistration", "check_preserved", "verify_all",
]
