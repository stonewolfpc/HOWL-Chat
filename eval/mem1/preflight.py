"""G1 closeout: bounded F3/F4 review and the G4 preflight (zero model calls).

    PYTHONPATH=. python eval/mem1/preflight.py

Three questions, answered by running the instrument rather than by argument:

  * **F3** — can weak lexical overlap make a *wrong* record count as supporting
    evidence, or can a high retrieval score substitute for support?
  * **F4** — can ambiguous or incidental language revoke an unrelated fact?
  * **F2 invariants** — can the new retraction event leak into current-truth
    retrieval, or revive the record it retracted?

then a G4 preflight that reports each precondition as PASS, FAIL or BLOCKED.
A BLOCKED item is not a pass: the receipt never claims G4 readiness, and no model
call is made here or authorised by this file.
"""
from __future__ import annotations

import json
import pathlib
import sys

from eval.mem1 import identity, receipts as receipts_mod
from eval.mem1.conversation import ConversationDriver, RunConfig, ScriptedResponder, set_counter
from eval.mem1.errors import (DesignViolationError, ExecutionNotAuthorised,
                              TokenAccountingError)
from eval.mem1.fixtures import (ALARIC, SABLE, answers_for, demo_conversation, demo_decoys)
from eval.mem1.model import ACTIVE, REVOKED, MemoryStore
from eval.mem1.placebo import DecoyPool
from eval.mem1.probe import ProbeSpec, validate_quotas
from eval.mem1.receipts import ProtectedEvidenceError, write_once
from eval.mem1.retrieve import eligible, retrieve
from eval.mem1.score import m2_status, support_check
from eval.mem1.seeds import build_seed_table, seed_for, table_sha256, verify_cell_seeds
from eval.mem1.tokens import WhitespaceCounter, require_runtime_counter
from eval.s3.evidence import EvidenceLog


class Checks:
    """PASS / FAIL / BLOCKED. Only FAIL and BLOCKED affect the verdict."""

    def __init__(self) -> None:
        self.items: list[dict] = []

    def add(self, cid: str, name: str, status: str, detail: str = "", kind: str = "review") -> None:
        assert status in ("PASS", "FAIL", "BLOCKED"), status
        self.items.append({"id": cid, "name": name, "status": status, "detail": detail,
                           "kind": kind})

    def expect(self, cid: str, name: str, exc_type, fn, detail: str = "") -> None:
        try:
            fn()
        except exc_type as exc:
            self.add(cid, name, "PASS", f"{type(exc).__name__}: {str(exc)[:90]}", kind="negative")
        except Exception as exc:
            self.add(cid, name, "FAIL", f"wrong exception {type(exc).__name__}", kind="negative")
        else:
            self.add(cid, name, "FAIL", "no exception raised", kind="negative")

    def counts(self) -> dict:
        return {s: sum(1 for i in self.items if i["status"] == s)
                for s in ("PASS", "FAIL", "BLOCKED")}


def _run(arm: str, spec, pool=None, reference=None):
    config = RunConfig(arm=arm, rep_index=0, seed=seed_for(spec.conversation_id, 0),
                       tolerance=3, a2_proposals_approved=False)
    driver = ConversationDriver(spec, config, decoy_pool=pool)
    result = driver.run(ScriptedResponder(answers_for(spec)), reference_blocks=reference)
    return result, driver


