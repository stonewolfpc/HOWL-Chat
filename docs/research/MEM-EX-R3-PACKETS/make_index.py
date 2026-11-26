"""Build Packet 00 (reading guide) and INDEX.md from the written packets.

The index is generated from each packet's own header, so its word counts and
omission notes cannot drift from the files it describes.
"""

from __future__ import annotations

import re
from pathlib import Path

HERE = Path(__file__).resolve().parent


def words(s: str) -> int:
    return len(s.split())


def field(text: str, label: str) -> str:
    m = re.search(rf"^\*\*{re.escape(label)}:\*\*\s*(.+)$", text, re.M)
    return m.group(1).strip() if m else ""


def sanitise(s: str) -> str:
    return s.replace("|", "/").replace("\n", " ").strip()


rows = []
for p in sorted(HERE.glob("PACKET-*.md")):
    if p.name.startswith("PACKET-00"):
        continue
    text = p.read_text(encoding="utf-8", errors="replace")
    first = text.split("\n", 1)[0]
    m = re.match(r"# Packet ([^ ]+) — (.+)$", first)
    pid = m.group(1) if m else p.stem
    title = m.group(2) if m else p.stem
    rows.append({
        "id": pid,
        "title": title,
        "sources": field(text, "Original sources"),
        "cases": field(text, "Case / run identifiers covered"),
        "verbatim": field(text, "Raw model output in this packet"),
        "omitted": field(text, "Omitted"),
        "words": words(text),
        "file": p.name,
    })


def key(r):
    m = re.match(r"(\d+)([A-Z]?)(?:-(\d+))?", r["id"])
    return (int(m.group(1)), m.group(2) or "", int(m.group(3) or 0))


rows.sort(key=key)

table = ["| Packet | Title | Sources | Cases / runs | Words | Raw output | Omitted |",
         "|---|---|---|---|---:|---|---|"]
for r in rows:
    verb = "verbatim" if r["verbatim"].lower().startswith("yes") else r["verbatim"]
    table.append(
        f"| {r['id']} | {sanitise(r['title'])} | {sanitise(r['sources'])} | "
        f"{sanitise(r['cases'])} | {r['words']:,} | {sanitise(verb)} | "
        f"{sanitise(r['omitted'])} |")
table = "\n".join(table)

by_id = {r["id"]: r for r in rows}
total = sum(r["words"] for r in rows)
over = [r for r in rows if r["words"] > 4200]

