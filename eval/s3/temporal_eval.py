"""E-TEMPORAL-SET — does bge consistently help distinguish was-true from is-true?

A = existing retrieval (lexical + subword RRF) over the case pool
B = same pool, final selection by bge with abstention at a threshold calibrated
    on this store's own held-out pairs (never on the cases)

Same rule as everywhere else in S3: the case set and its expectations are frozen
before the candidate runs, and the threshold is never tuned on the cases.

Run:  PYTHONPATH=. python eval/s3/temporal_eval.py --freeze
      PYTHONPATH=. python eval/s3/temporal_eval.py
"""
from __future__ import annotations

import hashlib
import json
import pathlib
import statistics
import sys
import time

from eval.s3.bge_candidate import BgeReranker, health
from eval.s3.hybrid import run_branches
from eval.s3.temporal_set import CALIBRATION, FROZEN_NOTE, build_cases, build_store

SPEC = pathlib.Path("eval/s3/temporal-spec.json")
OUT = pathlib.Path("eval/s3/temporal-receipt.json")
RRF_K = 60
SCORED = {"scope_gap_known"}


def canon(x) -> str:
    return json.dumps(x, sort_keys=True, separators=(",", ":"))


def payload(cases: list[dict]) -> list[dict]:
    return [{k: c[k] for k in ("id", "criterion", "owner", "query", "expect",
                               "expect_abstain", "candidates")} for c in cases]


def freeze() -> int:
    store = build_store()
    cases = payload(build_cases(store))
    spec = {
        "experiment": "E-TEMPORAL-SET", "version": "v2", "frozen_on": "2026-11-15",
        "supersedes": "v1 (cases 8c6f834a) — withdrawn before any result was used: it "
                      "compared record ids against expected values, and its scope probe "
                      "built a pool that bypassed the ownership filter it meant to test",
        "cases_sha256": hashlib.sha256(canon(cases).encode()).hexdigest(),
        "cases": cases,
        "selection_rule": "keep scores >= calibrated threshold; if none, abstain",
        "threshold_note": "calibrated on this store's held-out pairs at run time; "
                          "never fitted to these cases",
        "gap_cases": sorted(SCORED),
        "isolation_note": FROZEN_NOTE,
    }
    SPEC.write_text(json.dumps(spec, indent=2) + "\n", encoding="utf-8")
    print(f"froze E-TEMPORAL-SET {spec['version']}: {len(cases)} cases | "
          f"cases {spec['cases_sha256'][:16]} | {SPEC} ({SPEC.stat().st_size} B)")
    return 0


def load_spec(cases: list[dict]) -> dict | None:
    if not SPEC.exists():
        print("REFUSING: no frozen spec. Run --freeze first.")
        return None
    spec = json.loads(SPEC.read_text())
    live = hashlib.sha256(canon(cases).encode()).hexdigest()
    if live != spec["cases_sha256"]:
        print(f"REFUSING: cases changed since freeze ({spec['cases_sha256'][:16]} -> {live[:16]})")
        return None
    return spec


def calibrate(bge) -> dict:
    rel, irr = [], []
    for pid, q, text, is_rel in CALIBRATION:
        s = bge.raw(q, [text])[0]
        (rel if is_rel else irr).append((pid, round(s, 3)))
    min_rel = min(s for _, s in rel)
    max_irr = max(s for _, s in irr)
    sep = min_rel - max_irr
    return {"relevant": dict(rel), "irrelevant": dict(irr),
            "min_relevant": min_rel, "max_irrelevant": max_irr,
            "separation": round(sep, 3),
            "threshold": round((min_rel + max_irr) / 2, 3) if sep > 0 else None}


def pipeline(query: str, records) -> str | None:
    branches = run_branches(query, records)
    acc: dict[str, float] = {}
    for name, hits in branches.items():
        if name == "semantic":
            continue
        for rank, hit in enumerate(hits, start=1):
            acc[hit.record_id] = acc.get(hit.record_id, 0.0) + 1.0 / (RRF_K + rank)
    return max(acc.items(), key=lambda kv: (-kv[1], kv[0]))[0] if acc else None


