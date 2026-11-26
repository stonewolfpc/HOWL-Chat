"""Probes, quota enforcement, departure assertions and receipts.

Quotas and coverage are frozen (prereg §5), so they are *enforced* here rather
than described: a conversation that does not carry 4/2/2/2/1/1/1 probes and at
least three rollovers cannot produce a scored observation.

A probe that is scored must prove its supporting span left the window. If the
span is still present, the instrument raises `DroppedSpanStillPresentError` —
that probe is not a data point about memory, and treating it as one is the
failure the plan's threat table names first.
"""
from __future__ import annotations

import hashlib
from dataclasses import asdict, dataclass, field

from . import PROBE_QUOTAS, PROBES_PER_CONVERSATION
from .errors import DesignViolationError, DroppedSpanStillPresentError, ReceiptFieldError

PROBE_CATEGORIES = tuple(PROBE_QUOTAS)

# Categories whose correct answer requires evidence that has left the window.
MEMORY_RELEVANT = ("character_specific", "ambient", "restricted", "experiential",
                   "correction", "retraction")


@dataclass(frozen=True)
class ProbeSpec:
    probe_id: str
    turn: int
    category: str
    question: str
    owner_id: str
    key_terms: tuple[str, ...] = ()
    contradiction_terms: tuple[str, ...] = ()
    key_predicate: str = ""
    key_value: str = ""
    temporal: str = "current"          # current | historical | none
    answerable: bool = True
    acceptable: tuple[str, ...] = ()
    abstain_terms: tuple[str, ...] = ("i don't know", "i do not know", "not established",
                                      "no record", "cannot say", "i can't say")
    support_span_hint: str = ""        # expected evidence span, asserted absent
    context: str = ""                  # "" derives from category (world_state for ambient)
    # Which scenario the question is asked from. "" means no scenario context, so
    # scenario-owned records are not visible: one scenario cannot read another's.
    scenario_id: str = ""
    # Retraction probes name the predicate that was retracted, so support comes from
    # an explicit retraction event record rather than from the absence of a fact.
    retraction_of: str = ""

    def __post_init__(self) -> None:
        if self.category not in PROBE_CATEGORIES:
            raise DesignViolationError(f"unknown probe category {self.category!r}")
        if self.category == "unanswerable" and self.answerable:
            raise DesignViolationError(f"{self.probe_id}: unanswerable probe marked answerable")
        if self.category != "unanswerable" and not self.answerable:
            raise DesignViolationError(f"{self.probe_id}: answerable category marked unanswerable")

    @property
    def required_fields(self) -> tuple[str, ...]:
        return ("probe_id", "turn", "category", "question", "owner_id", "key_predicate")


def validate_quotas(probes: list[ProbeSpec]) -> dict:
    """Enforce prereg §5 exactly. No tolerance, no reclassification."""
    counts: dict[str, int] = {c: 0 for c in PROBE_CATEGORIES}
    for probe in probes:
        counts[probe.category] += 1
    if len(probes) != PROBES_PER_CONVERSATION:
        raise DesignViolationError(
            f"expected {PROBES_PER_CONVERSATION} probes, got {len(probes)}"
        )
    wrong = {c: (counts[c], PROBE_QUOTAS[c]) for c in PROBE_CATEGORIES if counts[c] != PROBE_QUOTAS[c]}
    if wrong:
        raise DesignViolationError(f"probe quota violated: {wrong}")
    turns = [p.turn for p in probes]
    if len(set(turns)) != len(turns):
        raise DesignViolationError("two probes share one turn; each probe needs its own turn")
    return {"counts": counts, "total": len(probes)}


def assert_span_departed(window, span_id: str, probe: ProbeSpec) -> dict:
    """Prove the evidence left the window before this probe is scored."""
    if not span_id:
        raise DesignViolationError(f"{probe.probe_id}: no supporting span declared")
    present = window.contains_span(span_id)
    if present:
        raise DroppedSpanStillPresentError(
            f"{probe.probe_id}: supporting span {span_id} is still in the active window"
        )
    dropped = any(span_id in r.dropped_spans for r in window.rollovers)
    if not dropped:
        raise DroppedSpanStillPresentError(
            f"{probe.probe_id}: span {span_id} is absent but no rollover records dropping it — "
            f"absence cannot be attributed to displacement"
        )
    return {"span_id": span_id, "absent": True,
            "rollover_index": next(r.index for r in window.rollovers if span_id in r.dropped_spans)}


@dataclass
class ProbeReceipt:
    conversation_id: str
    arm: str
    rep_index: int
    seed: int
    seed_honoured: bool
    probe_id: str
    turn: int
    category: str
    question: str
    response_text: str
    m1_category: str
    m1_credit: int
    m2_status: str
    selected_ids: list[str] = field(default_factory=list)
    abstained: bool = False
    retrieval_reason: str = ""
    supporting_span: str = ""
    span_absent: bool = False
    block_sha256: str = ""
    tokens: dict = field(default_factory=dict)
    flags: list[str] = field(default_factory=list)

    REQUIRED = ("conversation_id", "arm", "probe_id", "turn", "category", "response_text",
                "m1_category", "m2_status", "tokens")

    def as_dict(self) -> dict:
        return asdict(self)

    def validate(self) -> None:
        problems: list[str] = []
        for name in self.REQUIRED:
            value = getattr(self, name)
            if value is None or value == "" or (name == "tokens" and not value):
                problems.append(name)
        if self.selected_ids and not self.block_sha256:
            # A selection must name the block it was rendered into. An empty
            # selection is legitimate (abstention, or A0) and is flagged instead.
            problems.append("block_sha256")
        if self.m2_status not in ("SUPPORTED", "NOT_SUPPORTED", "N/A"):
            problems.append("m2_status")
        if problems:
            raise ReceiptFieldError(
                f"{self.probe_id} ({self.arm}): missing or invalid receipt fields: {sorted(set(problems))}"
            )


def block_sha256(text: str) -> str:
    return hashlib.sha256(text.encode("utf-8")).hexdigest()


__all__ = ["PROBE_CATEGORIES", "MEMORY_RELEVANT", "ProbeSpec", "validate_quotas",
           "assert_span_departed", "ProbeReceipt", "block_sha256"]
