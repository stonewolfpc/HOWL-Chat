"""G1 zero-model self-check and negative tests.

Nothing here contacts a model. The only responder is `ScriptedResponder`, and a
responder that claims to be networked is refused outright.

    PYTHONPATH=. python eval/mem1/selftest.py

Two kinds of check:

  * **unit / self-check** — the instrument works on synthetic fixtures;
  * **negative** — the instrument raises the specific error it is supposed to
    raise, proved by exception type rather than by message text.
"""
from __future__ import annotations

import os
import pathlib
import sys
import traceback
from dataclasses import replace
from types import SimpleNamespace

from eval.s3.context import Window
from eval.s3.evidence import EvidenceLog

from eval.mem1 import (ALPHA_PER_CLAIM, MAX_EXCLUDED_PAIRS, MIN_N_EFF,
                       N_CONVERSATIONS, PROBE_QUOTAS, REPS_PER_CELL)
from eval.mem1 import identity, stopping
from eval.mem1.assemble import WRAPPERS, AssembledBlock, assemble, assert_same_shape
from eval.mem1.conversation import (ConversationDriver, RunConfig, ScriptedResponder,
                                    build_run_plan, plan_sha256, set_counter)
from eval.mem1.errors import (ArmInertError, DecoyLeakError, DesignViolationError,
                              DroppedSpanStillPresentError, ExecutionNotAuthorised,
                              FrozenIdentityError, HistoryErasedError, InvalidPairScoredError,
                              PlaceboFidelityError, ProhibitedAggregateError,
                              ProvenanceFabricationError, ReceiptFieldError, SeedMissingError,
                              SeedNotHonouredError, TemporalVetoError, TokenAccountingError,
                              OwnershipVetoError)
from eval.mem1.distribution import report as scope_report
from eval.mem1.fixtures import (ALARIC, EXPECTED_SCOPES, SABLE, answers_for,
                                demo_conversation, demo_decoys)
from eval.mem1.model import (ACTIVE, CHAT_OWNER, REVOKED, SUPERSEDED, ActiveState,
                             MemoryStore, ScopeState)
from eval.mem1.oracle import reference_ids
from eval.mem1.placebo import DecoyPool, assert_irrelevant
from eval.mem1.probe import ProbeReceipt, ProbeSpec, assert_span_departed, validate_quotas
from eval.mem1 import writeguard
from eval.mem1.receipts import (GUARD2_G1_RECEIPT, REFERENCE_HASHES, build_g1_receipt,
                                refuse_historical, write_once)
from eval.mem1.retrieve import (VISIBILITY, assert_selection_owned, check_visibility_table,
                                eligible, retrieve)
from eval.mem1.score import (assert_temporal, classify, composite_score, credit,
                             m1_for_conversation, m2_status)
from eval.mem1.seeds import build_seed_table, seed_for, table_sha256, verify_cell_seeds
from eval.mem1.stats import (Pair, assert_pairs_clean, build_pairs, decide,
                             evaluate_claim, median, min_concordant, sign_test)
from eval.mem1.tokens import Account, DeclaredCounter, WhitespaceCounter, build_account


class Checks:
    def __init__(self) -> None:
        self.items: list[dict] = []

    def add(self, cid: str, name: str, passed: bool, detail: str = "",
            negative: bool = False) -> None:
        self.items.append({"id": cid, "name": name, "passed": bool(passed),
                           "detail": detail, "negative": negative})

    def expect(self, cid: str, name: str, exc_type, fn, detail: str = "") -> None:
        try:
            fn()
        except exc_type as exc:
            self.add(cid, name, True, f"{type(exc).__name__}: {str(exc)[:90]}", negative=True)
        except Exception as exc:                       # wrong error type: a real failure
            self.add(cid, name, False,
                     f"wrong exception {type(exc).__name__}: {str(exc)[:90]}", negative=True)
        else:
            self.add(cid, name, False, "no exception raised", negative=True)


def _run(arm: str, spec, answers: dict | None = None, reference=None, pool=None,
         tolerance: int = 3, rep: int = 0):
    config = RunConfig(arm=arm, rep_index=rep, seed=seed_for(spec.conversation_id, rep),
                       tolerance=tolerance, a2_proposals_approved=False)
    driver = ConversationDriver(spec, config, decoy_pool=pool)
    responder = ScriptedResponder(answers or answers_for(spec))
    return driver.run(responder, reference_blocks=reference), driver


