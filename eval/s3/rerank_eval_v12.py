"""Protocol v1.2 — score scale, ranking, and abstention kept apart.

v1.1 stays byte-immutable at eval/s3/rerank-spec-v1.1.json. This module is a new
protocol version, not an edit: the case set and its hash are inherited unchanged.

v1.2 rules:
  * a candidate declares its score scale and its transform before scoring;
  * RANKING is judged threshold-free (is the expected record ranked first?);
  * ABSTENTION is judged only with a threshold calibrated on held-out pairs;
  * if the calibration bands overlap, abstention is UNCALIBRATABLE — never fitted
    anyway, and never allowed to launder a wrong ranking into a pass.

Usage:
  PYTHONPATH=. python eval/s3/rerank_eval_v12.py --freeze
  PYTHONPATH=. python eval/s3/rerank_eval_v12.py
"""
from __future__ import annotations

import hashlib
import json
import pathlib
import sys

from eval.s3.bge_candidate import BgeReranker, health
from eval.s3.calibration_set import PAIRS, RULE
from eval.s3.evaluate import build_state
from eval.s3.rerank_eval import build_cases, cases_hash, harness_hash  # read-only use
from eval.s3.rerank_set import CRITERIA, VETOES

SPEC_V11 = pathlib.Path("eval/s3/rerank-spec-v1.1.json")
SPEC = pathlib.Path("eval/s3/rerank-spec-v1.2.json")
RECEIPT = pathlib.Path("eval/s3/rerank-v12-receipt.json")


def spec_hash() -> str:
    return hashlib.sha256(pathlib.Path(__file__).read_bytes()).hexdigest()


def freeze() -> int:
    v11 = json.loads(SPEC_V11.read_text(encoding="utf-8"))
    _, store, _, _ = build_state()
    payload = cases_hash(_payload(build_cases(store)))
    if payload != v11["cases_sha256"]:
        print("REFUSING: case set differs from v1.1; v1.2 must inherit it unchanged")
        return 2
    spec = {
        "version": "v1.2",
        "inherits": {"spec": "v1.1", "sha256": v11["harness_sha256"],
                     "cases_sha256": payload},
        "cases_sha256": payload,
        "score_scale": {
            "policy": "candidate declares raw|prob and its transform before scoring; "
                      "the transform is recorded, never silent",
            "transforms_allowed": ["identity", "sigmoid"],
        },
        "ranking_rule": "expected record ranked first among candidates; threshold-free",
        "abstention_rule": RULE,
        "abstention_note": "calibration uses eval/s3/calibration_set.py only; the nine "
                           "cases never inform the threshold",
        "veto_policy": "R4 and R5 are judged on the selection verdict; a ranking failure "
                       "is reported separately and is never offset by calibration",
        "conditions": [
            {"id": "C-SCHEMA", "status": "original condition, preserved",
             "detail": "records rendered as 'Subject predicate: value'"},
            {"id": "E-PROSE-RENDER", "status": "registered, not run",
             "detail": "natural-language rendering is a separate experiment and must "
                       "not replace C-SCHEMA"},
        ],
        "harness_sha256": spec_hash(),
        "frozen_on": "2026-11-15",
    }
    SPEC.write_text(json.dumps(spec, indent=2) + "\n", encoding="utf-8")
    print(f"froze v1.2 | cases {payload[:16]} (inherited unchanged) | "
          f"harness {spec['harness_sha256'][:16]}")
    return 0


def _payload(cases: list[dict]) -> list[dict]:
    return [{k: c[k] for k in ("id", "criterion", "owner", "query", "expect",
                               "expect_abstain", "why", "candidates")} for c in cases]


def load() -> dict | None:
    if not SPEC.exists():
        print("no v1.2 spec; run --freeze")
        return None
    spec = json.loads(SPEC.read_text(encoding="utf-8"))
    _, store, _, _ = build_state()
    live = cases_hash(_payload(build_cases(store)))
    if live != spec["cases_sha256"]:
        print("REFUSING: case set changed since v1.2 freeze")
        return None
    if spec_hash() != spec["harness_sha256"]:
        print("REFUSING: v1.2 harness changed since freeze")
        return None
    return spec


def calibrate(reranker) -> dict:
    rel, irr = [], []
    for pid, query, text, is_rel in PAIRS:
        s = reranker.raw(query, [text])[0]
        (rel if is_rel else irr).append((pid, s))
    min_rel = min(s for _, s in rel)
    max_irr = max(s for _, s in irr)
    separation = min_rel - max_irr
    ok = separation > 0
    return {
        "relevant": {p: round(s, 3) for p, s in rel},
        "irrelevant": {p: round(s, 3) for p, s in irr},
        "min_relevant": round(min_rel, 3), "max_irrelevant": round(max_irr, 3),
        "separation": round(separation, 3),
        "threshold": round((min_rel + max_irr) / 2, 3) if ok else None,
        "calibratable": ok,
    }


