"""Layer 3 — knowledge and state.

Typed records with one authoritative home, stable identity, and supersession.
Current state is derived; history is never erased. Ownership is carried on every
record and checked at retrieval time, not at retrieval author's discretion.
"""
from __future__ import annotations

from dataclasses import dataclass, field, asdict
import json

ACTIVE = "active"
SUPERSEDED = "superseded"
REJECTED = "rejected"

# Record kinds follow the destination taxonomy in OWNER-DECISION-MEMORY-FORMATION §3.
KINDS = ("fact", "preference", "commitment", "event", "state_change", "absence")


@dataclass
class Record:
    record_id: str
    kind: str
    owner_id: str
    owner_kind: str
    text: str
    subject: str
    predicate: str
    value: str
    source_span: str
    established_turn: int
    confidence: float
    status: str = ACTIVE
    supersedes: str | None = None
    superseded_by: str | None = None
    formation_method: str = "deterministic"

    def as_dict(self) -> dict:
        return asdict(self)


class RecordStore:
    def __init__(self) -> None:
        self._records: dict[str, Record] = {}
        self._next = 0

    def add(self, **kw) -> Record:
        self._next += 1
        rid = f"rec-{self._next:04d}"
        rec = Record(record_id=rid, **kw)
        assert rec.kind in KINDS, f"unknown kind {rec.kind}"
        assert rec.owner_kind in ("character", "chat"), "owner_kind must be stable kind"
        self._records[rid] = rec
        return rec

    def get(self, rid: str) -> Record:
        return self._records[rid]

    def all(self) -> list[Record]:
        return list(self._records.values())

    def supersede(self, old_id: str, new_id: str) -> None:
        old, new = self._records[old_id], self._records[new_id]
        old.status, old.superseded_by = SUPERSEDED, new_id
        new.supersedes = old_id

    def current(self, owner_id: str | None = None) -> list[Record]:
        return [r for r in self._records.values()
                if r.status == ACTIVE and (owner_id is None or r.owner_id == owner_id)]

    def history(self, owner_id: str, subject: str, predicate: str) -> list[Record]:
        return sorted((r for r in self._records.values()
                       if r.owner_id == owner_id and r.subject == subject
                       and r.predicate == predicate),
                      key=lambda r: r.established_turn)

    def count_by_status(self) -> dict[str, int]:
        out: dict[str, int] = {}
        for r in self._records.values():
            out[r.status] = out.get(r.status, 0) + 1
        return out

    def save(self, path) -> None:
        path.write_text(json.dumps([r.as_dict() for r in self._records.values()],
                                   indent=2, sort_keys=True) + "\n", encoding="utf-8")
