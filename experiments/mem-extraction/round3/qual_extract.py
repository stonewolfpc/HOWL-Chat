"""Build the qualitative evaluation track from preserved corpus runs.

Layers:
  qualitative/archive/case-Cxx.md  full verbatim: exact prompt, every draw,
                                   every reasoning block, reference material
  qualitative/models/<label>.md    one model at a time, all 18 cases
  qualitative/READING.md           compact side-by-side, draw 1 per model

Nothing here scores anything. It only reassembles what the harness stored.
"""
import json
import textwrap
from pathlib import Path

HERE = Path(__file__).resolve().parent
OUT = HERE / "qualitative"
CASES = json.loads((HERE / "corpus.json").read_text(encoding="utf-8"))["cases"]
RUNS = json.loads((HERE / "results" / "corpus.json").read_text(encoding="utf-8"))

# Reading order: strongest first, then the small ones.
ORDER = [
    ("corpus-4b-nonthink", "Qwen3.5-4B", "non-thinking", "T0.7 / top_p0.8 / top_k20 / min_p0.0 / presence 1.5"),
    ("corpus-4b-think-512", "Qwen3.5-4B", "thinking, runtime budget 512", "T0.6 / top_p0.95 / top_k20 / min_p0.0 / presence 1.5"),
    ("corpus-lfm3-think-512", "LFM3-1.7B", "thinking, runtime budget 512", "T0.1 / top_p0.1 / min_p0.0 / repeat 1.1"),
    ("corpus-lfm120-instruct", "LFM2.5-1.2B-Instruct", "non-thinking, publisher preset", "T0.1 / top_p0.1 / min_p0.0 / repeat 1.05"),
    ("corpus-lfm120-t08", "LFM2.5-1.2B-Instruct", "non-thinking, round-2 temperature", "T0.8 / top_p0.95 / min_p0.0 / repeat 1.05"),
    ("corpus-tinyllama", "TinyLlama-1.1B-Chat", "non-thinking", "T0.7 / top_p0.95 / top_k50 / repeat 1.0"),
]

BY_LABEL = {}
for r in RUNS:
    BY_LABEL.setdefault(r["label"], []).append(r)
for v in BY_LABEL.values():
    v.sort(key=lambda r: (r["case"], r["draw"]))

