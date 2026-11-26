"""Synthetic G1 fixtures — NOT experimental cases.

These conversations exist to exercise the instrument and prove its refusals.
They are deliberately not drawn from, and must never be scored as, the
experimental case set: G1 produces no results (user G1 §4).

The scenario carries every shape prereg §5 requires: 48 turns, 13 probes with
the frozen quota, planted facts of each class, a correction after rollover, a
retraction, a deliberate knowledge asymmetry, an unanswerable probe, and enough
displacement for at least three rollovers before the first probe.
"""
from __future__ import annotations

from dataclasses import dataclass, field, replace

from . import DEFAULT_CAPACITY, DEFAULT_RESERVE
from .distribution import ExpectedScope
from .errors import DesignViolationError
from .model import CHAT_OWNER
from .probe import ProbeSpec

SABLE, ALARIC = "char-001", "char-002"
OWNER_OF = {"Sable": SABLE, "Alaric": ALARIC}

TURN_TOKENS = 55          # declared per turn: enough to force rollovers
FILLER = [
    "The rain came in sideways off the moor and nobody wanted to speak first.",
    "We walked the north road until the light went grey and the horses refused.",
    "Someone had left bread and hard cheese on the table and it did not last.",
    "There is a particular silence in a room where two people are deciding not to argue.",
    "By morning the fire had burned down to a red line and the kettle was cold.",
    "The innkeeper counted coins twice, then a third time, and still looked unhappy.",
]


@dataclass
class TurnSpec:
    speaker: str
    text: str
    tokens: int = TURN_TOKENS


@dataclass
class ConversationSpec:
    conversation_id: str
    characters: dict[str, str]
    turns: list[TurnSpec]
    probes: list[ProbeSpec]
    capacity: int = DEFAULT_CAPACITY
    reserve: int = DEFAULT_RESERVE
    notes: str = ""
    extra: dict = field(default_factory=dict)
    # Declared world/scenario state, applied at the turn boundary in the driver.
    # Empty means no gate at all (D1): the conversation behaves as it always has.
    state_changes: list = field(default_factory=list)
    # Explicit scenario -> parent world map. Declared, never inferred from content,
    # and required for every scenario that appears in `state_changes`.
    scenario_parents: dict = field(default_factory=dict)

    def __post_init__(self) -> None:
        """Ownership referential integrity, checked before any turn can run.

        A character id that collides with the reserved world owner id, or a probe
        naming an owner the conversation never declared, would otherwise produce
        records and receipts whose ownership cannot be verified after the fact.
        """
        ids = list(self.characters.values())
        if not ids:
            raise DesignViolationError(
                f"{self.conversation_id}: no character ids declared")
        if len(set(ids)) != len(ids):
            raise DesignViolationError(
                f"{self.conversation_id}: duplicate character ids {ids}")
        if CHAT_OWNER in ids:
            raise DesignViolationError(
                f"{self.conversation_id}: character id {CHAT_OWNER!r} is the conversation "
                "container, not a character: narrator-derived records belong to the "
                "chat when no world is active and to the world when one is")
        known = set(ids)
        for probe in self.probes:
            if probe.owner_id not in known:
                raise DesignViolationError(
                    f"{self.conversation_id}/{probe.probe_id}: probe owner "
                    f"{probe.owner_id!r} is not a declared character {sorted(known)}")
        seen: set[tuple[int, str, str]] = set()
        for change in self.state_changes:
            if change.turn > len(self.turns):
                raise DesignViolationError(
                    f"{self.conversation_id}: state change at turn {change.turn} is "
                    f"beyond the {len(self.turns)} turns declared")
            key = (change.turn, change.owner_kind, change.owner_id)
            if key in seen:
                raise DesignViolationError(
                    f"{self.conversation_id}: duplicate state change {key}")
            seen.add(key)

        self._validate_scope_declarations()

    def _validate_scope_declarations(self) -> None:
        """Identity, parent mapping, and the single-active-world rule.

        Every scenario needs exactly one declared parent world, the parent must be
        a world this case actually declares, and no declared sequence may resolve
        to two active worlds at any turn boundary. Same-turn transitions are applied
        as one set, so a swap such as `A-, B+` is valid while `A+, B+` is not.
        """
        for scenario, world in self.scenario_parents.items():
            if not scenario or not world:
                raise DesignViolationError(
                    f"{self.conversation_id}: scenario parent mapping needs non-empty ids")
            if CHAT_OWNER in (scenario, world):
                raise DesignViolationError(
                    f"{self.conversation_id}: {CHAT_OWNER!r} is the chat container; "
                    "it cannot identify a scenario or world")
            if scenario == world:
                raise DesignViolationError(
                    f"{self.conversation_id}: scenario {scenario!r} cannot be its own "
                    "parent world")

        declared_worlds = {c.owner_id for c in self.state_changes
                           if c.owner_kind == "world"}
        declared_scenarios = {c.owner_id for c in self.state_changes
                              if c.owner_kind == "scenario"}
        collision = declared_worlds & declared_scenarios
        if collision:
            raise DesignViolationError(
                f"{self.conversation_id}: ids {sorted(collision)} are used as both a "
                "world and a scenario; ownership would be ambiguous")
        for scenario in sorted(declared_scenarios):
            if scenario not in self.scenario_parents:
                raise DesignViolationError(
                    f"{self.conversation_id}: scenario {scenario!r} has no declared "
                    "parent world")
        for scenario, world in sorted(self.scenario_parents.items()):
            if world not in declared_worlds:
                raise DesignViolationError(
                    f"{self.conversation_id}: parent world {world!r} of scenario "
                    f"{scenario!r} is not declared by this case")

        active: set[str] = set()
        for turn in sorted({c.turn for c in self.state_changes}):
            for change in [c for c in self.state_changes if c.turn == turn]:
                if change.owner_kind != "world":
                    continue
                if change.active:
                    active.add(change.owner_id)
                else:
                    active.discard(change.owner_id)
            if len(active) > 1:
                raise DesignViolationError(
                    f"{self.conversation_id}: turn {turn} resolves to {len(active)} "
                    f"active worlds {sorted(active)}; a declared sequence may have at "
                    "most one active world at any boundary")

    def owner_of(self) -> dict[str, str]:
        return dict(self.characters)

    def probe_at(self, turn: int) -> ProbeSpec | None:
        for probe in self.probes:
            if probe.turn == turn:
                return probe
        return None