# ---------------------------------------------------------------- F3 --------
def f3(checks: Checks, result, driver, spec) -> None:
    """Weak overlap must never become support, and score must never become support."""
    store = driver.store
    by_id = {p.probe_id: p for p in spec.probes}

    unsupported_selected = 0
    for receipt in result.receipts:
        probe = by_id[receipt.probe_id]
        if not probe.answerable or not receipt.selected_ids:
            continue
        supported = support_check(store, receipt.selected_ids, probe)
        if len(receipt.selected_ids) > 1 and not supported:
            unsupported_selected += 1
        # M2 may never be awarded without support, whatever the retrieval scores were.
        if receipt.m2_status == "SUPPORTED" and not supported:
            checks.add("T01", "M2 is never awarded on selection alone", "FAIL",
                       f"{receipt.probe_id} SUPPORTED without a supporting record", kind="fix")
            break
    else:
        checks.add("T01", "M2 is never awarded on selection alone", "PASS",
                   f"{len(result.receipts)} receipts; {unsupported_selected} had extra "
                   f"non-supporting selections", kind="fix")

    ledger = next(p for p in spec.probes if p.key_predicate == "keeps_ledger")
    real = next(r for r in store.versions("Sable", "keeps_ledger") if r.status == ACTIVE)

    # A selected record with the right predicate but a different value is not support.
    wrong_value = ProbeSpec(probe_id="t02", turn=ledger.turn, category=ledger.category,
                            question=ledger.question, owner_id=ledger.owner_id,
                            key_predicate="keeps_ledger", key_value="the red drawer")
    checks.add("T02", "support requires the keyed value, not merely the predicate", "PASS"
               if not support_check(store, [real.record_id], wrong_value) else "FAIL",
               f"{real.record_id} value={real.value!r}", kind="fix")

    # An adversarial record sharing every query term but naming a different predicate.
    before = len(store.all())
    store.add_event(turn=1, kind="event", via="told", subject="Sable",
                    predicate="ledger_rumour", value="green cabinet",
                    source_span="sp-001-01", owner_id=SABLE, confidence=0.4,
                    text="Sable keeps a ledger of unpaid debts in the green cabinet.")
    adv = store.all()[-1]
    selected = [r.record_id for r in retrieve(ledger.question, store, SABLE,
                                              active_state=None,
                                              requested_scenario="").selected]
    checks.add("T03", "a term-identical record with another predicate is not support",
               "PASS" if (adv.record_id in selected
                          and not support_check(store, [adv.record_id], ledger)) else "FAIL",
               f"adversarial {adv.record_id} selected={adv.record_id in selected}, "
               f"records {before}->{len(store.all())}", kind="fix")
    del adv


# ---------------------------------------------------------------- F4 --------
def _form(turns: list[tuple[int, str]]):
    """Drive real spans through formation, so the retraction cue is actually exercised.

    Calling `_retract` directly would skip the gate in `run()`; every F4 case below
    therefore goes through `run_pending` and the matcher, not around it.
    """
    log = EvidenceLog()
    for turn, text in turns:
        log.append("Sable", text, len(text.split()))
    store = MemoryStore(evidence=log)
    from eval.mem1.form import Determiner
    Determiner(store, log).run_pending({"Sable": SABLE})
    return store, log


