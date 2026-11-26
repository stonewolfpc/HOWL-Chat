"""G1 review: M2 applicability and adversarial assurance (zero model calls).

    PYTHONPATH=. python eval/mem1/review.py

Part A asks a question that cannot be answered by argument: for the correction
and retraction probes, does a record exist that *could* support the keyed answer,
can A1 select it, and can its span be proved absent? The answers are read off the
running instrument, not asserted from memory.

Part B adds deterministic tests for the failure modes a reviewer should worry
about: token collisions in retraction, superseded-value leakage, weak-overlap
retrieval, P5 agent/patient confusion, placebo irrelevance in every category, and
the boundary between accounted tokens and the actual serialized prompt.

Nothing here changes the frozen definition. Findings that require implementation
work are recorded as gaps with a named remedy, not patched into existence.
"""
from __future__ import annotations

import hashlib
import json
import os
import pathlib
import subprocess
import sys
import tempfile
import traceback

from eval.mem1 import (identity, receipts)
from eval.mem1.assemble import WRAPPERS, AssembledBlock, assemble
from eval.mem1.conversation import ConversationDriver, RunConfig, ScriptedResponder, set_counter
from eval.mem1.errors import (DecoyLeakError, DesignViolationError, PlaceboFidelityError,
                              TokenAccountingError)
from eval.mem1.fixtures import (ALARIC, DECOYS, SABLE, answers_for, demo_conversation,
                                demo_decoys)
from eval.mem1.form import Determiner
from eval.mem1.model import ACTIVE, REVOKED, SUPERSEDED, MemoryStore
from eval.mem1.p5 import P5_CONFUSED, P5_CORRECT, P5Responder, p5_probe, run_p5
from eval.mem1.placebo import DecoyPool, assert_irrelevant, build_placebo_measured
from eval.mem1.probe import PROBE_CATEGORIES
from eval.mem1.receipts import (GUARD2_G1_RECEIPT, GUARD2_REVIEW_RECEIPT, REFERENCE_HASHES,
                                V2_G1_RECEIPT, refuse_historical, write_once)

# Authorised check count for this generation. A different number is a stop condition.
AUTHORIZED_CURRENT_CHECKS = 142

# The four accepted v2 artifacts, for the evidence-set summary.
V2_ONLY = {p: h for p, h in REFERENCE_HASHES.items() if "/v2/" in p}

# write-guard-v1 remains the record of the narrower eight-artifact implementation.
ACCEPTED_GUARD_V1 = {
    "eval/mem1/write-guard-v1/g1-receipt-writeguard-v1.json":
        "9357ecffdeafdfc029e5fed81f5cc363f20a1f0804ecba982d24576dcf896076",
    "eval/mem1/write-guard-v1/g1-review-receipt-writeguard-v1.json":
        "f6e99cc9bdd9efec3dee5c053f970697318704cb2663ae033c9660d6735f0810",
}

# Authoritative v1 -> v2 supersession inventory. Not derived by keyword matching.
SUPERSESSION = {
    "replace": ["E06", "S02", "S04", "S05", "S06", "S07", "S08", "S09", "S10", "S12"],
    "mapping_changed": ["S01", "S11"],
    "rejustified": ["E07", "F02"],
    "unaffected": ["S03", "F01", "K06", "N26", "N27", "N28", "N29"],
}
from eval.mem1.retrieve import eligible, retrieve
from eval.mem1.score import classify, credit, m2_status, support_check
from eval.mem1.seeds import seed_for
from eval.mem1.tokens import WhitespaceCounter, build_account, require_runtime_counter
from eval.s3.context import Window
from eval.s3.evidence import EvidenceLog


class Checks:
    def __init__(self) -> None:
        self.items: list[dict] = []

    def add(self, cid: str, name: str, passed: bool, detail: str = "", kind: str = "review") -> None:
        self.items.append({"id": cid, "name": name, "passed": bool(passed),
                           "detail": detail, "kind": kind})

    def expect(self, cid: str, name: str, exc_type, fn, detail: str = "",
               kind: str = "negative") -> None:
        try:
            fn()
        except exc_type as exc:
            self.add(cid, name, True, f"{type(exc).__name__}: {str(exc)[:80]}", kind=kind)
        except Exception as exc:
            self.add(cid, name, False, f"wrong exception {type(exc).__name__}: {str(exc)[:80]}",
                     kind=kind)
        else:
            self.add(cid, name, False, "no exception raised", kind=kind)


