"""Build qualitative documents for the abliterated/uncensored variants only.

Writes to qualitative/ablation/, leaving the baseline qualitative track
untouched. Same structure as the baseline: exact prompt, every draw verbatim,
reference material, one model at a time.
"""
import json
from pathlib import Path

HERE = Path(__file__).resolve().parent
OUT = HERE / "qualitative" / "ablation"
RUNS = json.loads((HERE / "results" / "ablation_corpus_all.json").read_text(encoding="utf-8"))
CASES = json.loads((HERE / "corpus.json").read_text(encoding="utf-8"))["cases"]

NAMES = {
    "abl-4b-nonthink": ("Qwen3.5-4B HauhauCS Uncensored", "non-thinking"),
    "abl-2b-nonthink": ("Qwen3.5-2B HauhauCS Uncensored", "non-thinking"),
    "abl-lfm120": ("LFM2.5-1.2B Huihui abliterated", "non-thinking"),
}
ORDER = list(NAMES)

BY = {}
for r in RUNS:
    BY.setdefault(r["label"], []).append(r)
for v in BY.values():
    v.sort(key=lambda r: (r["case"], r["draw"]))

INSTR = (
    "You are the memory extractor for a long-running roleplay. Read the conversation and "
    "decide what belongs in long-term memory.\nRules:\n"
    "- KEEP: durable facts, preferences, relationships, world rules, or corrections that will "
    "matter in later sessions.\n"
    "- REJECT: transient chatter, scheduling, mood-of-the-moment, or restated facts with no new "
    "content.\n"
    "- HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.\n"
    "Output exactly one line per candidate in the form:\n"
    "ID | DECISION | MEMORY | REASON\n"
    "Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE\n"
    "Do not invent details that are not in the conversation."
)


def fence(t, lang="text"):
    return f"```{lang}\n{('' if t is None else t).rstrip()}\n```"


def reference(c):
    lines = []
    if c.get("required"):
        lines.append("| id | expected | what it is | key phrases |")
        lines.append("|---|---|---|---|")
        for it in c["required"]:
            keys = "; ".join(" / ".join(k) for k in it.get("keys", []))
            lines.append(f"| {it['id']} | **{it['decision']}** | {it.get('points','')} | {keys} |")
    else:
        lines.append("**No required items.** Correct answer is `NONE`.")
    if c.get("forbidden"):
        lines += ["", "Must not appear as stored canon: " +
                  ", ".join(f"`{x}`" for x in c["forbidden"])]
    if c.get("notes"):
        lines += ["", c["notes"]]
    return "\n".join(lines)


def cfg(r):
    s = r.get("sampling") or {}
    bits = [f"seed {r.get('seed')}", f"max_tokens {r.get('max_tokens')}",
            "sampling " + ", ".join(f"{k}={v}" for k, v in s.items()),
            f"finish `{r.get('finish_reason')}`", f"{r.get('completion_tokens')} tok",
            f"{r.get('wall_s'):.1f} s", f"run `{r.get('run_code')}`"]
    return " · ".join(bits)


def main():
    (OUT / "archive").mkdir(parents=True, exist_ok=True)
    (OUT / "models").mkdir(parents=True, exist_ok=True)
    reading = ["# Abliterated variants — case-by-case reading copy", "",
               "Generated from `results/ablation_corpus.json`, same corpus and prompts as the",
               "baseline track. Draw 1 shown; all draws in `archive/`.", ""]

    for c in CASES:
        cid = c["id"]
        prompt = INSTR + "\n\nConversation:\n" + c["conversation"]
        doc = [f"# {cid} — {c['category']} ({c['length']})", "",
               "## Source material", "", fence(c["conversation"]), "",
               "## Reference material", "", reference(c), "",
               "## Exact prompt", "", fence(prompt), "",
               "## Responses, all draws", ""]
        reading += [f"## {cid} — {c['category']}", "", fence(c["conversation"]), "",
                    "**Reference**", "", reference(c), ""]
        for label in ORDER:
            rows = [r for r in BY.get(label, []) if r["case"] == cid]
            if not rows:
                continue
            name, mode = NAMES[label]
            doc += [f"### {name} — {mode}", ""]
            for r in rows:
                doc += [f"**draw {r['draw']}** — {cfg(r)}", ""]
                if r.get("reasoning_content"):
                    doc += ["<details><summary>reasoning</summary>", "",
                            fence(r["reasoning_content"]), "", "</details>", ""]
                doc += [fence(r.get("content") or "(empty content)"), ""]
            r0 = rows[0]
            reading += [f"**{name}** — {r0.get('completion_tokens')} tok, {r0.get('finish_reason')}",
                        "", fence(r0.get("content") or "(empty content)"), ""]
        (OUT / "archive" / f"case-{cid}.md").write_text("\n".join(doc), encoding="utf-8")

    (OUT / "READING.md").write_text("\n".join(reading), encoding="utf-8")

    for label in ORDER:
        rows = BY.get(label, [])
        if not rows:
            continue
        name, mode = NAMES[label]
        doc = [f"# {name} — {mode}", "",
               f"{len(rows)} runs over {len({r['case'] for r in rows})} cases.", ""]
        for c in CASES:
            cr = [r for r in rows if r["case"] == c["id"]]
            if not cr:
                continue
            doc += [f"## {c['id']} — {c['category']}", ""]
            for r in cr:
                doc += [f"**draw {r['draw']}** — {cfg(r)}", ""]
                if r.get("reasoning_content"):
                    doc += ["<details><summary>reasoning</summary>", "",
                            fence(r["reasoning_content"]), "", "</details>", ""]
                doc += [fence(r.get("content") or "(empty content)"), ""]
        (OUT / "models" / f"{label}.md").write_text("\n".join(doc), encoding="utf-8")
    print("wrote", OUT, "cases", len(CASES))


if __name__ == "__main__":
    main()
