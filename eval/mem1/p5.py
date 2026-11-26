"""P5 as its own case, in every arm (plan §4.4, §7).

Gate B left P5 — *"Who trained the horse?"* — as an **open failure**, and the plan
forbids the subject inheriting it. So P5 is not folded into the 13-probe quota or
any aggregate: it is a standalone case whose result is reported on its own.

This module proves the instrument can keep it observable, including the specific
failure mode: **agent/patient confusion** (answering "the horse trained the smith").
That answer must be scored `INCORRECT` and must never be memory-supported.

The case is a deterministic fixture. It builds the knowledge record directly so
that the *reporting* path is under test rather than the extractor; formation is
covered separately in the main self-check. G4 must run the real P5 case through
the same path in every arm.
"""
from __future__ import annotations

from dataclasses import dataclass, field

from eval.s3.context import Window
from eval.s3.evidence import EvidenceLog

from .assemble import WRAPPERS, AssembledBlock
from .errors import DesignViolationError
from .fixtures import DECOYS, FILLER, TURN_TOKENS
from .model import MemoryStore
from .placebo import DecoyPool, build_placebo_measured
from .probe import ProbeSpec, assert_span_departed
from .retrieve import retrieve
from .score import classify, credit, m2_status, support_check

P5_QUERY = "Who trained the horse?"
P5_FACT = "The smith trained the horse at the ford."
P5_KEY = "smith"
P5_CONFUSED = "The horse trained the smith at the ford."
P5_CORRECT = "The smith trained the horse at the ford."
P5_PREDICATE = "trained_animal"
P5_VALUE = "horse trained by smith"
P5_OWNER = "char-001"
PROBE_TURN = 12


@dataclass
class P5Result:
    arm: str
    question: str
    response: str
    m1_category: str
    m1_credit: int
    m2_status: str
    selected_ids: list[str] = field(default_factory=list)
    supporting_span: str = ""
    span_absent: bool = False
    agent_patient_confusion: bool = False
    flags: list[str] = field(default_factory=lambda: ["p5_case", "reported_standalone"])

    def as_dict(self) -> dict:
        return {"arm": self.arm, "question": self.question, "response": self.response,
                "m1_category": self.m1_category, "m1_credit": self.m1_credit,
                "m2_status": self.m2_status, "selected_ids": list(self.selected_ids),
                "supporting_span": self.supporting_span, "span_absent": self.span_absent,
                "agent_patient_confusion": self.agent_patient_confusion,
                "flags": list(self.flags)}


def p5_probe() -> ProbeSpec:
    return ProbeSpec(
        probe_id="P5", turn=PROBE_TURN, category="character_specific", question=P5_QUERY,
        owner_id=P5_OWNER, key_terms=("trained", "horse"), key_predicate=P5_PREDICATE,
        key_value=P5_VALUE, acceptable=(P5_KEY,),
        contradiction_terms=("horse trained", "horse taught", "smith was trained"),
        support_span_hint="sp-002-01", context="world_state", temporal="current",
    )


def _log() -> EvidenceLog:
    log = EvidenceLog()
    log.append("You", "Has anyone been through the ford lately?", TURN_TOKENS)
    log.append("Narrator", P5_FACT, TURN_TOKENS)
    for i in range(3, PROBE_TURN):
        log.append(("You", "Sable", "Alaric")[i % 3], FILLER[i % len(FILLER)], TURN_TOKENS)
    return log


def _store(log: EvidenceLog) -> MemoryStore:
    store = MemoryStore(evidence=log, conversation_id="fixture-p5")
    # No world is declared for this fixture, so the narrator fact belongs to the
    # chat container (chat_state), not to a world. It remains visible in the
    # world_state context, which is what P5's claims depend on.
    store.add_event(turn=2, via="witnessed", owner_id="chat", owner_kind="chat",
                    kind="fact", subject="Narrator", predicate=P5_PREDICATE,
                    value=P5_VALUE, source_span="sp-002-01", confidence=0.9,
                    scope="chat_state", audience=())
    return store


def run_p5(arm: str, responder, counter, tolerance: int = 3,
           decoy_pool: DecoyPool | None = None) -> P5Result:
    """Run the P5 case standalone. No model call, no aggregate contribution."""
    probe = p5_probe()
    log = _log()
    window = Window(capacity=200, reserve=50)
    for turn in log.turns():
        window.add(turn)

    store = MemoryStore(evidence=log, conversation_id="fixture-p5") if arm != "A0" else \
        MemoryStore(evidence=None, conversation_id="fixture-p5")
    selected: list[str] = []
    block = AssembledBlock(text="", citation_ids=[], record_count=0, shape=(), arm=arm)

    if arm in ("A1", "A2", "C"):
        if arm == "A2":
            raise DesignViolationError("P5 must not be run in A2 while its proposals are unapproved")
        store = _store(log)
        retrieval = retrieve(P5_QUERY, store, P5_OWNER, context="world_state",
                             max_records=3, max_chars=600,
                             active_state=None, requested_scenario="")
        if not retrieval.abstained:
            from .assemble import assemble
            block = assemble(retrieval.ids(), store, WRAPPERS["labelled_evidence"], arm=arm)
            selected = list(block.citation_ids)
    elif arm == "D":
        # Build the reference through the same assembler A1 uses, so the shape and
        # token count compared are the real ones rather than a hand-written guess.
        store = _store(log)
        from .assemble import assemble as _assemble
        retrieval = retrieve(P5_QUERY, store, P5_OWNER, context="world_state",
                             max_records=3, max_chars=600,
                             active_state=None, requested_scenario="")
        reference = _assemble(retrieval.ids(), store, WRAPPERS["labelled_evidence"], arm="A1")
        pool = decoy_pool or DecoyPool([{"decoy_id": f"dec-p5-{i:02d}", "text": t,
                                         "tokens": counter.count(t)}
                                        for i, t in enumerate(DECOYS[:4], start=1)])
        block, _report = build_placebo_measured(reference, pool,
                                                WRAPPERS["labelled_evidence"], counter, tolerance)
        selected = list(block.citation_ids)

    departure = assert_span_departed(window, probe.support_span_hint, probe)
    response = responder.respond(probe, block.text)
    category = classify(response, probe)
    supported = False if arm == "D" else support_check(store, selected, probe)
    status = m2_status(arm, probe, category, supported, departure["absent"])
    confused = any(term in response.lower() for term in probe.contradiction_terms)

    return P5Result(arm=arm, question=P5_QUERY, response=response, m1_category=category,
                    m1_credit=credit(category, probe), m2_status=status,
                    selected_ids=selected, supporting_span=probe.support_span_hint,
                    span_absent=departure["absent"], agent_patient_confusion=confused)


class P5Responder:
    """Scripted: correct or deliberately agent/patient-confused. Never networked."""

    network = False

    def __init__(self, confused: bool = False) -> None:
        self.confused = confused
        self.calls = 0

    def respond(self, probe, prompt: str) -> str:
        self.calls += 1
        return P5_CONFUSED if self.confused else P5_CORRECT


__all__ = ["P5_QUERY", "P5_KEY", "P5_CORRECT", "P5_CONFUSED", "P5Result", "p5_probe",
           "run_p5", "P5Responder", "PROBE_TURN"]
