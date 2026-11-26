"""Round 2: Qwen2.5/3.5 small, FunctionGemma, and other small CPU models.

Same prompt, same rubric, same llama-server flags, same measurement code as
round 1 (imported from run_probe). Three default-settings draws per model.
CPU only (-ngl 0). Outputs live under round2/ and never overwrite round 1.
"""

from __future__ import annotations

import json
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parent
sys.path.insert(0, str(ROOT))
import run_probe as rp  # noqa: E402

rp.MODELS_DIR = ROOT / "round2" / "models"
rp.RAW = ROOT / "round2" / "raw"
RESULTS = ROOT / "round2" / "results.json"

MODELS = [
    ("Qwen2.5-0.5B-Instruct", "qwen2.5-0.5b-instruct-q4_k_m.gguf"),
    ("Qwen2.5-1.5B-Instruct", "qwen2.5-1.5b-instruct-q4_k_m.gguf"),
    ("Qwen3.5-0.8B", "Qwen3.5-0.8B-Q4_K_M.gguf"),
    ("Qwen3.5-2B", "Qwen3.5-2B-Q4_K_M.gguf"),
    ("FunctionGemma-270M-it", "functiongemma-270m-it-Q4_K_M.gguf"),
    ("Gemma-3-270M-it", "gemma-3-270m-it-Q4_K_M.gguf"),
    ("SmolLM2-1.7B-Instruct", "smollm2-1.7b-instruct-q4_k_m.gguf"),
    ("Granite-4.0-350M", "granite-4.0-350m-Q4_K_M.gguf"),
]


def main() -> int:
    missing = [f for _, f in MODELS if not (rp.MODELS_DIR / f).exists()]
    if missing:
        print("missing model files:", missing, file=sys.stderr)
        return 2
    out: list[dict] = []
    for label, filename in MODELS:
        for draw in (1, 2, 3):
            print(f"=== {label} draw {draw} ===", flush=True)
            try:
                rec = rp.run_one(label, filename, draw=draw)
            except Exception as exc:
                rec = {"model": label, "draw": draw, "file": filename,
                       "error": repr(exc)}
                print("FAILED:", exc, file=sys.stderr, flush=True)
            rec["round"] = 2
            out.append(rec)
            RESULTS.parent.mkdir(parents=True, exist_ok=True)
            RESULTS.write_text(json.dumps(out, indent=2), encoding="utf-8")
            print(json.dumps({k: v for k, v in rec.items()
                              if k not in ("content", "reasoning_content",
                                           "prompt", "command", "sampling")},
                             indent=2), flush=True)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
