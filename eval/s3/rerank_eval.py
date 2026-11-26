"""R1–R5 harness for cross-encoder candidates, against a FROZEN case set.

Frozen-before-candidate rules:
  * cases, expected outcomes, selection rule, threshold and veto logic live in
    eval/s3/rerank-spec.json, with a canonical SHA-256 over the case payload;
  * this harness refuses to run if the rebuilt cases do not hash to the frozen
    value, or if the harness's own scoring code changed since the freeze;
  * the threshold is fixed in the spec. It is never tuned on the evaluation set.

A reranker is any object with `score(query: str, candidate: dict) -> float` and a
`name`. It never receives the expected label, and that is checked.

Run:  PYTHONPATH=. python eval/s3/rerank_eval.py --freeze v1
      PYTHONPATH=. python eval/s3/rerank_eval.py
"""
from __future__ import annotations

import hashlib
import inspect
import json
import pathlib
import re
import sys

from eval.s3.evaluate import build_state
from eval.s3.rerank_set import CRITERIA, VETOES, build_cases

SPEC_PATH = pathlib.Path("eval/s3/rerank-spec.json")
RECEIPT_PATH = pathlib.Path("eval/s3/rerank-harness-receipt.json")
STOP = {"who", "was", "the", "a", "an", "is", "are", "did", "do", "you", "your",
        "again", "what", "where", "why", "how", "of", "in", "at", "it", "that",
        "kept", "think", "before", "now", "one", "there", "were", "be"}

# Existing system outputs, used as the comparison baseline. These are the frozen
# Gate B probes that correspond to each adversarial case — not recomputed here.
BASELINE = {
    "R1-horse": {"probe": "P5", "outcome": "FAIL", "detail": "returns the trainer record"},
    "R1-patient": {"probe": None, "outcome": "UNTESTED", "detail": "no equivalent probe"},
    "R2-vault-foreign": {"probe": "P3b", "outcome": "PASS", "detail": "abstains"},
    "R2-vault-own": {"probe": "P3a", "outcome": "PASS", "detail": "returns vault_code"},
    "R3-now": {"probe": "P4a", "outcome": "PASS", "detail": "current record"},
    "R3-before": {"probe": "P4b", "outcome": "PASS", "detail": "history reachable"},
    "R4-paraphrase": {"probe": "P1", "outcome": "PASS", "detail": "subword + semantic"},
    "R5-weather": {"probe": "P2", "outcome": "PASS", "detail": "abstains"},
    "R5-fair": {"probe": None, "outcome": "UNTESTED", "detail": "no equivalent probe"},
}


def stems(text: str) -> set[str]:
    words = re.findall(r"[a-z0-9]+", text.lower())
    return {w[:5] for w in words if w not in STOP and len(w) > 2}


class OverlapReranker:
    """Baseline analogue: topical overlap only — what a bi-encoder effectively does."""
    name = "overlap-baseline"

    def score(self, query: str, cand: dict) -> float:
        return len(stems(query) & stems(cand["text"])) / max(len(stems(query)), 1)


class TrapReranker:
    """Deliberately wrong: rewards any shared relation word, ignores the arguments."""
    name = "trap-mock"

    def score(self, query: str, cand: dict) -> float:
        return 1.0 if (stems(query) & stems(cand["text"])) else 0.0


class NoiseReranker:
    """Deliberately breaks abstention: scores everything equally."""
    name = "noise-mock"

    def score(self, query: str, cand: dict) -> float:
        return 0.5


class JointReranker:
    """Harness control only. Hand-fitted to the cases; NOT a performance baseline."""
    name = "joint-reference-mock"
    ENTITIES = {"horse", "weather", "fair", "prize", "vault", "code", "ledger",
                "corvin", "trainer", "tea"}

    def score(self, query: str, cand: dict) -> float:
        q = stems(query)
        text = cand["text"].lower()
        c = stems(text)
        if not (q & c):
            return 0.0
        for ent in re.findall(r"[a-z0-9]+", query.lower()):
            if ent in self.ENTITIES and ent not in text and ent[:5] not in c:
                return 0.0
        score = len(q & c) / max(len(q), 1)
        if "now" in query.lower():
            score += 0.10 if cand["status"] == "active" else -0.10
        elif any(w in query.lower() for w in ("before", "earlier")):
            score += 0.10 if cand["status"] != "active" else 0.0
        return score


def canonical(payload) -> str:
    return json.dumps(payload, sort_keys=True, separators=(",", ":"))


