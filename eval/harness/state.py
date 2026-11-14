"""Memory-state interface (row 478 §2-3).

Rules enforced here:
  * every mutation returns a hash;
  * snapshots are immutable and hashed canonically;
  * no silent normalisation, deduplication or fuzzy matching;
  * "deletion" marks removal with a turn, it never erases.
"""

from __future__ import annotations

import copy
import json
from dataclasses import dataclass
from typing import Any

from .adapter import canonical, sha256_text

ENTITY_KINDS = (
    "character",
    "persona",
    "world",
    "fact",
    "item",
    "equipment",
    "relationship",
    "event",
)

FACT_SOURCES = ("authored", "extracted_from_turn", "owner_edit")


class StateError(ValueError):
    """Raised on a malformed state operation."""


@dataclass(frozen=True)
class Snapshot:
    """Immutable: holds canonical JSON, not a mutable dict.

    A frozen dataclass around a dict is not immutable — the dict can still be
    mutated, which would let evidence change after it was hashed. Storing the
    serialised form and parsing on read closes that.
    """

    _json: str
    sha256: str

    @classmethod
    def from_data(cls, data: dict) -> "Snapshot":
        payload = canonical(data)
        return cls(_json=payload, sha256=sha256_text(payload))

    @property
    def data(self) -> dict:
        return json.loads(self._json)

    def entity_count(self) -> int:
        return sum(len(v) for v in self.data.values())


class State:
    """Mutable builder; `snapshot()` is the immutable, hashable evidence."""

    def __init__(self, data: dict | None = None) -> None:
        self._data: dict[str, dict[str, dict]] = {k: {} for k in ENTITY_KINDS}
        if data is not None:
            self._adopt(data)

    # -- internal ---------------------------------------------------------
    def _adopt(self, data: dict) -> None:
        unknown = set(data) - set(ENTITY_KINDS)
        if unknown:
            raise StateError(f"unknown entity kinds: {sorted(unknown)}")
        for kind, records in data.items():
            if not isinstance(records, dict):
                raise StateError(f"{kind} must map id -> record")
            for rid, rec in records.items():
                self._data[kind][rid] = copy.deepcopy(rec)
        self._validate_records()

    def _validate_records(self) -> None:
        for kind in ("fact",):
            for rid, rec in self._data[kind].items():
                src = rec.get("source")
                if src not in FACT_SOURCES:
                    raise StateError(f"{kind}/{rid} has invalid source {src!r}")

    def hash(self) -> str:
        return sha256_text(canonical(self._data))

    # -- operations -------------------------------------------------------
    def put(self, kind: str, rid: str, record: dict) -> str:
        if kind not in ENTITY_KINDS:
            raise StateError(f"unknown entity kind: {kind}")
        if not rid:
            raise StateError("entity id must not be empty")
        if not isinstance(record, dict):
            raise StateError("record must be a mapping")
        if kind == "fact" and record.get("source") not in FACT_SOURCES:
            raise StateError(f"fact {rid} requires a valid source, got {record.get('source')!r}")
        self._data[kind][rid] = copy.deepcopy(record)
        self._validate_records()
        return self.hash()

    def get(self, kind: str, rid: str) -> dict:
        if kind not in ENTITY_KINDS:
            raise StateError(f"unknown entity kind: {kind}")
        if rid not in self._data[kind]:
            raise StateError(f"no such entity: {kind}/{rid}")
        return copy.deepcopy(self._data[kind][rid])

    def has(self, kind: str, rid: str) -> bool:
        return rid in self._data.get(kind, {})

    def delete(self, kind: str, rid: str, turn: int) -> str:
        """Removal is recorded, never erased."""
        rec = self.get(kind, rid)
        rec["state"] = "removed"
        rec["removed_at_turn"] = turn
        return self.put(kind, rid, rec)

    def assert_fact(self, fact_id: str) -> tuple[bool, Any, Any]:
        """(present, source, established_at_turn) — the C2 check."""
        rec = self._data["fact"].get(fact_id)
        if rec is None:
            return (False, None, None)
        if rec.get("state") == "removed":
            return (False, rec.get("source"), rec.get("established_at_turn"))
        return (True, rec.get("source"), rec.get("established_at_turn"))

    def snapshot(self) -> Snapshot:
        return Snapshot.from_data(copy.deepcopy(self._data))

    def data(self) -> dict:
        return copy.deepcopy(self._data)


def diff(before: Snapshot, after: Snapshot) -> list[dict]:
    """Entity-level changes between two pinned snapshots."""
    changes: list[dict] = []
    b_all = before.data
    a_all = after.data
    for kind in ENTITY_KINDS:
        b = b_all.get(kind, {})
        a = a_all.get(kind, {})
        for rid in sorted(set(b) | set(a)):
            if rid not in b:
                changes.append({"kind": kind, "id": rid, "change": "added"})
            elif rid not in a:
                changes.append({"kind": kind, "id": rid, "change": "erased"})
            elif canonical(b[rid]) != canonical(a[rid]):
                if a[rid].get("state") == "removed" and b[rid].get("state") != "removed":
                    changes.append({"kind": kind, "id": rid, "change": "removed"})
                else:
                    changes.append({"kind": kind, "id": rid, "change": "changed"})
    return changes