def run_cases(reranker, cases, threshold) -> list[dict]:
    rows = []
    for case in cases:
        docs = [c["text"] for c in case["candidates"]]
        scores = reranker.raw(case["query"], docs)
        paired = sorted(zip(scores, case["candidates"]),
                        key=lambda t: (-t[0], t[1]["record_id"]))
        ranking_top = paired[0][1]["record_id"]
        keeps = [(s, c) for s, c in paired if s >= threshold] if threshold is not None else []
        abstained = not keeps
        sel_top = keeps[0][1]["record_id"] if keeps else None
        if case["expect_abstain"]:
            sel_ok = abstained
        else:
            sel_ok = (not abstained) and sel_top == case["expect"]
        rank_ok = (None if case["expect_abstain"]
                   else ranking_top == case["expect"])
        rows.append({
            "id": case["id"], "criterion": case["criterion"],
            "expect": "abstain" if case["expect_abstain"] else case["expect"],
            "ranking_top": ranking_top, "ranking_score": round(paired[0][0], 3),
            "ranking_ok": rank_ok,
            "selection_top": sel_top, "abstained": abstained, "selection_ok": sel_ok,
            "scores": {c["record_id"]: round(s, 3) for s, c in paired},
        })
    return rows


def main() -> int:
    if "--freeze" in sys.argv:
        return freeze()
    spec = load()
    if spec is None:
        return 2
    if not health():
        print("candidate server not reachable on 8092")
        return 2

    reranker = BgeReranker()
    cal = calibrate(reranker)
    _, store, _, _ = build_state()
    rows = run_cases(reranker, build_cases(store), cal["threshold"])

    print(f"v1.2 | cases {spec['cases_sha256'][:16]} | candidate {reranker.name}")
    print(f"declared scale: raw logits (identity transform)\n")
    print("=== calibration (held-out pairs, never the nine cases) ===")
    print(f"  min relevant {cal['min_relevant']} | max irrelevant {cal['max_irrelevant']} | "
          f"separation {cal['separation']}")
    print(f"  threshold {'%.3f' % cal['threshold'] if cal['threshold'] is not None else 'UNCALIBRATABLE'}\n")

    print("=== RANKING (threshold-free) ===")
    for r in rows:
        if r["ranking_ok"] is None:
            print(f"  {r['id']:<18} (abstain case) top={r['ranking_top']} {r['ranking_score']:>8.3f}")
        else:
            print(f"  {r['id']:<18} {r['ranking_top']:<10} {r['ranking_score']:>8.3f}  "
                  f"{'ok' if r['ranking_ok'] else 'WRONG'}")

    print("\n=== SELECTION (calibrated abstention) ===")
    for r in rows:
        print(f"  {r['id']:<18} top={str(r['selection_top']):<10} abstained={str(r['abstained']):<5} "
              f"{'PASS' if r['selection_ok'] else 'FAIL'}")

    per_crit = {}
    for crit in CRITERIA:
        sel = [r for r in rows if r["criterion"] == crit]
        per_crit[crit] = {
            "selection_ok": all(r["selection_ok"] for r in sel),
            "ranking_ok": all(r["ranking_ok"] is not False for r in sel),
            "veto": crit in VETOES,
        }
    sel_veto = [c for c in VETOES if not per_crit[c]["selection_ok"]]
    rank_veto = [c for c in VETOES if not per_crit[c]["ranking_ok"]]

    print("\n=== verdicts, kept separate ===")
    print(f"  ranking  : {sum(1 for r in rows if r['ranking_ok'])}/"
          f"{sum(1 for r in rows if r['ranking_ok'] is not None)} answerable cases correct"
          + (f" | R4 ranking {'FAIL' if rank_veto else 'ok'}"))
    print(f"  selection: {sum(1 for r in rows if r['selection_ok'])}/{len(rows)} "
          f"| veto broken: {sel_veto or 'none'}")
    print(f"  subject distinction (R1-horse): top={rows[0]['ranking_top']} "
          f"score={rows[0]['ranking_score']} -> "
          f"{'CONFUSION PERSISTS' if rows[0]['ranking_top'] == 'rec-0001' else 'distinguished'}")

    receipt = {
        "protocol": "v1.2", "candidate": reranker.name,
        "declared_scale": "raw_logit", "transform": "identity",
        "cases_sha256": spec["cases_sha256"], "calibration": cal,
        "rows": rows, "per_criterion": per_crit,
        "veto_broken_selection": sel_veto, "veto_broken_ranking": rank_veto,
        "verdicts": {
            "ranking": "FAIL" if rank_veto else "PASS",
            "selection": "FAIL" if sel_veto else "PASS",
            "subject_distinction": ("FAIL" if rows[0]["ranking_top"] == "rec-0001"
                                    else "PASS"),
        },
        "note": "ranking and abstention judged separately; calibration cannot convert a "
                "wrong ranking into a pass",
        "gate_b": "6/7 open",
    }
    RECEIPT.write_text(json.dumps(receipt, indent=2) + "\n", encoding="utf-8")
    print(f"\nreceipt: {RECEIPT} ({RECEIPT.stat().st_size} B)")
    return 0


if __name__ == "__main__":
    sys.exit(main())
