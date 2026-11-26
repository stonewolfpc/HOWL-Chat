"""Context-reset / job-isolation test.

Distinguishes (a) per-request history clearing from (b) actual inference-state
reset. A benign job A is run, then an adversarial job B, then A again in the
same server; a second A is also run on a fresh server. If A2 differs from A1,
or if prefix-cache counters show reuse after an erase, that is recorded.
"""
import json, sys
from pathlib import Path
sys.path.insert(0, str(Path(__file__).resolve().parent))
import harness
from configs import stage_a

CORPUS = json.loads((Path(__file__).resolve().parent / "corpus.json").read_text(encoding="utf-8"))["cases"]
C = {c["id"]: c for c in CORPUS}
INSTR = harness.EXTRACTION_INSTRUCTION

def run_job(spec, case_id, seed):
    return harness.chat(spec, INSTR + "\n\nConversation:\n" + C[case_id]["conversation"], seed=seed)

def main(label):
    spec = next(s for s in stage_a() if s.label == label) if label in [s.label for s in stage_a()] else None
    if spec is None:
        from configs import lfm3
        spec = lfm3(False)
    spec.max_tokens = 768
    out = {"label": spec.label, "jobs": []}
    erase_ok = True
    with harness.Server(spec, harness.RAW / "reset" / "server.log") as srv:
        for name, case_id in (("A1", "C03"), ("B_adversarial", "C17"), ("A2", "C03")):
            try:
                t_erase = srv.erase_slot()
            except Exception as exc:
                # this build answers 501 here; record it and continue
                erase_ok = False
                out["erase_error"] = repr(exc)
                t_erase = float("nan")
            r = run_job(spec, case_id, 4101)
            tl = (r.get("timings") or {})
            out["jobs"].append({"name": name, "erase_ms": round(t_erase, 2),
                                "case": case_id, "content": r["content"],
                                "reasoning_chars": len(r["reasoning"]),
                                "cache_n": tl.get("cache_n"),
                                "prompt_n": tl.get("prompt_n"),
                                "ttft_s": r["ttft_s"], "wall_s": r["wall_s"],
                                "finish": r["finish_reason"]})
            print(name, "erase_ms", round(t_erase, 2), "cache_n", tl.get("cache_n"),
                  "prompt_n", tl.get("prompt_n"), "wall", r["wall_s"], flush=True)
    # fresh-server reference for the same benign job
    out["erase_supported"] = erase_ok
    with harness.Server(spec, harness.RAW / "reset" / "server-fresh.log") as srv2:
        try:
            srv2.erase_slot()
        except Exception:
            pass
        r = run_job(spec, "C03", 4101)
        out["jobs"].append({"name": "A_fresh", "erase_ms": 0.0, "case": "C03",
                            "content": r["content"], "reasoning_chars": len(r["reasoning"]),
                            "cache_n": (r.get("timings") or {}).get("cache_n"),
                            "prompt_n": (r.get("timings") or {}).get("prompt_n"),
                            "ttft_s": r["ttft_s"], "wall_s": r["wall_s"],
                            "finish": r["finish_reason"]})
    a1 = next(j for j in out["jobs"] if j["name"] == "A1")["content"].strip()
    a2 = next(j for j in out["jobs"] if j["name"] == "A2")["content"].strip()
    out["a1_equals_a2"] = a1 == a2
    out["a1_equals_fresh"] = a1 == next(j for j in out["jobs"] if j["name"] == "A_fresh")["content"].strip()
    out["cache_reused_after_erase"] = any((j.get("cache_n") or 0) > 0 for j in out["jobs"])
    harness.RESULTS.mkdir(parents=True, exist_ok=True)
    (harness.RESULTS / "reset.json").write_text(json.dumps(out, indent=2), encoding="utf-8")
    print(json.dumps({k: v for k, v in out.items() if k != "jobs"}, indent=1))

if __name__ == "__main__":
    main(sys.argv[1] if len(sys.argv) > 1 else "lfm120-instruct-official")
