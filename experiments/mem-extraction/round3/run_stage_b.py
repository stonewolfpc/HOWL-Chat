"""Stage B: bounded, one-factor-at-a-time tuning around a Stage A starting point.

variants.json: [{"label":..., "base": <stage A label>, "sampling": {...},
                 "max_tokens":..., "budget":...}, ...]
Each variant inherits the base model, mode, system prompt and extra_body unless
overridden. This is a search, not a sweep of the whole parameter space.
"""
import json, sys
from pathlib import Path
sys.path.insert(0, str(Path(__file__).resolve().parent))
import harness
from configs import by_label

VAR = json.loads((Path(__file__).resolve().parent / "variants.json").read_text(encoding="utf-8"))
DRAWS = int(sys.argv[1]) if len(sys.argv) > 1 else 3
specs = []
for v in VAR:
    s = by_label(v["base"])
    s.label = v["label"]
    s.preset = v.get("preset", s.preset)
    s.sampling = {**s.sampling, **v.get("sampling", {})}
    if "max_tokens" in v:
        s.max_tokens = v["max_tokens"]
    if "budget" in v:
        s.reasoning_budget = v["budget"]
    s.notes = v.get("notes", "")
    specs.append(s)
harness.run_triage(specs, DRAWS, "stageB")
