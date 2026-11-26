"""Run experiment E-PROSE-RENDER against the frozen v1.2 protocol.

Only the rendering changes. Same model, runtime, nine cases, held-out calibration
discipline and separated ranking/abstention verdicts. Schema results are read from
the existing v1.2 receipt and never overwritten.

Usage: PYTHONPATH=. python eval/s3/prose_experiment.py
"""
from __future__ import annotations

import json
import pathlib
import sys

from eval.s3.bge_candidate import BgeReranker, health
from eval.s3.evaluate import build_state
from eval.s3.prose_render import CALIBRATION, render_cases
from eval.s3.rerank_eval import build_cases
from eval.s3.rerank_eval_v12 import load, run_cases

SCHEMA_RECEIPT = pathlib.Path("eval/s3/rerank-v12-receipt.json")
OUT = pathlib.Path("eval/s3/prose-experiment-receipt.json")


def calibrate_prose(reranker) -> dict:
    rel, irr = [], []
    for pid, q, text, is_rel in CALIBRATION:
        s = reranker.raw(q, [text])[0]
        (rel if is_rel else irr).append((pid, s))
    min_rel = min(s for _, s in rel)
    max_irr = max(s for _, s in irr)
    sep = min_rel - max_irr
    return {
        "relevant": {p: round(s, 3) for p, s in rel},
        "irrelevant": {p: round(s, 3) for p, s in irr},
        "min_relevant": round(min_rel, 3), "max_irrelevant": round(max_irr, 3),
        "separation": round(sep, 3),
        "threshold": round((min_rel + max_irr) / 2, 3) if sep > 0 else None,
        "calibratable": sep > 0,
    }


def main() -> int:
    spec = load()
    if spec is None:
        return 2
    if not health():
        print("candidate server not reachable on 8092")
        return 2

    r = BgeReranker()
    cal = calibrate_prose(r)
    _, store, _, _ = build_state()
    prose_cases = render_cases(build_cases(store))
    rows = run_cases(r, prose_cases, cal["threshold"])

    schema = json.loads(SCHEMA_RECEIPT.read_text(encoding="utf-8"))
    schema_rows = {row["id"]: row for row in schema["rows"]}

    print(f"E-PROSE-RENDER | candidate {r.name} | cases {spec['cases_sha256'][:16]}")
    print(f"calibration: min_rel {cal['min_relevant']} max_irr {cal['max_irrelevant']} "
          f"sep {cal['separation']} threshold {cal['threshold']}\n")
    print(f"{'case':<18} {'schema top':<12} {'schema':>8}  {'prose top':<12} {'prose':>8}  changed")
    for row in rows:
        s = schema_rows[row["id"]]
        changed = (s["ranking_top"] != row["ranking_top"])
        print(f"{row['id']:<18} {s['ranking_top']:<12} {s['ranking_score']:>8.3f}  "
              f"{row['ranking_top']:<12} {row['ranking_score']:>8.3f}  "
              f"{'YES' if changed else '-'}")

    rank_ok = sum(1 for x in rows if x["ranking_ok"])
    rank_total = sum(1 for x in rows if x["ranking_ok"] is not None)
    sel_ok = sum(1 for x in rows if x["selection_ok"])
    horse = rows[0]
    print(f"\nRANKING   (answerable, threshold-free): {rank_ok}/{rank_total}")
    print(f"SELECTION (calibrated abstention)     : {sel_ok}/{len(rows)}")
    print(f"SUBJECT DISTINCTION (R1-horse)        : "
          f"{'FAIL - trainer record first at %.3f' % horse['ranking_score'] if horse['ranking_top'] == 'rec-0001' else 'PASS - trainer record not first'}")
    print("\nThese three numbers stay separate. 5/5 ranking is not a subject-distinction pass.")

    receipt = {
        "experiment": "E-PROSE-RENDER",
        "condition": "natural prose; schema condition retained separately",
        "candidate": r.name, "cases_sha256": spec["cases_sha256"],
        "calibration": cal, "rows": rows,
        "metrics": {
            "ranking_correct": f"{rank_ok}/{rank_total}",
            "selection_correct": f"{sel_ok}/{len(rows)}",
            "subject_distinction": ("FAIL" if horse["ranking_top"] == "rec-0001" else "PASS"),
            "schema_ranking_correct": schema["verdicts"]["ranking"],
        },
        "comparison_note": "schema condition verdicts are in rerank-v12-receipt.json and "
                           "are not superseded by this experiment",
        "gate_b": "6/7 open",
    }
    OUT.write_text(json.dumps(receipt, indent=2) + "\n", encoding="utf-8")
    print(f"\nreceipt: {OUT} ({OUT.stat().st_size} B)")
    return 0


if __name__ == "__main__":
    sys.exit(main())
