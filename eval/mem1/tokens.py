"""Measured token accounting (plan §2.4).

Rules this module enforces:

  * counts come from a supplied counter — the runtime's own accounting — never
    from an estimate. There is deliberately no `len(text)//4` fallback: an
    estimate would let a capacity assertion pass for the wrong reason;
  * memory competes with history inside one prompt capacity, it does not extend
    it: `tokens.total <= capacity - reserve`, asserted every turn;
  * a missing or non-positive count is a hard failure, not a zero.

`DeclaredCounter` exists for offline fixtures: the counts are *declared* per
string exactly as a runtime would report them, and a string with no declared
count raises rather than being guessed.
"""
from __future__ import annotations

from dataclasses import asdict, dataclass, field

from .errors import TokenAccountingError


class TokenCounter:
    """Anything that can report a token count for a serialised string."""

    def count(self, text: str) -> int:      # pragma: no cover - interface
        raise NotImplementedError


class DeclaredCounter(TokenCounter):
    """Offline stand-in: counts are declared per string, never inferred."""

    def __init__(self, counts: dict[str, int] | None = None, default: int | None = None) -> None:
        self._counts = dict(counts or {})
        self._default = default

    def declare(self, text: str, count: int) -> None:
        if not isinstance(count, int) or count <= 0:
            raise TokenAccountingError(f"declared count must be a positive int, got {count!r}")
        self._counts[text] = count

    def count(self, text: str) -> int:
        if text == "":
            return 0
        if text in self._counts:
            value = self._counts[text]
        elif self._default is not None:
            value = self._default
        else:
            raise TokenAccountingError(
                "no token count supplied for this text: counts must come from the runtime, "
                "not be estimated here"
            )
        if not isinstance(value, int) or value <= 0:
            raise TokenAccountingError(f"token count must be a positive int, got {value!r}")
        return value


class WhitespaceCounter(TokenCounter):
    """G1 TEST ONLY. Whitespace words, not model tokens.

    Fixtures need *some* deterministic count to exercise matching and capacity
    logic. This counter is honest about what it is, and it must never be used to
    produce a G4 receipt: there, counts come from the runtime's own accounting
    (plan §2.4). The G4 receipt records which counter was used.
    """

    name = "whitespace-words-test-only"

    def count(self, text: str) -> int:
        if not text:
            return 0
        n = len(text.split())
        if n <= 0:
            raise TokenAccountingError("whitespace counter saw no tokens in non-empty text")
        return n


@dataclass
class Account:
    """Per-turn split, exactly the plan's §2.4 table."""

    system: int
    history: int
    memory: int
    overhead: int
    reserve: int
    capacity: int
    arm: str = ""
    turn: int = 0

    @property
    def total(self) -> int:
        return self.system + self.history + self.memory + self.overhead

    def as_dict(self) -> dict:
        out = asdict(self)
        out["total"] = self.total
        out["prompt_capacity"] = self.capacity - self.reserve
        return out

    def assert_fits(self) -> None:
        if self.reserve <= 0:
            raise TokenAccountingError("reply reserve must be positive")
        if self.total > self.capacity - self.reserve:
            raise TokenAccountingError(
                f"turn {self.turn} ({self.arm}): prompt {self.total} exceeds capacity "
                f"{self.capacity} - reserve {self.reserve} = {self.capacity - self.reserve}"
            )


def build_account(*, counter: TokenCounter, arm: str, turn: int, capacity: int, reserve: int,
                  system_text: str, history_tokens: int, memory_text: str = "",
                  wrapper_text: str = "") -> Account:
    """Assemble one turn's account from measured, non-overlapping parts.

    `memory_text` is the bullet lines only; `wrapper_text` is the header and any
    citation decoration around them. Keeping them separate is what stops the
    header being counted twice. `history_tokens` is the window's own bookkeeping
    (per-turn tokens supplied when the turn was appended), not a re-count.
    """
    if memory_text and wrapper_text and memory_text.startswith(wrapper_text.rstrip("\n")):
        raise TokenAccountingError(
            "memory_text appears to include the wrapper: pass bullets and header separately "
        )
    system = counter.count(system_text) if system_text else 0
    memory = counter.count(memory_text) if memory_text else 0
    overhead = counter.count(wrapper_text) if wrapper_text else 0
    account = Account(system=system, history=int(history_tokens), memory=memory,
                      overhead=overhead, reserve=reserve, capacity=capacity,
                      arm=arm, turn=turn)
    account.assert_fits()
    return account


def memory_delta(a: Account, b: Account) -> int:
    """Measured memory-token difference, reported per turn (never averaged away)."""
    return a.memory - b.memory


def token_advantage(accounts_by_arm: dict[str, list[Account]]) -> dict:
    """A systematic total-token advantage is a finding that blocks the capability claim.

    Reported, never smoothed: plan §2.4 says an A1 advantage in `tokens.total` is a
    finding about the experiment until explained.
    """
    def mean_total(arm: str) -> float:
        vals = [a.total for a in accounts_by_arm.get(arm, [])]
        return (sum(vals) / len(vals)) if vals else 0.0

    a1, a0, d = mean_total("A1"), mean_total("A0"), mean_total("D")
    return {"mean_total": {"A0": a0, "A1": a1, "D": d},
            "delta_A1_A0": a1 - a0, "delta_A1_D": a1 - d,
            "systematic_A1_advantage_over_A0": a1 > a0,
            "systematic_A1_advantage_over_D": a1 > d}


def require_runtime_counter(counter: TokenCounter) -> None:
    """G4 precondition: counts must come from the runtime, not a test stand-in.

    G1 fixtures may use `WhitespaceCounter`; a scored run may not. This is checked
    rather than documented, so a fixture counter cannot slip into a G4 receipt.
    """
    name = getattr(counter, "name", "")
    if not hasattr(counter, "count"):
        raise TokenAccountingError("counter has no count() method")
    if not name:
        raise TokenAccountingError(
            "counter has no declared name: a G4 receipt must record which tokenizer produced counts"
        )
    if name.endswith("test-only"):
        raise TokenAccountingError(
            f"counter {name!r} is a test instrument; G4 requires the runtime's own accounting"
        )


__all__ = ["TokenCounter", "DeclaredCounter", "WhitespaceCounter", "Account", "build_account",
           "memory_delta", "token_advantage", "require_runtime_counter"]