def cases_hash(cases: list[dict]) -> str:
    return hashlib.sha256(canonical(cases).encode()).hexdigest()


def harness_hash() -> str:
    return hashlib.sha256(pathlib.Path(__file__).read_bytes()).hexdigest()


def case_payload(cases: list[dict]) -> list[dict]:
    return [{k: c[k] for k in ("id", "criterion", "owner", "query", "expect",
                               "expect_abstain", "why", "candidates")} for c in cases]


def freeze(version: str) -> int:
    _, store, _, _ = build_state()
    payload = case_payload(build_cases(store))
    spec = {
        "version": version,
        "frozen_on": "2026-11-15",
        "threshold": 0.30,
        "threshold_note": "fixed before any candidate existed; never tuned on this set",
        "selection_rule": "keep candidates scoring >= threshold; if none remain, abstain",
        "criteria": CRITERIA,
        "vetoes": list(VETOES),
        "veto_rule": "a broken veto rejects the candidate outright; not averaged",
        "cases": payload,
        "cases_sha256": cases_hash(payload),
        "harness_sha256": harness_hash(),
        "baseline_system_as_is": BASELINE,
        "isolation": "scorer receives query and candidate fields only; never the label",
        "generalization_note": "9 cases test mechanism, not generalization; "
                               "the reference mock is a control, not a baseline",
    }
    SPEC_PATH.write_text(json.dumps(spec, indent=2) + "\n", encoding="utf-8")
    print(f"froze {version}: {len(payload)} cases | cases_sha256={spec['cases_sha256'][:16]} "
          f"| harness={spec['harness_sha256'][:16]}")
    print(f"written: {SPEC_PATH} ({SPEC_PATH.stat().st_size} B)")
    return 0


def load_spec(cases: list[dict]) -> dict | None:
    if not SPEC_PATH.exists():
        print(f"REFUSING: no frozen spec at {SPEC_PATH}. Run with --freeze <version> first.")
        return None
    spec = json.loads(SPEC_PATH.read_text(encoding="utf-8"))
    payload = case_payload(cases)
    live = cases_hash(payload)
    if live != spec["cases_sha256"]:
        print("REFUSING: the case set no longer matches the frozen spec.")
        print(f"  frozen {spec['cases_sha256'][:16]} | now {live[:16]}")
        print("  A changed case set is a new experiment: re-freeze under a new version.")
        return None
    live_h = harness_hash()
    if live_h != spec["harness_sha256"]:
        print("REFUSING: harness scoring code changed since the freeze.")
        print(f"  frozen {spec['harness_sha256'][:16]} | now {live_h[:16]}")
        print("  Re-freeze under a new version, or revert the harness.")
        return None
    return spec


def run_case(reranker, case: dict, threshold: float) -> dict:
    # Isolation: the scorer sees only query + candidate fields, never the label.
    assert "expect" not in canonical(case["candidates"]), "label leaked into candidates"
    # Ties break on record_id, never on input order: candidate ordering must not
    # be able to change a decision (checked by the order-invariance self-check).
    scored = sorted(((reranker.score(case["query"], dict(c)), dict(c))
                     for c in case["candidates"]),
                    key=lambda t: (-t[0], t[1]["record_id"]))
    kept = [c for s, c in scored if s >= threshold]
    abstained = not kept
    top = kept[0]["record_id"] if kept else None
    if case["expect_abstain"]:
        passed = abstained
        reason = None if abstained else f"selected {top} but should have abstained"
    elif abstained:
        passed = False
        reason = f"abstained; expected {case['expect']}"
    elif top != case["expect"]:
        passed = False
        reason = f"selected {top}; expected {case['expect']}"
    else:
        passed = True
        reason = None
    return {
        "id": case["id"], "criterion": case["criterion"], "query": case["query"],
        "expect": "abstain" if case["expect_abstain"] else case["expect"],
        "top": top, "abstained": abstained,
        "top_score": round(scored[0][0], 3) if scored else None,
        "scores": {c["record_id"]: round(s, 3) for s, c in scored},
        "pass": passed, "failure_reason": reason,
    }


