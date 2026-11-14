"""Call plan, caps, and the unassigned-model barrier (row 482 §3-5).

The plan cannot be turned into calls while any model slot is unassigned, and
the slot assignments are an owner decision recorded on the schedule — not a
harness default. This is a second barrier in front of execution, independent
of the grant file.
"""

from __future__ import annotations

from dataclasses import dataclass

from .adapter import AdapterConfigError, GenerateRequest
from .conditions import CONDITION_ASSISTED, CONDITION_UNASSISTED, Condition

SLOTS = ("S1-A", "S1-B", "S1-C")
REPS = 3
CONDITIONS = (CONDITION_ASSISTED, CONDITION_UNASSISTED)

PLANNED_CALLS = len(SLOTS) * len(CONDITIONS) * REPS  # 18
RETRY_ALLOWANCE = 12
HARD_CEILING = 30

UNASSIGNED = "UNASSIGNED"


class PlanIncomplete(RuntimeError):
    """Raised when the plan needs an owner decision it does not have."""


@dataclass(frozen=True)
class Slot:
    name: str
    model_id: str
    provider: str

    @property
    def assigned(self) -> bool:
        return self.model_id not in ("", None, UNASSIGNED) and self.provider not in ("", None, UNASSIGNED)


def default_slots() -> tuple[Slot, ...]:
    """P1 ships no model identities. Naming them is an owner decision."""
    return tuple(Slot(name=s, model_id=UNASSIGNED, provider=UNASSIGNED) for s in SLOTS)


def build_plan(slots: tuple[Slot, ...], conditions: tuple[Condition, Condition]) -> list[dict]:
    """Expand slots x conditions x reps. Refuses unassigned slots."""
    unassigned = [s.name for s in slots if not s.assigned]
    if unassigned:
        raise PlanIncomplete(
            "model slots are unassigned: "
            + ", ".join(unassigned)
            + ". Naming models and providers is an owner decision on the BEA-1 schedule."
        )
    if len(conditions) != len(CONDITIONS):
        raise PlanIncomplete("expected exactly two conditions")
    calls: list[dict] = []
    n = 0
    for slot in slots:
        for condition in conditions:
            for rep in range(1, REPS + 1):
                n += 1
                request = GenerateRequest(
                    messages=condition.messages,
                    model_id=slot.model_id,
                    provider=slot.provider,
                    params=condition.params,
                )
                calls.append(
                    {
                        "call": n,
                        "slot": slot.name,
                        "condition": condition.name,
                        "rep": rep,
                        "model_id": request.model_id,
                        "provider": request.provider,
                        "prompt_hash": request.prompt_hash(),
                        "receives_state": condition.receives_state,
                    }
                )
    if len(calls) != PLANNED_CALLS:
        raise PlanIncomplete(f"plan produced {len(calls)} calls, expected {PLANNED_CALLS}")
    return calls


def caps_summary() -> dict:
    return {
        "planned": PLANNED_CALLS,
        "retry_allowance": RETRY_ALLOWANCE,
        "hard_ceiling": HARD_CEILING,
        "note": "a failed call counts against the ceiling; retries buy no headroom",
    }


def check_caps(attempted: int, ceiling: int = HARD_CEILING) -> None:
    if attempted > ceiling:
        raise PlanIncomplete(f"{attempted} calls exceeds the hard ceiling of {ceiling}")
    if ceiling > HARD_CEILING:
        raise PlanIncomplete(
            f"ceiling {ceiling} exceeds the authorised maximum of {HARD_CEILING}; a higher ceiling needs a new grant"
        )


def validate_slot(slot: Slot) -> None:
    """A named slot must still satisfy T-1/T-2 (no auto, no routing suffix)."""
    if not slot.assigned:
        raise PlanIncomplete(f"{slot.name} unassigned")
    GenerateRequest(
        messages=({"role": "user", "content": "x"},),
        model_id=slot.model_id,
        provider=slot.provider,
        params={},
    )
