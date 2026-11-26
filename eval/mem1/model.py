"""A' model — records, acquisition events, and knowledge edges (plan §3.1).

Three streams:

  * `MemoryRecord` — what the system currently holds. Chat-owned, with an explicit
    `owner_id`, `scope`, status, supersession link and a `source_span`.
  * `AcquisitionEvent` — append-only and authoritative. `via` is one of the seven
    adopted channels; the event stream is never rewritten.
  * `KnowledgeEdge` — a projection of the events. **It carries no truth status**;
    marking one current or stale is a record/event matter, never an edge's.

Enforced here, because the plan makes them vetoes rather than preferences:

  * a record whose `source_span` does not resolve is *fabricated provenance*;
  * history is never erased — supersession and revocation keep every version;
  * nothing is silently normalised, deduplicated or fuzzy-matched.
"""
from __future__ import annotations

import copy
import hashlib
import json
from dataclasses import asdict, dataclass, field

from .errors import (DesignViolationError, HistoryErasedError, ProvenanceFabricationError,
                     UnresolvedSpanError)

ACTIVE = "active"
SUPERSEDED = "superseded"
REVOKED = "revoked"

VIAS = ("witnessed", "verified", "read", "told", "inferred", "narrator_grant", "revoked")

# Where a record may be read. Visibility only; never ownership.
SCOPES = ("owner", "restricted", "world", "chat_state", "scenario_state")

# Who a record belongs to. A character owns what it said; a world or scenario owns
# what the narrator establishes inside it; the chat owns conversation-level
# narrator facts that belong to no world. Ownership kind is not a visibility scope.
OWNER_KINDS = ("character", "chat", "world", "scenario")
# Owner kinds whose eligibility depends on declared active state (Stage 1, G1).
GATED_OWNER_KINDS = ("world", "scenario")
# The reserved owner id of the conversation container itself. Narrator-derived
# records belong to the active world when one is declared and to the chat
# otherwise; the id never names a character, world, or scenario.
CHAT_OWNER = "chat"

KINDS = ("fact", "preference", "commitment", "event", "state_change", "absence")


def canonical(obj) -> str:
    return json.dumps(obj, sort_keys=True, separators=(",", ":"), ensure_ascii=False)


@dataclass(frozen=True)
class ScopeState:
    """One declared activation state for a world or scenario, applied at a turn.

    State is a declaration, not a record: applying it changes eligibility only. It
    never creates, deletes, supersedes or revokes anything, and it never touches
    truth status. Preservation is retention; restoration is renewed eligibility.
    """

    turn: int
    owner_kind: str
    owner_id: str
    active: bool

    def __post_init__(self) -> None:
        if self.owner_kind not in GATED_OWNER_KINDS:
            raise DesignViolationError(
                f"scope state applies to {GATED_OWNER_KINDS}, not {self.owner_kind!r}")
        if not self.owner_id:
            raise DesignViolationError("scope state has no owner_id")
        if self.owner_id == CHAT_OWNER:
            raise DesignViolationError(
                f"{CHAT_OWNER!r} is the chat container, not a world or scenario id")
        if self.turn < 1:
            raise DesignViolationError(f"scope state turn must be >= 1, got {self.turn}")


@dataclass(frozen=True)
class ActiveState:
    """Resolved activation state at one turn boundary.

    `declared` separates "no state declared" (D1: no gate) from "declared, with
    nothing active". `parents` maps a scenario id to its single parent world id and
    is supplied explicitly, never inferred from record content.
    """

    declared: bool
    worlds: frozenset[str] = frozenset()
    scenarios: frozenset[str] = frozenset()
    parents: dict = field(default_factory=dict)

    def as_dict(self) -> dict:
        return {"declared": self.declared, "worlds": sorted(self.worlds),
                "scenarios": sorted(self.scenarios),
                "parents": dict(sorted(self.parents.items()))}


def record_is_active(rec, state: "ActiveState | None") -> bool:
    """Activation gate for one record. Undeclared state gates nothing (D1).

    World records follow their own world. A scenario record needs its scenario
    active AND its declared parent world active: the parent comes from the
    explicit map, so a missing mapping fails closed rather than defaulting.
    Character and chat records are never gated by world or scenario state.
    """
    if state is None or not state.declared:
        return True
    if rec.owner_kind == "world":
        return rec.owner_id in state.worlds
    if rec.owner_kind == "scenario":
        if rec.owner_id not in state.scenarios:
            return False
        parent = state.parents.get(rec.owner_id)
        if parent is None:
            return False
        return parent in state.worlds
    return True


@dataclass
class MemoryRecord:
    record_id: str
    kind: str
    owner_id: str
    owner_kind: str
    scope: str
    audience: tuple[str, ...]
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
    revoked_at_turn: int | None = None
    # A retraction event record points at the record it retracts. This is a link,
    # not the derived `revoked_turn` timestamp: the contract forbids the latter as a
    # field, while an append-only reference is how the event stream already works.
    retracts: str | None = None
    formation_method: str = "deterministic"

    def as_dict(self) -> dict:
        out = asdict(self)
        out["audience"] = list(self.audience)
        return out


