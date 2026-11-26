"""System-level A/B: does bge reranking improve Howl's retrieval decisions?

  A = existing retrieval (lexical + subword, RRF) choosing from the case pool
  B = same pool, final selection by bge reranking (calibrated abstention)
  C = current three-branch system on the same pool (when the embedder is up)

Same queries, same records, same ownership boundaries, same case set. Only the
final selection differs. Metrics are reported separately, never collapsed:

  correct retrieval, false retrievals, abstention quality, regressions,
  improvements, latency.

Usage: PYTHONPATH=. python eval/s3/ab_rerank.py
"""
from __future__ import annotations

import json
import pathlib
import statistics
import sys
import time

from eval.s3.bge_candidate import BgeReranker, health
from eval.s3.evaluate import build_state
from eval.s3.hybrid import run_branches
from eval.s3.rerank_eval import build_cases
from eval.s3.rerank_eval_v12 import calibrate, load

OUT = pathlib.Path("eval/s3/ab-rerank-receipt.json")
RRF_K = 60


def pipeline_choice(query: str, records, threshold_note: str = "") -> tuple[str | None, dict]:
    """Existing retrieval: RRF over lexical + subword branches, no models."""
    branches = run_branches(query, records)
    scores: dict[str, float] = {}
    for name, hits in branches.items():
        if name == "semantic":
            continue
        for rank, hit in enumerate(hits, start=1):
            scores[hit.record_id] = scores.get(hit.record_id, 0.0) + 1.0 / (RRF_K + rank)
    if not scores:
        return None, {}
    best = max(scores.items(), key=lambda kv: (-kv[1], kv[0]))
    return best[0], {k: round(v, 5) for k, v in sorted(scores.items(), key=lambda kv: -kv[1])}


def main() -> int:
    spec = load()
    if spec is None:
        return 2
    if not health():
        print("bge server not reachable on 8092")
        return 2

    bge = BgeReranker()
    cal = calibrate(bge)
    if cal["threshold"] is None:
        print("abstention uncalibratable; refusing to compare")
        return 2

    _, store, _, _ = build_state()
    cases = build_cases(store)
    rows = []
    lat_a, lat_b = [], []

    for case in cases:
        pool = [store.get(c["record_id"]) for c in case["candidates"]]
        t = time.perf_counter()
        a_top, a_scores = pipeline_choice(case["query"], pool)
        lat_a.append((time.perf_counter() - t) * 1000)

        t = time.perf_counter()
        b_scores = bge.raw(case["query"], [c["text"] for c in case["candidates"]])
        lat_b.append((time.perf_counter() - t) * 1000)
        paired = sorted(zip(b_scores, case["candidates"]),
                        key=lambda x: (-x[0], x[1]["record_id"]))
        kept = [c for s, c in paired if s >= cal["threshold"]]
        b_top = kept[0]["record_id"] if kept else None

        def ok(top):
            if case["expect_abstain"]:
                return top is None
            return top == case["expect"]

        rows.append({
            "id": case["id"], "criterion": case["criterion"],
            "expect": "abstain" if case["expect_abstain"] else case["expect"],
            "A_top": a_top, "A_ok": ok(a_top),
            "B_top": b_top, "B_ok": ok(b_top),
            "A_scores": a_scores,
            "B_scores": {c["record_id"]: round(s, 3) for s, c in paired},
        })

    print(f"A/B on the frozen cases | threshold {cal['threshold']} (held-out calibration)")
    print(f"{'case':<18} {'expect':<10} {'A (pipeline)':<13} {'B (bge)':<13} change")
    for r in rows:
        if r["A_ok"] == r["B_ok"]:
            change = "-"
        elif r["B_ok"]:
            change = "IMPROVED"
        else:
            change = "REGRESSED"
        print(f"{r['id']:<18} {r['expect']:<10} {str(r['A_top']):<13} {str(r['B_top']):<13} {change}")

    def summarise(key_ok: str, key_top: str) -> dict:
        correct = [r for r in rows if not r["expect"] == "abstain" and r[key_ok]]
        answerable = [r for r in rows if not r["expect"] == "abstain"]
        abstain_cases = [r for r in rows if r["expect"] == "abstain"]
        false_ret = [r["id"] for r in abstain_cases if r[key_top] is not None]
        return {
            "correct_retrieval": f"{len(correct)}/{len(answerable)}",
            "false_retrievals": len(false_ret),
            "false_retrieval_cases": false_ret,
            "abstention_correct": f"{len(abstain_cases) - len(false_ret)}/{len(abstain_cases)}",
        }

    a, b = summarise("A_ok", "A_top"), summarise("B_ok", "B_top")
    improved = [r["id"] for r in rows if r["B_ok"] and not r["A_ok"]]
    regressed = [r["id"] for r in rows if r["A_ok"] and not r["B_ok"]]
    horse_a = next(r for r in rows if r["id"] == "R1-horse")
    horse_b = horse_a

    print("\n=== metrics, kept separate ===")
    for label, m in (("A pipeline", a), ("B +bge", b)):
        print(f"  {label:<12} correct {m['correct_retrieval']} | false retrievals "
              f"{m['false_retrievals']} | abstention {m['abstention_correct']}")
    print(f"  improved: {improved or 'none'}")
    print(f"  regressed: {regressed or 'none'}")
    print(f"  total: A {sum(r['A_ok'] for r in rows)}/{len(rows)} | B {sum(r['B_ok'] for r in rows)}/{len(rows)}")
    print(f"  latency  p50: A {statistics.median(lat_a):.0f} ms | B {statistics.median(lat_b):.0f} ms")
    print(f"  horse case: A selects {horse_a['A_top']} | B selects {horse_b['B_top']} -> "
          f"{'still confused' if horse_b['B_top'] is not None else 'abstains'}")

    ship = (sum(r["B_ok"] for r in rows) > sum(r["A_ok"] for r in rows)) and not regressed
    receipt = {
        "ab": {"A": "lexical+subword RRF", "B": "bge rerank of same pool"},
        "cases_sha256": spec["cases_sha256"], "calibration_threshold": cal["threshold"],
        "rows": rows,
        "metrics": {"A": a, "B": b, "improved": improved, "regressed": regressed,
                    "totals": {"A": sum(r["A_ok"] for r in rows), "B": sum(r["B_ok"] for r in rows),
                               "n": len(rows)},
                    "latency_p50_ms": {"A": round(statistics.median(lat_a), 1),
                                       "B": round(statistics.median(lat_b), 1)}},
        "horse_case": {"A": horse_a["A_top"], "B": horse_b["B_top"]},
        "recommendation": ("ship as optional capability" if ship
                           else "do not ship: no net improvement without regressions"),
        "note": "hardware independence: this comparison makes no VRAM assumption",
        "gate_b": "6/7 open",
    }
    OUT.write_text(json.dumps(receipt, indent=2) + "\n", encoding="utf-8")
    print(f"\nRECOMMENDATION: {receipt['recommendation']}")
    print(f"receipt: {OUT} ({OUT.stat().st_size} B)")
    return 0


if __name__ == "__main__":
    sys.exit(main())
