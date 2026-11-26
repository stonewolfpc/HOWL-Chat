"""Run the 18-case corpus per an explicit plan, appending into results/corpus.json.

plan.json: [{"label":..., "base":..., "draws":2, "budget":512|null,
             "max_tokens":..., "sampling":{...}}]
Configs run sequentially; one server per config; slot erased before every job.
"""
import json, sys, time
from pathlib import Path
sys.path.insert(0, str(Path(__file__).resolve().parent))
import harness
from configs import by_label, lfm3

PLAN = json.loads((Path(__file__).resolve().parent / "plan_tinyllama.json").read_text(encoding="utf-8"))
CORPUS = json.loads((Path(__file__).resolve().parent / "corpus.json").read_text(encoding="utf-8"))["cases"]
OUT = harness.RESULTS / "corpus.json"
out = json.loads(OUT.read_text(encoding="utf-8")) if OUT.exists() else []
done = {(r["label"], r["case"], r["draw"]) for r in out if "error" not in r}

def build(p):
    lab = p["base"]
    if lab.startswith("lfm3"):
        s = lfm3("think" in lab)
    else:
        s = by_label(lab)
    s.label = p["label"]
    s.max_tokens = p.get("max_tokens", 768)
    if p.get("budget") is not None:
        s.reasoning_budget = p["budget"]
    if p.get("sampling"):
        s.sampling = {**s.sampling, **p["sampling"]}
    s.notes = p.get("notes", "")
    return s

for p in PLAN:
    s = build(p)
    if not s.path.exists():
        print("MISSING", s.path, flush=True); continue
    todo = [c for c in CORPUS if any((p["label"], c["id"], d) not in done for d in range(1, p["draws"] + 1))]
    if not todo:
        print("SKIP complete", p["label"], flush=True); continue
    print(f"=== corpus {p['label']} cases={len(todo)} draws={p['draws']} budget={s.reasoning_budget} ===", flush=True)
    try:
        with harness.Server(s, harness.RAW / f"corpus-{p['label']}" / "server.log") as srv:
            for case in todo:
                for draw in range(1, p["draws"] + 1):
                    if (p["label"], case["id"], draw) in done:
                        continue
                    try:
                        srv.erase_slot()
                    except Exception:
                        pass
                    user = harness.EXTRACTION_INSTRUCTION + "\n\nConversation:\n" + case["conversation"]
                    r = harness.chat(s, user, seed=4100 + draw)
                    rec = {"run_code": harness._q(), "label": p["label"], "base": p["base"],
                           "case": case["id"], "category": case["category"],
                           "length": case["length"], "draw": draw, "seed": 4100 + draw,
                           "sampling": s.sampling, "max_tokens": s.max_tokens,
                           "reasoning_budget": s.reasoning_budget, "ngl": s.ngl,
                           "file": s.file, "kind": "corpus",
                           "content": r["content"], "reasoning_content": r["reasoning"],
                           "finish_reason": r["finish_reason"],
                           "prompt_tokens": r["usage"].get("prompt_tokens"),
                           "completion_tokens": r["usage"].get("completion_tokens"),
                           "ttft_s": r["ttft_s"], "thinking_s": r["thinking_s"],
                           "answer_s": r["answer_s"], "wall_s": r["wall_s"],
                           "peak_rss_mb": round(srv.peak["rss"] / 1e6, 1),
                           "peak_vram_mb": srv.peak["vram"],
                           "load_s": round(srv.load_s, 2)}
                    out.append(rec)
                    OUT.write_text(json.dumps(out, indent=2), encoding="utf-8")
                    print(f"  {case['id']} d{draw} {r['wall_s']:6.1f}s fin={r['finish_reason']} "
                          f"tok={rec['completion_tokens']}", flush=True)
    except Exception as exc:
        print("FAILED", p["label"], repr(exc), flush=True)
    time.sleep(1)
print("TOTAL", len(out), flush=True)
