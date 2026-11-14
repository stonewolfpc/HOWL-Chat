"""P2 execution gate.

P1 boundary 4: no P2 execution or behavioural evaluation. This module makes
that a property of the code rather than a promise: execution requires a grant
file (`eval/runs/P2-GRANT.json`) carrying the owner's signature fields, and no
such file exists or may be created under P1.

A grant file alone is not sufficient: the plan must also be free of unassigned
model slots (plan.build_plan), so two independent barriers stand in front of a
single hosted call.
"""

from __future__ import annotations

import json
from dataclasses import dataclass
from pathlib import Path

from .adapter import ExecutionNotAuthorised, canonical, sha256_text
from .plan import HARD_CEILING, SLOTS

EVAL_ROOT = Path(__file__).resolve().parents[1]
GRANT_PATH = EVAL_ROOT / "runs" / "P2-GRANT.json"

# L3 is a *different* authorisation from P2: one local call against an identified
# server and model, never a hosted route. It gets its own file so an L3 run cannot
# be mistaken for, or reused as, a behavioural-evaluation grant.
L3_GRANT_PATH = EVAL_ROOT / "runs" / "L3-GRANT.json"
L3_CAPS_FIELDS = ("date", "call_ceiling", "cost_ceiling_usd", "endpoint", "model_path", "model_sha256")
L3_REQUIRED = ("report_id", "owner_signature") + L3_CAPS_FIELDS

REQUIRED_SIGNATURE_FIELDS = (
    "owner_signature",
    "date",
    "call_ceiling",
    "cost_ceiling_usd",
    "models",
    "providers",
    "caps_sha256",
)

# Fields the signature covers. Editing any of them after signing breaks the hash.
CAPS_FIELDS = ("date", "call_ceiling", "cost_ceiling_usd", "models", "providers")

MAX_PAID_GRANT_USD = 0.0


def caps_sha256(payload: dict) -> str:
    """Hash of exactly the fields a signature is meant to cover."""
    return sha256_text(canonical({k: payload.get(k) for k in CAPS_FIELDS}))


@dataclass(frozen=True)
class ExecutionGrant:
    owner_signature: str
    date: str
    call_ceiling: int
    cost_ceiling_usd: float
    models: dict
    providers: dict
    grants_execution: bool = True


def grant_status() -> str:
    return "PRESENT" if GRANT_PATH.exists() else "ABSENT"


def load_grant(path: Path | None = None) -> ExecutionGrant:
    target = path or GRANT_PATH
    if not target.exists():
        raise ExecutionNotAuthorised(
            "P2 is not granted. BEA-1 authorises P1 only; execution needs a separate "
            "owner decision after P1 inspection. No hosted call may be made."
        )
    try:
        payload = json.loads(target.read_text(encoding="utf-8"))
    except Exception as exc:  # malformed grant is not permission
        raise ExecutionNotAuthorised(f"P2 grant file is not readable as JSON: {exc}") from exc
    if not isinstance(payload, dict):
        raise ExecutionNotAuthorised("P2 grant must be a JSON object")
    missing = [f for f in REQUIRED_SIGNATURE_FIELDS if f not in payload]
    if missing:
        raise ExecutionNotAuthorised(f"P2 grant is incomplete: missing {missing}")

    signature = payload.get("owner_signature")
    if not isinstance(signature, str) or len(signature.strip()) < 8:
        raise ExecutionNotAuthorised("P2 grant carries no usable owner signature")

    date = payload.get("date")
    if not isinstance(date, str) or len(date) != 10 or date[4] != "-" or date[7] != "-":
        raise ExecutionNotAuthorised(f"grant date must be YYYY-MM-DD, got {date!r}")

    ceiling = payload.get("call_ceiling")
    if isinstance(ceiling, bool) or not isinstance(ceiling, int):
        raise ExecutionNotAuthorised(f"call_ceiling must be an integer, got {type(ceiling).__name__}")
    if ceiling < 1 or ceiling > HARD_CEILING:
        raise ExecutionNotAuthorised(
            f"grant ceiling {ceiling} is outside 1..{HARD_CEILING}"
        )

    cost = payload.get("cost_ceiling_usd")
    if isinstance(cost, bool) or not isinstance(cost, (int, float)):
        raise ExecutionNotAuthorised("cost_ceiling_usd must be a number")
    if float(cost) > MAX_PAID_GRANT_USD:
        raise ExecutionNotAuthorised(
            f"cost ceiling {cost} exceeds the authorised {MAX_PAID_GRANT_USD:.2f}; no grant may raise it"
        )

    for key in ("models", "providers"):
        value = payload.get(key)
        if not isinstance(value, dict):
            raise ExecutionNotAuthorised(f"{key} must be an object mapping slots to values")
        if set(value) != set(SLOTS):
            raise ExecutionNotAuthorised(f"{key} must name exactly the slots {list(SLOTS)}")
        if any(not isinstance(v, str) or v.strip() in ("", "UNASSIGNED") for v in value.values()):
            raise ExecutionNotAuthorised(f"{key} contains an unassigned slot")

    if any(v.strip().lower() == "auto" for v in payload["providers"].values()):
        raise ExecutionNotAuthorised("grant names provider 'auto' (T-1 prohibits it)")

    stated = payload.get("caps_sha256")
    expected = caps_sha256(payload)
    if not isinstance(stated, str) or stated != expected:
        raise ExecutionNotAuthorised(
            "caps_sha256 does not match the caps it claims to sign: the grant was edited "
            "after signing, or never signed"
        )

    return ExecutionGrant(
        owner_signature=signature.strip(),
        date=date,
        call_ceiling=ceiling,
        cost_ceiling_usd=float(cost),
        models={k: str(v) for k, v in payload["models"].items()},
        providers={k: str(v) for k, v in payload["providers"].items()},
    )


