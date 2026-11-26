import sys
from pathlib import Path
sys.path.insert(0, str(Path(__file__).resolve().parent))
import harness
from configs import tinyllama
harness.run_triage([tinyllama()], 3, "tinyllama")
