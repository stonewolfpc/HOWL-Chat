"""Stage A: publisher-recommended configuration, original six-item prompt."""
import json, sys
from pathlib import Path
sys.path.insert(0, str(Path(__file__).resolve().parent))
import harness
from configs import stage_a

DRAWS = int(sys.argv[2]) if len(sys.argv) > 2 else 3
TAG = sys.argv[1] if len(sys.argv) > 1 else "stageA"
specs = stage_a()
if TAG.endswith("-no9b"):
    specs = [s for s in specs if "9b" not in s.label]
if TAG.endswith("-9b"):
    specs = [s for s in specs if "9b" in s.label]
print("specs:", [s.label for s in specs], flush=True)
harness.run_triage(specs, DRAWS, TAG)
