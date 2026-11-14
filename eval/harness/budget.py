"""Inference route policy: free allowance or local. There is no third route.

Owner direction (2026-11-13): **no payment method is on file, and none will be
added.** So paid inference is not "prohibited pending authorisation" — it is
*impossible*. That removes the need for a balance API, a spend authorisation
flow, or any credential that could incur a charge.

The rule reduces to three lines:

1. Hugging Face is usable **only** while it serves the request inside the free
   allowance.
2. If it asks for payment, or the allowance is exhausted, **stop** — do not
   retry, do not top up, do not switch automatically.
3. The fallback is local (KoboldCpp, SillyTavern, another compatible server),
   chosen by a human, not by the harness.

Verification stays lightweight on purpose: the provider's own response is the
signal. A successful call was free; a payment demand is a stop condition.
"""

from __future__ import annotations

from dataclasses import dataclass

# Owner statement, 2026-11-13: no card, no credits purchased, none planned.
PAYMENT_METHOD_ON_FILE = False

# Hugging Face removed the free-account monthly allowance (docs change merged
# 2026-10-07): free accounts get *no* included credits. With no payment method
# there is therefore no coverage to verify and no route to free hosted use.
# See MEMORY-09 P1 inspection addendum A2, section 9.4.
FREE_ALLOWANCE_EXISTS = False

HOSTED_REFUSAL = (
    "hosted inference is refused. A free Hugging Face account has no included monthly "
    "allowance, and neither the applicable price nor remaining coverage can be established "
    "before dispatch. Zero out-of-pocket spend cannot be guaranteed, so the call is not made. "
    "Use a local route (local:koboldcpp, local:sillytavern, ...)."
)

MAX_PAID_SPEND_USD = 0.00
AUTO_TOPUP_PERMITTED = False
AUTOMATIC_FALLBACK_PERMITTED = False

LOCAL_PREFIX = "local:"

ROUTE_PERMISSION = {
    "huggingface_free": "refused: a free account has no included monthly allowance",
    "huggingface_paid": "unreachable: no payment method on file",
    "automatic_paid_fallback": "no permission",
    "local": "preferred fallback; no per-call charge",
    "other_free": "eligible after a free-tier check",
}

PAYMENT_SIGNALS = (
    "402",
    "payment required",
    "insufficient credit",
    "insufficient_quota",
    "quota exceeded",
    "credit balance",
    "billing",
    "upgrade your plan",
    "add a payment method",
)


class BudgetRefusal(RuntimeError):
    """Raised when a route would require payment."""


class RouteChangeRefused(BudgetRefusal):
    """Raised when something would silently move a call to another route."""


class PaymentDemanded(BudgetRefusal):
    """The provider asked for money. Stop; do not retry or substitute."""


def is_local(provider: str) -> bool:
    return isinstance(provider, str) and provider.strip().lower().startswith(LOCAL_PREFIX)


def classify(provider: str) -> str:
    return "local" if is_local(provider) else "hosted"


def check_provider(provider: str, payment_method_on_file: bool | None = None) -> str:
    """Return the permitted route label, or raise.

    Hosted inference is permitted precisely because nothing can be charged. If a
    payment method ever exists, this refuses until the owner re-states the rule —
    the harness never assumes a charge is acceptable.
    """
    if not isinstance(provider, str) or not provider.strip():
        raise BudgetRefusal("no provider named")
    if is_local(provider):
        return "local"
    on_file = PAYMENT_METHOD_ON_FILE if payment_method_on_file is None else payment_method_on_file
    if on_file:
        raise BudgetRefusal(
            "a payment method is on file: paid inference has no permission ($0.00 ceiling). "
            "Remove the payment method or use a local route."
        )
    # No payment method is not enough: with no free allowance there is no
    # coverage to verify, so hosted is refused rather than assumed free.
    raise BudgetRefusal(HOSTED_REFUSAL)


def check_spend(usd: float) -> None:
    if usd > MAX_PAID_SPEND_USD:
        raise BudgetRefusal(
            f"spend {usd:.2f} exceeds the authorised ceiling {MAX_PAID_SPEND_USD:.2f}"
        )


def is_payment_demand(exc_or_text: object) -> bool:
    """Does this failure mean 'pay us'? Checked on status codes and message text."""
    status = getattr(exc_or_text, "status_code", None) or getattr(exc_or_text, "status", None)
    if status == 402:
        return True
    text = str(exc_or_text).lower()
    return any(signal in text for signal in PAYMENT_SIGNALS)


def check_payment_demand(exc: object) -> None:
    """Turn a payment-related failure into a stop, never a retry or a switch."""
    if not is_payment_demand(exc):
        return
    message = (
        "the provider demanded payment or reported an exhausted allowance. "
        "Stopping: do not retry, top up, or switch providers automatically. "
        "Move to a local route by explicit decision."
    )
    if isinstance(exc, BaseException):
        raise PaymentDemanded(message) from exc
    raise PaymentDemanded(message)


def check_no_route_change(requested: str, actual: str) -> None:
    if classify(requested) != classify(actual) or requested != actual:
        raise RouteChangeRefused(
            f"route changed from {requested!r} to {actual!r}; automatic fallback is not permitted"
        )


def permission_table() -> dict:
    return dict(ROUTE_PERMISSION)


@dataclass(frozen=True)
class FreeOnlyState:
    """Minimal record of why a hosted route is considered free-only."""

    payment_method_on_file: bool
    stated_by: str
    stated_on: str

    def describe(self) -> str:
        if self.payment_method_on_file:
            return "payment method present: hosted routes refused"
        return f"no payment method on file (stated by {self.stated_by}, {self.stated_on})"


def owner_state() -> FreeOnlyState:
    return FreeOnlyState(
        payment_method_on_file=PAYMENT_METHOD_ON_FILE,
        stated_by="owner",
        stated_on="2026-11-13",
    )
