import json, sys
from pathlib import Path
sys.path.insert(0, str(Path(__file__).resolve().parent))
import harness
from configs import by_label
VAR = json.loads((Path(__file__).resolve().parent / "variants2.json").read_text(encoding="utf-8"))
specs = []
for v in VAR:
    s = by_label(v["base"]); s.label = v["label"]; s.preset = v["preset"]
    s.sampling = {**s.sampling, **v["sampling"]}
    s.max_tokens = v.get("max_tokens", s.max_tokens)
    s.ctx = v.get("ctx", s.ctx)
    specs.append(s)
harness.run_triage(specs, 3, "stageC")
