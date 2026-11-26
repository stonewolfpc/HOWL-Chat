"""Blind corpus scoring, then unblind for aggregation."""
import json, random, statistics as st, collections
from pathlib import Path
from score_corpus import main as score_main, summarise, GROUNDS

R = Path(__file__).resolve().parent
runs = json.loads((R/"results/corpus.json").read_text(encoding="utf-8"))
# 1. score blind (only run_code/case/content are read)
rows = score_main(R/"results/corpus.json", R/"corpus.json", R/"results/corpus_scores.json")
# 2. code map, written after scoring
cmap = {r["run_code"]: {"label": r["label"], "case": r["case"], "draw": r["draw"]} for r in runs if "error" not in r}
(R/"results/code_map.json").write_text(json.dumps(cmap, indent=2), encoding="utf-8")
summary = summarise(rows, cmap)

def boot(xs, n=4000, seed=11):
    if len(xs) < 4: return None
    rnd = random.Random(seed)
    m = sorted(st.mean(rnd.choices(xs, k=len(xs))) for _ in range(n))
    return round(m[int(.025*n)],1), round(m[int(.975*n)],1)

table = ["| Config | Cases | Draws | Mean /100 | 95% CI | Median | Min | Max | SD cases | SD draws | Identical decisions | Catastrophic runs | Empty | Truncated |",
         "|---|---|---|---|---|---|---|---|---|---|---|---|---|---|"]
for s in summary:
    per = collections.defaultdict(list)
    for r in rows:
        if cmap[r["run_code"]]["label"] == s["label"]:
            per[r["case"]].append(r["score"])
    means = [st.mean(v) for v in per.values()]
    ci = boot(means)
    table.append(f"| {s['label']} | {s['cases']} | {s['draws']} | **{s['mean']}** | {ci} | {s['median']} | "
                 f"{s['min']} | {s['max']} | {s['sd_cases']} | {s['sd_draws_within_case']} | "
                 f"{s['identical_decisions_frac']} | {s['catastrophic_runs']} | {s['empty_runs']} | {s['truncated_runs']} |")
(R/"results/corpus_table.md").write_text("\n".join(table), encoding="utf-8")
print("\n".join(table))

# per-category for the best config
best = summary[0]["label"]
print("\nPer-case means for", best)
per = collections.defaultdict(list)
for r in rows:
    if cmap[r["run_code"]]["label"] == best:
        per[r["case"]].append(r["score"])
for k in sorted(per):
    print(f"  {k} {round(st.mean(per[k]),1)} {[round(x) for x in per[k]]}")
