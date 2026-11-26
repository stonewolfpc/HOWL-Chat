"""Long-conversation driver (plan §4.4, prereg §2/§4/§5).

The driver is the thing that would execute the frozen design. It builds no
results by itself: it needs a `Responder`, and during G1 the only responder is
scripted. Anything claiming to be a network responder raises
`ExecutionNotAuthorised` — G4 is not authorized.

Enforced here:

  * one prompt capacity for every arm; memory is paid for out of a declared
    share, so a larger memory block shrinks history rather than growing the
    prompt (plan §2.4);
  * per-turn measured accounting, asserted every turn;
  * a probe cannot be scored unless its supporting span has verifiably left the
    window through a recorded rollover;
  * arm D is built from A1's block for the *same* turn, matched on count, tokens
    and shape, with decoys proved irrelevant to the probe;
  * A2 refuses to run at all while its proposals are unapproved;
  * N = 16 and 3 reps are enforced when a run plan is built, not when results
    are convenient.
"""
from __future__ import annotations

import hashlib
from dataclasses import dataclass, field

from eval.s3.context import Window
from eval.s3.evidence import EvidenceLog

from . import ARMS, N_CONVERSATIONS, PLACEBO_TOKEN_TOLERANCE, REPS_PER_CELL
from .assemble import WRAPPERS, AssembledBlock, DEFAULT_WRAPPER, assemble
from .errors import (ArmInertError, DesignViolationError, ExecutionNotAuthorised,
                     PlaceboFidelityError, TokenAccountingError)
from .form import Determiner
from .model import ActiveState, MemoryStore
from .placebo import DecoyPool, assert_irrelevant, build_placebo_measured
from .probe import ProbeReceipt, assert_span_departed, block_sha256, validate_quotas
from .retrieve import DEFAULT_CONTEXT, assert_selection_owned, retrieve
from .score import assert_temporal, classify, credit, m2_status, support_check
from .seeds import seed_for
from .tokens import build_account

MEMORY_BUDGET = 220          # declared share of the prompt that memory may occupy
# Declared assembly budget, identical in every arm (plan §2.2). Fixed here, before
# any run, because "same max records" must be a constant of the instrument.
DECLARED_MAX_RECORDS = 3
DECLARED_MAX_CHARS = 600


@dataclass
class RunConfig:
    arm: str
    rep_index: int
    seed: int
    seed_honoured: bool = True
    tolerance: int = PLACEBO_TOKEN_TOLERANCE
    wrapper: str = DEFAULT_WRAPPER
    a2_proposals_approved: bool = False
    memory_budget: int = MEMORY_BUDGET


class ScriptedResponder:
    """G1 only: canned answers. Never contacts anything."""

    network = False

    def __init__(self, answers: dict[str, str]) -> None:
        self.answers = dict(answers)
        self.calls = 0

    def respond(self, probe, prompt: str) -> str:
        self.calls += 1
        return self.answers.get(probe.probe_id, "")


@dataclass
class ArmResult:
    conversation_id: str
    arm: str
    rep_index: int
    seed: int
    seed_honoured: bool
    receipts: list[ProbeReceipt] = field(default_factory=list)
    accounts: list = field(default_factory=list)
    rollover_count: int = 0
    rollovers: list = field(default_factory=list)
    store_sha256: str = ""
    formation: dict = field(default_factory=dict)
    conversation_sha256: str = ""
    model_calls: int = 0
    vetoes: list[str] = field(default_factory=list)
    flags: list[str] = field(default_factory=list)
    state_transitions: list = field(default_factory=list)
    blocks: dict = field(default_factory=dict)     # turn -> AssembledBlock (for D matching)

    def as_dict(self) -> dict:
        return {"conversation_id": self.conversation_id, "arm": self.arm,
                "rep_index": self.rep_index, "seed": self.seed,
                "seed_honoured": self.seed_honoured,
                "receipts": [r.as_dict() for r in self.receipts],
                "rollover_count": self.rollover_count,
                "store_sha256": self.store_sha256, "formation": self.formation,
                "conversation_sha256": self.conversation_sha256,
                "model_calls": self.model_calls, "vetoes": list(self.vetoes),
                "flags": list(self.flags),
                "state_transitions": list(self.state_transitions)}


