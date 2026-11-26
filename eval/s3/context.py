"""Context window management.

Displacement is asserted from the window's own bookkeeping, never inferred from
turn count. Each rollover records exactly which turns and spans left, so a probe
can prove that the evidence it should retrieve was no longer in the request.
"""
from __future__ import annotations

from dataclasses import dataclass, field

from eval.s3.evidence import Turn


@dataclass
class Rollover:
    index: int
    at_turn: int
    tokens_before: int
    dropped_turns: list[int]
    dropped_spans: list[str]
    retained_turns: list[int]

    def as_dict(self) -> dict:
        return {"index": self.index, "at_turn": self.at_turn,
                "tokens_before": self.tokens_before,
                "dropped_turns": self.dropped_turns,
                "dropped_spans": self.dropped_spans,
                "retained_turns": self.retained_turns}


class Window:
    """Sliding prompt window with an explicit reply reserve."""

    def __init__(self, capacity: int, reserve: int) -> None:
        assert capacity > reserve > 0, "capacity must leave room for the reply"
        self.capacity = capacity
        self.reserve = reserve
        self.kept: list[Turn] = []
        self.rollovers: list[Rollover] = []

    @property
    def prompt_capacity(self) -> int:
        return self.capacity - self.reserve

    def tokens_used(self) -> int:
        return sum(t.tokens for t in self.kept)

    def add(self, turn: Turn) -> Rollover | None:
        self.kept.append(turn)
        if self.tokens_used() <= self.prompt_capacity:
            return None
        # Drop oldest turns until the prompt fits. Evidence leaves here, by definition.
        dropped_turns: list[int] = []
        dropped_spans: list[str] = []
        while self.kept and self.tokens_used() > self.prompt_capacity:
            victim = self.kept.pop(0)
            dropped_turns.append(victim.turn)
            dropped_spans.extend(victim.span_ids)
        roll = Rollover(index=len(self.rollovers) + 1, at_turn=turn.turn,
                        tokens_before=sum(t.tokens for t in self.kept) + sum(
                            next(t.tokens for t in [turn]) for _ in [0]),
                        dropped_turns=dropped_turns, dropped_spans=dropped_spans,
                        retained_turns=[t.turn for t in self.kept])
        self.rollovers.append(roll)
        return roll

    def present_turns(self) -> list[int]:
        return [t.turn for t in self.kept]

    def present_spans(self) -> set[str]:
        return {s for t in self.kept for s in t.span_ids}

    def contains_span(self, span_id: str) -> bool:
        return span_id in self.present_spans()

    def as_dict(self) -> dict:
        return {"capacity": self.capacity, "reserve": self.reserve,
                "prompt_capacity": self.prompt_capacity,
                "tokens_used": self.tokens_used(),
                "present_turns": self.present_turns(),
                "rollovers": [r.as_dict() for r in self.rollovers]}
