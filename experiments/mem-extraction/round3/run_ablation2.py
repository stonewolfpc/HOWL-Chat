import json, sys
from pathlib import Path
sys.path.insert(0, str(Path(__file__).resolve().parent))
import harness
from configs_ablation import ablation
from configs import by_label

specs = {s.label: s for s in ablation()}
for lb in ("abl-4b-nonthink", "abl-2b-nonthink"):
    print("== rerun", lb, flush=True)
    harness.run_corpus([specs[lb]], 2, f"ablation_corpus_{lb}",
                       Path(__file__).resolve().parent / "corpus.json")

print("== control: stock Qwen3.5-2B corpus ==", flush=True)
harness.run_corpus([by_label("qwen35-2b-nonthink")], 2, "control_corpus_2b",
                   Path(__file__).resolve().parent / "corpus.json")
print("DONE2", flush=True)