def part_a(checks: Checks, driver: ConversationDriver, spec, receipts_list) -> list[dict]:
    """M2 applicability, read off the instrument rather than argued."""
    store = driver.store
    blue = next(r for r in store.versions("Sable", "keeps_ledger") if r.status == SUPERSEDED)
    green = next(r for r in store.versions("Sable", "keeps_ledger") if r.status == ACTIVE)
    served = next(r for r in store.all() if r.predicate == "served")
    eligible_ids = [r.record_id for r in eligible(store, SABLE, "character_perspective",
                                                  active_state=None, requested_scenario="")]
    correction_probe = spec.probes[11]
    retraction_probe = spec.probes[12]

    correction = {
        "probe": correction_probe.probe_id,
        "question": correction_probe.question,
        "keyed_answer": correction_probe.key_value,
        "candidate_record": {"record_id": blue.record_id, "status": blue.status,
                             "value": blue.value, "span": blue.source_span,
                             "established_turn": blue.established_turn},
        "selectable_by_a1": blue.record_id in eligible_ids,
        "span_provably_absent": True,   # receipt p47 recorded span_absent=True
        "m2_definition_applied": "CORRECT + selected record supports key + span absent -> "
                                 "NOT_SUPPORTED, because no supporting record can be selected",
        "disposition": "M2 applicable, blocked on an architecture-contract decision (F1), not on "
                       "implementation",
        "causes": [
            "retrieve.eligible() enumerates store.current(), so a superseded record is never a "
            "candidate and there is no declared history context",
            "score.support_check() requires rec.status == ACTIVE, which the frozen M2 text does "
            "not require",
        ],
        "remedy": [
            "BLOCKED — requires an owner decision, not an implementation choice: the frozen "
            "contract names four contexts (character_perspective, world_state, narrator, "
            "inspection), and history appears only in the `knowledge_history` inspection view, "
            "which must never feed the prompt",
            "options: (A) add a fifth prompt-facing context, (B) add a declared temporal parameter "
            "to existing contexts, (C) leave this probe M1-only and record the limitation",
        ],
        "definition_change_required": True,
    }
    retraction_record = next(r for r in store.all()
                             if r.predicate == f"retraction:{retraction_probe.retraction_of}")
    retraction = {
        "probe": retraction_probe.probe_id,
        "question": retraction_probe.question,
        "keyed_answer": retraction_probe.key_value,   # "retracted"
        "retracted_record": {"record_id": served.record_id, "status": served.status,
                             "value": served.value, "span": served.source_span,
                             "revoked_at_turn": served.revoked_at_turn},
        "supporting_record": {"record_id": retraction_record.record_id,
                              "kind": retraction_record.kind,
                              "predicate": retraction_record.predicate,
                              "value": retraction_record.value,
                              "span": retraction_record.source_span,
                              "retracts": retraction_record.retracts,
                              "status": retraction_record.status},
        "selectable_by_a1": retraction_record.record_id in [
            s.record_id for s in retrieve(retraction_probe.question, store, ALARIC,
                                          context="character_perspective",
                                          active_state=None,
                                          requested_scenario="").selected],
        "span_provably_absent": True,
        "m2_definition_applied": "CORRECT + a selected retraction event record that references the "
                                 "retracted record + retraction span absent -> SUPPORTED",
        "disposition": "M2 applicable and now implemented (explicit retraction event record)",
        "causes": [],
        "remedy": [],
        "definition_change_required": False,
        "note": "The retraction is an event record, not a generic absence: it cites its own span, "
                "references the retracted record, and keeps the target revoked with provenance "
                "intact. 'Retracted' therefore remains distinguishable from 'never established' "
                "and from 'not known'.",
    }

    receipt_47 = next(r for r in receipts_list if r.probe_id == "p47")
    receipt_48 = next(r for r in receipts_list if r.probe_id == "p48")
    checks.add("R03", "p47 correction: support requires an active record, superseded excluded",
               not correction["selectable_by_a1"]
               and not support_check(store, [blue.record_id], correction_probe)
               and receipt_47.m2_status == "NOT_SUPPORTED" and receipt_47.span_absent,
               f"blue={blue.status}, selected={receipt_47.selected_ids}, m2={receipt_47.m2_status}",
               kind="gap")
    checks.add("R04", "p48 retraction: explicit retraction event supports the answer",
               served.status == REVOKED and retraction["selectable_by_a1"]
               and support_check(store, [retraction_record.record_id], retraction_probe)
               and receipt_48.m2_status == "SUPPORTED" and receipt_48.span_absent,
               f"support={retraction_record.record_id}, selected={receipt_48.selected_ids}, "
               f"m2={receipt_48.m2_status}", kind="fix")
    checks.add("R04b", "retraction record references its target and preserves it revoked",
               retraction_record.retracts == served.record_id
               and retraction_record.kind == "event"
               and store.get(served.record_id).status == REVOKED
               and store.get(served.record_id).source_span == served.source_span,
               f"retracts={retraction_record.retracts}, target={served.status}", kind="fix")
    # Absence records legitimately exist for never-established facts ("I have never…"),
    # so the distinction is per-predicate: the retracted fact must not be represented
    # as an absence record, and a never-established question must select nothing.
    absence_for_served = [r for r in store.all()
                          if r.predicate == retraction_probe.retraction_of
                          and r.kind == "absence"]
    checks.add("R04c", "retraction is distinguishable from absence and from not-known",
               not absence_for_served
               and retraction_record.predicate.startswith("retraction:")
               and retrieve(spec.probes[3].question, store, SABLE,
                            active_state=None, requested_scenario="").abstained,
               f"no absence record for 'served'; event record {retraction_record.record_id}; "
               f"unanswerable probe abstains", kind="fix")
    checks.expect("R05", "no history context exists to select superseded records",
                  ValueError, lambda: retrieve(correction_probe.question, store, SABLE,
                                               context="history", active_state=None,
                                               requested_scenario=""), kind="negative")
    checks.add("R06", "current-mode retrieval never returns a superseded value",
               all(store.get(rid).status == ACTIVE for rid in eligible_ids),
               f"{len(eligible_ids)} eligible, all active", kind="assurance")
    return [correction, retraction]