INSTR_HEAD = (
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


def fence(text: str, lang: str = "text") -> str:
    text = "" if text is None else text
    return f"```{lang}\n{text.rstrip()}\n```"


def reference(c: dict) -> str:
    lines = []
    if c.get("required"):
        lines.append("| id | expected | what it is | key phrases |")
        lines.append("|---|---|---|---|")
        for it in c["required"]:
            keys = "; ".join(" / ".join(k) for k in it.get("keys", []))
            lines.append(f"| {it['id']} | **{it['decision']}** | {it.get('points','')} | {keys} |")
    else:
        lines.append("**No required items.** The correct answer is `NONE` (or, at most, a "
                     "single low-confidence note).")
    if c.get("forbidden"):
        lines.append("")
        lines.append("Must not appear as stored canon: " +
                     ", ".join(f"`{x}`" for x in c["forbidden"]))
    if c.get("notes"):
        lines.append("")
        lines.append(c["notes"])
    return "\n".join(lines)


def cfg_line(r: dict, mode: str) -> str:
    s = r.get("sampling") or {}
    bits = [f"`{r['label']}`", mode,
            f"seed {r.get('seed')}", f"max_tokens {r.get('max_tokens')}",
            f"ctx {r.get('ctx', '?')}"]
    if s:
        bits.append("sampling " + ", ".join(f"{k}={v}" for k, v in s.items()))
    if r.get("reasoning_budget"):
        bits.append(f"runtime budget {r['reasoning_budget']}")
    bits.append(f"finish `{r.get('finish_reason')}`")
    bits.append(f"{r.get('completion_tokens')} tok")
    if r.get("wall_s"):
        bits.append(f"{r['wall_s']:.1f} s")
    bits.append(f"run `{r.get('run_code')}`")
    return " · ".join(bits)


def main() -> None:
    (OUT / "archive").mkdir(parents=True, exist_ok=True)
    (OUT / "models").mkdir(parents=True, exist_ok=True)

    reading = ["# Case-by-case reading copy", "",
               "Draw 1 of each configuration, verbatim, with the reference material for the case.",
               "All draws and all reasoning are in `archive/case-Cxx.md`;",
               "per-model files are in `models/`.", ""]

    for c in CASES:
        cid = c["id"]
        prompt = INSTR_HEAD + "\n\nConversation:\n" + c["conversation"]
        doc = [f"# {cid} — {c['category']} ({c['length']})", "",
               "## Source material (exactly what the model received after the instruction)", "",
               fence(c["conversation"]), "",
               "## Reference material", "", reference(c), "",
               "## Exact prompt", "", fence(prompt), "",
               "## Responses, all draws", ""]
        for label, name, mode, preset in ORDER:
            rows = [r for r in BY_LABEL.get(label, []) if r["case"] == cid]
            if not rows:
                continue
            doc.append(f"### {name} — {mode}")
            doc.append("")
            doc.append(f"*{preset}*")
            doc.append("")
            for r in rows:
                doc.append(f"**draw {r['draw']}** — {cfg_line(r, mode)}")
                doc.append("")
                if r.get("reasoning_content"):
                    doc.append("<details><summary>reasoning (verbatim)</summary>")
                    doc.append("")
                    doc.append(fence(r["reasoning_content"]))
                    doc.append("")
                    doc.append("</details>")
                    doc.append("")
                doc.append(fence(r.get("content") or "(empty content)"))
                doc.append("")
        (OUT / "archive" / f"case-{cid}.md").write_text("\n".join(doc), encoding="utf-8")

        reading += [f"## {cid} — {c['category']} ({c['length']})", "",
                    "**Source**", "", fence(c["conversation"]), "",
                    "**Reference**", "", reference(c), ""]
        for label, name, mode, preset in ORDER:
            rows = [r for r in BY_LABEL.get(label, []) if r["case"] == cid]
            if not rows:
                continue
            r = rows[0]
            reading.append(f"**{name}** ({mode}) — {r.get('completion_tokens')} tok, "
                           f"{r.get('finish_reason')}")
            reading.append("")
            reading.append(fence(r.get("content") or "(empty content)"))
            if len(rows) > 1:
                other = "; ".join(f"d{x['draw']}: {x.get('finish_reason')}/{x.get('completion_tokens')}tok"
                                  for x in rows[1:])
                reading.append(f"*also: {other}*")
            reading.append("")

    (OUT / "READING.md").write_text("\n".join(reading), encoding="utf-8")

    for label, name, mode, preset in ORDER:
        rows = BY_LABEL.get(label, [])
        if not rows:
            continue
        doc = [f"# {name} — {mode}", "",
               f"Preset: *{preset}*", "",
               f"{len(rows)} runs over {len({r['case'] for r in rows})} cases. "
               f"Prompt, reference material and other models' answers for each case: "
               f"`../archive/case-Cxx.md`.", ""]
        for c in CASES:
            cid = c["id"]
            cr = [r for r in rows if r["case"] == cid]
            if not cr:
                continue
            doc += [f"## {cid} — {c['category']}", ""]
            for r in cr:
                doc.append(f"**draw {r['draw']}** — {cfg_line(r, mode)}")
                doc.append("")
                if r.get("reasoning_content"):
                    doc += ["<details><summary>reasoning (verbatim)</summary>", "",
                            fence(r["reasoning_content"]), "", "</details>", ""]
                doc.append(fence(r.get("content") or "(empty content)"))
                doc.append("")
        slug = label.replace("corpus-", "")
        (OUT / "models" / f"{slug}.md").write_text("\n".join(doc), encoding="utf-8")

    print("cases:", len(CASES), "reading KB:", round((OUT / "READING.md").stat().st_size / 1024))
    for p in sorted((OUT / "models").glob("*.md")):
        print(" ", p.name, round(p.stat().st_size / 1024, 1), "KB")


if __name__ == "__main__":
    main()
