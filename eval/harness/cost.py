"""Pricing pre-flight (row 484 requirement 3; P2 prerequisite).

No pricing, no spend. With a $0.00 ceiling, the only quote that passes is a
verified free-tier quote whose worst-case cost is zero. A quote is owner-supplied
metadata, never a harness guess.
"""

from __future__ import annotations

from dataclasses import dataclass

from .budget import MAX_PAID_SPEND_USD

# Rough prompt-token allowance per call, used only for a worst-case bound.
ASSUMED_INPUT_TOKENS = 1024


class CostRefusal(RuntimeError):
    """Raised when a call cannot be priced, or would exceed the ceiling."""


@dataclass(frozen=True)
class PriceQuote:
    model: str
    provider: str
    currency: str
    input_per_1k: float
    output_per_1k: float
    retrieved_at: str
    source: str
    is_free_tier: bool

    def validate(self) -> None:
        if not self.model or not self.provider:
            raise CostRefusal("quote does not identify model and provider")
        if not self.currency:
            raise CostRefusal("quote has no currency")
        if self.input_per_1k < 0 or self.output_per_1k < 0:
            raise CostRefusal("negative price in quote")
        if not self.retrieved_at or not self.source:
            raise CostRefusal("quote has no retrieval date or source")


def worst_case_usd(quote: PriceQuote, planned_calls: int, max_output_tokens: int) -> float:
    per_call = (
        quote.input_per_1k * ASSUMED_INPUT_TOKENS + quote.output_per_1k * max_output_tokens
    ) / 1000.0
    return round(per_call * planned_calls, 6)


def preflight(
    quote: PriceQuote | None, planned_calls: int, max_output_tokens: int
) -> dict:
    """Return the cost decision, or refuse. Never returns an over-ceiling plan."""
    if quote is None:
        raise CostRefusal(
            "no pricing quote: refusing. A plan cannot be priced by assumption."
        )
    quote.validate()
    if not quote.is_free_tier:
        raise CostRefusal("quote is not a free tier: paid inference has no permission")
    worst = worst_case_usd(quote, planned_calls, max_output_tokens)
    if worst > MAX_PAID_SPEND_USD:
        raise CostRefusal(
            f"worst-case cost {worst:.6f} {quote.currency} exceeds the ceiling "
            f"{MAX_PAID_SPEND_USD:.2f}"
        )
    if planned_calls <= 0:
        raise CostRefusal("no calls to price")
    return {
        "model": quote.model,
        "provider": quote.provider,
        "planned_calls": planned_calls,
        "max_output_tokens": max_output_tokens,
        "worst_case": worst,
        "currency": quote.currency,
        "ceiling": MAX_PAID_SPEND_USD,
        "verdict": "priced within ceiling",
    }
