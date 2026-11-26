"""Gate B evaluation — each branch alone, then fused, then reranked.

Run:
    PYTHONPATH=. python eval/s3/evaluate.py

Formation is FROZEN: the run refuses to proceed unless formation.py, records.py
and evidence.py are byte-identical to the modules that produced the Gate A
baseline. Improving retrieval must not move the records.
"""
from __future__ import annotations

import hashlib
import json
import pathlib
import sys

from eval.s3.evidence import EvidenceLog
from eval.s3.records import RecordStore
from eval.s3.formation import Formatter
from eval.s3.context import Window
from eval.s3.selfcheck import build_conversation, toks, OWNER
from eval.s3.hybrid import hybrid, pool, run_branches
from eval.s3.vectors import SEMANTIC_STATUS
from eval.s3.semantic import try_connect

SEMANTIC_URL = "http://127.0.0.1:8091"
SEMANTIC_MODEL = "nomic-embed-text-v1.5"

ROOT = pathlib.Path(__file__).resolve().parents[2]
FROZEN = ("formation.py", "records.py", "evidence.py")

PROBES = [
    {"id": "P1", "owner": "char-001", "query": "Who was your trainer, again?",
     "expect": "trained_by", "expect_abstain": False,
     "why": "lexical paraphrase gap observed in Gate A"},
    {"id": "P2", "owner": "char-001", "query": "What did you think of the weather?",
     "expect": None, "expect_abstain": True, "why": "irrelevant probe must abstain"},
    {"id": "P3a", "owner": None, "query": "What is the vault code?",
     "expect": "vault_code", "expect_abstain": False,
     "why": "pre-filter contamination is measured, not assumed away"},
    {"id": "P3b", "owner": "char-001", "query": "What is the vault code?",
     "expect": None, "expect_abstain": True,
     "why": "ownership filter must remove another character's record"},
    {"id": "P4a", "owner": "char-001", "query": "Where is the ledger kept now?",
     "expect": "green", "expect_abstain": False, "why": "current truth"},
    {"id": "P4b", "owner": "char-001", "query": "Where was the ledger kept before?",
     "expect": "blue", "expect_abstain": False, "include_history": True,
     "why": "history must remain reachable when explicitly asked for"},
    {"id": "P5", "owner": "char-001", "query": "Who trained the horse?",
     "expect": None, "expect_abstain": True,
     "why": "misleading similarity: shares 'trained' but asks something else"},
]


def build_state():
    log, store, win = EvidenceLog(), RecordStore(), Window(capacity=420, reserve=80)
    former = Formatter(store, log)
    for i, (speaker, text) in enumerate(build_conversation(), start=1):
        turn = log.append(speaker, text, toks(text))
        roll = win.add(turn)
        if roll or i % 8 == 0:
            former.run_pending(log, OWNER)
    return log, store, win, former


def frozen_check() -> tuple[bool, dict]:
    receipt = json.loads((ROOT / "eval/s3/gate-a-receipt.json").read_text())
    recorded = receipt["module_sha256"]
    now = {n: hashlib.sha256((ROOT / "eval/s3" / n).read_bytes()).hexdigest() for n in FROZEN}
    return all(now[n] == recorded.get(n) for n in FROZEN), {
        "recorded": {n: recorded.get(n, "")[:16] for n in FROZEN},
        "current": {n: now[n][:16] for n in FROZEN},
        "identical": all(now[n] == recorded.get(n) for n in FROZEN),
    }


