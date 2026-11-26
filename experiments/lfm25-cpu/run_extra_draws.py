"""Additional default-settings draws for the LFM2.5 probe.

Draw 1 is preserved exactly as first recorded in results/results.json and
raw/<model>/ (before this script existed). This appends draws 2 and 3 per
model, same prompt, same sampling, same runtime. All draws are reported;
none replaces another.
"""

from __future__ import annotations

import json
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))
import run_probe as rp  # noqa: E402

RESULTS = rp.RESULTS / "results.json"


def main() -> int:
    existing = json.loads(RESULTS.read_text(encoding="utf-8"))
    for rec in existing:
        rec.setdefault("draw", 1)
    if any(r.get("draw", 1) != 1 for r in existing):
        print("results.json already contains extra draws; refusing to duplicate",
              file=sys.stderr)
        return 2
    for draw in (2, 3):
        for label, filename in rp.MODELS:
            print(f"=== {label} draw {draw} ===", flush=True)
            try:
                rec = rp.run_one(label, filename, draw=draw)
            except Exception as exc:
                rec = {"model": label, "draw": draw, "file": filename,
                       "error": repr(exc)}
                print("FAILED:", exc, file=sys.stderr, flush=True)
            print(json.dumps({k: v for k, v in rec.items()
                              if k not in ("content", "prompt")}, indent=2),
                  flush=True)
            existing.append(rec)
            RESULTS.write_text(json.dumps(existing, indent=2), encoding="utf-8")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