def _script() -> dict[int, tuple[str, str]]:
    """The planted material, by turn. Everything here is synthetic."""
    return {
        1: ("You", "The road south was worse than I remembered. My boots are still wet."),
        2: ("Sable", "Then you came the long way. Sit down before you fall down."),
        3: ("You", "Long way or not, I am here. Has it always been this quiet?"),
        4: ("Sable", "I was trained by Corvin Ashgrave of Ossary."),
        5: ("Sable", "I like tea more than I like surprises."),
        6: ("Sable", "I keep a ledger of unpaid debts in the blue cabinet."),
        7: ("You", "Noted. I will try to stay off that ledger."),
        8: ("Alaric", "The vault code is seven-four-nine-one."),
        9: ("Alaric", "I served four years in the Ninth Company."),
        10: ("Narrator", "The north road is flooded."),
        11: ("You", "Flooded, and I walked it. Of course it is."),
        12: ("Sable", "Between us, the cellar key is under the third stair."),
        13: ("You", "I will pretend I did not hear that."),
        14: ("Sable", "I have never trusted a map that was drawn by someone who has not walked the ground."),
        15: ("Alaric", "Maps lie more often than people do."),
        16: ("Narrator", "The mill wheel is broken."),
        17: ("You", "Then the bread will be late again."),
        18: ("Sable", "I promise to keep the north gate locked at dusk."),
        19: ("Alaric", "Promises are cheaper than gates."),
        20: ("Sable", "Actually, I keep a ledger of debts in the green cabinet."),
        21: ("You", "Green. Understood."),
        22: ("Alaric", "Forget the Ninth Company — I was wrong to say it."),
        23: ("Sable", "That is the first sensible thing you have said all evening."),
    }