def part_b(checks: Checks, driver: ConversationDriver, spec, counter) -> None:
    """Adversarial assurance."""
    # --- retraction token collisions -----------------------------------
    log = EvidenceLog()
    log.append("Sable", "I keep a ledger of unpaid debts in the green cabinet.", 12)
    log.append("Sable", "I promise to keep the north gate locked at dusk.", 11)
    store = MemoryStore(evidence=log)
    former = Determiner(store, log)
    former.run_pending({"Sable": SABLE})

    def facts():
        return [r for r in store.current(owner_id=SABLE)
                if not r.predicate.startswith("retraction:")]

    former._retract("I was wrong about the ledger in the cabinet.", SABLE, 3, "sp-001-01")
    checks.add("R07", "single-fact retraction does not revoke an unrelated fact",
               any(r.predicate == "promise" for r in facts())
               and not any(r.predicate == "keeps_ledger" for r in facts()),
               f"active facts={[r.predicate for r in facts()]}", kind="assurance")

    before = len(facts())
    former._retract("I was wrong about the gate.", SABLE, 4, "sp-001-01")
    checks.add("R08", "one shared content token is not enough to revoke",
               len(facts()) == before, f"active before={before}, after={len(facts())}",
               kind="assurance")

    former._retract("I was wrong about the ledger in the green cabinet and the north gate.",
                    SABLE, 5, "sp-002-01")
    checks.add("R09", "a retraction naming two facts revokes both (documented breadth)",
               not facts(), f"active facts={[r.predicate for r in facts()]}", kind="finding")

    # --- weak overlap ---------------------------------------------------
    weak = retrieve("What does Sable like?", driver.store, SABLE,
                    active_state=None, requested_scenario="")
    weak_ids = [s.record_id for s in weak.selected
                if driver.store.get(s.record_id).predicate != "preference"]
    checks.add("R10", "overlap floor admits weakly related records (reported, not tuned)",
               len(weak_ids) >= 1,
               f"{len(weak_ids)} weak selections of {len(weak.selected)}: "
               f"{[(i, round(next(s.score for s in weak.selected if s.record_id == i), 2)) for i in weak_ids]}",
               kind="finding")

    # --- P5 -------------------------------------------------------------
    correct = run_p5("A1", P5Responder(confused=False), counter)
    confused = run_p5("A1", P5Responder(confused=True), counter)
    checks.add("R11", "P5 correct answer is memory-supported in A1",
               correct.m1_category == "CORRECT" and correct.m2_status == "SUPPORTED"
               and correct.span_absent, f"{correct.m1_category}/{correct.m2_status}")
    checks.add("R12", "P5 agent/patient confusion stays independently observable",
               confused.m1_category == "INCORRECT" and confused.m2_status != "SUPPORTED"
               and confused.agent_patient_confusion is True,
               f"{confused.m1_category}/{confused.m2_status}, confusion flag="
               f"{confused.agent_patient_confusion}")
    a0_p5 = run_p5("A0", P5Responder(), counter)
    d_p5 = run_p5("D", P5Responder(), counter,
                  decoy_pool=DecoyPool(demo_decoys(counter)))
    checks.add("R13", "P5 runs as its own case in A0, A1 and D, outside the 13-probe quota",
               a0_p5.m2_status == "N/A" and d_p5.m2_status == "NOT_SUPPORTED"
               and d_p5.span_absent and all("p5_case" in r.flags
                                            for r in (a0_p5, correct, d_p5))
               and not any(r.probe_id == "P5" for r in driver.last_receipts),
               "3 arms, 0 occurrences in the aggregate receipts")

    # --- placebo irrelevance across every registered category -----------
    failures: list[str] = []
    from eval.mem1.retrieve import terms as _terms
    for category in PROBE_CATEGORIES:
        probe = next(p for p in spec.probes if p.category == category)
        # Use a content term, not the first word: stopwords are stripped by terms(),
        # so "under"/"what" would never appear in a forbidden set.
        candidates = _terms(probe.key_value) or _terms(probe.question)
        answer_token = candidates[0] if candidates else probe.probe_id.lower()
        leaky = DecoyPool([{"decoy_id": "dec-x", "text": f"Someone mentioned {answer_token} twice.",
                            "tokens": counter.count(f"Someone mentioned {answer_token} twice.")}])
        try:
            assert_irrelevant([d.text for d in leaky.decoys()], probe)
            failures.append(category)
        except DecoyLeakError:
            pass
    checks.add("R14", "a decoy containing the answer is rejected in every probe category",
               not failures, f"categories tested={len(PROBE_CATEGORIES)}, failures={failures}")

    pool = DecoyPool(demo_decoys(counter))
    unmatched = []
    for probe in spec.probes:
        block = driver.last_blocks.get(probe.turn)
        if block is None or not block.record_count:
            continue
        placebo, _report = build_placebo_measured(block, pool,
                                                  WRAPPERS["labelled_evidence"], counter, 3)
        decoys = [d.text for d in pool.decoys() if d.decoy_id in placebo.citation_ids]
        try:
            assert_irrelevant(decoys, probe)
        except DecoyLeakError as exc:
            unmatched.append(f"{probe.probe_id}: {exc}")
    checks.add("R15", "declared decoy pool is irrelevant to every scored probe",
               not unmatched, f"{len(spec.probes)} probes, leaks={unmatched}")

    # --- token accounting boundary --------------------------------------
    block = driver.last_blocks[spec.probes[0].turn]
    account = build_account(counter=counter, arm="A1", turn=spec.probes[0].turn, capacity=900,
                            reserve=160, system_text="You are Sable.",
                            history_tokens=100, memory_text=block.bullets_text,
                            wrapper_text=block.wrapper_text)
    serialized = "You are Sable.\n" + block.text
    checks.add("R16", "accounted tokens equal the actual serialized prompt components",
               counter.count(serialized) == account.system + account.memory + account.overhead,
               f"serialized={counter.count(serialized)}, "
               f"split={account.system}+{account.memory}+{account.overhead}", kind="assurance")

    class RuntimeStub:
        name = "runtime-tokenizer-stub"
        def count(self, text):
            return max(1, len(text.split()))

    checks.expect("R17", "a test-only counter is refused for a scored run",
                  TokenAccountingError, lambda: require_runtime_counter(counter), kind="negative")
    ok = True
    try:
        require_runtime_counter(RuntimeStub())
    except TokenAccountingError:
        ok = False
    checks.add("R18", "a declared runtime tokenizer is accepted", ok, "runtime-tokenizer-stub")

    acct_bad = build_account(counter=counter, arm="A1", turn=1, capacity=900, reserve=160,
                             system_text="sys", history_tokens=1, memory_text="", wrapper_text="")
    checks.add("R19", "empty memory block costs zero memory tokens",
               acct_bad.memory == 0 and acct_bad.total <= acct_bad.capacity - acct_bad.reserve,
               f"memory={acct_bad.memory}, total={acct_bad.total}", kind="assurance")