def main() -> int:
    if "--freeze" in sys.argv:
        return freeze()
    store = build_store()
    cases = build_cases(store)
    spec = load_spec(payload(cases))
    if spec is None:
        return 2
    if not health():
        print("bge server not reachable on 8092")
        return 2

    bge = BgeReranker()
    cal = calibrate(bge)
    if cal["threshold"] is None:
        print("REFUSING: abstention uncalibratable on this store")
        return 2

    rows = []
    lat_a, lat_b = [], []
    for case in cases:
        recs = [store.get(c["record_id"]) for c in case["candidates"]]
        t = time.perf_counter(); a_top = pipeline(case["query"], recs)
        lat_a.append((time.perf_counter() - t) * 1000)
        t = time.perf_counter(); scores = bge.raw(case["query"], [c["text"] for c in case["candidates"]])
        lat_b.append((time.perf_counter() - t) * 1000)
        paired = sorted(zip(scores, case["candidates"]), key=lambda x: (-x[0], x[1]["record_id"]))
        kept = [c for s, c in paired if s >= cal["threshold"]]
        b_top = kept[0]["record_id"] if kept else None

        def ok(top):
            # expectations are values; tops are record ids
            if case["expect_abstain"]:
                return top is None
            return top is not None and store.get(top).value == case["expect"]

        rows.append({"id": case["id"], "criterion": case["criterion"], "gap": case["gap"],
                     "expect": "abstain" if case["expect_abstain"] else case["expect"],
                     "A_top": a_top, "A_value": store.get(a_top).value if a_top else None,
                     "A_ok": ok(a_top), "B_top": b_top,
                     "B_value": store.get(b_top).value if b_top else None,
                     "B_ok": ok(b_top),
                     "B_scores": {c["record_id"]: round(s, 3) for s, c in paired},
                     "conflict": "conflicting_timestamps" == case["criterion"]})

    print(f"E-TEMPORAL-SET {spec['version']} | cases {len(rows)} | threshold {cal['threshold']} "
          f"(separation {cal['separation']}, held-out)\n")
    print(f"{'case':<26} {'expect':<16} {'A':<20} {'B':<20} change")
    for r in rows:
        ch = "-" if r["A_ok"] == r["B_ok"] else ("IMPROVED" if r["B_ok"] else "REGRESSED")
        tag = " [gap]" if r["gap"] else ""
        print(f"{r['id']:<26} {str(r['expect']):<16} {str(r['A_value']):<20} "
              f"{str(r['B_value']):<20} {ch}{tag}")

    scored = [r for r in rows if not r["gap"]]
    a_ok = sum(r["A_ok"] for r in scored); b_ok = sum(r["B_ok"] for r in scored)
    improved = [r["id"] for r in scored if r["B_ok"] and not r["A_ok"]]
    regressed = [r["id"] for r in scored if r["A_ok"] and not r["B_ok"]]
    hist = [r for r in scored if r["criterion"].startswith(("explicit_history", "correction_history"))]
    cur = [r for r in scored if "current" in r["criterion"] or "without_temporal" in r["criterion"]]
    own = [r for r in scored if r["criterion"].startswith("ownership")]

    def tally(group):
        return f"A {sum(r['A_ok'] for r in group)}/{len(group)} | B {sum(r['B_ok'] for r in group)}/{len(group)}"

    print("\n=== metrics, separate by construction ===")
    print(f"  scored cases        : {len(scored)}  A {a_ok}/{len(scored)} | B {b_ok}/{len(scored)}")
    print(f"  current-state       : {tally(cur)}")
    print(f"  historical-state    : {tally(hist)}")
    print(f"  ownership           : {tally(own)}")
    print(f"  improved            : {improved or 'none'}")
    print(f"  regressed           : {regressed or 'none'}")
    print(f"  latency p50         : A {statistics.median(lat_a):.0f} ms | B {statistics.median(lat_b):.0f} ms")
    gap_rows = [r for r in rows if r["gap"]]
    print(f"  scope gap (unscored): world fact under a character filter -> "
          f"A {'reachable' if gap_rows[0]['A_top'] else 'unreachable'}, "
          f"B {'reachable' if gap_rows[0]['B_top'] else 'unreachable'}")

    receipt = {
        "experiment": "E-TEMPORAL-SET", "version": spec["version"],
        "cases_sha256": spec["cases_sha256"], "calibration": cal,
        "rows": rows,
        "metrics": {
            "scored": len(scored), "A_correct": a_ok, "B_correct": b_ok,
            "improved": improved, "regressed": regressed,
            "current_state": {"A": sum(r["A_ok"] for r in cur), "B": sum(r["B_ok"] for r in cur), "n": len(cur)},
            "historical_state": {"A": sum(r["A_ok"] for r in hist), "B": sum(r["B_ok"] for r in hist), "n": len(hist)},
            "ownership": {"A": sum(r["A_ok"] for r in own), "B": sum(r["B_ok"] for r in own), "n": len(own)},
            "latency_p50_ms": {"A": round(statistics.median(lat_a), 1), "B": round(statistics.median(lat_b), 1)},
        },
        "scope_gap": {"note": "chat-scope records are unreachable under a character filter; "
                              "scope/ownership modelling is a parked architectural workstream",
                      "A": gap_rows[0]["A_top"], "B": gap_rows[0]["B_top"]},
        "verdict": "bge improves current-state resolution without regressions on this set"
                   if b_ok > a_ok and not regressed else
                   ("bge matches the pipeline on this set" if b_ok == a_ok and not regressed else
                    "bge is not consistently better on this set"),
        "frozen_note": FROZEN_NOTE,
    }
    OUT.write_text(json.dumps(receipt, indent=2) + "\n", encoding="utf-8")
    print(f"\nVERDICT: {receipt['verdict']}")
    print(f"receipt: {OUT} ({OUT.stat().st_size} B)")
    return 0


if __name__ == "__main__":
    sys.exit(main())
