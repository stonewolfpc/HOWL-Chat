"""Stage A2: official sampler preset for the remaining models.

For thinking models the *unbounded* official condition was attempted first and
failed to terminate (Qwen3.5-2B: 3072/3072 tokens, empty content). Those runs
are kept in stageA.json. Here the official sampler chain is unchanged but a
runtime reasoning budget of 1024 is imposed, which is a harness-side cap, not a
model effort level. Non-thinking models get no budget.
"""
import sys
from pathlib import Path
sys.path.insert(0, str(Path(__file__).resolve().parent))
import harness
from configs import stage_a, lfm3

BUDGETED = {"qwen35-2b-think", "qwen35-4b-think", "qwen35-9b-think",
            "qwen3-17b-think", "smollm3-3b-think", "lfm120-think-official",
            "lfm3-17b-think"}
WANT = sys.argv[1].split(",") if len(sys.argv) > 1 else []
TAG = sys.argv[2] if len(sys.argv) > 2 else "stageA2"
DRAWS = int(sys.argv[3]) if len(sys.argv) > 3 else 3

pool = {s.label: s for s in stage_a()}
pool["lfm3-17b-think"] = lfm3(True)
pool["lfm3-17b-nonthink"] = lfm3(False)
specs = []
for lab in WANT:
    s = pool[lab]
    if lab in BUDGETED:
        s.reasoning_budget = 1024
        s.max_tokens = max(s.max_tokens, 2048)
        s.notes = (s.notes + " | runtime budget 1024 (unbounded did not terminate)").strip(" |")
        s.label = f"{lab}-rb1024"
    specs.append(s)
print([s.label for s in specs], flush=True)
harness.run_triage(specs, DRAWS, TAG)
