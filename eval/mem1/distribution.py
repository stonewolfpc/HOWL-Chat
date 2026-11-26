"""Declared expectation versus assigned scope (Stage 1, gap G4).

Reporting only. This module reads a finished store and a list of declared
expectations and counts what agrees; it never forms, selects or scores anything.
It is deliberately outside the probe receipt: the report is an artefact derived
from the store alone, so adding it changes no receipt schema and no receipt hash.

`ExpectedScope` is a case declaration, not an inference. The expected scope of an
event comes from the case design, never from the record's own content — a report
that derived the expectation from the record would always agree with itself.

`report` is reproducible by construction: given the same store and the same
expectations it returns the same table. It reports three things separately:

  * distribution — how many records exist per assigned scope;
  * agreement — declared versus assigned over the declared events only;
  * coverage — declared events with no record, and records with no declaration.

Coverage is reported rather than hidden: a silent mismatch between the declared
event inventory and what formation produced is exactly what this table exists to
make visible.
"""
from __future__ import annotations

from dataclasses import asdict, dataclass

from .errors import DesignViolationError
from .model import SCOPES


@dataclass(frozen=True)
class ExpectedScope:
    """One declared event: where it was said, and the scope it should land in."""

    turn: int
    subject: str
    predicate: str
    scope: str

    def __post_init__(self) -> None:
        if self.scope not in SCOPES:
            raise DesignViolationError(
                f"expected scope {self.scope!r} is not one of {SCOPES}")
        if self.turn < 1:
            raise DesignViolationError(f"expected event turn must be >= 1, got {self.turn}")
        if not self.subject or not self.predicate:
            raise DesignViolationError("expected event needs a subject and a predicate")

    def as_dict(self) -> dict:
        return asdict(self)


def _key(turn: int, subject: str, predicate: str) -> tuple[int, str, str]:
    return (turn, subject, predicate)


def report(store, expectations=()) -> dict:
    """Declared-versus-assigned scope table. Pure function of store + declarations."""
    declared: dict[tuple[int, str, str], ExpectedScope] = {}
    for expectation in expectations:
        key = _key(expectation.turn, expectation.subject, expectation.predicate)
        if key in declared:
            raise DesignViolationError(f"duplicate expectation for {key}")
        declared[key] = expectation

    by_key: dict[tuple[int, str, str], list] = {}
    per_scope = {scope: 0 for scope in SCOPES}
    for rec in store.all():
        per_scope[rec.scope] = per_scope.get(rec.scope, 0) + 1
        by_key.setdefault(_key(rec.established_turn, rec.subject, rec.predicate),
                          []).append(rec)

    confusion: dict[str, int] = {}
    unmatched: list[dict] = []
    matched = 0
    mismatched = 0
    mismatches: list[dict] = []

    for key in sorted(declared):
        expectation = declared[key]
        hits = by_key.get(key, [])
        if not hits:
            unmatched.append({"turn": key[0], "subject": key[1], "predicate": key[2],
                              "expected_scope": expectation.scope})
            continue
        for rec in hits:
            matched += 1
            bucket = f"{expectation.scope}->{rec.scope}"
            confusion[bucket] = confusion.get(bucket, 0) + 1
            if rec.scope != expectation.scope:
                mismatched += 1
                mismatches.append({"record_id": rec.record_id, "turn": key[0],
                                   "subject": key[1], "predicate": key[2],
                                   "expected_scope": expectation.scope,
                                   "assigned_scope": rec.scope})

    undeclared = [rec.record_id for rec in store.all()
                  if _key(rec.established_turn, rec.subject, rec.predicate) not in declared]

    return {
        "records_total": len(store.all()),
        "per_scope": per_scope,
        "declared_total": len(declared),
        "matched_records": matched,
        "mismatched_records": mismatched,
        "agreement_rate": ((matched - mismatched) / matched) if matched else None,
        "confusion": dict(sorted(confusion.items())),
        "mismatches": mismatches,
        "unmatched_expectations": unmatched,
        "undeclared_record_ids": undeclared,
        "declared_scopes": sorted({e.scope for e in expectations}),
        "declarations": [declared[k].as_dict() for k in sorted(declared)],
    }


__all__ = ["ExpectedScope", "report"]
