"""CLI for instrument v2.

Writes only to paths you name explicitly, refuses to write inside `results/` or
`raw/`, and refuses to overwrite an existing file without --force. Historical
round-3 artifacts are inputs here and are never opened for writing.
"""

from __future__ import annotations

import argparse
import hashlib
import json
import sys
from datetime import datetime, timezone
from pathlib import Path

import reporting
import rubric_v2

VERSION = "MEMX-EVAL 2.0.0"
HERE = Path(__file__).resolve().parent
FORBIDDEN_OUT_DIRS = {"results", "raw", "models"}


def sha256(path: Path) -> str:
    h = hashlib.sha256()
    h.update(path.read_bytes())
    return h.hexdigest()


def guard_out(path: Path, force: bool) -> None:
    if path.resolve().parent.name in FORBIDDEN_OUT_DIRS:
        raise SystemExit(f"refusing to write into {path.parent}: historical "
                         "evidence directories are read-only for v2")
    if path.exists() and not force:
        raise SystemExit(f"refusing to overwrite {path} without --force")


def main(argv: list[str] | None = None) -> int:
    ap = argparse.ArgumentParser(description="MEMX-EVAL 2.0.0")
    ap.add_argument("--runs", required=True)
    ap.add_argument("--corpus", required=True)
    ap.add_argument("--out", required=True)
    ap.add_argument("--force", action="store_true")
    args = ap.parse_args(argv)

    runs_path, corpus_path = Path(args.runs), Path(args.corpus)
    out_path = Path(args.out)
    guard_out(out_path, args.force)

    runs = json.loads(runs_path.read_text(encoding="utf-8"))
    corpus = {c["id"]: c for c in
              json.loads(corpus_path.read_text(encoding="utf-8"))["cases"]}

    scored = []
    for r in runs:
        if r.get("kind") not in (None, "corpus"):
            continue
        case = corpus.get(r.get("case"))
        if case is None:
            scored.append({"case": r.get("case"), "run_code": r.get("run_code"),
                           "label": r.get("label"), "status": "UNKNOWN_CASE",
                           "semantic": None, "contract": None, "events": []})
            continue
        scored.append(rubric_v2.score_run(r, case))

    agg = rubric_v2.aggregate(scored)
    out = {
        "instrument": VERSION,
        "generated_utc": datetime.now(timezone.utc).isoformat(timespec="seconds"),
        "inference_run": False,
        "inputs": {"runs": str(runs_path), "runs_sha256": sha256(runs_path),
                   "corpus": str(corpus_path), "corpus_sha256": sha256(corpus_path)},
        "runs": scored,
        "summary": agg,
        "reporting_notes": {
            label: reporting.sample_qualifier(v["draws_per_case"])
            for label, v in agg.items()
        },
    }
    out_path.parent.mkdir(parents=True, exist_ok=True)
    out_path.write_text(json.dumps(out, indent=2), encoding="utf-8")

    print(f"{VERSION} -> {out_path}")
    for label, v in sorted(agg.items(), key=lambda kv: -(kv[1]["semantic_mean_eligible"] or 0)):
        print(f"  {label:34s} semantic={v['semantic_mean_eligible']} "
              f"contract={v['contract_mean']} n={v['draws_per_case']} "
              f"truncated={v['truncated_runs']} "
              f"[{reporting.sample_qualifier(v['draws_per_case'])}]")
    return 0


if __name__ == "__main__":
    sys.exit(main())