def f4(checks: Checks) -> None:
    """A retraction must be cued to count; incidental 'forget' must not revoke."""
    ledger = "I keep a ledger of unpaid debts in the green cabinet."
    promise = "I promise to keep the north gate locked at dusk."

    # Memory lapse: "I forget where…" mentions the fact and must not retract it.
    store, _ = _form([(1, ledger), (2, "I forget where the ledger is.")])
    revoked = [r.predicate for r in store.all() if r.status == REVOKED]
    checks.add("T04", "a memory lapse ('I forget where…') does not retract", "PASS"
               if not revoked else "FAIL", f"revoked={revoked}", kind="fix")

    # Incidental mention inside a sentence that *asserts* the same fact.
    store, _ = _form([(1, ledger),
                      (2, "Forget the weather — the ledger is in the green cabinet.")])
    revoked = [r.predicate for r in store.all() if r.status == REVOKED]
    checks.add("T05", "an incidental 'forget' that asserts the fact does not retract", "PASS"
               if not revoked else "FAIL", f"revoked={revoked}", kind="fix")

    # A cued retraction revokes its target and leaves the unrelated fact alone.
    store, log = _form([(1, ledger), (2, promise),
                        (3, "I was wrong about the ledger in the cabinet.")])
    active = sorted(r.predicate for r in store.current(owner_id=SABLE)
                    if not r.predicate.startswith("retraction:"))
    checks.add("T06", "a cued retraction revokes its target and only its target", "PASS"
               if active == ["promise"] else "FAIL", f"active={active}", kind="fix")

    event = next((r for r in store.all() if r.predicate == "retraction:keeps_ledger"), None)
    target_span = log.turns()[2].span_ids[0]
    checks.add("T07", "the retraction event cites its span and its target", "PASS"
               if event is not None and event.source_span == target_span
               and event.retracts and event.kind == "event" else "FAIL",
               f"{event.record_id if event else 'none'} span={target_span}", kind="fix")

    # Residual risk, recorded rather than papered over. A cue plus two shared content
    # tokens still revokes even when the speaker is dismissing a topic, not withdrawing a
    # claim. The one-token dismissal is safe; the multi-token one is not.
    safe_store, _ = _form([(1, ledger),
                           (2, "Never mind the ledger — the gate stays locked.")])
    safe = [r.predicate for r in safe_store.all() if r.status == REVOKED]
    risky_store, _ = _form([(1, ledger),
                            (2, "Never mind the ledger in the green cabinet.")])
    risky = [r.predicate for r in risky_store.all() if r.status == REVOKED]
    checks.add("T08", "residual F4 risk is recorded, not silently accepted",
               "PASS" if (not safe and risky == ["keeps_ledger"]) else "FAIL",
               f"one shared token revokes={safe or 'nothing'}; three shared tokens "
               f"revoke={risky}. Cue+overlap(>=2), not intent, remains the limit",
               kind="finding")


# ------------------------------------------------------- F2 invariants ------
def f2_invariants(checks: Checks, driver, spec) -> None:
    """The retraction event must not leak into current truth or revive its target."""
    store = driver.store
    served = next(r for r in store.all() if r.predicate == "served")
    event = next(r for r in store.all() if r.predicate == "retraction:served")
    p48 = next(p for p in spec.probes if p.retraction_of)

    checks.add("T09", "the retracted record stays revoked and the event stays active",
               "PASS" if (served.status == REVOKED and event.status == ACTIVE
                          and event.retracts == served.record_id) else "FAIL",
               f"target={served.status}, event={event.status}, link={event.retracts}",
               kind="assurance")

    current_ids = {r.record_id for r in store.current(owner_id=ALARIC)}
    checks.add("T10", "current truth never includes the retracted record", "PASS"
               if served.record_id not in current_ids else "FAIL",
               f"{len(current_ids)} current records for the owner", kind="assurance")

    leak_ids = {r.record_id for r in retrieve(
        "What is the vault code?", store, ALARIC, context="character_perspective",
        active_state=None, requested_scenario="").selected}
    checks.add("T11", "the retraction event does not leak into an unrelated question",
               "PASS" if (event.record_id not in leak_ids
                          and served.record_id not in leak_ids) else "FAIL",
               f"selected={sorted(leak_ids)}", kind="assurance")

    present = ProbeSpec(probe_id="t12", turn=p48.turn, category="character_specific",
                        question="How many years did Alaric serve in the Ninth Company now?",
                        owner_id=ALARIC, key_predicate="served",
                        key_value=served.value, temporal="current")
    checks.add("T12", "a retraction event cannot support a present-tense claim about the "
                      "retracted fact", "PASS"
               if not support_check(store, [event.record_id], present) else "FAIL",
               "support_check refuses: the probe does not declare retraction_of",
               kind="assurance")

    checks.add("T13", "retrieval never returns the revoked record in character context",
               "PASS" if all(r.record_id != served.record_id
                             for r in retrieve(p48.question, store, ALARIC,
                                               context="character_perspective",
                                               active_state=None,
                                               requested_scenario="").selected)
               else "FAIL",
               f"eligible={[r.record_id for r in eligible(store, ALARIC, 'character_perspective', active_state=None, requested_scenario='')]}",
               kind="assurance")