PROBE_TURNS = {
    36: ("character_specific", SABLE, "Who trained you?", "Corvin Ashgrave of Ossary",
         "trained_by", "Corvin Ashgrave of Ossary"),
    37: ("character_specific", SABLE, "Where is the ledger kept now?", "green cabinet",
         "keeps_ledger", "ledger in the green cabinet"),
    38: ("character_specific", ALARIC, "What is the vault code?", "seven-four-nine-one",
         "vault_code", "seven-four-nine-one"),
    39: ("unanswerable", SABLE, "What is the name of your sister?", "", "", ""),
    40: ("character_specific", SABLE, "What does Sable like?", "tea", "preference",
         "tea more than I like surprises"),
    41: ("ambient", SABLE, "Is the north road passable?", "north road is flooded",
         "world_state", "north road is flooded"),
    42: ("ambient", SABLE, "Is the mill working?", "mill wheel is broken",
         "world_state", "mill wheel is broken"),
    43: ("restricted", SABLE, "Where is the cellar key?", "under the third stair",
         "cellar_key_location", "under the third stair"),
    44: ("restricted", SABLE, "Remind me where the cellar key is kept.", "under the third stair",
         "cellar_key_location", "under the third stair"),
    45: ("experiential", SABLE, "How do you feel about a map?", "never trusted a map",
         "never", "trusted a map that was drawn by someone who has not walked the ground"),
    46: ("experiential", SABLE, "What did you promise about the gate?", "keep the north gate locked",
         "promise", "to keep the north gate locked at dusk"),
    47: ("correction", SABLE, "Where was the ledger kept before?", "blue cabinet",
         "keeps_ledger", "ledger in the blue cabinet"),
    48: ("retraction", ALARIC, "How long did you serve in the Ninth Company?", "retracted", "served",
         "four years in the Ninth Company"),
}

# Span hints are resolved at run time from the formation store; these are turns.
SUPPORT_TURNS = {36: 4, 37: 20, 38: 8, 40: 5, 41: 10, 42: 16, 43: 12, 44: 12,
                 45: 14, 46: 18, 47: 6, 48: 9}

# Declared expected scope for every event this fixture plants (Stage 1, gap G4).
# Declarations come from the case design above, never from inspecting the store.
EXPECTED_SCOPES = (
    ExpectedScope(4, "Sable", "trained_by", "owner"),
    ExpectedScope(5, "Sable", "preference", "owner"),
    ExpectedScope(6, "Sable", "keeps_ledger", "owner"),
    ExpectedScope(8, "Alaric", "vault_code", "owner"),
    ExpectedScope(9, "Alaric", "served", "owner"),
    # No world is declared for this fixture, so narrator facts belong to the chat
    # container and are visible as chat_state, not as world state.
    ExpectedScope(10, "Narrator", "world_state", "chat_state"),
    ExpectedScope(12, "Sable", "cellar_key_location", "restricted"),
    ExpectedScope(14, "Sable", "never", "owner"),
    ExpectedScope(16, "Narrator", "world_state", "chat_state"),
    ExpectedScope(18, "Sable", "promise", "owner"),
    ExpectedScope(20, "Sable", "keeps_ledger", "owner"),
    ExpectedScope(22, "Alaric", "retraction:served", "owner"),
)


def demo_conversation(conversation_id: str = "fixture-demo-01") -> ConversationSpec:
    script = _script()
    turns: list[TurnSpec] = []
    filler_index = 0
    for turn_no in range(1, 49):
        if turn_no in script:
            speaker, text = script[turn_no]
        elif turn_no in PROBE_TURNS:
            _cat, _owner, question, *_rest = PROBE_TURNS[turn_no]
            speaker, text = "You", question
        else:
            speaker = ("You", "Sable", "Alaric")[filler_index % 3]
            text = FILLER[filler_index % len(FILLER)]
            filler_index += 1
        turns.append(TurnSpec(speaker=speaker, text=text))

    probes: list[ProbeSpec] = []
    for turn_no, (category, owner, question, key_value, predicate, record_value) in PROBE_TURNS.items():
        if category == "unanswerable":
            probe = ProbeSpec(probe_id=f"p{turn_no:02d}", turn=turn_no, category=category,
                              question=question, owner_id=owner, answerable=False,
                              key_predicate="", key_value="",
                              contradiction_terms=("sister", "her name is"))
        elif category == "retraction":
            probe = ProbeSpec(probe_id=f"p{turn_no:02d}", turn=turn_no, category=category,
                              question=question, owner_id=owner, key_predicate=predicate,
                              key_value=key_value, acceptable=("retracted", "withdrawn",
                                                                "no longer", "was wrong"),
                              contradiction_terms=("four years",), temporal="none",
                              retraction_of=predicate)
        elif category == "correction":
            probe = ProbeSpec(probe_id=f"p{turn_no:02d}", turn=turn_no, category=category,
                              question=question, owner_id=owner, key_predicate=predicate,
                              key_value=key_value, acceptable=("blue", "blue cabinet"),
                              contradiction_terms=("green",), temporal="historical")
        elif category == "ambient":
            probe = ProbeSpec(probe_id=f"p{turn_no:02d}", turn=turn_no, category=category,
                              question=question, owner_id=owner, key_predicate=predicate,
                              key_value=key_value, acceptable=(key_value,),
                              contradiction_terms=(("passable",) if "road" in question
                                                   else ("working",)),
                              context="world_state", temporal="current")
        elif category == "restricted":
            probe = ProbeSpec(probe_id=f"p{turn_no:02d}", turn=turn_no, category=category,
                              question=question, owner_id=owner, key_predicate=predicate,
                              key_value=key_value, acceptable=(key_value,),
                              contradiction_terms=("first stair", "cellar door"),
                              temporal="current")
        elif category == "experiential":
            probe = ProbeSpec(probe_id=f"p{turn_no:02d}", turn=turn_no, category=category,
                              question=question, owner_id=owner, key_predicate=predicate,
                              key_value=key_value, acceptable=(key_value, "distrust",
                                                               "does not trust"),
                              contradiction_terms=(), temporal="current")
        else:
            probe = ProbeSpec(probe_id=f"p{turn_no:02d}", turn=turn_no, category=category,
                              question=question, owner_id=owner, key_predicate=predicate,
                              key_value=key_value, acceptable=(key_value,),
                              contradiction_terms=(), temporal="current")
        support_turn = SUPPORT_TURNS.get(turn_no)
        if support_turn is not None:
            # Span ids are deterministic from the evidence log (turn, sentence index).
            probe = replace(probe, support_span_hint=f"sp-{support_turn:03d}-01")
        probes.append(probe)

    return ConversationSpec(
        conversation_id=conversation_id,
        characters=OWNER_OF,
        turns=turns,
        probes=probes,
        notes="synthetic G1 fixture: exercises quotas, rollover displacement, correction, "
              "retraction, ownership, and abstention. Not an experimental case.",
    )