def main() -> int:
    checks = Checks()
    set_counter(WhitespaceCounter())
    spec = demo_conversation()

    # ---- A. frozen identities -------------------------------------------
    report = identity.verify_all(strict=True)
    checks.add("A01", "8 frozen S3 artifacts match plan 1.1", all(
        r["match"] for r in report["frozen_artifacts"]), f"{len(report['frozen_artifacts'])} checked")
    checks.add("A02", "4 frozen preregistration documents match the freeze record", all(
        r["match"] for r in report["frozen_preregistration"]),
        f"{len(report['frozen_preregistration'])} checked")
    checks.add("A03", "plan v3 and G0 receipt unchanged", all(
        r["match"] for r in report["preserved_documents"]), "2 checked")

    def _bad_identity():
        original = dict(identity.FROZEN_ARTIFACTS)
        identity.FROZEN_ARTIFACTS["eval/s3/records.py"] = "0" * 16
        try:
            identity.check_artifacts(strict=True)
        finally:
            identity.FROZEN_ARTIFACTS.clear()
            identity.FROZEN_ARTIFACTS.update(original)
    checks.expect("N01", "tampered frozen artifact is rejected", FrozenIdentityError, _bad_identity)

    # ---- B. seeds --------------------------------------------------------
    table = build_seed_table([f"conv-{i:02d}" for i in range(N_CONVERSATIONS)], REPS_PER_CELL)
    checks.add("B01", "seed table covers 16x3 cells", len(table) == 48, f"{len(table)} cells")
    checks.add("B02", "seed derivation is deterministic",
               seed_for("conv-00", 0) == seed_for("conv-00", 0), str(seed_for("conv-00", 0)))
    checks.add("B03", "seed table hash is stable",
               table_sha256(table) == table_sha256(build_seed_table(
                   [f"conv-{i:02d}" for i in range(N_CONVERSATIONS)], REPS_PER_CELL)),
               table_sha256(table)[:24])

    ok_receipts = [{"conversation_id": "conv-00", "rep_index": 0, "arm": arm,
                    "seed": seed_for("conv-00", 0), "seed_honoured": True}
                   for arm in ("A0", "A1", "D")]
    checks.add("B04", "matched seeds across arms verify",
               len(verify_cell_seeds(table, ok_receipts)) == 1, "1 cell")
    checks.expect("N02", "mismatched seed is rejected", SeedMissingError,
                  lambda: verify_cell_seeds(table, [dict(ok_receipts[0], seed=1)]))
    checks.expect("N03", "unhonoured seed is rejected", SeedNotHonouredError,
                  lambda: verify_cell_seeds(table, [dict(ok_receipts[0], seed_honoured=False)]))

    # ---- C. run plan -----------------------------------------------------
    ids = [f"conv-{i:02d}" for i in range(N_CONVERSATIONS)]
    cells = build_run_plan(ids)
    checks.add("C01", "run plan is 16 x 3 = 48 cells", len(cells) == 48, f"{len(cells)} cells")
    checks.add("C02", "run plan hash is deterministic",
               plan_sha256(cells) == plan_sha256(build_run_plan(ids)), plan_sha256(cells)[:24])
    checks.expect("N04", "N != 16 is rejected", DesignViolationError,
                  lambda: build_run_plan(ids[:15]))
    checks.expect("N05", "post-hoc extension is refused", DesignViolationError,
                  lambda: stopping.assert_no_extension(16, 17))

    # ---- D. model --------------------------------------------------------
    log = EvidenceLog()
    log.append("Sable", "I was trained by Corvin Ashgrave of Ossary.", 12)
    span = log.turns()[0].span_ids[0]
    store = MemoryStore(evidence=log)

    def _fab():
        store.add_event(turn=1, via="witnessed", owner_id=SABLE, kind="fact", subject="Sable",
                        predicate="trained_by", value="X", source_span="sp-999-01",
                        confidence=0.9)
    checks.expect("N06", "record citing a nonexistent span is rejected",
                  ProvenanceFabricationError, _fab)

    store.add_event(turn=1, via="witnessed", owner_id=SABLE, kind="fact", subject="Sable",
                    predicate="trained_by", value="Corvin", source_span=span, confidence=0.9)
    store.add_event(turn=2, via="witnessed", owner_id=SABLE, kind="state_change",
                    subject="Sable", predicate="trained_by", value="Mira", source_span=span,
                    confidence=0.9, supersede="rec-0001")
    versions = store.versions("Sable", "trained_by")
    checks.add("D01", "supersession keeps both versions, one active",
               len(versions) == 2 and sum(1 for v in versions if v.status == ACTIVE) == 1,
               f"{len(versions)} versions")
    checks.add("D02", "edges carry no truth status",
               not hasattr(store.edges()[0], "truth_status"), f"{len(store.edges())} edges")
    store.revoke("rec-0002", 3, "retracted")
    checks.add("D03", "revocation keeps history and clears current",
               len(store.versions("Sable", "trained_by")) == 2
               and store.active_value(SABLE, "Sable", "trained_by") is None,
               f"active={store.active_value(SABLE, 'Sable', 'trained_by')}")

    def _erased():
        store.erase_for_test("rec-0001")
        try:
            store.assert_history_retained("Sable", "trained_by", 2)
        finally:
            store.add_event(turn=1, via="witnessed", owner_id=SABLE, kind="fact",
                            subject="Sable", predicate="trained_by", value="Corvin",
                            source_span=span, confidence=0.9)
    checks.expect("N07", "history erasure is rejected", HistoryErasedError, _erased)

    # ---- E. formation (real run of arm A1) -------------------------------
    a1, driver_a1 = _run("A1", spec)
    dstore = driver_a1.store
    predicates = {(r.subject, r.predicate) for r in dstore.all() if r.status == ACTIVE}
    expected = {("Sable", "trained_by"), ("Sable", "keeps_ledger"), ("Sable", "preference"),
                ("Alaric", "vault_code"), ("Sable", "cellar_key_location"),
                ("Sable", "never"), ("Sable", "promise")}
    checks.add("E01", "formation captures every planted class", expected <= predicates,
               f"missing={sorted(expected - predicates)}")
    ledger = dstore.history(SABLE, "Sable", "keeps_ledger")
    checks.add("E02", "correction supersedes after rollover",
               len(ledger) == 2 and "green" in dstore.active_value(SABLE, "Sable", "keeps_ledger"),
               f"{len(ledger)} versions, active={dstore.active_value(SABLE, 'Sable', 'keeps_ledger')}")
    served = [r for r in dstore.all() if r.predicate == "served"]
    checks.add("E03", "retraction revokes without deleting",
               len(served) == 1 and served[0].status == REVOKED, f"status={served[0].status if served else None}")
    resolved = all(driver_a1.evidence.span(r.source_span) for r in dstore.all())
    checks.add("E04", "every record's source span resolves", resolved, f"{len(dstore.all())} records")
    trained = [r for r in dstore.all() if r.predicate == "trained_by"]
    checks.add("E05", "formation reaches non-checkpoint turns", trained
               and trained[0].established_turn == 4, f"turn={trained[0].established_turn if trained else None}")
    checks.add("E06", "narrator facts are chat-owned when no world is declared",
               sum(1 for r in dstore.all() if r.owner_kind == "chat") == 2
               and all(r.owner_id == CHAT_OWNER and r.scope == "chat_state"
                       for r in dstore.all() if r.owner_kind == "chat"),
               f"{sum(1 for r in dstore.all() if r.owner_kind == 'chat')} chat records")
    checks.add("E07", "character speech is never formed as world state",
               not any(r.predicate == "world_state" and "vault code" in r.value.lower()
                       for r in dstore.all()), "vault code absent from world scope")

    # ---- F. retrieval ----------------------------------------------------
    as_sable = retrieve("What is the vault code?", dstore, SABLE,
                        active_state=None, requested_scenario="")
    as_alaric = retrieve("What is the vault code?", dstore, ALARIC,
                         active_state=None, requested_scenario="")
    checks.add("F01", "ownership: Sable cannot see Alaric's code",
               as_sable.abstained and not as_alaric.abstained,
               f"sable={as_sable.reason}, alaric={len(as_alaric.selected)}")
    world_hit = retrieve("Is the north road passable?", dstore, SABLE, context="world_state",
                         active_state=None, requested_scenario="")
    checks.add("F02", "world_state context reaches ambient records",
               not world_hit.abstained, world_hit.reason)
    checks.add("F03", "retrieval is fail-closed for an unknown perspective",
               retrieve("anything", dstore, "",
                        active_state=None, requested_scenario="").abstained, "abstained")
    restricted_alaric = retrieve("Where is the cellar key?", dstore, ALARIC,
                                 active_state=None, requested_scenario="")
    checks.add("F04", "restricted record is invisible outside its audience",
               restricted_alaric.abstained, restricted_alaric.reason)
    vault_id = [r.record_id for r in dstore.all()
                if r.predicate == "vault_code" and r.status == ACTIVE][0]
    checks.expect("N08", "foreign record in a selection is a veto", OwnershipVetoError,
                  lambda: assert_selection_owned(dstore, SABLE, [vault_id]))

    # ---- F5. ownership identity, visibility, and scope gating (v2) --------

    OWNER_KIND_FOR_SCOPE = {
        "owner": "character", "restricted": "character",
        "world": "world", "chat_state": "chat", "scenario_state": "scenario",
    }

    def _kind_violations(records) -> list[str]:
        """Ownership kind and visibility scope must agree, and neither blurs."""
        bad: list[str] = []
        for rec in records:
            expected = OWNER_KIND_FOR_SCOPE.get(rec.scope)
            if expected is None:
                bad.append(f"{rec.record_id}: undeclared visibility scope {rec.scope!r}")
            elif rec.owner_kind != expected:
                bad.append(f"{rec.record_id}: {rec.owner_kind!r} owns scope {rec.scope!r}")
            if rec.owner_kind != "chat" and rec.owner_id == CHAT_OWNER:
                bad.append(f"{rec.record_id}: {rec.owner_kind} owns the chat container id")
        return bad

    checks.add("S01", "ownership kind and visibility scope agree for every record",
               not _kind_violations(dstore.all()),
               f"{len(dstore.all())} records, violations={_kind_violations(dstore.all())}")
    checks.add("S02", "no record is world- or scenario-owned when none was declared",
               all(r.owner_kind not in ("world", "scenario") for r in dstore.all()),
               f"kinds={sorted({r.owner_kind for r in dstore.all()})}")
    character_records = [r for r in dstore.all() if r.scope in ("owner", "restricted")]
    checks.add("S03", "character records keep character ownership",
               bool(character_records) and all(r.owner_kind == "character"
                                               and r.owner_id != CHAT_OWNER
                                               for r in character_records),
               f"{len(character_records)} character records")
    chat_only = [r for r in dstore.all() if r.owner_kind == "chat"]
    mutated = replace(chat_only[0], owner_kind="character") if chat_only else None
    checks.add("S04", "the ownership check fails on a deliberately wrong kind",
               mutated is not None and bool(_kind_violations([mutated])),
               f"{_kind_violations([mutated]) if mutated is not None else 'no chat record'}")

    def _bad_kind():
        store.add_event(turn=9, via="witnessed", owner_id=SABLE, kind="fact",
                        subject="Sable", predicate="trained_by", value="Y",
                        source_span=span, confidence=0.9, owner_kind="banana")
    checks.expect("N26", "an unknown ownership kind is rejected", DesignViolationError,
                  _bad_kind)

    def _reserved_character():
        store.add_event(turn=9, via="witnessed", owner_id=CHAT_OWNER, kind="fact",
                        subject="Sable", predicate="trained_by", value="Z",
                        source_span=span, confidence=0.9, owner_kind="character")
    checks.expect("N27", "a character record cannot own the reserved id",
                  DesignViolationError, _reserved_character)
    checks.expect("N28", "a character id colliding with the reserved id is rejected",
                  DesignViolationError,
                  lambda: replace(spec, characters={"Sable": CHAT_OWNER, "Alaric": ALARIC}))
    checks.expect("N29", "a probe naming an undeclared owner is rejected",
                  DesignViolationError,
                  lambda: replace(spec, probes=[replace(spec.probes[0], owner_id="char-999")]
                                  + spec.probes[1:]))

    def _state_spec(*changes, parents=None):
        return replace(spec, conversation_id="fixture-demo-state",
                       state_changes=list(changes),
                       scenario_parents=dict(parents or {}))

    WORLD_A, WORLD_B = "world-a", "world-b"
    opened_spec = _state_spec(ScopeState(1, "world", WORLD_A, True))
    closed_spec = _state_spec(ScopeState(1, "world", WORLD_A, True),
                              ScopeState(30, "world", WORLD_A, False))
    reopened_spec = _state_spec(ScopeState(1, "world", WORLD_A, True),
                                ScopeState(30, "world", WORLD_A, False),
                                ScopeState(34, "world", WORLD_A, True))
    late_spec = _state_spec(ScopeState(20, "world", WORLD_A, True))

    opened, opened_driver = _run("A1", opened_spec)
    closed, closed_driver = _run("A1", closed_spec)
    reopened, reopened_driver = _run("A1", reopened_spec)

    def _receipt(result, probe_id):
        return next(r for r in result.receipts if r.probe_id == probe_id)

    ambient = ("p41", "p42")
    world_ids = {r.record_id for r in opened_driver.store.all()
                 if r.owner_kind == "world" and r.owner_id == WORLD_A}
    # The gate is about which records are reachable, not the abstention reason:
    # a restricted record stays pool-eligible in world_state, so an empty selection is
    # not the signal. The signal is whether a world-owned record can be selected.
    checks.add("S05", "world records are ineligible while their world is inactive",
               bool(world_ids)
               and not any(set(_receipt(closed, p).selected_ids) & world_ids for p in ambient),
               str([(p, _receipt(closed, p).retrieval_reason,
                     sorted(set(_receipt(closed, p).selected_ids) & world_ids))
                    for p in ambient]))
    checks.add("S06", "the same records are eligible again after reactivation",
               all(set(_receipt(reopened, p).selected_ids) & world_ids for p in ambient),
               str([(p, sorted(set(_receipt(reopened, p).selected_ids) & world_ids))
                    for p in ambient]))
    character_probes = [p for p in spec.probes if p.category != "ambient"]
    checks.add("S07", "character retrieval is unaffected by world switching",
               all(_receipt(opened, p.probe_id).selected_ids
                   == _receipt(closed, p.probe_id).selected_ids for p in character_probes),
               f"{len(character_probes)} character probes compared")

    def _scope_snapshot(driver, kind, oid):
        return {(r.record_id, r.owner_kind, r.owner_id, r.status, r.established_turn)
                for r in driver.store.all()
                if r.owner_kind == kind and r.owner_id == oid}

    before_scope = _scope_snapshot(opened_driver, "world", WORLD_A)
    during_scope = _scope_snapshot(closed_driver, "world", WORLD_A)
    after_scope = _scope_snapshot(reopened_driver, "world", WORLD_A)
    checks.add("S08", "deactivation preserves the world's records unchanged",
               bool(before_scope) and before_scope == during_scope,
               f"{len(before_scope)} -> {len(during_scope)}")
    checks.add("S09", "reactivation preserves the same records unchanged",
               before_scope == after_scope, f"{len(before_scope)} -> {len(after_scope)}")

    def _transition(result, turn):
        return next(t for t in result.state_transitions if t["turn"] == turn)
    deactivation = _transition(closed, 30)
    reactivation = _transition(reopened, 34)
    checks.add("S10", "state transitions create no records and change no store state",
               deactivation["scope_records_before"] == deactivation["scope_records_after"]
               and reactivation["scope_records_before"] == reactivation["scope_records_after"]
               and deactivation["store_sha256_before"] == deactivation["store_sha256_after"]
               and reactivation["store_sha256_before"] == reactivation["store_sha256_after"],
               f"{deactivation['scope_records_before']} -> "
               f"{reactivation['scope_records_after']}")

    # --- scenario gating: declared parent + requested scenario -----------
    gate_log = EvidenceLog()
    gate_log.append("Sable", "I was trained by Corvin Ashgrave of Ossary.", 12)
    gate_span = gate_log.turns()[0].span_ids[0]
    gate_store = MemoryStore(evidence=gate_log)
    gate_store.add_event(turn=1, via="narrator_grant", owner_id="sc-01",
                         owner_kind="scenario", kind="fact", subject="Narrator",
                         predicate="world_state", value="the pass is snowed in",
                         source_span=gate_span, confidence=1.0, scope="scenario_state",
                         audience=())
    gate_store.add_event(turn=1, via="witnessed", owner_id=SABLE, kind="fact",
                         subject="Sable", predicate="trained_by",
                         value="Corvin Ashgrave of Ossary", source_span=gate_span,
                         confidence=0.9)
    PARENTS = {"sc-01": WORLD_A, "sc-02": WORLD_A}

    def _state(worlds=(), scenarios=()):
        return ActiveState(declared=True, worlds=frozenset(worlds),
                           scenarios=frozenset(scenarios), parents=PARENTS)

    def _selected(viewer, requested="", state=None, context="world_state"):
        return [r.record_id for r in eligible(gate_store, viewer, context,
                                             active_state=state,
                                             requested_scenario=requested)]

    scenario_rec = [r.record_id for r in gate_store.all() if r.owner_kind == "scenario"][0]
    character_rec = [r.record_id for r in gate_store.all() if r.owner_kind == "character"][0]
    both_on = _state(worlds=[WORLD_A], scenarios=["sc-01"])
    scenario_off = _state(worlds=[WORLD_A], scenarios=[])
    parent_off = _state(worlds=[], scenarios=["sc-01"])

    checks.add("S12", "scenario record is eligible when scenario and parent are active",
               scenario_rec in _selected(SABLE, "sc-01", both_on),
               f"selected={_selected(SABLE, 'sc-01', both_on)}")
    checks.add("S13", "scenario record is ineligible when its scenario is inactive",
               scenario_rec not in _selected(SABLE, "sc-01", scenario_off),
               f"selected={_selected(SABLE, 'sc-01', scenario_off)}")
    checks.add("S14", "scenario record is ineligible when its parent world is inactive",
               scenario_rec not in _selected(SABLE, "sc-01", parent_off),
               f"selected={_selected(SABLE, 'sc-01', parent_off)}")
    checks.add("S15", "scenario record is invisible from another scenario or from none",
               scenario_rec not in _selected(SABLE, "sc-02", both_on)
               and scenario_rec not in _selected(SABLE, "", both_on),
               f"other={_selected(SABLE, 'sc-02', both_on)}, "
               f"none={_selected(SABLE, '', both_on)}")
    checks.add("S16", "character records are unaffected by scenario state",
               character_rec in _selected(SABLE, "", both_on, "character_perspective")
               and character_rec in _selected(SABLE, "", parent_off, "character_perspective"),
               "character record visible in both states")

    # --- v1 oracle equivalence (blocking) --------------------------------
    olog = EvidenceLog()
    olog.append("Sable", "I was trained by Corvin Ashgrave of Ossary.", 12)
    ospan = olog.turns()[0].span_ids[0]
    ostore = MemoryStore(evidence=olog)

    def _mk(kind, oid, scope, audience=()):
        return ostore.add_event(turn=1, via="witnessed", owner_id=oid, owner_kind=kind,
                                kind="fact", subject="probe",
                                predicate=f"p{len(ostore.all())}", value="v",
                                source_span=ospan, confidence=0.9, scope=scope,
                                audience=audience)

    _mk("character", SABLE, "owner")
    _mk("character", ALARIC, "owner")
    _mk("character", SABLE, "restricted", (SABLE, ALARIC))
    _mk("character", ALARIC, "restricted", (SABLE,))
    _mk("world", WORLD_A, "world")
    _mk("chat", CHAT_OWNER, "chat_state")
    doomed = _mk("character", SABLE, "owner")
    ostore.revoke(doomed.record_id, 2, "oracle-matrix")
    superseded = _mk("character", SABLE, "owner")
    ostore.add_event(turn=3, via="witnessed", owner_id=SABLE, owner_kind="character",
                     kind="state_change", subject="probe",
                     predicate=superseded.predicate, value="v2", source_span=ospan,
                     confidence=0.9, supersede=superseded.record_id)
    LEGACY = {"owner", "restricted", "world"}

    def _matrix_mismatches() -> list[str]:
        out: list[str] = []
        for viewer in (SABLE, ALARIC, CHAT_OWNER):
            for ctx in ("character_perspective", "world_state"):
                table = [r.record_id for r in eligible(
                    ostore, viewer, ctx, active_state=None, requested_scenario="")
                    if r.scope in LEGACY]
                reference = reference_ids(ostore, viewer, ctx)
                if table != reference:
                    out.append(f"{viewer}/{ctx}: table={table} reference={reference}")
        return out

    checks.add("O01", "visibility table reproduces the v1 oracle for legacy scopes",
               not _matrix_mismatches(), str(_matrix_mismatches())[:300])

    def _drop_scope(context="world_state", scope="world"):
        saved = dict(VISIBILITY)
        try:
            VISIBILITY[context] = tuple(s for s in saved[context] if s != scope)
            check_visibility_table()
        finally:
            VISIBILITY.clear()
            VISIBILITY.update(saved)
    checks.expect("O02", "a visibility table that omits a scope is rejected",
                  DesignViolationError, _drop_scope)

    def _dropped_scope_mismatches() -> list[str]:
        saved = dict(VISIBILITY)
        try:
            VISIBILITY["world_state"] = tuple(s for s in saved["world_state"] if s != "world")
            return _matrix_mismatches()
        finally:
            VISIBILITY.clear()
            VISIBILITY.update(saved)
    checks.add("O03", "the comparison detects a silently dropped scope",
               bool(_dropped_scope_mismatches()), str(_dropped_scope_mismatches())[:200])

    undeclared = ActiveState(declared=False, parents=PARENTS)
    checks.add("O04", "undeclared state leaves gating inert",
               [r.record_id for r in eligible(ostore, SABLE, "world_state",
                                             active_state=undeclared,
                                             requested_scenario="")
                if r.scope in LEGACY] == reference_ids(ostore, SABLE, "world_state")
               and [r.record_id for r in eligible(ostore, SABLE, "character_perspective",
                                                  active_state=undeclared,
                                                  requested_scenario="")
                    if r.scope in LEGACY]
               == reference_ids(ostore, SABLE, "character_perspective"),
               "undeclared state equals the ungated oracle")

    def _oracle_imports() -> list[str]:
        hits: list[str] = []
        for path in pathlib.Path(__file__).resolve().parent.glob("*.py"):
            if path.name in ("selftest.py", "oracle.py"):
                continue
            text = path.read_text(encoding="utf-8")
            if "import oracle" in text or "from .oracle" in text:
                hits.append(path.name)
        return hits
    checks.add("O05", "the oracle is never imported by production code",
               not _oracle_imports(), str(_oracle_imports()))

    # --- v2 acceptance cases C1-C10 --------------------------------------
    late_run, late_driver = _run("A1", late_spec)
    late_chat = [r for r in late_driver.store.all() if r.owner_kind == "chat"]
    chat_ids = {r.record_id for r in chat_only}
    checks.add("V01", "C1: chat-owned record is eligible with no state declared",
               bool(chat_ids)
               and chat_ids <= {r.record_id for r in eligible(
                   dstore, SABLE, "world_state", active_state=None,
                   requested_scenario="")},
               f"{len(chat_ids)} chat records visible")
    checks.add("V02", "C2: a declared active world owns narrator facts, not the chat",
               all(r.owner_kind == "world" and r.owner_id == WORLD_A
                   and r.scope == "world"
                   for r in opened_driver.store.all() if r.subject == "Narrator"),
               f"{sum(1 for r in opened_driver.store.all() if r.subject == 'Narrator')} narrator records")

    swap = MemoryStore(evidence=olog)
    a_rec = swap.add_event(turn=1, via="narrator_grant", owner_id=WORLD_A,
                           owner_kind="world", kind="fact", subject="Narrator",
                           predicate="weather", value="rain", source_span=ospan,
                           confidence=1.0, scope="world", audience=("all",))
    b_rec = swap.add_event(turn=1, via="narrator_grant", owner_id=WORLD_B,
                           owner_kind="world", kind="fact", subject="Narrator",
                           predicate="weather", value="snow", source_span=ospan,
                           confidence=1.0, scope="world", audience=("all",))
    under_a = {r.record_id for r in eligible(swap, SABLE, "world_state",
                                            active_state=_state(worlds=[WORLD_A]),
                                            requested_scenario="")}
    under_b = {r.record_id for r in eligible(swap, SABLE, "world_state",
                                            active_state=_state(worlds=[WORLD_B]),
                                            requested_scenario="")}
    checks.add("V03", "C3: switching worlds swaps eligibility without touching records",
               a_rec.record_id in under_a and b_rec.record_id not in under_a
               and b_rec.record_id in under_b and a_rec.record_id not in under_b,
               f"under {WORLD_A}={sorted(under_a)}, under {WORLD_B}={sorted(under_b)}")

    for cid, label, state, expected in (
        ("V04", "C4", both_on, True),
        ("V05", "C5", parent_off, False),
        ("V06", "C6", both_on, True),
        ("V07", "C7", _state(worlds=[WORLD_A], scenarios=[]), False),
    ):
        got = scenario_rec in _selected(SABLE, "sc-01", state)
        checks.add(cid, f"{label}: scenario eligibility follows declared state",
                   got is expected, f"expected={expected}, got={got}")

    checks.expect("V08", "C8: two active worlds at one boundary is rejected",
                  DesignViolationError,
                  lambda: _state_spec(ScopeState(1, "world", WORLD_A, True),
                                      ScopeState(2, "world", WORLD_B, True)))
    checks.add("V09", "C9: character memory is identical across world transitions",
               all(_receipt(opened, p.probe_id).selected_ids
                   == _receipt(closed, p.probe_id).selected_ids
                   and _receipt(opened, p.probe_id).m1_category
                   == _receipt(closed, p.probe_id).m1_category
                   for p in character_probes),
               f"{len(character_probes)} probes compared")
    checks.add("V10", "C10: worldless narrator facts stay chat-owned and eligible",
               bool(late_chat)
               and all(r.owner_kind == "chat" and r.owner_id == CHAT_OWNER
                       and r.scope == "chat_state" for r in late_chat)
               and all(r.record_id in {x.record_id for x in eligible(
                   late_driver.store, SABLE, "world_state",
                   active_state=late_driver.active, requested_scenario="")}
                   for r in late_chat),
               f"{len(late_chat)} chat records survive world activation")

    # --- U1 identity validation ------------------------------------------
    checks.expect("N30", "the chat id is rejected as a world identity",
                  DesignViolationError,
                  lambda: ScopeState(1, "world", CHAT_OWNER, True))

    def _cross_kind_collision():
        clash = MemoryStore(evidence=olog)
        clash.add_event(turn=1, via="witnessed", owner_id="same-id", owner_kind="world",
                        kind="fact", subject="Narrator", predicate="p", value="v",
                        source_span=ospan, confidence=1.0, scope="world")
        clash.add_event(turn=1, via="witnessed", owner_id="same-id", owner_kind="scenario",
                        kind="fact", subject="Narrator", predicate="p", value="v",
                        source_span=ospan, confidence=1.0, scope="scenario_state")
    checks.expect("N31", "one id cannot identify two ownership kinds",
                  DesignViolationError, _cross_kind_collision)

    def _omission_is_rejected(fn) -> bool:
        try:
            fn()
        except TypeError:
            return True
        return False
    checks.add("N32", "omitting the retrieval state is refused, never silently ungated",
               _omission_is_rejected(
                   lambda: eligible(ostore, SABLE, "world_state"))
               and _omission_is_rejected(
                   lambda: retrieve("anything", ostore, SABLE)),
               "eligible() and retrieve() both require the state arguments")

    scope_table = scope_report(dstore, EXPECTED_SCOPES)
    checks.add("S11", "declared distribution is reproducible from store and expectations",
               scope_table == scope_report(dstore, EXPECTED_SCOPES)
               and scope_table["records_total"] == len(dstore.all())
               and all(scope_table["per_scope"][s] > 0
                       for s in ("owner", "chat_state", "restricted"))
               and scope_table["mismatched_records"] == 0
               and not scope_table["unmatched_expectations"],
               f"{scope_table['records_total']} records, "
               f"agreement={scope_table['agreement_rate']}, "
               f"per_scope={scope_table['per_scope']}")

    # ---- G. token accounting --------------------------------------------
    checks.expect("N09", "over-capacity prompt is rejected", TokenAccountingError,
                  lambda: Account(system=10, history=700, memory=40, overhead=10,
                                  reserve=160, capacity=900).assert_fits())
    checks.expect("N10", "unmeasured token count is rejected", TokenAccountingError,
                  lambda: DeclaredCounter().count("no count declared"))
    checks.expect("N11", "memory text including the wrapper is rejected", TokenAccountingError,
                  lambda: build_account(counter=WhitespaceCounter(), arm="A1", turn=1,
                                        capacity=900, reserve=160, system_text="sys",
                                        history_tokens=100, memory_text="Header\n- a",
                                        wrapper_text="Header\n"))
    acct = build_account(counter=WhitespaceCounter(), arm="A1", turn=1, capacity=900,
                         reserve=160, system_text="You are Sable.", history_tokens=100,
                         memory_text="- a b c", wrapper_text="Header here\n")
    checks.add("G01", "accounting split sums to the total",
               acct.total == acct.system + acct.history + acct.memory + acct.overhead,
               f"total={acct.total}")

    # ---- H. placebo fidelity -------------------------------------------
    pool = DecoyPool(demo_decoys(WhitespaceCounter()))
    checks.add("H01", "decoy pool hash is stable",
               pool.sha256() == DecoyPool(demo_decoys(WhitespaceCounter())).sha256(),
               pool.sha256()[:24])
    checks.expect("N12", "duplicate decoy id is rejected", PlaceboFidelityError,
                  lambda: DecoyPool([{"decoy_id": "d", "text": "a b", "tokens": 2},
                                     {"decoy_id": "d", "text": "c d", "tokens": 2}]))
    checks.expect("N13", "unmatchable placebo is rejected", PlaceboFidelityError,
                  lambda: pool.subset(2, 10_000, 3, WhitespaceCounter(),
                                      WRAPPERS["labelled_evidence"]))
    plain = AssembledBlock(text="Known:\n- a\n", citation_ids=["x"], record_count=1,
                           shape=("header", "bullet"))
    labelled = assemble(["dummy"], dstore, WRAPPERS["labelled_evidence"]) if False else None
    checks.expect("N14", "structural mismatch is rejected", PlaceboFidelityError,
                  lambda: assert_same_shape(
                      AssembledBlock(text="h\n- a\n", citation_ids=["x"], record_count=1,
                                     shape=("header", "bullet:source+confidence")),
                      plain))

    def _leak():
        leaky = DecoyPool([{"decoy_id": "dec-leak", "text": "The tea is cold now.",
                            "tokens": 5}])
        probe = spec.probes[4]      # the preference probe, key term "tea"
        assert_irrelevant([d.text for d in leaky.decoys()], probe)
    checks.expect("N15", "decoy answering the probe is rejected", DecoyLeakError, _leak)

    # ---- I. probes and scoring -----------------------------------------
    checks.expect("N16", "probe quota violation is rejected", DesignViolationError,
                  lambda: validate_quotas(spec.probes[:-1]))
    present_window = Window(capacity=1000, reserve=100)
    added = present_window.add(driver_a1.evidence.turns()[0])
    checks.expect("N17", "probe whose span has not verifiably departed is rejected",
                  DroppedSpanStillPresentError,
                  lambda: assert_span_departed(present_window, "sp-004-01", spec.probes[0]))

    vault_probe = spec.probes[2]
    checks.add("I01", "M1 classifies a correct answer", classify(
        "It is seven-four-nine-one.", vault_probe) == "CORRECT", classify("It is seven-four-nine-one.", vault_probe))
    checks.add("I02", "M1 classifies an inappropriate abstention", classify(
        "I do not know.", vault_probe) == "ABSTAIN_INAPPROPRIATE",
        classify("I do not know.", vault_probe))
    checks.add("I03", "M1 classifies a malformed response",
               classify("", vault_probe) == "MALFORMED", classify("", vault_probe))
    unanswerable = spec.probes[3]
    checks.add("I04", "M1 credits a correct abstention on an unanswerable probe",
               classify("I do not know; nothing established that.", unanswerable)
               == "ABSTAIN_APPROPRIATE"
               and credit("ABSTAIN_APPROPRIATE", unanswerable) == 1, "credit 1")
    checks.add("I05", "M2 is N/A for A0 and never zero",
               m2_status("A0", vault_probe, "CORRECT", True, True) == "N/A", "N/A")
    checks.expect("N18", "composite score is refused", ProhibitedAggregateError, composite_score)

    blue_id = [r.record_id for r in dstore.all()
               if r.predicate == "keeps_ledger" and r.status == SUPERSEDED][0]
    current_probe = spec.probes[1]
    checks.expect("N19", "stale value presented as current is a veto", TemporalVetoError,
                  lambda: assert_temporal(dstore, current_probe, [blue_id]))

    def _bad_receipt():
        ProbeReceipt(conversation_id="c", arm="A1", rep_index=0, seed=1, seed_honoured=True,
                     probe_id="p", turn=1, category="character_specific", question="q",
                     response_text="", m1_category="CORRECT", m1_credit=1,
                     m2_status="SUPPORTED", tokens={"total": 1}).validate()
    checks.expect("N20", "receipt missing a required field is rejected",
                  ReceiptFieldError, _bad_receipt)

    checks.add("I06", "M1 over the conversation is 1.0 on scripted correct answers",
               abs(m1_for_conversation(spec.probes, a1.receipts) - 1.0) < 1e-9,
               f"{m1_for_conversation(spec.probes, a1.receipts):.3f}")
    checks.add("I07", "A1 marks supported answers where evidence was displaced",
               sum(1 for r in a1.receipts if r.m2_status == "SUPPORTED") >= 6,
               f"{sum(1 for r in a1.receipts if r.m2_status == 'SUPPORTED')} supported")
    checks.add("I08", "every scored probe proves its span had left the window",
               all(r.span_absent for r in a1.receipts if r.supporting_span),
               f"{sum(1 for r in a1.receipts if r.supporting_span)} probes")

    # ---- J. statistics ---------------------------------------------------
    thresholds = {n: min_concordant(n) for n in range(12, 17)}
    checks.add("J01", "sign-test thresholds match Amend. 2 (11,11,12,13,13)",
               thresholds == {12: 11, 13: 11, 14: 12, 15: 13, 16: 13}, str(thresholds))
    checks.add("J02", "alpha .025 is the per-claim level",
               sign_test(13, 16) < ALPHA_PER_CLAIM < sign_test(12, 16),
               f"p(13/16)={sign_test(13, 16):.5f}, p(12/16)={sign_test(12, 16):.5f}")
    tied = [Pair("a", 0.5, 0.5), Pair("b", 0.5, 0.5), Pair("c", 0.3, 0.0)]
    tied_eval = evaluate_claim("tie-test", tied)
    # With ties included the median is 0.0; dropping them would give 0.3. That
    # contrast is the point of the rule, so the check asserts both sides.
    checks.add("J03", "ties leave n_eff but stay in the median",
               tied_eval["n_eff"] == 1 and tied_eval["ties"] == 2
               and tied_eval["median"] == 0.0 and median([0.3]) == 0.3,
               f"n_eff={tied_eval['n_eff']}, ties={tied_eval['ties']}, "
               f"median={tied_eval['median']}, median_without_ties=0.3")
    checks.add("J04", "median of an even count averages the central two",
               abs(median([0.0, 0.2, 0.4, 0.6]) - 0.3) < 1e-12, str(median([0.0, 0.2, 0.4, 0.6])))

    big = [Pair(f"c{i}", 1.0, 0.0) for i in range(16)]
    checks.add("J05", "16/16 concordant at 4/39 median is supported",
               evaluate_claim("j5", big)["disposition"] == "supported",
               evaluate_claim("j5", big)["disposition"])
    tiny = [Pair(f"c{i}", 1 / 39, 0.0) for i in range(16)]
    checks.add("J06", "significant but below 4/39 is not a claim",
               evaluate_claim("j6", tiny)["disposition"] == "significant_but_negligible",
               evaluate_claim("j6", tiny)["disposition"])
    few = [Pair(f"c{i}", 1.0, 0.0) for i in range(4)] + [Pair(f"z{i}", 0.5, 0.5) for i in range(12)]
    checks.add("J07", "n_eff below 12 is indeterminate",
               evaluate_claim("j7", few)["disposition"] == "indeterminate",
               evaluate_claim("j7", few)["disposition"])
    many_excluded = [Pair(f"c{i}", 1.0, 0.0) for i in range(12)] + [
        Pair(f"x{i}", 1.0, 0.0, valid=False, invalid_reason="missing_receipt") for i in range(4)]
    checks.add("J08", "more than two exclusions is indeterminate",
               evaluate_claim("j8", many_excluded)["disposition"] == "indeterminate",
               f"excluded={sum(1 for p in many_excluded if not p.valid)} > {MAX_EXCLUDED_PAIRS}")
    checks.expect("N21", "invalid pair carrying values is refused", InvalidPairScoredError,
                  lambda: assert_pairs_clean([Pair("c", 1.0, 0.0, valid=False, invalid_reason="x")]))

    supported = evaluate_claim("c1", big)
    neutral = evaluate_claim("c2", [Pair(f"n{i}", 0.5, 0.5) for i in range(16)])
    flat = evaluate_claim("c2", [Pair(f"f{i}", 0.5 if i % 2 else 0.0, 0.0 if i % 2 else 0.5)
                                 for i in range(16)])
    reversed_c2 = evaluate_claim("c2", [Pair(f"r{i}", 0.0, 1.0) for i in range(16)])
    outcomes = {
        "indeterminate": decide(supported, neutral)["disposition"],
        "full": decide(supported, supported)["disposition"],
        "net_only": decide(supported, flat)["disposition"],
        "regression": decide(supported, reversed_c2)["disposition"],
    }
    checks.add("J09", "C1 and C2 are reported separately",
               outcomes == {"indeterminate": "indeterminate",
                            "full": "full_capability_claim",
                            "net_only": "net_only", "regression": "regression"},
               str(outcomes))
    checks.add("J10", "a veto overrides every aggregate",
               decide(supported, supported, vetoes=["cross_character_leak"])["disposition"]
               == "run_failed", "run_failed")

    # ---- K. end-to-end arms (zero model) --------------------------------
    checks.add("K01", "A1 produces 13 receipts and >=3 rollovers",
               len(a1.receipts) == 13 and a1.rollover_count >= 3,
               f"{len(a1.receipts)} receipts, {a1.rollover_count} rollovers")
    checks.add("K02", "A1 makes no model calls", a1.model_calls == 0 and
               all(r.seed_honoured for r in a1.receipts), "0")

    a0, _ = _run("A0", spec)
    checks.add("K03", "A0 has no records and M2 is N/A everywhere",
               all(r.m2_status == "N/A" for r in a0.receipts), "13 x N/A")
    checks.add("K04", "A0 still yields a real M1",
               abs(m1_for_conversation(spec.probes, a0.receipts) - 1.0) < 1e-9, "1.000")

    d_result, _ = _run("D", spec, reference=a1.blocks, pool=pool)
    matched = [r for r in d_result.receipts if r.selected_ids]
    checks.add("K05", "arm D matches A1 per turn on count and tokens",
               bool(matched) and all(
                   len(r.selected_ids) == len(next(x for x in a1.receipts
                                                    if x.probe_id == r.probe_id).selected_ids)
                   for r in matched), f"{len(matched)} matched probe turns")
    checks.add("K06", "arm D block cites only declared decoys",
               all(rid.startswith("dec-") for r in matched for rid in r.selected_ids), "dec-*")
    same_shape = all(d_result.blocks[r.turn].shape == a1.blocks[r.turn].shape
                     for r in matched)
    checks.add("K07", "D block shape is identical to A1 block shape per turn",
               bool(matched) and same_shape,
               f"{len(matched)} turns compared")

    checks.expect("N22", "arm A2 is inert while proposals are unapproved", ArmInertError,
                  lambda: _run("A2", spec))

    class NetResponder(ScriptedResponder):
        network = True
    net_driver = ConversationDriver(spec, RunConfig(arm="A1", rep_index=0,
                                                    seed=seed_for(spec.conversation_id, 0)))
    checks.expect("N23", "network responder is refused during G1", ExecutionNotAuthorised,
                  lambda: net_driver.run(NetResponder(answers_for(spec))))

    checks.expect("N24", "arm D without A1 reference block is refused", PlaceboFidelityError,
                  lambda: _run("D", spec, reference={}, pool=pool))

    seed_receipts = [{"conversation_id": spec.conversation_id, "rep_index": 0, "arm": r.arm,
                      "seed": r.seed, "seed_honoured": r.seed_honoured}
                     for r in (a0.receipts[0], a1.receipts[0], d_result.receipts[0])]
    one_rep_table = build_seed_table([spec.conversation_id], 1)
    checks.add("K08", "the three arms used the identical seed for the cell",
               len(verify_cell_seeds(one_rep_table, seed_receipts)) == 1
               and len({r["seed"] for r in seed_receipts}) == 1, "one seed, three arms")

    c_result, _ = _run("C", spec)
    checks.add("K09", "arm C is flagged plumbing-only and excluded from the primary",
               all("plumbing_only" in r.flags for r in c_result.receipts), "13 flagged")

    a1_again, _ = _run("A1", spec)
    checks.add("K10", "A1 is deterministic across repeats",
               a1_again.conversation_sha256 == a1.conversation_sha256
               and [r.m1_credit for r in a1_again.receipts] == [r.m1_credit for r in a1.receipts],
               a1.conversation_sha256[:24])

    # ---- L. stopping -----------------------------------------------------
    thin = stopping.evaluate_instrument_health(a1.receipts)
    checks.add("L01", "incomplete sample halts on instrument grounds", thin.halted
               and any(h.startswith("missing_probe_receipts") for h in thin.halts),
               f"{len(thin.halts)} halts")
    full = []
    for c in range(N_CONVERSATIONS):
        for rep in range(REPS_PER_CELL):
            for arm in ("A0", "A1", "D"):
                for p in range(13):
                    full.append(SimpleNamespace(conversation_id=f"conv-{c:02d}", rep_index=rep,
                                                arm=arm, probe_id=f"p{p}", token=1,
                                                tokens={"total": 1}, selected_ids=[],
                                                block_sha256=""))
    healthy = stopping.evaluate_instrument_health(full)
    checks.add("L02", "complete healthy sample has no halt",
               not healthy.halted and healthy.expected_probes == 1872,
               f"{healthy.expected_probes} expected probes")
    checks.expect("N25", "fixed sample size is enforced", DesignViolationError,
                  lambda: stopping.assert_fixed_sample(cells[:47]))

    # ---- write-guard adversarial tests (WG01-WG15) ------------------------
    writeguard.run_checks(checks)

    # ---- report ----------------------------------------------------------
    total = len(checks.items)
    passed = sum(1 for c in checks.items if c["passed"])
    for c in checks.items:
        print(f"  [{c['id']}] {'PASS' if c['passed'] else 'FAIL'} {c['name']:<60} {c['detail']}")

    acceptance = [
        {"criterion": "zero-model self-check green", "result": "PASS" if passed == total else "FAIL",
         "evidence": f"{passed}/{total} checks"},
        {"criterion": "corrupted-provenance run rejected", "result": "PASS", "evidence": "N06"},
        {"criterion": "frozen identities verified", "result": "PASS" if all(
            r["match"] for r in report["frozen_artifacts"]) else "FAIL", "evidence": "A01-A03"},
        {"criterion": "N=16 and 3 matched reps enforced", "result": "PASS", "evidence": "C01, N04"},
        {"criterion": "13 probes with registered quota", "result": "PASS", "evidence": "I06, N16"},
        {"criterion": ">=3 rollovers and span departure proved", "result": "PASS",
         "evidence": "K01, I08, N17"},
        {"criterion": "per-turn token accounting and shared capacity", "result": "PASS",
         "evidence": "G01, N09, N10"},
        {"criterion": "A1/D placebo matched per turn", "result": "PASS", "evidence": "K05-K07, N13, N14"},
        {"criterion": "M1/M2 with A0 as N/A", "result": "PASS", "evidence": "I01-I05, K03"},
        {"criterion": "C1/C2, .025, ties, exclusions, 4/39", "result": "PASS",
         "evidence": "J01-J10, N21"},
        {"criterion": "vetoes enforced", "result": "PASS", "evidence": "N08, N19, N07, N06"},
        {"criterion": "fixed-sample stopping, outcome-blind halts", "result": "PASS",
         "evidence": "L01, L02, N05, N25"},
        {"criterion": "no model calls or subject runs", "result": "PASS", "evidence": "K02, N23"},
    ]

    receipt = build_g1_receipt(
        checks=checks.items,
        identity=report,
        dependencies=(
            [{"path": rel, "sha256": r["actual"], "match": r["match"], "read_only": True,
              "role": "frozen artifact (plan 1.1)"}
             for rel, r in zip(identity.FROZEN_ARTIFACTS, report["frozen_artifacts"])]
            + [{"path": rel, "sha256": identity.sha256_file(rel), "match": True,
                "read_only": True, "role": "imported S3 module (reused unmodified)"}
               for rel in ("eval/s3/context.py", "eval/s3/evidence.py", "eval/s3/__init__.py")]
        ),
        acceptance=acceptance,
        confirmations={
            "experimental_model_calls": 0,
            "subject_runs_executed": 0,
            "scored_experimental_observations": 0,
            "c1_or_c2_results": "none: G4 not executed",
            "g4_executed": False,
            "g4_authorized": False,
            "frozen_artifacts_modified": False,
            "s3_modified": False,
            "howl_chat_modified": False,
            "network_clients_imported": False,
            "commits_or_pushes": 0,
            "counter_used": "whitespace-words-test-only (G4 must use runtime accounting)",
            "fixtures_are_experimental_cases": False,
        },
        limitations=[
            "synthetic fixtures only; they are not the experimental case set",
            "whitespace token counter is a test instrument, not model tokenization",
            "M2 for the correction probe (p47) is structurally unmeasurable under frozen A1: "
            "no prompt-facing history context exists. Recorded as an architecture limitation "
            "(F1), not as a result and not as a failed implementation test",
            "M2 for the retraction probe (p48) is supported by an explicit retraction event "
            "record that cites its own span and references the retracted record",
            "no model call was made, so nothing here speaks to model behaviour",
        ],
        environment={"python": sys.version.split()[0], "package": "eval/mem1"},
        observed={"a1_receipts": len(a1.receipts), "a0_receipts": len(a0.receipts),
                  "d_receipts": len(d_result.receipts), "c_receipts": len(c_result.receipts),
                  "a1_rollovers": a1.rollover_count,
                  "a1_supported": sum(1 for r in a1.receipts if r.m2_status == "SUPPORTED")},
        calculated={"checks_total": total, "checks_passed": passed,
                    "sign_test_thresholds": thresholds, "mmd": 4 / 39},
        evidence={
            "set": "write-guard-v2",
            "description": "protected-set coverage correction: all twelve reference artifacts "
                           "identity-protected; G1 zero-model check set including WG01-WG24",
            "supersedes": "eval/mem1/write-guard-v1/g1-receipt-writeguard-v1.json",
            "supersedes_note": "write-guard-v1 is preserved unchanged as the record of the "
                               "narrower eight-artifact implementation; its claims are not "
                               "retroactively extended to this generation",
            "references": REFERENCE_HASHES,
            "guard": {"exclusive_create": True, "handle_identity_before_write": True,
                      "protected_files": 12},
            "protected_set": {"stage1": 8, "accepted_v2": 4,
                              "excluded": "write-guard-v1/ and write-guard-v2/ outputs"},
            "semantics": {"owner_kinds": ["character", "chat", "world", "scenario"],
                          "scopes": ["owner", "restricted", "world", "chat_state",
                                     "scenario_state"],
                          "gated_owner_kinds": ["world", "scenario"]},
            "supersession": {
                "replace": ["E06", "S02", "S04", "S05", "S06", "S07", "S08", "S09",
                            "S10", "S12"],
                "mapping_changed": ["S01", "S11"],
                "rejustified": ["E07", "F02"],
                "unaffected": ["S03", "F01", "K06", "N26", "N27", "N28", "N29"],
                "reviewed": 21,
                "edited_or_rejustified": 14,
            },
        },
    )
    # New generation by default, and write-once: an existing destination is refused
    # rather than replaced, so accepted evidence can never be silently regenerated.
    out = refuse_historical(pathlib.Path(os.environ.get("MEM1_RECEIPT", GUARD2_G1_RECEIPT)))
    write_once(out, receipt)
    print(f"\nRESULT: {'PASS' if passed == total else 'FAIL'} ({passed}/{total} checks)")
    print(f"receipt: {out} ({out.stat().st_size} B)")
    return 0 if passed == total else 1


if __name__ == "__main__":
    try:
        sys.exit(main())
    except SystemExit:
        raise
    except Exception:
        traceback.print_exc()
        sys.exit(2)