def require_grant() -> ExecutionGrant:
    return load_grant()


@dataclass(frozen=True)
class LocalGrant:
    """A single-call local grant (L3). `grants_execution` is the adapter's key."""

    report_id: str
    owner_signature: str
    date: str
    call_ceiling: int
    cost_ceiling_usd: float
    endpoint: str
    model_path: str
    model_sha256: str
    grants_execution: bool = True


def l3_caps_sha256(payload: dict) -> str:
    return sha256_text(canonical({k: payload.get(k) for k in L3_CAPS_FIELDS}))


def load_l3_grant(path: Path | None = None) -> LocalGrant:
    """Validate the L3 grant. Absent, edited or over-ceiling means no permission."""
    target = path or L3_GRANT_PATH
    if not target.exists():
        raise ExecutionNotAuthorised(
            "L3 is not granted: local inference needs its own owner authorisation, "
            "and P2's grant does not cover it"
        )
    try:
        payload = json.loads(target.read_text(encoding="utf-8"))
    except Exception as exc:
        raise ExecutionNotAuthorised(f"L3 grant is not readable as JSON: {exc}") from exc
    if not isinstance(payload, dict):
        raise ExecutionNotAuthorised("L3 grant must be a JSON object")
    missing = [f for f in L3_REQUIRED if f not in payload]
    if missing:
        raise ExecutionNotAuthorised(f"L3 grant is incomplete: missing {missing}")

    signature = payload.get("owner_signature")
    if not isinstance(signature, str) or len(signature.strip()) < 8:
        raise ExecutionNotAuthorised("L3 grant carries no usable owner signature")

    ceiling = payload.get("call_ceiling")
    if isinstance(ceiling, bool) or not isinstance(ceiling, int) or not 1 <= ceiling <= HARD_CEILING:
        raise ExecutionNotAuthorised(f"L3 call_ceiling {ceiling!r} is outside 1..{HARD_CEILING}")

    cost = payload.get("cost_ceiling_usd")
    if isinstance(cost, bool) or not isinstance(cost, (int, float)) or float(cost) > MAX_PAID_GRANT_USD:
        raise ExecutionNotAuthorised("L3 grant may not authorise any spend")

    from .local_adapter import LOOPBACK_HOSTS
    from urllib.parse import urlparse

    endpoint = payload.get("endpoint")
    if not isinstance(endpoint, str) or (urlparse(endpoint).hostname or "").lower() not in LOOPBACK_HOSTS:
        raise ExecutionNotAuthorised(f"L3 endpoint {endpoint!r} is not loopback")

    digest = payload.get("model_sha256")
    if not isinstance(digest, str) or len(digest) != 64 or any(c not in "0123456789abcdef" for c in digest):
        raise ExecutionNotAuthorised("L3 grant needs a lowercase hex sha256 for the model")

    stated = payload.get("caps_sha256")
    if not isinstance(stated, str) or stated != l3_caps_sha256(payload):
        raise ExecutionNotAuthorised("L3 caps_sha256 does not match the caps it claims to sign")

    return LocalGrant(
        report_id=str(payload["report_id"]),
        owner_signature=signature.strip(),
        date=str(payload["date"]),
        call_ceiling=ceiling,
        cost_ceiling_usd=float(cost),
        endpoint=endpoint,
        model_path=str(payload["model_path"]),
        model_sha256=digest,
    )