def answers_for(spec: ConversationSpec, wrong: tuple[str, ...] = ()) -> dict[str, str]:
    """Canned responses for the scripted responder: correct by default."""
    out: dict[str, str] = {}
    for probe in spec.probes:
        if probe.probe_id in wrong:
            out[probe.probe_id] = "I could not possibly comment on any of that."
        elif probe.category == "unanswerable":
            out[probe.probe_id] = "I do not know; no such thing was established here."
        elif probe.category == "retraction":
            out[probe.probe_id] = "That claim was retracted, so I have nothing current to give."
        elif probe.category == "correction":
            out[probe.probe_id] = "Before the change it was kept in the blue cabinet."
        elif probe.category == "ambient":
            out[probe.probe_id] = f"As far as the world goes, the {probe.key_value}."
        elif probe.category == "restricted":
            out[probe.probe_id] = "It is under the third stair, and that stays between us."
        else:
            out[probe.probe_id] = f"It is {probe.key_value}."
    return out


# Short and irrelevant by construction: no probe key term appears here, and each is
# close in length to a bullet so per-turn token matching can succeed.
# Deliberately varied in length: arm D must match A1 on *every* turn, including
# turns where A1 selected two records. No probe key term appears in any decoy.
DECOYS = [
    "The stable boy counts oats twice.",
    "The well rope was replaced last spring.",
    "A grey cat sleeps by the stove all afternoon.",
    "Two travellers argued about the price of a room until the innkeeper intervened.",
    "The chandler orders wax in barrels and keeps the tallow cold through winter.",
    "A dog barked at the eastern fence for an hour and then gave up and slept.",
    "Bread was left cooling on the sill and the butter had already softened in the heat.",
    "The orchard path is muddy but shorter than the road in dry weather, and nobody uses it.",
    "A kettle hissed, then went quiet, and the room settled into the silence after an argument.",
    "Someone mended the eastern fence with wire and a great deal of patience, though it may not hold.",
]


def demo_decoys(counter=None) -> list[dict]:
    """Declared decoy pool. Counts come from the supplied counter, never estimated here."""
    if counter is None:
        from .tokens import WhitespaceCounter
        counter = WhitespaceCounter()
    return [{"decoy_id": f"dec-{i:04d}", "text": text, "tokens": counter.count(text)}
            for i, text in enumerate(DECOYS, start=1)]


__all__ = ["SABLE", "ALARIC", "OWNER_OF", "TURN_TOKENS", "TurnSpec", "ConversationSpec",
           "PROBE_TURNS", "SUPPORT_TURNS", "EXPECTED_SCOPES", "demo_conversation",
           "answers_for", "DECOYS", "demo_decoys"]