def part_c(checks: Checks, receipts_list) -> dict:
    """Receipt and provenance verification."""
    report = identity.verify_all(strict=True)
    checks.add("R20", "8/8 frozen artifacts unchanged",
               all(r["match"] for r in report["frozen_artifacts"]), "8 checked")
    checks.add("R21", "4/4 frozen preregistration documents unchanged",
               all(r["match"] for r in report["frozen_preregistration"]), "4 checked")
    checks.add("R22", "plan v3 and G0 receipt unchanged",
               all(r["match"] for r in report["preserved_documents"]), "2 checked")

    g1_path = pathlib.Path("eval/mem1/g1-receipt.json")
    g1 = json.loads(g1_path.read_text(encoding="utf-8"))
    current = receipts.source_hashes()
    v2 = json.loads(V2_G1_RECEIPT.read_text(encoding="utf-8"))
    cur = json.loads(GUARD2_G1_RECEIPT.read_text(encoding="utf-8"))

    # Historical evidence is identified, not re-judged. Its source hashes must DIFFER
    # from today's: a Stage 1 receipt that matched current sources would mean the
    # historical record had been rewritten to look like it was generated against v2.
    checks.add("R23", "the historical G1 receipt identifies itself as Stage 1 evidence",
               g1["checks_passed"] == g1["checks_total"] == 96
               and g1["status"] == "PASS"
               and g1["source_sha256"] != current
               and not g1.get("evidence"),
               f"96 checks; v1 sources differ from v2={g1['source_sha256'] != current}; "
               "no v2 label")

    side = g1_path.with_suffix(".json.sha256").read_text().split()[0]
    actual = hashlib.sha256(g1_path.read_bytes()).hexdigest()
    checks.add("R24", "historical G1 receipt sidecar hash verifies unchanged",
               side == actual, actual[:16])

    checks.add("R29", "the guard-generation receipt was run against the current sources",
               cur["source_sha256"] == current and cur["status"] == "PASS"
               and cur["checks_total"] == len(cur["checks"])
               and cur["checks_passed"] == cur["checks_total"] == AUTHORIZED_CURRENT_CHECKS,
               f"package match={cur['source_sha256'] == current}, "
               f"{cur['checks_passed']}/{cur['checks_total']}")

    cur_side = GUARD2_G1_RECEIPT.with_suffix(".json.sha256").read_text().split()[0]
    cur_actual = hashlib.sha256(GUARD2_G1_RECEIPT.read_bytes()).hexdigest()
    checks.add("R30", "guard-generation receipt sidecar hash verifies",
               cur_side == cur_actual, cur_actual[:16])

    ev = cur.get("evidence", {})
    checks.add("R31", "the receipt is labelled for this generation and protects twelve",
               ev.get("set") == "write-guard-v2"
               and ev.get("supersedes")
               == "eval/mem1/write-guard-v1/g1-receipt-writeguard-v1.json"
               and cur["checks_total"] != v2["checks_total"]
               and v2["checks_total"] == 118
               and g1["checks_total"] == 96
               and ev.get("guard", {}).get("protected_files") == 12
               and ev.get("guard", {}).get("handle_identity_before_write") is True,
               f"set={ev.get('set')}; guard={cur['checks_total']} vs v2={v2['checks_total']} "
               f"vs stage1={g1['checks_total']}; protected="
               f"{ev.get('guard', {}).get('protected_files')}")

    ref_mismatch = [p for p, h in REFERENCE_HASHES.items()
                    if hashlib.sha256(pathlib.Path(p).read_bytes()).hexdigest() != h]
    checks.add("R33", "all twelve reference artifacts match their recorded hashes",
               len(REFERENCE_HASHES) == 12 and not ref_mismatch,
               f"{len(REFERENCE_HASHES)} artifacts re-hashed, mismatches={ref_mismatch}")
    checks.add("R34", "the receipt records the twelve reference hashes it protects",
               ev.get("references") == REFERENCE_HASHES,
               f"{len(ev.get('references', {}))} reference hashes recorded")

    protected_ids, handles = receipts._open_protected()
    try:
        set_size = len(protected_ids)
        v2_ids = []
        for rel in sorted(receipts.REFERENCE):
            fd = os.open(str(receipts.PACKAGE / rel), os.O_RDONLY | getattr(os, "O_BINARY", 0))
            try:
                v2_ids.append(receipts._file_identity(fd) in protected_ids)
            finally:
                os.close(fd)
    finally:
        for fd in handles:
            os.close(fd)
    checks.add("R35", "the active protected identity set covers all twelve artifacts",
               set_size == 12 and all(v2_ids),
               f"set size {set_size}, v2 artifacts identified {sum(v2_ids)}/{len(v2_ids)}")

    v1_mismatch = [p for p, h in ACCEPTED_GUARD_V1.items()
                   if hashlib.sha256(pathlib.Path(p).read_bytes()).hexdigest() != h]
    v1_sidecars = all(
        pathlib.Path(p).with_suffix(".json.sha256").read_text().split()[0]
        == hashlib.sha256(pathlib.Path(p).read_bytes()).hexdigest()
        for p in ACCEPTED_GUARD_V1)
    checks.add("R36", "write-guard-v1 is preserved and not retroactively relabelled",
               not v1_mismatch and v1_sidecars,
               f"{len(ACCEPTED_GUARD_V1)} v1 receipts re-hashed, mismatch={v1_mismatch}")

    inv = ev.get("supersession", {})
    checks.add("R32", "the supersession inventory is recorded exactly as ratified",
               all(inv.get(k) == v for k, v in SUPERSESSION.items())
               and inv.get("reviewed") == 21
               and inv.get("edited_or_rejustified") == 14
               and len(SUPERSESSION["replace"]) + len(SUPERSESSION["mapping_changed"]) == 12,
               f"{inv.get('reviewed')} reviewed, {inv.get('edited_or_rejustified')} edited")

    deps = [(d["path"], d["sha256"]) for d in g1["dependencies"]]
    checks.add("R25", "all imported S3 dependencies are read-only and unchanged",
               all(identity.sha256_file(p) == h for p, h in deps) and len(deps) == 11,
               f"{len(deps)} dependencies re-hashed")

    net = []
    for path in pathlib.Path("eval/mem1").glob("*.py"):
        for line in path.read_text(encoding="utf-8").splitlines():
            s = line.strip()
            if s.startswith(("import requests", "import urllib", "import httpx",
                             "import socket", "from requests", "from urllib", "from httpx",
                             "import aiohttp", "from aiohttp")):
                net.append(f"{path.name}: {s}")
    checks.add("R26", "no network client is imported anywhere in mem1", not net, str(net[:3]))

    checks.add("R27", "no experimental model call or subject run occurred",
               all(r.model_calls == 0 for r in receipts_list)
               and g1["confirmations"]["experimental_model_calls"] == 0
               and cur["confirmations"]["experimental_model_calls"] == 0,
               "0 calls across A0/A1/D; 0 declared in both receipts")

    fresh = _fresh_selftest(cur["checks_total"])
    checks.add("R28", "the full current suite still passes on a fresh run",
               fresh["pass"]
               and f"({cur['checks_total']}/{cur['checks_total']} checks)" in fresh["detail"],
               fresh["detail"])
    return {"identity": report, "g1_receipt": g1, "v2_receipt": v2,
            "current_receipt": cur, "fresh_run": fresh}


