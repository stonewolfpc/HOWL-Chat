"""Thinking-budget sweep (runtime-imposed via llama-server --reasoning-budget).

Budgets are experimental levels, not claims that the model implements an effort
control. LFM2.5-Thinking has no thinking control at all, so its budget rows are
explicitly labelled as forced truncation of an always-on trace.
"""
import sys
from pathlib import Path
sys.path.insert(0, str(Path(__file__).resolve().parent))
import harness
from configs import by_label, lfm3

BUDGETS = [int(x) for x in (sys.argv[2].split(",") if len(sys.argv) > 2 else ["512", "1024", "2048"])]
DRAWS = int(sys.argv[3]) if len(sys.argv) > 3 else 2
labels = (sys.argv[1].split(",") if len(sys.argv) > 1 else ["qwen35-4b-think"])
specs = []
for lab in labels:
    for b in BUDGETS:
        s = lfm3(True) if lab.startswith("lfm3") else by_label(lab)
        s.reasoning_budget = b
        s.label = f"{lab}-rb{b}"
        s.notes = (s.notes + " | budget forced by runtime, not a model effort level").strip(" |")
        specs.append(s)
harness.run_triage(specs, DRAWS, "budget")