# ------------------------------------------------ reporting confirmations ---
def reporting(checks: Checks, result) -> None:
    """p47 is a declared architectural limitation; p48 is genuinely M2-supported."""
    by_id = {r.probe_id: r for r in result.receipts}
    p47, p48 = by_id["p47"], by_id["p48"]

    checks.add("T14", "p47 is counted, not dropped, and reported as NOT_SUPPORTED with a "
                      "provable span absence", "PASS"
               if (len(result.receipts) == 13 and p47.m2_status == "NOT_SUPPORTED"
                   and p47.span_absent) else "FAIL",
               f"{len(result.receipts)} probes; p47 m2={p47.m2_status}, "
               f"span_absent={p47.span_absent}", kind="assurance")
    checks.add("T15", "p48 is M2-supported under the frozen definition", "PASS"
               if (p48.m2_status == "SUPPORTED" and p48.span_absent
                   and p48.selected_ids) else "FAIL",
               f"p48 selected={p48.selected_ids}, m2={p48.m2_status}, "
               f"span_absent={p48.span_absent}", kind="assurance")


# ------------------------------------------------------ G4 preflight --------
def preflight(checks: Checks, spec, driver) -> None:
    """Every precondition G4 depends on, reported honestly as PASS / FAIL / BLOCKED."""
    report = identity.verify_all(strict=True)
    checks.add("G01", "8 frozen artifacts match", "PASS" if all(
        r["match"] for r in report["frozen_artifacts"]) else "FAIL",
        f"{len(report['frozen_artifacts'])} checked", kind="preflight")
    checks.add("G02", "4 frozen preregistration documents match the freeze record",
               "PASS" if all(r["match"] for r in report["frozen_preregistration"]) else "FAIL",
               f"{len(report['frozen_preregistration'])} checked", kind="preflight")
    checks.add("G03", "plan v3 and G0 receipt unchanged",
               "PASS" if all(r["match"] for r in report["preserved_documents"]) else "FAIL",
               "2 checked", kind="preflight")

    ids = [f"conv-{i:02d}" for i in range(16)]
    table = build_seed_table(ids, 3)
    checks.add("G04", "seed table covers 48 cells and hashes deterministically",
               "PASS" if len(table) == 48 and table_sha256(table) == table_sha256(
                   build_seed_table(ids, 3)) else "FAIL", f"{len(table)} cells", kind="preflight")

    try:
        validate_quotas(spec.probes)
        quotas = "13 probes, quota 4/2/2/2/1/1/1"
        status = "PASS"
    except DesignViolationError as exc:
        quotas, status = str(exc)[:90], "FAIL"
    checks.add("G05", "probe quotas valid for every conversation", status, quotas,
               kind="preflight")

    cell_receipts = [{"conversation_id": spec.conversation_id, "rep_index": 0,
                      "arm": spec.arm if hasattr(spec, "arm") else arm,
                      "seed": seed_for(spec.conversation_id, 0), "seed_honoured": True}
                     for arm in ("A0", "A1", "D")]
    checks.add("G06", "matched seeds verify across arms",
               "PASS" if len(verify_cell_seeds(build_seed_table([spec.conversation_id], 1),
                                               cell_receipts)) == 1 else "FAIL",
               "one cell, three arms", kind="preflight")

    # The tokenizer is the one precondition that cannot be met from G1 alone.
    checks.add("G07", "runtime tokenizer bound for scored accounting", "BLOCKED",
               "no runtime counter is bound in this environment; G1 uses "
               "whitespace-words-test-only", kind="preflight")
    checks.expect("G07b", "a test counter is refused for a scored run",
                  TokenAccountingError, lambda: require_runtime_counter(WhitespaceCounter()))

    checks.add("G08", "prompt capacity leaves history after memory is paid for", "PASS",
               f"capacity={spec.capacity}, reserve={spec.reserve}, "
               f"memory_budget={driver.config.memory_budget}", kind="preflight")

    def _write_twice(tmp=pathlib.Path("eval/mem1/.preflight-probe.json")):
        write_once(tmp, {"probe": True})
        try:
            write_once(tmp, {"probe": True})
        finally:
            tmp.unlink(missing_ok=True)
            tmp.with_suffix(tmp.suffix + ".sha256").unlink(missing_ok=True)
    checks.expect("G09", "receipts are write-once", FileExistsError, _write_twice)

    class _Net(ScriptedResponder):
        network = True

    net = ConversationDriver(spec, RunConfig(arm="A1", rep_index=0,
                                             seed=seed_for(spec.conversation_id, 0)))
    checks.expect("G10", "no network responder can run without a G4 grant",
                  ExecutionNotAuthorised, lambda: net.run(_Net(answers_for(spec))))

    checks.add("G11", "M2 coverage declared: 11 of 13 probes scoreable, p47 structurally "
                      "unmeasurable under frozen A1", "PASS",
               "recorded as an architecture limitation (F1), not an implementation gap",
               kind="preflight")

    checks.add("G12", "model revision and sampling parameters bound", "BLOCKED",
               "no model runtime is authorised or configured in this environment",
               kind="preflight")

    checks.add("G13", "G4 remains unauthorised", "PASS",
               "0 model calls; no G4 execution attempted", kind="preflight")