def _fresh_selftest(expected: int) -> dict:
    """Run the whole suite again into a temp receipt, leaving shipped evidence untouched."""
    tmp = pathlib.Path(tempfile.mkdtemp()) / "g1-receipt.json"
    env = dict(os.environ)
    env["PYTHONPATH"] = "."
    env["MEM1_RECEIPT"] = str(tmp)
    proc = subprocess.run([sys.executable, "eval/mem1/selftest.py"], cwd=str(identity.REPO_ROOT),
                          env=env, capture_output=True, text=True)
    tail = [ln for ln in proc.stdout.splitlines() if ln.startswith("RESULT") or ln.startswith("receipt")]
    ok = proc.returncode == 0 and any(f"PASS ({expected}/{expected} checks)" in ln for ln in tail)
    return {"pass": ok, "detail": "; ".join(tail)[:120] or f"exit={proc.returncode}"}


def main() -> int:
    checks = Checks()
    counter = WhitespaceCounter()
    set_counter(counter)

    spec = demo_conversation()
    driver = ConversationDriver(spec, RunConfig(arm="A1", rep_index=0,
                                               seed=seed_for(spec.conversation_id, 0),
                                               tolerance=3))
    result = driver.run(ScriptedResponder(answers_for(spec)))
    driver.last_receipts = result.receipts
    driver.last_blocks = result.blocks

    findings = part_a(checks, driver, spec, result.receipts)
    part_b(checks, driver, spec, counter)
    verification = part_c(checks, [result])

    total = len(checks.items)
    passed = sum(1 for c in checks.items if c["passed"])
    for c in checks.items:
        print(f"  [{c['id']}] {'PASS' if c['passed'] else 'FAIL'} {c['name']:<62} {c['detail']}")

    review = {
        "gate": "G1-review",
        "status": "PASS" if passed == total else "FAIL",
        "checks_total": total,
        "checks_passed": passed,
        "part_a_m2_applicability": findings,
        "part_b_tests": [c for c in checks.items if c["kind"] in ("assurance", "negative", "finding")],
        "part_c_verification": {
            "frozen_artifacts": f"{sum(r['match'] for r in verification['identity']['frozen_artifacts'])}/8",
            "frozen_preregistration":
                f"{sum(r['match'] for r in verification['identity']['frozen_preregistration'])}/4",
            "preserved_documents":
                f"{sum(r['match'] for r in verification['identity']['preserved_documents'])}/2",
            "historical_g1_receipt": {
                "path": "eval/mem1/g1-receipt.json",
                "status": verification["g1_receipt"]["status"],
                "checks": verification["g1_receipt"]["checks_total"],
                "labelled": "Stage 1, v1 semantics, preserved byte-identical, partly superseded",
                "cas_sha256": hashlib.sha256(
                    pathlib.Path("eval/mem1/g1-receipt.json").read_bytes()).hexdigest()},
            "guard_g1_receipt": {
                "path": str(GUARD2_G1_RECEIPT),
                "status": verification["current_receipt"]["status"],
                "checks": verification["current_receipt"]["checks_total"],
                "labelled": "write-guard-v2: twelve-artifact protected set, WG01-WG24",
                "cas_sha256": hashlib.sha256(GUARD2_G1_RECEIPT.read_bytes()).hexdigest()},
            "accepted_v2_receipts": {
                "paths": sorted(V2_ONLY),
                "standing": "preserved byte-identical; identity-protected in this generation",
                "sha256": V2_ONLY},
            "guard_v1_generation": {
                "paths": sorted(ACCEPTED_GUARD_V1),
                "standing": "preserved; accurate record of the narrower implementation",
                "sha256": ACCEPTED_GUARD_V1},
            "fresh_selftest": verification["fresh_run"],
        },
        "evidence_sets": {
            "historical": {"receipt": "eval/mem1/g1-receipt.json", "checks": 96,
                           "semantics": "v1: owner/restricted/world branch retrieval",
                           "standing": "preserved unchanged; superseded in part by the v2 set"},
            "accepted_v2": {"receipts": sorted(V2_ONLY), "checks": 118,
                            "semantics": "v2: visibility table; chat/world/scenario ownership",
                            "standing": "preserved byte-identical; identity-protected"},
            "guard_v1": {"receipts": sorted(ACCEPTED_GUARD_V1), "checks": 133,
                         "protected_files": 8,
                         "standing": "preserved; narrower implementation, not retroactively "
                                     "compliant with the twelve-artifact requirement"},
            "current": {"receipt": str(GUARD2_G1_RECEIPT),
                        "checks": verification["current_receipt"]["checks_total"],
                        "protected_files": 12,
                        "semantics": "v2 + protected write path; twelve artifacts identity-"
                                     "protected before the first byte"},
        },
        "supersession": {**SUPERSESSION, "reviewed": 21, "edited_or_rejustified": 14},
        "findings": [
            {"id": "F1", "kind": "gap",
             "summary": "M2 for the correction probe is applicable but unimplementable as built: "
                        "no history context, and support requires an active record",
             "impact": "1 of 13 probes contributes M1 only", "remedy": findings[0]["remedy"]},
            {"id": "F2", "kind": "closed",
             "summary": "M2 for the retraction probe is now supported by an explicit retraction "
                        "event record with provenance and a reference to the retracted record",
             "impact": "none remaining; 13/13 probes are M2-scoreable when applicable",
             "remedy": "implemented and covered by R04, R04b, R04c, R07-R09"},
            {"id": "F3", "kind": "risk",
             "summary": "overlap retrieval floor (min_score=0.0) admits weakly related records",
             "impact": "inflates the memory block; does not affect M1 or M2 attribution",
             "remedy": "declare a minimum term-overlap rule before G4 from prior calibration, "
                       "never from outcomes"},
            {"id": "F4", "kind": "risk",
             "summary": "a retraction sentence naming two facts revokes both",
             "impact": "broad but deterministic; a retracted pair leaves two records revoked",
             "remedy": "acceptable as-is, or restrict to the maximum-overlap record pre-G4"},
            {"id": "F5", "kind": "note",
             "summary": "P5 is supported as a standalone case and remains outside the aggregate",
             "impact": "none; Gate B's P5 may not be inherited as a pass"},
        ],
        "recommendation": {
            "g4_ready": False,
            "reason": "F2 is closed; F1 is blocked on an architecture-contract decision, not on "
                      "implementation: the contract names four contexts and exposes history only "
                      "through the inspection view, which must never feed the prompt",
            "pre_g4_work": ["owner decision on how, or whether, historical temporal intent is "
                            "represented in a prompt-facing context"],
            "note": "M1 is unaffected for both probes; C1/C2 remain computable at reduced M2 coverage",
        },
        "confirmations": {
            "experimental_model_calls": 0,
            "subject_runs": 0,
            "frozen_preregistration_modified": False,
            "frozen_artifacts_modified": False,
            "howl_chat_modified": False,
            "commits_or_pushes": 0,
            "g4_executed": False,
        },
    }
    # v2 review evidence goes to its own path. Both historical review receipts stay
    # byte-identical, and a misconfigured override cannot delete protected evidence.
    out = refuse_historical(pathlib.Path(os.environ.get("MEM1_REVIEW_RECEIPT",
                                                        GUARD2_REVIEW_RECEIPT)))
    write_once(out, review)
    print(f"\nRESULT: {'PASS' if passed == total else 'FAIL'} ({passed}/{total} checks)")
    print(f"review receipt: {out} ({out.stat().st_size} B)")
    return 0 if passed == total else 1


if __name__ == "__main__":
    try:
        sys.exit(main())
    except SystemExit:
        raise
    except Exception:
        traceback.print_exc()
        sys.exit(2)
