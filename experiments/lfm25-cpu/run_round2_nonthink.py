"""Round 2b: Qwen3.5 small models with thinking disabled.

Separate, clearly labelled condition. The default-mode runs are preserved in
round2/results.json; these runs send
`chat_template_kwargs: {"enable_thinking": false}` and are stored separately.
Same prompt, sampling, flags, and scoring rubric.
"""

from __future__ import annotations

import json
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parent
sys.path.insert(0, str(ROOT))
import run_probe as rp  # noqa: E402

rp.MODELS_DIR = ROOT / "round2" / "models"
rp.RAW = ROOT / "round2" / "raw-nonthink"
RESULTS = ROOT / "round2" / "nonthink_results.json"

MODELS = [
    ("Qwen3.5-0.8B", "Qwen3.5-0.8B-Q4_K_M.gguf"),
    ("Qwen3.5-2B", "Qwen3.5-2B-Q4_K_M.gguf"),
]
EXTRA = {"chat_template_kwargs": {"enable_thinking": False}}


def main() -> int:
    out: list[dict] = []
    for label, filename in MODELS:
        for draw in (1, 2, 3):
            print(f"=== {label} non-thinking draw {draw} ===", flush=True)
            try:
                rec = rp.run_one(label, filename, draw=draw, extra_body=EXTRA)
            except Exception as exc:
                rec = {"model": label, "draw": draw, "file": filename,
                       "error": repr(exc)}
                print("FAILED:", exc, file=sys.stderr, flush=True)
            rec["condition"] = "non-thinking (enable_thinking=false)"
            rec["round"] = "2b"
            out.append(rec)
            RESULTS.write_text(json.dumps(out, indent=2), encoding="utf-8")
            print(json.dumps({k: v for k, v in rec.items()
                              if k not in ("content", "reasoning_content",
                                           "prompt", "command", "sampling")},
                             indent=2), flush=True)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
