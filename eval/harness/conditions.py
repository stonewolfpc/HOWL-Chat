"""Test-condition construction and separation (row 479 §6; row 482 §2).

A = assisted (canonical fact present). B = unassisted (fact absent).
The two conditions differ in **one line of the system message and nothing
else**: identical user turn, identical params, identical model, same reps.

`assert_separated` refuses any construction where the difference is wider,
because a wider difference is a different experiment (row 470 §4 disclosure).
"""

from __future__ import annotations

from dataclasses import dataclass

from .adapter import canonical, sha256_text
from .fixtures import Fixture

# Frozen for S1. These are experiment settings, not claims about defaults.
FROZEN_PARAMS: dict = {"temperature": 0.8, "top_p": 0.95, "max_tokens": 256}

CONDITION_ASSISTED = "A"
CONDITION_UNASSISTED = "B"


class SeparationError(ValueError):
    """Raised when conditions are not a clean single-line contrast."""


@dataclass(frozen=True)
class Condition:
    name: str
    messages: tuple[dict, ...]
    params: dict
    receives_state: bool

    def system_text(self) -> str:
        return self.messages[0]["content"]

    def user_text(self) -> str:
        return self.messages[-1]["content"]

    def prompt_hash(self) -> str:
        """Same formula as GenerateRequest.prompt_hash, usable before a model is named."""
        return sha256_text(canonical({"messages": [dict(m) for m in self.messages]}))


def base_system_text(fixture: Fixture) -> str:
    return f"You are {fixture.character_name}."


def fact_block(fixture: Fixture) -> str:
    return f"Known: {fixture.fact_text}"


def build_conditions(fixture: Fixture) -> tuple[Condition, Condition]:
    base = base_system_text(fixture)
    user = {"role": "user", "content": fixture.probe.user_message}
    assisted = Condition(
        name=CONDITION_ASSISTED,
        messages=({"role": "system", "content": base + "\n" + fact_block(fixture)}, user),
        params=dict(FROZEN_PARAMS),
        receives_state=True,
    )
    unassisted = Condition(
        name=CONDITION_UNASSISTED,
        messages=({"role": "system", "content": base}, user),
        params=dict(FROZEN_PARAMS),
        receives_state=False,
    )
    assert_separated(assisted, unassisted, fixture)
    return assisted, unassisted


def diff_conditions(a: Condition, b: Condition) -> dict:
    """Structural difference, for the record."""
    out: dict = {"differing_indices": [], "params_equal": a.params == b.params}
    for i in range(max(len(a.messages), len(b.messages))):
        ma = a.messages[i] if i < len(a.messages) else None
        mb = b.messages[i] if i < len(b.messages) else None
        if ma != mb:
            out["differing_indices"].append(
                {
                    "index": i,
                    "a_role": (ma or {}).get("role"),
                    "b_role": (mb or {}).get("role"),
                    "a_content": (ma or {}).get("content"),
                    "b_content": (mb or {}).get("content"),
                }
            )
    return out


def assert_separated(a: Condition, b: Condition, fixture: Fixture) -> None:
    """Only the state block may differ. Anything else is a different test."""
    if len(a.messages) != len(b.messages):
        raise SeparationError("conditions must have the same number of messages")
    if a.params != b.params:
        raise SeparationError("conditions must share identical sampling params")
    if [m["role"] for m in a.messages] != [m["role"] for m in b.messages]:
        raise SeparationError("conditions must share identical role sequence")
    # User turn identical.
    if a.user_text() != b.user_text():
        raise SeparationError("conditions must share an identical user turn")
    if a.user_text() != fixture.probe.user_message:
        raise SeparationError("user turn drifted from the fixture probe")
    # System message: B must equal A with the fact block removed, exactly.
    expected_a = base_system_text(fixture) + "\n" + fact_block(fixture)
    if a.system_text() != expected_a:
        raise SeparationError("assisted system message is not base + fact block")
    if b.system_text() != base_system_text(fixture):
        raise SeparationError("unassisted system message must be the bare base text")
    if not a.receives_state or b.receives_state:
        raise SeparationError("condition state flags are inverted")