def evaluate(reranker, cases, threshold: float) -> dict:
    rows = [run_case(reranker, c, threshold) for c in cases]
    # Order invariance: reversing candidate order must not change any decision.
    reversed_rows = [run_case(reranker, {**c, "candidates": list(reversed(c["candidates"]))},
                              threshold) for c in cases]
    order_stable = all(a["top"] == b["top"] and a["abstained"] == b["abstained"]
                       for a, b in zip(rows, reversed_rows))
    per_criterion = {}
    for crit in CRITERIA:
        sel = [r for r in rows if r["criterion"] == crit]
        per_criterion[crit] = {
            "label": CRITERIA[crit], "passed": sum(1 for r in sel if r["pass"]),
            "total": len(sel), "veto": crit in VETOES,
            "ok": all(r["pass"] for r in sel),
        }
    veto_broken = [c for c in VETOES if not per_criterion[c]["ok"]]
    return {
        "reranker": reranker.name, "cases": rows, "per_criterion": per_criterion,
        "veto_broken": veto_broken, "order_stable": order_stable,
        "passed": sum(1 for r in rows if r["pass"]), "total": len(rows),
        "verdict": "REJECTED (veto)" if veto_broken else
                   ("ACCEPTABLE" if all(v["ok"] for v in per_criterion.values()) else "PARTIAL"),
    }


def main() -> int:
    if "--freeze" in sys.argv:
        idx = sys.argv.index("--freeze")
        version = sys.argv[idx + 1] if len(sys.argv) > idx + 1 else "v1"
        return freeze(version)

    _, store, _, _ = build_state()
    cases = build_cases(store)
    spec = load_spec(cases)
    if spec is None:
        return 2
    threshold = spec["threshold"]

    results = [evaluate(r, cases, threshold) for r in
               (OverlapReranker(), TrapReranker(), NoiseReranker(), JointReranker())]

    print(f"spec {spec['version']} | cases {len(cases)} | threshold {threshold} "
          f"(frozen {spec['frozen_on']}, cases {spec['cases_sha256'][:12]})\n")
    print(f"{'reranker':<24} " + " ".join(f"{c:<7}" for c in CRITERIA) + "  veto  verdict")
    for res in results:
        cells = [f"{res['per_criterion'][c]['passed']}/{res['per_criterion'][c]['total']:<5}"
                 for c in CRITERIA]
        print(f"{res['reranker']:<24} " + " ".join(f"{c:<7}" for c in cells) +
              f"  {','.join(res['veto_broken']) or '-':<5} {res['verdict']}")
    print()
    for res in results:
        print(f"  {res['reranker']}: {res['passed']}/{res['total']}")
        for c in res["cases"]:
            if not c["pass"]:
                print(f"      {c['id']:<18} {c['failure_reason']}  (top={c['top']}, score={c['top_score']})")

    by_name = {r["reranker"]: r for r in results}
    scorer_params = len(inspect.signature(OverlapReranker.score).parameters)
    checks = {
        "trap_not_acceptable": by_name["trap-mock"]["verdict"] != "ACCEPTABLE",
        "baseline_fails_R1": not by_name["overlap-baseline"]["per_criterion"]["R1"]["ok"],
        "reference_passes_all": all(v["ok"] for v in by_name["joint-reference-mock"]["per_criterion"].values()),
        "veto_path_fires_and_rejects": by_name["noise-mock"]["verdict"].startswith("REJECTED")
                                       and by_name["noise-mock"]["veto_broken"] != [],
        "scorer_receives_no_label": scorer_params == 3 and
                                    all("expect" not in canonical(c["candidates"]) for c in cases),
        "order_invariance": all(r["order_stable"] for r in results),
    }
    print()
    for k, v in checks.items():
        print(f"  self-check {k}: {'PASS' if v else 'FAIL'}")

    receipt = {
        "suite": "s3.rerank_eval", "spec_version": spec["version"],
        "cases_sha256": spec["cases_sha256"], "harness_sha256": harness_hash(),
        "threshold": threshold, "criteria": CRITERIA, "vetoes": list(VETOES),
        "baseline_system_as_is": spec["baseline_system_as_is"],
        "results": results, "self_checks": checks,
        "scope_note": "mock rerankers only; no cross-encoder fetched or run. "
                      "Mocks prove the harness discriminates; they are not candidates. "
                      "Nine cases test mechanism, not generalization.",
    }
    RECEIPT_PATH.write_text(json.dumps(receipt, indent=2) + "\n", encoding="utf-8")
    ok = all(checks.values())
    print(f"\nRESULT: {'PASS' if ok else 'FAIL'} (harness self-checks) | "
          f"receipt: {RECEIPT_PATH} ({RECEIPT_PATH.stat().st_size} B)")
    return 0 if ok else 1


if __name__ == "__main__":
    sys.exit(main())
