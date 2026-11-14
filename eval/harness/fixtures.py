"""Synthetic fixtures (row 478 §4).

These are engineering fixtures, not corpus items. They score nothing, belong
to no family, and must never be cited as MEMORY-09 units.

`F-ST-01-S1` is the single P1 fixture. Its canonical fact is *invented*, so
correct-by-priors (C3) is excluded by construction rather than by judgement:
no model can know the name from training.
"""

from __future__ import annotations

from dataclasses import dataclass
from typing import Any

from .adapter import canonical, sha256_text
from .state import State

FIXTURE_ID = "F-ST-01-S1"
FIXTURE_VERSION = "1"
DIMENSION = "character knowledge separation (invented fact)"

# Invented, unguessable. Any correct mention is state-supplied.
CANONICAL_FACT_ID = "fact-001"
CANONICAL_FACT_TEXT = "The person who trained me was Corvin Ashgrave of Ossary."

CHARACTER_ID = "char-001"
CHARACTER_NAME = "Sable"
PERSONA_ID = "persona-001"
WORLD_ID = "world-001"


class FixtureError(ValueError):
    """Raised when a fixture is malformed or altered."""


@dataclass(frozen=True)
class Probe:
    user_message: str
    must_include: tuple[str, ...]
    must_exclude: tuple[str, ...]
    failure_signature: str


@dataclass(frozen=True)
class Fixture:
    fixture_id: str
    version: str
    dimension: str
    state: State
    character_id: str
    character_name: str
    persona_id: str
    fact_id: str
    fact_text: str
    probe: Probe

    def to_dict(self) -> dict:
        return {
            "fixture_id": self.fixture_id,
            "version": self.version,
            "dimension": self.dimension,
            "state": self.state.data(),
            "character_id": self.character_id,
            "character_name": self.character_name,
            "persona_id": self.persona_id,
            "fact_id": self.fact_id,
            "fact_text": self.fact_text,
            "probe": {
                "user_message": self.probe.user_message,
                "must_include": list(self.probe.must_include),
                "must_exclude": list(self.probe.must_exclude),
                "failure_signature": self.probe.failure_signature,
            },
        }

    def sha256(self) -> str:
        return sha256_text(canonical(self.to_dict()))


def build_f_st_01_s1() -> Fixture:
    """Deterministic ids so diffs stay readable and hashes stable."""
    st = State()
    st.put(
        "character",
        CHARACTER_ID,
        {"id": CHARACTER_ID, "name": CHARACTER_NAME, "known_facts": [CANONICAL_FACT_ID], "voice_notes": ""},
    )
    st.put("persona", PERSONA_ID, {"id": PERSONA_ID, "name": "the player", "facts": []})
    st.put("world", WORLD_ID, {"id": WORLD_ID, "canon_facts": ["Ossary is a river town."]})
    st.put(
        "fact",
        CANONICAL_FACT_ID,
        {
            "id": CANONICAL_FACT_ID,
            "text": CANONICAL_FACT_TEXT,
            "subject": CHARACTER_ID,
            "source": "authored",
            "established_at_turn": 0,
            "confidence": "canonical",
        },
    )
    probe = Probe(
        user_message="Remind me — who trained you?",
        must_include=("Corvin Ashgrave", "Ossary"),
        must_exclude=(),
        failure_signature=(
            "names a different mentor, hedges as unknown, or (in the unassisted "
            "condition) produces the invented name"
        ),
    )
    return Fixture(
        fixture_id=FIXTURE_ID,
        version=FIXTURE_VERSION,
        dimension=DIMENSION,
        state=st,
        character_id=CHARACTER_ID,
        character_name=CHARACTER_NAME,
        persona_id=PERSONA_ID,
        fact_id=CANONICAL_FACT_ID,
        fact_text=CANONICAL_FACT_TEXT,
        probe=probe,
    )


def from_dict(payload: dict[str, Any]) -> Fixture:
    """Rebuild a fixture from JSON; used to prove round-trip hash stability."""
    if payload.get("fixture_id") != FIXTURE_ID:
        raise FixtureError(f"unexpected fixture_id: {payload.get('fixture_id')!r}")
    state = State(payload["state"])
    probe = payload["probe"]
    return Fixture(
        fixture_id=payload["fixture_id"],
        version=payload["version"],
        dimension=payload["dimension"],
        state=state,
        character_id=payload["character_id"],
        character_name=payload["character_name"],
        persona_id=payload["persona_id"],
        fact_id=payload["fact_id"],
        fact_text=payload["fact_text"],
        probe=Probe(
            user_message=probe["user_message"],
            must_include=tuple(probe["must_include"]),
            must_exclude=tuple(probe["must_exclude"]),
            failure_signature=probe["failure_signature"],
        ),
    )


def assert_declared(fixture: Fixture) -> None:
    """Guard rails: the fixture must still be what the schedule authorised."""
    if fixture.fixture_id != FIXTURE_ID or fixture.version != FIXTURE_VERSION:
        raise FixtureError("fixture identity/version drifted from the authorised schedule")
    present, source, turn = fixture.state.assert_fact(fixture.fact_id)
    if not present:
        raise FixtureError("canonical fact absent from fixture state")
    if source != "authored":
        raise FixtureError(f"canonical fact must be authored, got {source!r}")
    if turn != 0:
        raise FixtureError("canonical fact must be established at turn 0")
    if "Corvin Ashgrave" not in fixture.fact_text:
        raise FixtureError("canonical fact text no longer contains the invented name")
