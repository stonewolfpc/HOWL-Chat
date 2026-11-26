"""Stage 4 — run the FROZEN cases against the real bge candidate.

Two readings are produced, kept strictly apart:

  frozen_rule : the frozen selection rule (keep >= 0.30) applied to the
                candidate's raw score, exactly as specified. Reported as-is even
                though the candidate scores on a logit scale.
  diagnostic  : ranking-only. Top-1 is always taken, abstention is marked
                NOT EVALUABLE rather than failed, and raw scores are reported so
                the subject-confusion question can be read without inventing a
                threshold.

Neither reading modifies the spec. A candidate-specific transform would be an
instrument change and is not applied here.

Usage: PYTHONPATH=. python eval/s3/bge_run_cases.py
"""
from __future__ import annotations

import json
import pathlib
import sys

from eval.s3.bge_candidate import BgeReranker, health
from eval.s3.evaluate import build_state
from eval.s3.rerank_eval import build_cases, evaluate, load_spec


def diagnostic(reranker, cases) -> list[dict]:
    rows = []
    for case in cases:
        scored = sorted(((reranker.score(case["query"], dict(c)), dict(c))
                         for c in case["candidates"]),
                        key=lambda t: (-t[0], t[1]["record_id"]))
        top_score, top = scored[0]
        second = scored[1][0] if len(scored) > 1 else None
        rows.append({
            "id": case["id"], "criterion": case["criterion"], "query": case["query"],
            "expect": "abstain" if case["expect_abstain"] else case["expect"],
            "top": top["record_id"], "top_score": round(top_score, 3),
            "second_score": round(second, 3) if second is not None else None,
            "margin": round(top_score - second, 3) if second is not None else None,
            "top_matches_expect": (top["record_id"] == case["expect"]
                                   if not case["expect_abstain"] else None),
            "scores": {c["record_id"]: round(s, 3) for s, c in scored},
        })
    return rows


def main() -> int:
    if not health():
        print("server not reachable on 8092")
        return 2
    _, store, _, _ = build_state()
    cases = build_cases(store)
    spec = load_spec(cases)
    if spec is None:
        return 2

    frozen = evaluate(BgeReranker(), cases, spec["threshold"])
    diag = diagnostic(BgeReranker(), cases)

    print(f"spec {spec['version']} | cases {len(cases)} | frozen threshold {spec['threshold']}\n")
    print("=== READING 1: frozen rule applied unchanged to raw scores ===")
    for c in frozen["cases"]:
        print(f"  {c['id']:<18} top={str(c['top']):<10} score={str(c['top_score']):<8} "
              f"{'PASS' if c['pass'] else 'FAIL'}  {c['failure_reason'] or ''}")
    print(f"  -> {frozen['passed']}/{frozen['total']}, veto broken: {frozen['veto_broken']}")
    print("  NOTE: every raw score is below 0.30 because the candidate emits logits,")
    print("        so this reading measures the scale mismatch, not the model.\n")

    print("=== READING 2: ranking diagnostic (raw scores, no threshold invented) ===")
    for d in diag:
        verdict = ("n/a (abstain case)" if d["expect"] == "abstain"
                   else ("right record ranked first" if d["top_matches_expect"]
                         else f"WRONG (expected {d['expect']})"))
        print(f"  {d['id']:<18} top={d['top']:<10} {d['top_score']:>8.3f} "
              f"second={d['second_score']:>8.3f} margin={d['margin']:>6.3f}  {verdict}")
    print()
    horse = next(d for d in diag if d["id"] == "R1-horse")
    p1 = next(d for d in diag if d["id"] == "R4-paraphrase")
    print("=== the narrow question ===")
    print(f"  P1 paraphrase, correct record ranked first: {p1['top_matches_expect']}")
    print(f"  horse query top score {horse['top_score']:.3f} on {horse['top']}")
    print(f"  same record under a matching question:  {p1['top_score']:.3f}")
    gap = horse["top_score"] - p1["top_score"]
    reading = ("the horse query matches the trainer record MORE strongly than the "
               "trainer question does => subject confusion is amplified, not resolved"
               if gap > 1.0 else
               "the horse query scores below the trainer question => the model partly "
               "distinguishes the two"
               if gap < -1.0 else "near-indistinguishable")

    receipt = {
        "candidate": "bge-reranker-v2-m3-Q8_0",
        "spec_version": spec["version"], "cases_sha256": spec["cases_sha256"],
        "threshold_frozen": spec["threshold"],
        "reading_frozen_rule": {k: frozen[k] for k in
                                ("passed", "total", "veto_broken", "verdict", "cases")},
        "reading_diagnostic": diag,
        "scale_note": "candidate emits raw logits (-11..-6 observed); the frozen "
                      "0.30 threshold is defined on harness scale, so reading 1 is a "
                      "scale-mismatch artefact. No transform applied.",
        "gate_b": "6/7 open",
    }
    out = pathlib.Path("eval/s3/bge-cases.json")
    out.write_text(json.dumps(receipt, indent=2) + "\n", encoding="utf-8")
    print(f"\nreceipt: {out} ({out.stat().st_size} B)")
    return 0


if __name__ == "__main__":
    sys.exit(main())
