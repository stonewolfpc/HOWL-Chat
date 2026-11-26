"""Embedder resource envelope — baseline capture and degraded-mode guarantee.

Runs with NO embedding model: records this machine's baseline, then proves the
second guarantee from the work order — memory still works when the embedder is
absent, and reports itself degraded rather than faking semantics.

    PYTHONPATH=. python eval/s3/resources.py
"""
from __future__ import annotations

import json
import pathlib
import re
import subprocess
import sys
import time

from eval.s3.selfcheck import build_conversation, toks, OWNER
from eval.s3.evidence import EvidenceLog
from eval.s3.records import RecordStore
from eval.s3.formation import Formatter
from eval.s3.context import Window
from eval.s3.hybrid import hybrid
from eval.s3.vectors import SEMANTIC_STATUS

MODEL = pathlib.Path("eval/s3/models/nomic-embed-text-v1.5.Q8_0.gguf")


def gpu() -> dict:
    out = subprocess.run(
        ["nvidia-smi", "--query-gpu=name,memory.total,memory.used,memory.free,driver_version,compute_cap",
         "--format=csv,noheader"], capture_output=True, text=True)
    if out.returncode != 0:
        return {"available": False}
    parts = [p.strip() for p in out.stdout.strip().split(",")]
    return {"available": True, "name": parts[0], "total_mib": int(re.sub(r"\D", "", parts[1])),
            "used_mib": int(re.sub(r"\D", "", parts[2])), "free_mib": int(re.sub(r"\D", "", parts[3])),
            "driver": parts[4], "compute_cap": parts[5] if len(parts) > 5 else None}


def ram() -> dict:
    import ctypes

    class M(ctypes.Structure):
        _fields_ = [("dwLength", ctypes.c_ulong), ("dwMemoryLoad", ctypes.c_ulong),
                    ("ullTotalPhys", ctypes.c_ulonglong), ("ullAvailPhys", ctypes.c_ulonglong),
                    ("ullTotalPageFile", ctypes.c_ulonglong), ("ullAvailPageFile", ctypes.c_ulonglong),
                    ("ullTotalVirtual", ctypes.c_ulonglong), ("ullAvailVirtual", ctypes.c_ulonglong),
                    ("ullAvailExtendedVirtual", ctypes.c_ulonglong)]

    m = M()
    m.dwLength = ctypes.sizeof(M)
    ctypes.windll.kernel32.GlobalMemoryStatusEx(ctypes.byref(m))
    return {"total_gib": round(m.ullTotalPhys / 2**30, 1),
            "available_gib": round(m.ullAvailPhys / 2**30, 1),
            "cores": __import__("os").cpu_count()}


def build_state():
    log, store, win = EvidenceLog(), RecordStore(), Window(capacity=420, reserve=80)
    former = Formatter(store, log)
    for i, (speaker, text) in enumerate(build_conversation(), start=1):
        turn = log.append(speaker, text, toks(text))
        roll = win.add(turn)
        if roll or i % 8 == 0:
            former.run_pending(log, OWNER)
    return log, store


def degraded_mode() -> dict:
    """Guarantee 2: lexical retrieval keeps working with no embedder present."""
    log, store = build_state()
    probes = [
        ("P1", "Who was your trainer, again?", "char-001", False, "trained_by"),
        ("P2", "What did you think of the weather?", "char-001", True, None),
        ("P3b", "What is the vault code?", "char-001", True, None),
        ("P4a", "Where is the ledger kept now?", "char-001", False, "green"),
    ]
    rows = []
    for pid, q, owner, expect_abstain, expect in probes:
        r = hybrid(q, store, owner)
        hay = [s["value"].lower() + " " + store.get(s["record_id"]).predicate.lower()
               for s in r.selected]
        ok = r.abstained if expect_abstain else (not r.abstained and any(expect in h for h in hay))
        rows.append({"id": pid, "abstained": r.abstained, "pass": ok,
                     "branches_available": r.branches_available,
                     "semantic": r.branches_unavailable})
    return {"semantic_available": SEMANTIC_STATUS["available"],
            "probes": rows, "passed": sum(1 for r in rows if r["pass"]), "total": len(rows)}


def envelope() -> dict:
    """Pre-registered budgets. The model is measured against these, not vice versa."""
    return {
        "vram_overhead_mib_max": 384,
        "vram_overhead_target": 0,
        "cold_load_seconds_max": 3.0,
        "unload_residual_mib_max": 64,
        "query_latency_p50_ms_max": 40,
        "query_latency_p95_ms_max": 120,
        "index_build_100_records_seconds_max": 2.0,
        "cpu_rss_gib_max": 1.0,
        "preferred_execution": "CPU/RAM; GPU only if CPU misses the latency budget",
        "matrix": ["chat model absent/present", "CPU/GPU", "cold/warm"],
    }


def main() -> int:
    g, r = gpu(), ram()
    degraded = degraded_mode()
    model_present = MODEL.exists()
    print(f"gpu: {g.get('name')} {g.get('total_mib')} MiB, used {g.get('used_mib')} MiB, "
          f"free {g.get('free_mib')} MiB, cc {g.get('compute_cap')}")
    print(f"ram: {r['total_gib']} GiB total, {r['available_gib']} GiB available, {r['cores']} cores")
    print(f"embedding model present: {model_present} ({MODEL})")
    print(f"degraded mode: {degraded['passed']}/{degraded['total']} probes pass with semantic "
          f"unavailable={not degraded['semantic_available']}")
    for p in degraded["probes"]:
        print(f"  {p['id']:<4} pass={p['pass']} abstained={p['abstained']}")
    out = pathlib.Path("eval/s3/resource-baseline.json")
    out.write_text(json.dumps({
        "captured_at": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()),
        "gpu": g, "ram": r,
        "embedding_model_present": model_present,
        "envelope": envelope(),
        "degraded_mode": degraded,
        "note": "baseline only; no embedding model loaded, nothing downloaded",
    }, indent=2, sort_keys=True) + "\n", encoding="utf-8")
    print(f"wrote {out} ({out.stat().st_size} B)")
    return 0 if degraded["passed"] == degraded["total"] else 1


if __name__ == "__main__":
    sys.exit(main())