guide = f"""# Packet 00 — Reading guide and evidence index

**Series:** Round 3 memory-extraction evidence, source-preserving reading packet
**Prepared:** 2026-11-19 · **For:** independent review (ChatGPT)
**Packets:** {len(rows)} files · **Total:** {total:,} words
**Master table:** the same table is repeated in `INDEX.md`.

## What this series is

Round 3 compared local language models on one narrow task: given a roleplay
conversation, produce memory candidates as `ID | DECISION | MEMORY | REASON`.
The series contains the raw generations, the references they were scored
against, the scoring code, the hardware measurements, and two independent
readings of the result. It contains no production decision.

## Source hierarchy

Read in this order of authority; later layers never override earlier ones.

1. **Raw generations** — `results/corpus.json` (235 runs), `results/all_triage.json`
   (51 runs), `results/ablation_corpus_all.json` (108 runs). These are the
   evidence. Nothing in the series rewrites them.
2. **Execution records** — the `command`, `sampling`, `seed`, `finish_reason`,
   token counts and timings stored beside each generation; `raw/**/server.log`.
3. **Frozen references** — `corpus.json` (18 cases with expected decisions and
   forbidden strings), `RUBRIC.md`, `SOURCES.md`.
4. **Derived scores** — `results/*_scores.json` (historical, v1) and
   `instrument_v2/out/corpus_v2.json` (separate corrected instrument, v2).
5. **Interpretation** — `docs/research/MEM-EXTRACTION-ROUND3-2026-11-19.md`,
   `qualitative/EVALUATION-BUFFY.md`, Muse's audit, Packet 09.

A derived score is not evidence that the model did something. It is evidence
that a scorer computed something. Where the two disagree, the raw text wins.

## Identifiers

- **`run_code`** — random 8-character id per generation, used for blind scoring.
  It is the join key between every packet and every JSON file.
- **`label`** — the configuration name (`corpus-4b-nonthink`, …). Written to the
  run record; *not* read by the v1 scorer, which is what made scoring blind.
- **`tag`** — the run batch that produced the file.
- **Case id** — `C01`–`C18` in `corpus.json`; `M1`–`M6` in the six-item prompt.
- **Item id** — `C1`, `C2`, … inside a response, one per candidate memory.
- **`draw`** — repetition index within the same case and configuration; seeds
  are `4100 + draw`.
- **`finish_reason`** — `stop` (model ended) or `length` (cap reached).
- **Status (v2)** — `VALID_KEEP`, `VALID_HOLD`, `VALID_REJECT`, `VALID_NONE`,
  `VALID_EMPTY`, `MALFORMED`, `UNPARSEABLE`, `MISSING`, `INCOMPLETE`,
  `EXECUTION_FAILURE`.
- **Event codes (v2)** — `FORBIDDEN_CLAIM_STORED`, `ATTRIBUTION_ERROR`,
  `PROMOTION`, `INHERITANCE`, `CONTRADICTION_COLLAPSE`, `SCENARIO_LEAK`,
  `ATOMICITY_MERGED_DISTINCT`, `EXECUTION_FAILURE`.

## Labels used in the packets

`OBSERVED` directly present in a stored response · `SOURCE-DERIVED` stated in a
source or config · `CALCULATED` explicit arithmetic · `INFERRED` a reading ·
`UNESTABLISHED` not established · `SUPERSEDED` replaced by later evidence.

**Whose interpretation:** fenced text is the model's. Packet 08 is Muse's.
Packet 09 and any unfenced commentary are Buffy's. The project owner has not
written a reading in this series.

## Raw, derived, interpreted — how to tell

| You are looking at | Layer | Trust it as |
|---|---|---|
| A fenced block after "verbatim `content`" | raw | what the model emitted |
| A `v1 score` in a run header | derived | the historical 100-point rubric |
| A `v2` line in a run header | derived | a different, corrected instrument |
| Packet 02B code excerpt | source | the actual scoring implementation |
| Packet 08 | interpretation | one adversarial reviewer's judgment |
| Packet 09 | interpretation | questions, not findings |

## Known missing evidence

1. **The Muse audit original.** No repository copy exists; Packet 08 is a
   transcription of the delivered message.
2. **Full `run_code` uniqueness proof** across the four ablation JSON files
   (audit F-13, open).
3. **Primary-source card re-fetch** at pinned revisions — sampler claims rest on
   `SOURCES.md` URLs and retrieval dates.
4. **Per-model template rendering dumps** for the Qwen default-mode conflict and
   the LFM3 off-switch.
5. **A truncation-free LFM3 corpus run** to separate quality from the 512-token
   cap.
6. **Inter-rater reliability** — one scorer, no second human.
7. **Concurrency / cross-seed isolation matrix** — never run.
8. **The DeepSeek dossier and the owner's adversarial review** — separate tracks,
   not produced.
9. **9B corpus run** — never performed.
10. **Prompt-token accounting per tokenizer** — recorded `prompt_tokens` not
    cross-checked against each tokenizer.

## Caveats before you read any score

- Draw counts are 1–3. A packet header saying `n=1` means one observation, not a
  ranking.
- The v1 and v2 totals are different instruments. Do not compare them directly
  or substitute one for the other.
- v2's tests prove the evaluator follows its own rules. They do not prove the
  rules track human judgment.
- The corpus has one author and 18 cases. Confidence intervals, where present,
  describe this corpus only.
- Nothing here is MEMORY-09 evidence or a production recommendation.

## Master index

{table}

## Recommended order for a first read

1. **00** (this packet) and **01** — what was asked and what was measured.
2. **02A** and **02B** — the rules and their known defects.
3. **05A** — C01, where a 100-point record is four rules in one line.
4. **08A** and **08B** — Muse's audit, including the missing-evidence list.
5. **04A**–**04R** as needed when a specific case is disputed.

Send **00, 01, 02A, 02B, 05A, 08A, 08B** first. That is about
{sum(by_id[i]['words'] for i in ('00','01','02A','02B','05A','08A','08B') if i in by_id):,}
words and is enough for an informed first verdict; the remaining packets are the
evidence behind individual claims.

## Size note

Target was 2,500–4,000 words per packet. {len(over)} packet(s) exceed 4,200
words after splitting at run boundaries, because a single generation with several
thousand tokens of reasoning cannot be cut without breaking the verbatim
guarantee. They are listed in `INDEX.md` with their sizes.
"""

(HERE / "PACKET-00-Reading-guide-and-evidence-index.md").write_text(
    guide, encoding="utf-8")

index = f"""# Master index — Round 3 source-preserving reading packets

Generated from each packet's own header. Word counts are approximate.

{table}

**Total:** {total:,} words across {len(rows)} packets.

**First send:** 00, 01, 02A, 02B, 05A, 08A, 08B.

**Verbatim guarantee:** every model response in Packets 03–05 was inserted
directly from `results/*.json` by `make_packets.py`, which asserts that each
stored `content` string appears in the series. Zero responses were missing at
generation time. Nothing was normalised, shortened or reworded.

**Oversize packets:** {", ".join(f"{r['id']} ({r['words']:,} words)" for r in over) or "none"}.

## Evidence that could not be included or recovered

1. **The Muse audit original file.** Produced in a chat session; no repository
   copy existed. Packet 08 is a transcription of the delivered message, and is
   the only copy in the workspace.
2. **Full `run_code` uniqueness proof across the ablation files** — never
   performed (audit F-13). The records themselves are present; the proof is not.
3. **Primary-source card text at pinned revisions** — citations and retrieval
   dates are in `SOURCES.md`; the pages were not re-fetched for this series.
4. **Chat-template renderings per model** — not dumped; the default-mode
   conflict and the LFM3 off-switch are inferred from observed reasoning fields.
5. **A truncation-free LFM3 corpus** — never run.
6. **Human scores and inter-rater reliability** — never collected; one scorer.
7. **Concurrency isolation matrix** — never run.
8. **The DeepSeek dossier and the owner's adversarial review** — separate tracks,
   not produced, so they are not in this series.
9. **The 9B corpus run** — never performed; only triage and a ceiling probe exist.
10. **Per-tokenizer prompt accounting** — recorded `prompt_tokens` were not
    cross-checked against each model's tokenizer.
11. **Full server logs** — only reset/coexistence logs are quoted; the rest are
    cited by path.
12. **108 ablation generations** — summarised in Packet 07 rather than reproduced;
    the file is cited and the omission is stated in that packet's header.
"""

(HERE / "INDEX.md").write_text(index, encoding="utf-8")

print(f"wrote PACKET-00 ({words(guide):,} words) and INDEX.md; {len(rows)} packets, {total:,} words")
print("oversize:", [(r["id"], r["words"]) for r in over])
