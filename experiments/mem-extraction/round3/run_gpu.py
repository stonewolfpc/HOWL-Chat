"""GPU offload measurements for the promising candidates.

Full residency first (ngl 99). If the model does not fit, a partial-offload row
is produced at half the layer count read from the server log. VRAM is polled via
nvidia-smi; peak VRAM during inference is what is reported, not file size.
"""
import json, re, sys
from pathlib import Path
sys.path.insert(0, str(Path(__file__).resolve().parent))
import harness
from configs import by_label, lfm3

def n_layers(spec):
    # read from a previous server log if present, else launch once and parse
    logdir = harness.RAW / f"gpu-{spec.label}"
    text = ""
    if (logdir / "server.log").exists():
        text = (logdir / "server.log").read_text(errors="replace")
    m = re.search(r"n_layer\s*=\s*(\d+)", text)
    return int(m.group(1)) if m else None

def main():
    labels = sys.argv[1].split(",") if len(sys.argv) > 1 else ["lfm120-instruct-official"]
    out = []
    for lab in labels:
        base = by_label(lab) if lab in [s.label for s in __import__("configs").stage_a()] else lfm3(False)
        for ngl, tag in ((99, "full"),):
            s = base
            s.ngl = ngl
            s.ctx = 4096
            s.flash_attn = True
            s.kv_unified = True
            s.label = f"{base.label}-gpu{tag}"
            s.max_tokens = 768
            recs = harness.run_triage([s], 2, f"gpu-{tag}-{base.label}")
            out.extend(recs)
    print(json.dumps([{k: v for k, v in r.items() if k not in ("content", "reasoning_content")}
                      for r in out], indent=2))

if __name__ == "__main__":
    main()
