"""Merge every triage run and print the unified quality/performance table.

Reads raw run files, applies the legacy 12-point scorer, and writes:
  results/all_triage.json     merged runs
  results/triage_table.md     one row per configuration
  results/corpus_summary.md   corpus rankings (delegates to score_corpus)
"""
from __future__ import annotations

import json
import statistics as st
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parent
sys.path.insert(0, str(ROOT))
from score_triage import score  # noqa: E402

R = ROOT / "results"
FILES = ["stageA.json", "stageA2a.json", "stageA2b.json", "stageA2b2.json",
         "budget.json", "stageB.json"]


def main():
    runs = []
    for f in FILES:
        p = R / f
        if p.exists():
            runs.extend([r for r in json.loads(p.read_text(encoding="utf-8"))
                         if r.get("kind") == "triage"])
    (R / "all_triage.json").write_text(json.dumps(runs, indent=2), encoding="utf-8")

    rows = [score(r) for r in runs if "error" not in r]
    by = {}
    for s in rows:
        by.setdefault(s["label"], []).append(s)

    lines = ["| Configuration | Admissible draws | 12-pt mean | Draws | Failures | Out tok | Wall s | Peak RSS MB |",
             "|---|---|---|---|---|---|---|---|"]
    ranked = []
    for lab, ss in by.items():
        ok = [s for s in ss if s["admissible"]]
        bad = [s for s in ss if not s["admissible"]]
        if not ok:
            ranked.append((lab, -1, ss, 0, len(bad)))
            continue
        ranked.append((lab, st.mean(s["score12"] for s in ok), ss, len(ok), len(bad)))
    ranked.sort(key=lambda x: -x[1])
    for lab, m, ss, n, nb in ranked:
        r0 = runs[[r["label"] for r in runs].index(lab)] if lab in [r["label"] for r in runs] else {}
        walls = [s["wall_s"] for s in ss if s.get("wall_s")]
        toks = [s["completion_tokens"] for s in ss if s.get("completion_tokens")]
        rss = [r.get("peak_rss_mb") for r in runs if r.get("label") == lab and r.get("peak_rss_mb")]
        lines.append(f"| {lab} | {n} | {m if m >= 0 else 'n/a'} | "
                     f"{', '.join(str(s['score12']) for s in ss)} | {nb} | "
                     f"{round(st.mean(toks)) if toks else ''} | "
                     f"{round(st.mean(walls), 1) if walls else ''} | "
                     f"{round(st.mean(rss)) if rss else ''} |")
    (R / "triage_table.md").write_text("\n".join(lines), encoding="utf-8")
    print("\n".join(lines))


if __name__ == "__main__":
    main()
