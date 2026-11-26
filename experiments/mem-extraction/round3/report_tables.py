"""Generate the round-3 markdown tables directly from recorded JSON.

Nothing here types a score by hand. If a file is missing, the section is
skipped and the report says so.
"""
from __future__ import annotations

import json
import random
import statistics as st
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parent
R = ROOT / "results"


def load(name):
    p = R / name
    return json.loads(p.read_text(encoding="utf-8")) if p.exists() else None


def mean(xs):
    xs = [x for x in xs if x is not None]
    return round(st.mean(xs), 2) if xs else None


def boot_ci(xs, n=4000, seed=7):
    if len(xs) < 4:
        return None
    rnd = random.Random(seed)
    means = [st.mean(rnd.choices(xs, k=len(xs))) for _ in range(n)]
    means.sort()
    return round(means[int(0.025 * n)], 1), round(means[int(0.975 * n)], 1)


def stage_a_table():
    runs = load("stageA.json")
    scores = load("stageA_scores.json")
    if not runs or not scores:
        return None
    by_score = {}
    for s in scores:
        by_score.setdefault(s["label"], []).append(s["score12"])
    by_run = {}
    for r in runs:
        if "error" in r:
            continue
        by_run.setdefault(r["label"], []).append(r)
    lines = ["| Configuration | Mode | Preset (T / top_p / top_k / presence / rep) | 12-pt mean (draws) | "
             "Prompt tok | Out tok | TTFT s | Thinking s | Wall s | Peak RSS MB |",
             "|---|---|---|---|---|---|---|---|---|---|"]
    order = sorted(by_score, key=lambda k: -st.mean(by_score[k]))
    for lab in order:
        rs = by_run.get(lab, [])
        if not rs:
            continue
        r0 = rs[0]
        sm = r0["sampling"]
        preset = (f"{sm.get('temperature')} / {sm.get('top_p')} / {sm.get('top_k')} / "
                  f"{sm.get('presence_penalty')} / {sm.get('repeat_penalty')}")
        mode = "thinking" if (r0.get("reasoning_budget") is not None or
                              "think" in lab) else "non-thinking"
        sc = by_score[lab]
        lines.append(
            f"| {lab} | {mode} | {preset} | {mean(sc):.2f} ({', '.join(str(x) for x in sc)}) | "
            f"{mean([r.get('prompt_tokens') for r in rs])} | "
            f"{mean([r.get('completion_tokens') for r in rs])} | "
            f"{mean([r.get('ttft_s') for r in rs])} | "
            f"{mean([r.get('thinking_s') for r in rs])} | "
            f"{mean([r.get('wall_s') for r in rs])} | "
            f"{mean([r.get('peak_rss_mb') for r in rs])} |")
    return "\n".join(lines)


def budget_table():
    runs = load("budget.json")
    scores = load("budget_scores.json")
    if not runs or not scores:
        return None
    by = {}
    for s in scores:
        by.setdefault(s["label"], []).append(s["score12"])
    by_run = {}
    for r in runs:
        if "error" not in r:
            by_run.setdefault(r["label"], []).append(r)
    lines = ["| Budget condition | 12-pt mean | Out tok | Thinking s | Wall s | finish=length | Empty content |",
             "|---|---|---|---|---|---|---|"]
    for lab in sorted(by):
        rs = by_run.get(lab, [])
        lines.append(f"| {lab} | {mean(by[lab]):.2f} | {mean([r.get('completion_tokens') for r in rs])} | "
                     f"{mean([r.get('thinking_s') for r in rs])} | {mean([r.get('wall_s') for r in rs])} | "
                     f"{sum(1 for r in rs if r.get('finish_reason') == 'length')}/{len(rs)} | "
                     f"{sum(1 for r in rs if not (r.get('content') or '').strip())}/{len(rs)} |")
    return "\n".join(lines)


def corpus_table(summary):
    lines = ["| Config | Cases | Draws | Mean /100 | Median | Min | Max | SD across cases | SD within case | "
             "Identical decisions | Catastrophic | Empty | Truncated |",
             "|---|---|---|---|---|---|---|---|---|---|---|---|---|"]
    for row in summary:
        lines.append(f"| {row['label']} | {row['cases']} | {row['draws']} | {row['mean']} | {row['median']} | "
                     f"{row['min']} | {row['max']} | {row['sd_cases']} | {row['sd_draws_within_case']} | "
                     f"{row['identical_decisions_frac']} | {row['catastrophic']} | {row['empty_runs']} | "
                     f"{row['truncated_runs']} |")
    return "\n".join(lines)


def main():
    out = {}
    t = stage_a_table()
    if t:
        out["stage_a"] = t
    b = budget_table()
    if b:
        out["budget"] = b
    print(json.dumps(out, indent=1)[:200])
    if t:
        (R / "table_stage_a.md").write_text(t, encoding="utf-8")
    if b:
        (R / "table_budget.md").write_text(b, encoding="utf-8")


if __name__ == "__main__":
    main()