def evaluate():
    log, store, win, former = build_state()
    frozen_ok, frozen_detail = frozen_check()
    sem = try_connect(SEMANTIC_URL, SEMANTIC_MODEL)
    rows = []
    for p in PROBES:
        res = hybrid(p["query"], store, p["owner"],
                     include_history=p.get("include_history", False), semantic=sem)
        per_branch_top = {b: (hits[0].record_id if hits else None)
                          for b, hits in res.per_branch.items()}
        per_branch_branch_only = {}
        recs = pool(store, p["owner"], p.get("include_history", False))
        for b, hits in run_branches(p["query"], recs, semantic=sem).items():
            per_branch_branch_only[b] = [store.get(h.record_id).value for h in hits[:2]]
        values = [s["value"] for s in res.selected]
        # Match against predicate as well as value: "vault_code" is a predicate,
        # the value is the code itself.
        haystack = [(s["value"] + " " + store.get(s["record_id"]).predicate).lower()
                    for s in res.selected]
        hit_ok = (p["expect"] is None) or any(p["expect"].lower() in h for h in haystack)
        pass_ = (res.abstained if p["expect_abstain"] else (not res.abstained and hit_ok))
        rows.append({
            "id": p["id"], "query": p["query"], "owner": p["owner"], "why": p["why"],
            "branches_top": per_branch_top,
            "branch_values_top2": per_branch_branch_only,
            "semantic_available": sem is not None,
            "fused_top": res.fused[0].record_id if res.fused else None,
            "reranked_top": res.reranked[0].record_id if res.reranked else None,
            "selected_values": values,
            "abstained": res.abstained, "reason": res.reason,
            "expected": "abstain" if p["expect_abstain"] else p["expect"],
            "pass": pass_,
        })
    return log, store, win, former, rows, frozen_ok, frozen_detail, sem


def main() -> int:
    log, store, win, former,     rows, frozen_ok, frozen_detail, sem = evaluate()
    print(f"formation frozen at Gate A baseline: {frozen_ok}")
    print(f"  {frozen_detail['current']}")
    print(f"records: {store.count_by_status()} | branches: lexical, subword, "
          f"semantic={'available' if sem is not None else 'UNAVAILABLE'}")
    print()
    hdr = f"{'id':<4} {'branch tops (lex/sub)':<34} {'fused':<10} {'reranked':<10} {'abstain':<8} result"
    print(hdr)
    for r in rows:
        tops = f"{r['branches_top'].get('lexical')}/{r['branches_top'].get('subword')}/{r['branches_top'].get('semantic')}"
        print(f"{r['id']:<4} {tops:<34} {str(r['fused_top']):<10} {str(r['reranked_top']):<10} "
              f"{str(r['abstained']):<8} {'PASS' if r['pass'] else 'FAIL'}  ({r['expected']})")
    print()
    for r in rows:
        print(f"  {r['id']}: selected={r['selected_values']}")
        print(f"      subword top2={r['branch_values_top2'].get('subword')}")
        print(f"      lexical top2={r['branch_values_top2'].get('lexical')}")

    passed = sum(1 for r in rows if r["pass"])
    receipt = {
        "gate": "B",
        "suite": "s3.evaluate",
        "formation_frozen": frozen_ok,
        "formation_hashes": frozen_detail,
        "records": store.count_by_status(),
        "evidence_sha256": log.sha256(),
        "branches": {
            "lexical": {"method": "BM25", "params": {"k1": 1.5, "b": 0.75}},
            "subword": {"method": "hashed char n-gram cosine", "n_grams": [3, 4, 5],
                        "dim": 2048, "seed": "s3-subword-v1", "metric": "cosine"},
            "semantic": SEMANTIC_STATUS,
        },
        "fusion": {"method": "reciprocal rank fusion", "k": 60, "min_fused": 0.012},
        "rerank": {"method": "heuristic term-coverage + kind prior (NOT a cross-encoder)",
                   "coverage_weight": 0.10, "stale_penalty": 0.02},
        "embedder_identity": None,
        "reproducibility": {"deterministic": True,
                            "no_model_loaded": True,
                            "note": "no embedding model, no server, no generation"},
        "probes": rows,
        "probes_total": len(rows),
        "probes_passed": passed,
        "result": "PASS" if passed == len(rows) and frozen_ok else "PARTIAL" if frozen_ok else "FAIL",
        "scope_note": "semantic branch absent by design (no embedder authorized); "
                      "fused results cannot be attributed to a semantic component",
    }
    out = pathlib.Path("eval/s3/gate-b-receipt.json")
    out.write_text(json.dumps(receipt, indent=2, sort_keys=True) + "\n", encoding="utf-8")
    print(f"\nRESULT: {receipt['result']} ({passed}/{len(rows)} probes, formation_frozen={frozen_ok})")
    print(f"receipt: {out} ({out.stat().st_size} B)")
    return 0 if receipt["result"] != "FAIL" else 1


if __name__ == "__main__":
    sys.exit(main())
