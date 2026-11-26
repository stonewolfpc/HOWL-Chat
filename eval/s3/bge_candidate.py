"""bge-reranker-v2-m3 candidate adapter + instrument checks.

This file is NOT part of the frozen apparatus. It imports the frozen harness
read-only and reports the candidate's own score scale without altering the
frozen selection rule.

Usage: PYTHONPATH=. python eval/s3/bge_candidate.py instrument
"""
from __future__ import annotations

import json
import math
import statistics
import sys
import time
import urllib.request

BASE = "http://127.0.0.1:8092"


class BgeReranker:
    name = "bge-reranker-v2-m3-Q8_0"

    def __init__(self, base: str = BASE, timeout: float = 120.0) -> None:
        self.base, self.timeout = base.rstrip("/"), timeout
        self.calls = 0

    def raw(self, query: str, documents: list[str]) -> list[float]:
        body = json.dumps({"query": query, "documents": documents}).encode()
        req = urllib.request.Request(f"{self.base}/v1/rerank", data=body,
                                     headers={"Content-Type": "application/json"})
        with urllib.request.urlopen(req, timeout=self.timeout) as r:
            payload = json.loads(r.read().decode())
        self.calls += 1
        order = {d["index"]: d["relevance_score"] for d in payload["results"]}
        return [order[i] for i in range(len(documents))]

    def score(self, query: str, cand: dict) -> float:
        return self.raw(query, [cand["text"]])[0]


def sigmoid(x: float) -> float:
    return 1.0 / (1.0 + math.exp(-x))


def health(base: str = BASE) -> bool:
    try:
        with urllib.request.urlopen(f"{base}/health", timeout=3) as r:
            return r.status == 200
    except Exception:
        return False


TRAINER = "Sable trained_by: Corvin Ashgrave of Ossary, and quiet was not part of it"
TEA = "Sable preference: tea more than I like surprises"


def instrument_checks() -> dict:
    r = BgeReranker()
    out: dict[str, dict] = {}

    # 1. Joint scoring: the same document must score differently under different
    #    queries. A document-only scorer would return identical values.
    q_a = "Who was your trainer, again?"
    q_b = "Do you prefer tea or coffee?"
    a1, b1 = r.raw(q_a, [TRAINER])[0], r.raw(q_b, [TRAINER])[0]
    out["joint_scoring"] = {
        "pass": abs(a1 - b1) > 1.0,
        "trainer_doc_under_trainer_query": round(a1, 3),
        "trainer_doc_under_tea_query": round(b1, 3),
        "delta": round(a1 - b1, 3),
        "note": "delta ~0 would mean the runtime is not scoring the pair jointly",
    }

    # 2. Rank discrimination: relevant above misleading.
    scores = r.raw(q_a, [TRAINER, TEA])
    out["rank_discrimination"] = {
        "pass": scores[0] > scores[1],
        "relevant": round(scores[0], 3), "irrelevant": round(scores[1], 3),
    }

    # 3. Reproducibility: repeated calls reproduce the ranking and the values.
    reps = [r.raw(q_a, [TRAINER, TEA]) for _ in range(3)]
    same_rank = all(s[0] > s[1] for s in reps)
    spread = max(abs(s[0] - reps[0][0]) for s in reps)
    out["reproducibility"] = {
        "pass": same_rank and spread < 0.01,
        "values": [[round(x, 4) for x in s] for s in reps],
        "max_drift": round(spread, 6),
    }

    # 4. Order invariance: reversing documents must not change the ranking.
    rev = r.raw(q_a, [TEA, TRAINER])
    out["order_invariance"] = {
        "pass": (rev[1] > rev[0]) == (scores[0] > scores[1]),
        "forward": [round(x, 3) for x in scores],
        "reversed": [round(x, 3) for x in rev],
    }

    # 5. Degenerate output: no NaN/inf, not constant across different documents.
    probe = r.raw(q_a, [TRAINER, TEA, "Alaric vault_code: seven-four-nine-one"])
    finite = all(math.isfinite(x) for x in probe)
    nonconstant = len({round(x, 6) for x in probe}) > 1
    out["nondegenerate"] = {
        "pass": finite and nonconstant,
        "values": [round(x, 3) for x in probe],
    }

    # 6. Score scale: reported, not corrected.
    out["score_scale"] = {
        "raw_range_observed": [round(min(probe), 3), round(max(probe), 3)],
        "sigmoid_of_examples": [round(sigmoid(x), 6) for x in probe],
        "interpretation": "raw logits, not probabilities; frozen threshold 0.30 "
                          "is not on this scale",
    }
    return out


def main() -> int:
    if not health():
        print("server not reachable on 8092")
        return 2
    checks = instrument_checks()
    print(json.dumps(checks, indent=2))
    core = ("joint_scoring", "rank_discrimination", "reproducibility",
            "order_invariance", "nondegenerate")
    ok = all(checks[k]["pass"] for k in core)
    print(f"\nSTAGE3 INSTRUMENT SANITY: {'PASS' if ok else 'FAIL'}")
    out = {"checks": checks, "stage3": "PASS" if ok else "FAIL"}
    with open("eval/s3/bge-instrument.json", "w", encoding="utf-8") as fh:
        json.dump(out, fh, indent=2)
    return 0 if ok else 1


if __name__ == "__main__":
    sys.exit(main())
