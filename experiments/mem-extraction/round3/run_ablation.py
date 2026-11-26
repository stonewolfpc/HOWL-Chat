import json, sys
from pathlib import Path
sys.path.insert(0, str(Path(__file__).resolve().parent))
import harness
from configs_ablation import ablation, CORPUS_SPECS

specs = ablation()
print("== triage, 3 draws ==", flush=True)
harness.run_triage(specs, 3, "ablation_triage")

print("== corpus, 2 draws (non-thinking configs) ==", flush=True)
for s in specs:
    if s.label in CORPUS_SPECS:
        harness.run_corpus([s], 2, "ablation_corpus",
                           Path(__file__).resolve().parent / "corpus.json")
print("DONE", flush=True)