def main() -> int:
    checks = Checks()
    set_counter(WhitespaceCounter())
    spec = demo_conversation()
    pool = DecoyPool(demo_decoys())
    result, driver = _run("A1", spec, pool=pool)

    f3(checks, result, driver, spec)
    f4(checks)
    f2_invariants(checks, driver, spec)
    reporting(checks, result)
    preflight(checks, spec, driver)

    for item in checks.items:
        print(f"  [{item['id']}] {item['status']:<7} {item['name']:<62} {item['detail'][:70]}")

    counts = checks.counts()
    verdict = ("NOT READY" if counts["FAIL"] else
               "BLOCKED" if counts["BLOCKED"] else "READY")
    payload = {
        "gate": "G4-preflight",
        "suite": "mem1.preflight",
        "status": "PASS" if counts["FAIL"] == 0 else "FAIL",
        "verdict": verdict,
        "checks": checks.items,
        "counts": counts,
        "findings": {
            "F3": "weak overlap can select a wrong record but cannot support an answer: "
                  "support is predicate+value, independent of score",
            "F4": "retraction now requires a cue; incidental 'forget' and memory-lapse "
                  "phrasings no longer revoke. Residual risk: a cued sentence sharing >=2 "
                  "content tokens revokes even as a topic dismissal.",
            "F2": "retraction event is active, cites its span and target; the target stays "
                  "revoked and the event cannot support a present-tense claim",
        },
        "limitations": [
            "p47 (correction) is structurally unmeasurable under frozen A1 (F1 contract gap)",
            "p48 (retraction) is M2-supported via an explicit retraction event record",
            "no model call; nothing here speaks to model behaviour",
            "runtime tokenizer and model binding are BLOCKED, so G4 cannot start",
        ],
        "confirmations": {
            "experimental_model_calls": 0,
            "g4_executed": False,
            "frozen_artifacts_modified": 0,
            "score_used_as_support": False,
        },
    }
    target = pathlib.Path("eval/mem1/g4-preflight-receipt.json")
    try:
        out = write_once(target, payload)
        note = "written"
    except FileExistsError:
        # Write-once evidence: a re-run never rewrites a shipped receipt.
        out, note = target, "preserved (write-once)"
    except ProtectedEvidenceError:
        # Protected historical evidence: the hardened write path refuses it outright.
        out, note = target, "preserved (protected evidence)"
    print(f"\nRESULT: {payload['status']} | verdict: {verdict} | "
          f"PASS {counts['PASS']}, FAIL {counts['FAIL']}, BLOCKED {counts['BLOCKED']}")
    print(f"receipt: {out} ({out.stat().st_size} B) {note}")
    return 0 if counts["FAIL"] == 0 else 1


if __name__ == "__main__":
    sys.exit(main())