@dataclass(frozen=True)
class AcquisitionEvent:
    event_id: str
    turn: int
    via: str
    owner_id: str
    kind: str
    subject: str
    predicate: str
    value: str
    source_span: str
    confidence: float
    record_id: str
    supersedes: str | None = None

    def as_dict(self) -> dict:
        return asdict(self)


@dataclass(frozen=True)
class KnowledgeEdge:
    """Projection only. `truth_status` is deliberately absent."""

    owner_id: str
    subject: str
    predicate: str
    value: str
    record_id: str
    via: str
    turn: int

    def as_dict(self) -> dict:
        return asdict(self)


class MemoryStore:
    """Append-only event stream with a derived, supersession-aware projection."""

    def __init__(self, evidence=None, conversation_id: str = "") -> None:
        self.evidence = evidence          # optional EvidenceLog for provenance checks
        self.conversation_id = conversation_id
        self._events: list[AcquisitionEvent] = []
        self._records: dict[str, MemoryRecord] = {}
        self._next_record = 0
        self._revocations: dict[str, dict] = {}

    # -- writes -----------------------------------------------------------
    def add_event(
        self,
        *,
        turn: int,
        via: str,
        owner_id: str,
        owner_kind: str = "character",
        kind: str,
        subject: str,
        predicate: str,
        value: str,
        source_span: str,
        confidence: float,
        scope: str = "owner",
        audience: tuple[str, ...] = (),
        supersede: str | None = None,
        formation_method: str = "deterministic",
        retracts: str | None = None,
        text: str | None = None,
    ) -> MemoryRecord:
        if via not in VIAS:
            raise ValueError(f"unknown acquisition channel: {via!r}")
        if kind not in KINDS:
            raise ValueError(f"unknown record kind: {kind!r}")
        if scope not in SCOPES:
            raise ValueError(f"unknown scope: {scope!r}")
        if owner_kind not in OWNER_KINDS:
            raise DesignViolationError(f"unknown ownership kind: {owner_kind!r}")
        if not owner_id:
            raise DesignViolationError("record has no owner_id")
        if owner_id == CHAT_OWNER and owner_kind != "chat":
            raise DesignViolationError(
                f"owner id {CHAT_OWNER!r} identifies the chat container; it cannot "
                f"identify a {owner_kind!r}")
        for other in self._records.values():
            if other.owner_id == owner_id and other.owner_kind != owner_kind:
                raise DesignViolationError(
                    f"owner id {owner_id!r} already identifies a {other.owner_kind!r} "
                    f"record ({other.record_id}); ownership would be ambiguous")
        if not source_span:
            raise ProvenanceFabricationError("record has no source_span")
        if self.evidence is not None:
            try:
                self.evidence.span(source_span)
            except KeyError:
                raise ProvenanceFabricationError(
                    f"source_span {source_span!r} does not resolve in this conversation"
                ) from None

        self._next_record += 1
        rid = f"rec-{self._next_record:04d}"
        rec = MemoryRecord(
            record_id=rid, kind=kind, owner_id=owner_id, owner_kind=owner_kind,
            scope=scope, audience=tuple(audience),
            text=text if text is not None else f"{subject} {predicate}: {value}",
            subject=subject, predicate=predicate, value=value, source_span=source_span,
            established_turn=turn, confidence=confidence, retracts=retracts,
            formation_method=formation_method,
        )
        if supersede:
            old = self._records.get(supersede)
            if old is None:
                raise HistoryErasedError(f"supersede target {supersede!r} does not exist")
            old.status, old.superseded_by = SUPERSEDED, rid
            rec.supersedes = supersede
            self._append_event(turn, "witnessed", rec, supersede=supersede)
        else:
            self._append_event(turn, via, rec)
        self._records[rid] = rec
        return rec

    def _append_event(self, turn: int, via: str, rec: MemoryRecord,
                      supersede: str | None = None) -> None:
        self._events.append(AcquisitionEvent(
            event_id=f"ev-{len(self._events) + 1:05d}", turn=turn, via=via,
            owner_id=rec.owner_id, kind=rec.kind, subject=rec.subject,
            predicate=rec.predicate, value=rec.value, source_span=rec.source_span,
            confidence=rec.confidence, record_id=rec.record_id, supersedes=supersede,
        ))

    def retract(self, record_id: str, turn: int, source_span: str, reason: str) -> MemoryRecord:
        """Retraction as an explicit, selectable event record.

        A generic absence record would not be evidence that something was retracted —
        it conflates "retracted" with "never established" and with "not known". So a
        retraction emits a record that (a) cites its own span, (b) references the
        retracted record, (c) carries the retracted value in its text for retrieval,
        and (d) does not present that value as current. The target record remains,
        revoked, with its provenance intact.
        """
        target = self._records[record_id]
        self.revoke(record_id, turn, reason)
        return self.add_event(
            turn=turn, via="revoked", owner_id=target.owner_id,
            owner_kind=target.owner_kind, kind="event",
            subject=target.subject, predicate=f"retraction:{target.predicate}",
            value="retracted", source_span=source_span, confidence=1.0,
            scope=target.scope, audience=target.audience, retracts=record_id,
            text=f"{target.subject} retracted: {target.value}",
        )

    def revoke(self, record_id: str, turn: int, reason: str) -> None:
        """Retraction. The record stops being current; it is never deleted."""
        rec = self._records[record_id]
        rec.status = REVOKED
        rec.revoked_at_turn = turn
        self._revocations[record_id] = {"turn": turn, "reason": reason}
        self._events.append(AcquisitionEvent(
            event_id=f"ev-{len(self._events) + 1:05d}", turn=turn, via="revoked",
            owner_id=rec.owner_id, kind=rec.kind, subject=rec.subject,
            predicate=rec.predicate, value=rec.value, source_span=rec.source_span,
            confidence=rec.confidence, record_id=record_id,
        ))

    # -- reads ------------------------------------------------------------
    def get(self, record_id: str) -> MemoryRecord:
        return self._records[record_id]

    def all(self) -> list[MemoryRecord]:
        return list(self._records.values())

    def current(self, owner_id: str | None = None) -> list[MemoryRecord]:
        return [r for r in self._records.values()
                if r.status == ACTIVE and (owner_id is None or r.owner_id == owner_id)]

    def history(self, owner_id: str, subject: str, predicate: str) -> list[MemoryRecord]:
        return sorted((r for r in self._records.values()
                       if r.owner_id == owner_id and r.subject == subject
                       and r.predicate == predicate),
                      key=lambda r: r.established_turn)

    def versions(self, subject: str, predicate: str) -> list[MemoryRecord]:
        return sorted((r for r in self._records.values()
                       if r.subject == subject and r.predicate == predicate),
                      key=lambda r: r.established_turn)

    def active_value(self, owner_id: str, subject: str, predicate: str) -> str | None:
        for rec in self.history(owner_id, subject, predicate):
            if rec.status == ACTIVE:
                return rec.value
        return None

    def events(self) -> list[AcquisitionEvent]:
        return list(self._events)

    def edges(self) -> list[KnowledgeEdge]:
        return [KnowledgeEdge(owner_id=e.owner_id, subject=e.subject, predicate=e.predicate,
                              value=e.value, record_id=e.record_id, via=e.via, turn=e.turn)
                for e in self._events]

    # -- integrity --------------------------------------------------------
    def verify_history(self) -> list[str]:
        """Every supersession link must resolve, and no version may be missing."""
        problems: list[str] = []
        for rec in self._records.values():
            if rec.superseded_by and rec.superseded_by not in self._records:
                problems.append(f"{rec.record_id} -> missing {rec.superseded_by}")
            if rec.supersedes and rec.supersedes not in self._records:
                problems.append(f"{rec.record_id} supersedes missing {rec.supersedes}")
        if problems:
            raise HistoryErasedError("; ".join(problems))
        return problems

    def assert_history_retained(self, subject: str, predicate: str, minimum: int) -> None:
        versions = self.versions(subject, predicate)
        if len(versions) < minimum:
            raise HistoryErasedError(
                f"{subject}/{predicate}: expected >= {minimum} retained versions, got {len(versions)}"
            )
        for v in versions:
            if v.superseded_by and v.status == ACTIVE:
                raise HistoryErasedError(f"{v.record_id} superseded but still active")

    def erase_for_test(self, record_id: str) -> None:
        """Test hook ONLY: simulates history erasure so the detector can be shown to bite."""
        self._records.pop(record_id, None)

    def snapshot(self) -> dict:
        return {"records": [r.as_dict() for r in self._records.values()],
                "events": [e.as_dict() for e in self._events]}

    def snapshot_sha256(self) -> str:
        return hashlib.sha256(canonical(self.snapshot()).encode("utf-8")).hexdigest()

    def copy(self) -> "MemoryStore":
        other = MemoryStore(evidence=self.evidence, conversation_id=self.conversation_id)
        other._records = {k: copy.deepcopy(v) for k, v in self._records.items()}
        other._events = list(self._events)
        other._next_record = self._next_record
        other._revocations = copy.deepcopy(self._revocations)
        return other


__all__ = ["ACTIVE", "SUPERSEDED", "REVOKED", "VIAS", "SCOPES", "KINDS",
           "OWNER_KINDS", "GATED_OWNER_KINDS", "CHAT_OWNER", "ScopeState",
           "ActiveState", "record_is_active",
           "MemoryRecord", "AcquisitionEvent", "KnowledgeEdge", "MemoryStore"]
