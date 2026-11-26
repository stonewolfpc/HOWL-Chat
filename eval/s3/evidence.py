"""Layer 1 — original evidence.

Immutable, addressable conversation evidence. Every later layer cites spans
from here; nothing in this module is ever rewritten by formation or retrieval.

Tokens are supplied per turn by the caller (measured by the runtime or, offline,
stated explicitly), never estimated here. An estimate would let a displacement
assertion pass for the wrong reason.
"""
from __future__ import annotations

from dataclasses import dataclass, field, asdict
import hashlib
import json
import re
import pathlib


@dataclass(frozen=True)
class Span:
    span_id: str
    turn: int
    speaker: str
    text: str

    def as_dict(self) -> dict:
        return asdict(self)


@dataclass
class Turn:
    turn: int
    speaker: str
    text: str
    tokens: int
    span_ids: list[str] = field(default_factory=list)

    def as_dict(self) -> dict:
        return asdict(self)


_SENTENCE = re.compile(r"[^.!?]+[.!?]?")


def _segments(text: str) -> list[str]:
    return [m.group(0).strip() for m in _SENTENCE.finditer(text) if m.group(0).strip()]


class EvidenceLog:
    """Append-only evidence. Spans are stable and citation-addressable."""

    def __init__(self) -> None:
        self._turns: list[Turn] = []
        self._spans: dict[str, Span] = {}
        self._sealed = False

    def append(self, speaker: str, text: str, tokens: int) -> Turn:
        assert not self._sealed, "evidence log is sealed"
        assert tokens > 0, "turn tokens must be supplied, not guessed"
        turn = Turn(turn=len(self._turns) + 1, speaker=speaker, text=text, tokens=int(tokens))
        for i, seg in enumerate(_segments(text), start=1):
            sid = f"sp-{turn.turn:03d}-{i:02d}"
            self._spans[sid] = Span(span_id=sid, turn=turn.turn, speaker=speaker, text=seg)
            turn.span_ids.append(sid)
        self._turns.append(turn)
        return turn

    def seal(self) -> str:
        self._sealed = True
        return self.sha256()

    def sha256(self) -> str:
        blob = json.dumps([t.as_dict() for t in self._turns], sort_keys=True).encode()
        return hashlib.sha256(blob).hexdigest()

    def turns(self) -> list[Turn]:
        return list(self._turns)

    def span(self, span_id: str) -> Span:
        return self._spans[span_id]

    def spans_for_turns(self, turns: list[int]) -> list[str]:
        keep = set(turns)
        return [s for s in self._spans if self._spans[s].turn in keep]

    def text_for(self, span_ids: list[str]) -> str:
        return " ".join(self._spans[s].text for s in span_ids)

    def save(self, path: pathlib.Path) -> None:
        path.write_text(json.dumps({
            "sha256": self.sha256(),
            "turns": [t.as_dict() for t in self._turns],
            "spans": {k: v.as_dict() for k, v in self._spans.items()},
        }, indent=2, sort_keys=True) + "\n", encoding="utf-8")
