"""The dispatch boundary: the only place a call may be counted or attempted.

The owner's inspection point was exact — a ceiling enforced when a *plan* is
built is not enforced when a *call* is dispatched. `CallLedger.reserve` is the
single choke point: it increments before the adapter is touched, so a failed
call still consumes a slot (row 482 §6: failures count against the ceiling).

P1 cannot dispatch anything: every adapter raises. What P1 can prove is that the
counter, the ceiling and the ordering are correct.
"""

from __future__ import annotations

from dataclasses import dataclass, field

from .adapter import GenerateRequest, GenerateResponse
from .budget import check_payment_demand
from .plan import HARD_CEILING


class CallCeilingExceeded(RuntimeError):
    """Raised at dispatch when the authorised call count is exhausted."""


class DispatchRefused(RuntimeError):
    """Raised when a dispatch is attempted outside the authorised conditions."""


@dataclass
class CallLedger:
    ceiling: int = HARD_CEILING
    attempted: int = 0
    per_slot: dict = field(default_factory=dict)
    consumed_by_failure: int = 0

    def __post_init__(self) -> None:
        if self.ceiling > HARD_CEILING:
            raise DispatchRefused(
                f"ledger ceiling {self.ceiling} exceeds the authorised maximum {HARD_CEILING}"
            )

    def remaining(self) -> int:
        return max(0, self.ceiling - self.attempted)

    def reserve(self, slot: str) -> int:
        """Count a call *before* it happens. Returns the 1-based call number."""
        if self.attempted + 1 > self.ceiling:
            raise CallCeilingExceeded(
                f"call {self.attempted + 1} exceeds the ceiling of {self.ceiling}; stop"
            )
        self.attempted += 1
        self.per_slot[slot] = self.per_slot.get(slot, 0) + 1
        return self.attempted

    def note_failure(self) -> None:
        self.consumed_by_failure += 1


class Dispatcher:
    """Counts, then calls. No other component may call an adapter's generate."""

    def __init__(self, ledger: CallLedger, adapter: object) -> None:
        if not hasattr(adapter, "generate"):
            raise DispatchRefused("dispatcher needs an adapter with generate()")
        self._ledger = ledger
        self._adapter = adapter

    @property
    def ledger(self) -> CallLedger:
        return self._ledger

    def dispatch(self, slot: str, request: GenerateRequest) -> GenerateResponse:
        if not isinstance(request, GenerateRequest):
            raise DispatchRefused("dispatch takes a GenerateRequest")
        self._ledger.reserve(slot)
        try:
            return self._adapter.generate(request)
        except Exception as exc:
            self._ledger.note_failure()
            # A payment demand is a stop condition, not a retryable error: it is
            # re-raised as one so no caller can treat it as transient.
            check_payment_demand(exc)
            raise