def build_run_plan(conversation_ids: list[str]) -> list[dict]:
    """N = 16 paired conversations x 3 matched reps, seeds committed up front."""
    if len(conversation_ids) != N_CONVERSATIONS:
        raise DesignViolationError(
            f"prereg fixes N={N_CONVERSATIONS}; got {len(conversation_ids)} conversations"
        )
    if len(set(conversation_ids)) != len(conversation_ids):
        raise DesignViolationError("conversation ids must be distinct independent units")
    cells = [{"conversation_id": cid, "rep_index": rep, "seed": seed_for(cid, rep)}
             for cid in conversation_ids for rep in range(REPS_PER_CELL)]
    return cells


def plan_sha256(cells: list[dict]) -> str:
    blob = "\n".join(f"{c['conversation_id']}|{c['rep_index']}|{c['seed']}" for c in cells)
    return hashlib.sha256(blob.encode("utf-8")).hexdigest()


class ConversationDriver:
    def __init__(self, spec, config: RunConfig, decoy_pool: DecoyPool | None = None) -> None:
        if config.arm not in ARMS:
            raise DesignViolationError(f"unknown arm {config.arm!r}")
        validate_quotas(spec.probes)
        self.spec = spec
        self.config = config
        self.pool = decoy_pool
        self.wrapper = WRAPPERS[config.wrapper]
        self.evidence = EvidenceLog()
        self.store = self._make_store()
        self.former = Determiner(self.store, self.evidence) if config.arm in ("A1", "A2") else None
        # History gets what memory does not: memory competes, it does not extend.
        self.window = Window(capacity=spec.capacity - config.memory_budget, reserve=spec.reserve)
        # D1: no declared state means no gate. Declared state starts empty, so a
        # scope is inactive until a declared transition activates it (no implicit
        # activation inside a declared sequence).
        self.state_declared = bool(getattr(spec, "state_changes", ()))
        self.parents = dict(getattr(spec, "scenario_parents", {}) or {})
        self.active: ActiveState | None = (
            ActiveState(declared=True, parents=self.parents) if self.state_declared else None)

    def _make_store(self) -> MemoryStore:
        arm = self.config.arm
        if arm == "A0":
            return MemoryStore(evidence=None, conversation_id=self.spec.conversation_id)
        if arm in ("A1", "A2"):
            return MemoryStore(evidence=self.evidence, conversation_id=self.spec.conversation_id)
        if arm == "C":
            store = MemoryStore(evidence=None, conversation_id=self.spec.conversation_id)
            self._inject_known(store)
            return store
        return MemoryStore(evidence=None, conversation_id=self.spec.conversation_id)  # D

    def _inject_known(self, store: MemoryStore) -> None:
        """Arm C: prewritten records. Plumbing only — never counted as formation."""
        injected = [
            ("Sable", "trained_by", "Corvin Ashgrave of Ossary", 0.99),
            ("Sable", "keeps_ledger", "ledger in the blue cabinet", 0.99),
            ("Alaric", "vault_code", "seven-four-nine-one", 0.99),
        ]
        for subject, predicate, value, conf in injected:
            store.add_event(turn=0, via="narrator_grant", owner_id=self.spec.characters[subject],
                            kind="fact", subject=subject, predicate=predicate, value=value,
                            source_span=f"injected:{subject}:{predicate}:{value[:12]}",
                            confidence=conf, formation_method="injected")

    # -- main loop --------------------------------------------------------
    def run(self, responder, reference_blocks: dict[int, AssembledBlock] | None = None) -> ArmResult:
        if getattr(responder, "network", False):
            raise ExecutionNotAuthorised(
                "a network responder requires G4 authorization, which has not been granted"
            )
        arm = self.config.arm
        if arm == "A2" and not self.config.a2_proposals_approved:
            raise ArmInertError(
                "arm A2 is inert until its proposals are approved; it cannot run as validation"
            )

        result = ArmResult(conversation_id=self.spec.conversation_id, arm=arm,
                           rep_index=self.config.rep_index, seed=self.config.seed,
                           seed_honoured=self.config.seed_honoured)
        owner_of = self.spec.owner_of()

        for turn_no, turn_spec in enumerate(self.spec.turns, start=1):
            # Declared state applies at the turn boundary, before this turn's
            # formation and before any probe on this turn is evaluated.
            self._apply_transitions(turn_no, result)
            turn = self.evidence.append(turn_spec.speaker, turn_spec.text, turn_spec.tokens)
            rollover = self.window.add(turn)
            if rollover is not None:
                result.rollovers.append(rollover.as_dict())

            if self.former is not None:
                self.former.run_pending(owner_of, self._current_world())

            probe = self.spec.probe_at(turn_no)
            perspective = probe.owner_id if probe else self.spec.characters["Sable"]
            context = _context_for(probe)
            query = turn_spec.text

            block, retr = self._block_for_turn(query, perspective, context, arm,
                                               reference_blocks, turn_no, probe)
            result.blocks[turn_no] = block

            system_text = f"You are {_name_of(self.spec, perspective)}."
            account = build_account(
                counter=_counter(), arm=arm, turn=turn_no, capacity=self.spec.capacity,
                reserve=self.spec.reserve, system_text=system_text,
                history_tokens=self.window.tokens_used(),
                memory_text=block.bullets_text, wrapper_text=block.wrapper_text,
            )
            share = account.system + account.memory + account.overhead
            if share > self.config.memory_budget:
                raise TokenAccountingError(
                    f"turn {turn_no} ({arm}): memory+system+overhead {share} exceeds the declared "
                    f"share {self.config.memory_budget}: memory would extend the prompt"
                )
            result.accounts.append(account)

            if probe is None:
                continue

            span_absent = False
            if probe.support_span_hint:
                departure = assert_span_departed(self.window, probe.support_span_hint, probe)
                span_absent = departure["absent"]

            # Arm D cites decoys, not records: ownership/temporal vetoes and the
            # support check apply to records, so they are skipped for D only.
            if arm != "D":
                assert_selection_owned(self.store, perspective, block.citation_ids)
                assert_temporal(self.store, probe, block.citation_ids)

            response = responder.respond(probe, block.text)
            category = classify(response, probe)
            supported = False if arm == "D" else support_check(
                self.store, block.citation_ids, probe)
            status = m2_status(arm, probe, category, supported, span_absent)

            receipt = ProbeReceipt(
                conversation_id=self.spec.conversation_id, arm=arm,
                rep_index=self.config.rep_index, seed=self.config.seed,
                seed_honoured=self.config.seed_honoured, probe_id=probe.probe_id,
                turn=turn_no, category=probe.category, question=probe.question,
                response_text=response, m1_category=category,
                m1_credit=credit(category, probe), m2_status=status,
                selected_ids=list(block.citation_ids), abstained=(retr is None or retr.abstained),
                retrieval_reason=(retr.reason if retr else "no_memory_arm"),
                supporting_span=probe.support_span_hint, span_absent=span_absent,
                block_sha256=block_sha256(block.text) if block.text else "",
                tokens=account.as_dict(),
                flags=(["plumbing_only", "excluded_from_primary"] if arm == "C" else []),
            )
            if arm in ("A1", "A2") and block.text and not block.citation_ids:
                receipt.flags.append("block_without_citation")
            receipt.validate()
            result.receipts.append(receipt)

        if len(result.rollovers) < 3:
            raise DesignViolationError(
                f"{self.spec.conversation_id}: {len(result.rollovers)} rollovers; prereg requires >= 3"
            )
        result.rollover_count = len(result.rollovers)
        result.store_sha256 = self.store.snapshot_sha256()
        result.formation = self.former.stats.as_dict() if self.former else {}
        result.conversation_sha256 = self.evidence.sha256()
        result.model_calls = 0
        return result

    def _current_world(self) -> str | None:
        """The single active world at this boundary, or None.

        Used for narrator attribution. The single-active-world rule guarantees at
        most one; if a caller ever hands the driver an inconsistent state, this
        returns None rather than picking one arbitrarily.
        """
        if self.active is None or not self.active.declared:
            return None
        worlds = sorted(self.active.worlds)
        return worlds[0] if len(worlds) == 1 else None

    def _apply_transitions(self, turn: int, result: ArmResult) -> None:
        """Apply all state scheduled for this turn as one atomic transition.

        Transitions are declarations about eligibility, so they must leave the
        record stream untouched: the store hash is compared before and after, and a
        transition that changed it is a design violation, not a side effect. The
        whole turn's changes are resolved before the state is committed, so a swap
        such as `A-, B+` never passes through a two-world intermediate, and no probe
        or formation step can observe a partial transition.
        """
        changes = [c for c in getattr(self.spec, "state_changes", ()) if c.turn == turn]
        if not changes:
            return
        before_sha = self.store.snapshot_sha256()

        def _counts(kind: str, oid: str) -> int:
            return sum(1 for r in self.store.all()
                       if r.owner_kind == kind and r.owner_id == oid)

        before_counts = {f"{c.owner_kind}/{c.owner_id}": _counts(c.owner_kind, c.owner_id)
                         for c in changes}
        worlds = set(self.active.worlds) if self.active else set()
        scenarios = set(self.active.scenarios) if self.active else set()
        for change in changes:
            target = worlds if change.owner_kind == "world" else scenarios
            if change.active:
                target.add(change.owner_id)
            else:
                target.discard(change.owner_id)
        if len(worlds) > 1:
            raise DesignViolationError(
                f"turn {turn}: transition resolves to {len(worlds)} active worlds "
                f"{sorted(worlds)}; at most one may be active")

        # Preservation is asserted before the new state is committed: a transition
        # that touched the record stream must leave the previous state in place, so
        # the check never runs against an already-committed state.
        after_counts = {f"{c.owner_kind}/{c.owner_id}": _counts(c.owner_kind, c.owner_id)
                        for c in changes}
        after_sha = self.store.snapshot_sha256()
        if after_sha != before_sha:
            raise DesignViolationError(
                f"turn {turn}: scope transition changed the memory store")

        self.active = ActiveState(declared=True, worlds=frozenset(worlds),
                                  scenarios=frozenset(scenarios), parents=self.parents)
        result.state_transitions.append({
            "turn": turn,
            "changes": [{"owner_kind": c.owner_kind, "owner_id": c.owner_id,
                         "active": c.active} for c in changes],
            "resolved": self.active.as_dict(),
            "scope_records_before": before_counts,
            "scope_records_after": after_counts,
            "store_sha256_before": before_sha, "store_sha256_after": after_sha,
        })

    def _block_for_turn(self, query: str, perspective: str, context: str, arm: str,
                        reference_blocks: dict[int, AssembledBlock] | None,
                        turn_no: int, probe) -> tuple[AssembledBlock, object]:
        if arm == "A0":
            return AssembledBlock(text="", citation_ids=[], record_count=0, shape=(), arm="A0"), None
        if arm == "D":
            if not reference_blocks or turn_no not in reference_blocks:
                raise PlaceboFidelityError(
                    f"arm D needs A1's block for turn {turn_no} to match per turn"
                )
            reference = reference_blocks[turn_no]
            if reference.record_count == 0:
                # A1 abstained here, so the faithful placebo is also empty: matching
                # decoys to nothing would be a different input, not a matched one.
                return AssembledBlock(text="", citation_ids=[], record_count=0,
                                      shape=(), arm="D"), None
            if self.pool is None:
                raise PlaceboFidelityError("arm D requires a declared decoy pool")
            block, report = build_placebo_measured(reference_blocks[turn_no], self.pool,
                                                   self.wrapper, _counter(),
                                                   self.config.tolerance)
            if probe is not None:
                decoy_texts = [d.text for d in self.pool.decoys()
                               if d.decoy_id in block.citation_ids]
                assert_irrelevant(decoy_texts, probe)
            return block, None

        requested = getattr(probe, "scenario_id", "") if probe is not None else ""
        retr = retrieve(query, self.store, perspective, context=context,
                        max_records=DECLARED_MAX_RECORDS, max_chars=DECLARED_MAX_CHARS,
                        active_state=self.active, requested_scenario=requested)
        if arm == "C":
            retr.owner_id = perspective
        block = assemble(retr.ids(), self.store, self.wrapper, arm=arm) if not retr.abstained \
            else AssembledBlock(text="", citation_ids=[], record_count=0, shape=(), arm=arm)
        return block, retr


_COUNTER = None


def set_counter(counter) -> None:
    """The measured counter for this process: runtime accounting in G4, fixture in G1."""
    global _COUNTER
    _COUNTER = counter


def _counter():
    if _COUNTER is None:
        raise TokenAccountingError("no token counter configured: counts must be measured")
    return _COUNTER


def _context_for(probe) -> str:
    if probe is None:
        return DEFAULT_CONTEXT
    if probe.context:
        return probe.context
    return "world_state" if probe.category == "ambient" else DEFAULT_CONTEXT


def _name_of(spec, owner_id: str) -> str:
    for name, oid in spec.characters.items():
        if oid == owner_id:
            return name
    return "the character"


__all__ = ["RunConfig", "ScriptedResponder", "ArmResult", "build_run_plan", "plan_sha256",
           "ConversationDriver", "set_counter", "MEMORY_BUDGET"]
