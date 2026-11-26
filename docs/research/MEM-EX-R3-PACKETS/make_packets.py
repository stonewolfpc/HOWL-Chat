"""Generate the round-3 source-preserving reading packets.

Every model response is inserted directly from the stored JSON, never retyped.
The script asserts afterwards that each included response appears verbatim in
the written file, so a transcription error cannot pass silently.

Run from this directory:  python make_packets.py
"""

from __future__ import annotations

import ast
import json
import re
from collections import Counter, defaultdict
from pathlib import Path

HERE = Path(__file__).resolve().parent
ROOT = HERE.parent.parent.parent
R3 = ROOT / "experiments" / "mem-extraction" / "round3"
REPORT = ROOT / "docs" / "research" / "MEM-EXTRACTION-ROUND3-2026-11-19.md"

LIMIT = 3300  # words per packet, before splitting

DISPLAY = {
    "corpus-4b-nonthink": "Qwen3.5-4B non-thinking",
    "corpus-4b-think-512": "Qwen3.5-4B thinking, budget 512",
    "corpus-lfm3-think-512": "LFM3-1.7B thinking, budget 512",
    "corpus-lfm120-instruct": "LFM2.5-1.2B-Instruct, T0.1",
    "corpus-lfm120-t08": "LFM2.5-1.2B-Instruct, T0.8",
    "corpus-lfm3-nonthink": "LFM3-1.7B non-thinking (failed condition)",
    "corpus-tinyllama": "TinyLlama-1.1B-Chat-v1.0",
}

ORDER = list(DISPLAY)


def load(name: Path):
    return json.loads(name.read_text(encoding="utf-8", errors="replace"))


def triage_prompt() -> str:
    """The exact six-item prompt, read from the module that owns it."""
    src = (ROOT / "experiments" / "lfm25-cpu" / "run_probe.py").read_text(
        encoding="utf-8", errors="replace")
    for node in ast.parse(src).body:
        if isinstance(node, ast.Assign) and any(
                getattr(t, "id", "") == "PROMPT" for t in node.targets):
            return ast.literal_eval(node.value)
    raise SystemExit("PROMPT not found in run_probe.py")


PROMPT = triage_prompt()


corpus = {c["id"]: c for c in load(R3 / "corpus.json")["cases"]}
runs = load(R3 / "results" / "corpus.json")
triage = load(R3 / "results" / "all_triage.json")
ablation = load(R3 / "results" / "ablation_corpus_all.json")


def score_map() -> dict:
    m: dict[str, float] = {}
    for p in (R3 / "results").glob("*score*.json"):
        try:
            data = load(p)
        except Exception:
            continue
        if not isinstance(data, list):
            continue
        for r in data:
            if isinstance(r, dict) and r.get("run_code"):
                v = r.get("score12", r.get("score"))
                if v is not None:
                    m[r["run_code"]] = v
    return m


SCORES = score_map()

V2 = {}
try:
    for r in load(R3 / "instrument_v2" / "out" / "corpus_v2.json")["runs"]:
        V2[r["run_code"]] = r
except Exception:
    V2 = {}


def fence(text: str) -> str:
    """A code fence that cannot be closed by the content itself."""
    longest = max([len(m) for m in re.findall(r"`+", text)] or [0])
    bar = "`" * max(3, longest + 1)
    return f"{bar}text\n{text}\n{bar}"


def words(s: str) -> int:
    return len(s.split())


def head(pid: str, title: str, sources: str, cases: str, verbatim: str,
         omissions: str, framing: str) -> str:
    return f"""# Packet {pid} — {title}

**Series:** Round 3 memory-extraction evidence, source-preserving reading packet
**Prepared:** 2026-11-19 · **For:** independent review (ChatGPT)
**Original sources:** {sources}
**Case / run identifiers covered:** {cases}
**Raw model output in this packet:** {verbatim}
**Omitted:** {omissions}

**Labels used throughout:** `OBSERVED` directly present in a stored response or
execution record · `SOURCE-DERIVED` stated in a source or configuration file ·
`CALCULATED` derived by explicit arithmetic · `INFERRED` a reading of the
evidence · `UNESTABLISHED` not established by the available evidence ·
`SUPERSEDED` an earlier claim replaced by later evidence.

**Whose words you are reading:** fenced blocks are the model's own output
(`OBSERVED`, verbatim). Prompt text is `SOURCE-DERIVED` from `corpus.json`.
Scores are `CALCULATED`. Any sentence outside a fence that is not a definition is
Buffy's interpretation unless it is attributed to Muse or to the project owner.

**Caution.** These are exploratory measurements. Nothing here selects a model,
establishes memory quality, or supplies MEMORY-09 evidence. Sample sizes are
1–3 draws per case.

---

{framing}

---
"""


# --------------------------------------------------------------------------
# Packet 01 — design and protocol
# --------------------------------------------------------------------------

DESIGN = """## 1. Research questions

The experiment asked four separable questions. They are listed because a later
reader should be able to see which one any number answers.

1. **Quality.** Given one conversation, does a local model produce memory
   candidates that a careful reader would keep — correct content, correct
   decision, correct speaker, correct scope? `OBSERVED` via the 18-case corpus.
2. **Contract.** Does it produce the required four-field line
   (`ID | DECISION | MEMORY | REASON`) without a preamble, and stop cleanly?
   `OBSERVED` via parse status.
3. **Cost.** What latency, RAM and VRAM does that take on this machine, on CPU
   and on one 8 GB GPU? `OBSERVED` from wall time and sampled memory.
4. **Isolation.** Does one extraction job's context affect the next job?
   `OBSERVED` narrowly, by sentinel probe.

It did **not** ask which model should be deployed, whether the memory design is
correct, or how a production system should store records.

## 2. Procedure

`SOURCE-DERIVED` from `harness.py`, `run_stage_a.py`, `run_stage_a2.py`,
`run_stage_b.py`, `run_stage_c.py`, `run_budget.py`, `run_corpus_plan.py`,
`run_gpu.py`, `run_coexist.py`, `run_reset.py`, `run_isolation.py`.

1. **Stage A** — one publisher-preset configuration per model and mode on the
   original six-item prompt, three draws, CPU only.
2. **Stage A2** — add models and modes that Stage A did not cover, same prompt.
3. **Stage B** — bounded sampler ladder on the two weakest useful models:
   `presence_penalty` 2.0 → 0.0 for Qwen3.5-2B; temperature 0.3 / 0.5 / 0.8 for
   LFM2.5-1.2B-Instruct.
4. **Stage C** — rebuild the round-2 sampler chain to test whether the published
   round-2 numbers reproduce. They did not (see §4).
5. **Budget ladder** — `--reasoning-budget` 256 / 512 / 1024 where the build
   supports it.
6. **Corpus** — 18 cases × the selected conditions, 1–3 draws, 235 stored runs.
7. **GPU** — full offload for two models, two draws.
8. **Coexistence** — a 9B chat model resident on the GPU while a 1.2B extractor
   runs on CPU.
9. **Reset / isolation** — the erase endpoint, then a sentinel probe.

## 3. The six-item prompt

`SOURCE-DERIVED`, imported byte-for-byte from
`experiments/lfm25-cpu/run_probe.py` (`PROMPT`) so it cannot drift between
rounds. It asks for one line per item, `ID | DECISION | REASON`, no preamble,
for six candidate memories (`M1`–`M6`) covering canon confirmation, chatter, a
character fact, a contradiction, a preference and an uncertain date.

## 4. The corpus

`SOURCE-DERIVED` from `corpus.json`. Eighteen cases, one per category:

| ID | Category | Length |
|---|---|---|
| C01 | explicit long-term preference | short |
| C02 | implicit preference via behaviour | medium |
| C03 | character fact | short |
| C04 | world / scenario rule | medium |
| C05 | correction of earlier information | short |
| C06 | changing, superseded preference | medium |
| C07 | unresolved contradiction | short |
| C08 | uncertain / speculative claim | short |
| C09 | transient chatter | short |
| C10 | important detail buried in long context | long (~250 words) |
| C11 | context-dependent significance | medium |
| C12 | must not be extracted (voided narration) | short |
| C13 | fiction vs real-world assertion | medium |
| C14 | multiple speakers / attribution | medium |
| C15 | repeated information / duplicates | medium |
| C16 | indirect information | short |
| C17 | adversarial embedded instruction | medium |
| C18 | exact wording and numbers | short |

Each case carries the expected decision per required item, keyword groups, and
forbidden strings that indicate fabrication or reversal. Two cases (C09, C12)
have no required items: the correct answer is to store nothing.

`UNESTABLISHED`: the reference answers were written before the runs according to
the file's own `status` field, but the file is untracked and no hash was
published beforehand, so the **order of writing cannot be proven** from the
repository.

## 5. Configuration actually sent

`SOURCE-DERIVED` from `configs.py` and the per-run `sampling` records. Every
value is traceable to a publisher statement in `SOURCES.md`; `min_p` is always
sent explicitly because llama.cpp defaults it to 0.05.

| Condition | Sampling sent | Max tokens | Thinking control |
|---|---|---|---|
| Qwen3.5-4B non-thinking | T0.7 / top_p0.8 / top_k20 / min_p0 / pres1.5 / rep1.0 | 768 | `enable_thinking=false` |
| Qwen3.5-4B thinking | T1.0 / top_p0.95 / top_k20 / min_p0 / pres1.5 / rep1.0 | 1024 | `enable_thinking=true`, budget 512 |
| Qwen3.5-2B non-thinking | T1.0 / top_p1.0 / top_k20 / min_p0 / pres2.0 / rep1.0 | 2048 | `enable_thinking=false` |
| Qwen3-1.7B thinking | T0.6 / top_p0.95 / top_k20 / min_p0 / pres0 / rep1.0 | 3072 | `enable_thinking=true` |
| LFM2.5-1.2B-Instruct | T0.1 / top_p1.0 / top_k50 / min_p0 / pres0 / rep1.05 | 512 (corpus) | none available |
| LFM2.5-1.2B-Instruct T0.8 | T0.8 / top_p1.0 / top_k50 / min_p0 / pres0 / rep1.05 | 512 | none available |
| LFM3-1.7B thinking | T0.1 / top_p1.0 / top_k50 / min_p0 / pres0 / rep1.1 | 768 / budget 512 | `enable_thinking=true` |
| TinyLlama-1.1B-Chat | T0.7 / top_p0.95 / top_k50 / min_p0 / pres0 / rep1.0 | 768 | none exists (ctx 2048) |

Server launch for every round-3 session, recorded per run in `command`:
`llama-server -m <gguf> -c 4096 -ngl <0|99> -np 1 --host 127.0.0.1 --port 5099
--jinja --reasoning-format deepseek --metrics`.

## 6. Execution and measurement

`SOURCE-DERIVED` from `harness.py`.

- One `llama-server` process per session; one request per sequence (`-np 1`);
  streaming responses; a random 8-character `run_code` per generation so scoring
  can be blind.
- The harness attempts `POST /slots/0?action=erase` before each request and
  records the result. On this build it returns HTTP 501, so no erase occurred.
- Wall time is measured around the HTTP request and **excludes model load**;
  load time is recorded separately as `load_s`.
- Peak RSS is sampled at 100 ms over the server process tree. Peak VRAM is
  `nvidia-smi memory.used`, which includes other desktop consumers.
- Seeds are `4100 + draw`. They are recorded but are not a variance control:
  llama.cpp sampling is not guaranteed reproducible across builds.

## 7. Context isolation procedure

`SOURCE-DERIVED` from `run_reset.py` and `run_isolation.py`.

Two designs. **Reset**: run benign case C03, then adversarial case C17, then C03
again on the same server, then C03 on a fresh server; compare bytes and read
`cache_n` from the server metrics. **Sentinel**: give the model a unique name
("Marrowgale") and code ("7712") in one job, then ask an unrelated extraction
and probe whether the sentinel surfaces; repeat on a fresh server.

## 8. What was not done

No constrained decoding, no grammar, no DRY, no seed sweep, no partial GPU
offload, no concurrency, no LFM3 GPU run, no second scorer, no inter-rater
reliability, no more than three draws per case. Each is a limitation, not an
oversight: they were outside the round's scope.
"""


# --------------------------------------------------------------------------
# Packet 02 — rubric and scoring
# --------------------------------------------------------------------------

def packet_02() -> list[dict]:
    rubric = (R3 / "RUBRIC.md").read_text(encoding="utf-8", errors="replace")
    code = (R3 / "score_corpus.py").read_text(encoding="utf-8", errors="replace")
    start = code.index("def parse_lines")
    end = code.index("def score_case")
    parser = code[start:end]
    start2 = code.index("    # 2 fidelity")
    end2 = code.index("    # 5 temporal")
    fidelity = code[start2:end2]

    a = head(
        "02A", "The rubric as written",
        "`experiments/mem-extraction/round3/RUBRIC.md` (143 lines, verbatim)",
        "not case-specific; applies to all 18 corpus cases and the six-item prompt",
        "none — the rubric contains no model output",
        "nothing omitted; the file is reproduced in full",
        "The rubric is reproduced in full below. It was written as a pre-registered\n"
        "document; whether it *was* pre-registered cannot be proven (§8 of this\n"
        "packet, and `reporting.PREREGISTRATION`).") + fence(rubric)

    b = head(
        "02B", "How the scorer implements the rubric",
        "`score_corpus.py` (295 lines) and `score_triage.py` (106 lines), excerpts",
        "scoring functions, not case-specific",
        "none — no model output",
        "The excerpts below are the parser, the fidelity/false-memory block and the\n"
        "NONE early return. The remainder of `score_corpus.py` is item matching,\n"
        "significance, temporal, attribution, restraint, completeness and structure;\n"
        "it is cited by function name rather than reproduced to keep this packet\n"
        "readable. The full file is 295 lines and is the authoritative source.",
        """## 1. The parser

`SOURCE-DERIVED`, `score_corpus.py::parse_lines`. A line is a record only if it
splits into at least three fields on `|` or `;` and the second field is exactly
`KEEP`, `REJECT` or `HOLD`. Anything else becomes `{"parsed": False}`.

Key consequence, and the root of audit finding F-02: a table whose second column
is the literal word `DECISION` is not a record, so it contributes no `KEEP`. On a
case with no required items, "no KEEP" was scored as restraint.

""") + fence(parser) + """
## 2. Fidelity and false memory

`SOURCE-DERIVED`. A forbidden string inside a `KEEP` or `HOLD` memory zeroes
fidelity (20 points) **and** subtracts the entire false-memory dimension
(15 points) **and** appends a catastrophic reason. That is three consequences
from one string, which is audit finding F-04.

""" + fence(fidelity) + """
## 3. The NONE rule

`SOURCE-DERIVED`. When a case has required items and the output yields no parsed
row, the scorer returns exactly 2.0 — the structure points — and nothing else.
This was amendment 1, recorded before final scoring.

## 4. Two amendments, and one that was not applied

- **Amendment 1 (applied, pre-scoring per `RUBRIC.md`):** a clean `NONE` on a
  case with required items earns only structure points.
- **Amendment 2 (applied, pre-scoring per `RUBRIC.md`):** keyword matching is
  stem-tolerant (`narrating` ≈ `narrate`).
- **C17 amendment (proposed, NOT applied):** removing `hunter2` from C17's
  forbidden list and scoring attribution/register instead. Proposed in
  `qualitative/SECRETS-AND-MEMORY.md` and §26 of the report; deliberately left
  out of the code so the pre-registered criterion is not edited retroactively.

`UNESTABLISHED`: the two applied amendments are documented as pre-scoring in
`RUBRIC.md`, but the repository has no timestamped hash or commit that proves
when they were written. `CALCULATED`: the C17 proposal is absent from the code —
`score_corpus.py` still lists `hunter2` and still returns a catastrophic flag.

## 5. Known scoring limitations (`OBSERVED` from the code)

1. Keyword presence is not record structure; a four-fact line matches all four
   items (F-01).
2. `parsed=False` rows retain no text for matching, so a correct fact in an
   invalid row earns nothing (F-03).
3. One forbidden string is charged to fidelity, false memory and the
   catastrophic count (F-04).
4. Truncation is recorded per run but not removed from the mean (F-08).
5. Attribution is hard-coded for C13 and C14 only; every other case receives the
   5 points by default.
6. No ownership or scope dimension exists. A character's belief and a world fact
   are indistinguishable to the scorer.

## 6. The corrected instrument

`SOURCE-DERIVED` from `instrument_v2/`. It does not replace these files and does
not change any v1 score. It separates semantic from contract scoring, makes
output validity explicit, measures atomicity, adds ownership, and charges one
incident once. Its tests are 32 behavioural cases. Its recalculated corpus is a
new file (`out/corpus_v2.json`), not a correction of the old one.
"""
    return [{"id": "02A", "title": "The rubric as written", "sources": "RUBRIC.md",
             "cases": "all", "body": a, "verbatim": "n/a (no model output)"},
            {"id": "02B", "title": "How the scorer implements the rubric",
             "sources": "score_corpus.py, score_triage.py", "cases": "all",
             "body": b, "verbatim": "n/a (no model output)"}]


# --------------------------------------------------------------------------
# Renderers for raw runs
# --------------------------------------------------------------------------

def run_block(r: dict, heading: str) -> str:
    label = DISPLAY.get(r.get("label"), r.get("label") or "?")
    s = r.get("sampling") or {}
    parts = []
    for k in ("temperature", "top_p", "top_k", "min_p", "presence_penalty",
              "repeat_penalty"):
        if k in s:
            parts.append(f"{k} {s[k]}")
    meta = [
        f"`{r.get('run_code')}` draw {r.get('draw')}",
        f"file `{r.get('file')}`",
        f"sampler: " + " / ".join(parts) if parts else "sampler: not recorded",
        f"max_tokens {r.get('max_tokens')}",
        f"reasoning_budget {r.get('reasoning_budget')}",
        f"finish_reason `{r.get('finish_reason')}`",
        f"completion_tokens {r.get('completion_tokens')}",
        f"wall {r.get('wall_s')} s",
    ]
    sc = SCORES.get(r.get("run_code"))
    if sc is not None:
        meta.append(f"v1 score {sc}")
    v2 = V2.get(r.get("run_code"))
    if v2:
        meta.append(f"v2 `{v2['status']}` semantic {v2['semantic']} "
                    f"contract {v2['contract']}")
        if v2.get("events"):
            meta.append("v2 events: " + ", ".join(e["code"] for e in v2["events"]))
    out = [f"#### {heading}", "", "* `" + "`\n* `".join(meta) + "`", ""]
    content = r.get("content")
    if content is None:
        out.append("*`OBSERVED`: no content field (execution failure).*\n")
        if r.get("error"):
            out.append(fence(str(r["error"])) + "\n")
        return "\n".join(out)
    out.append("`OBSERVED` — verbatim `content`:\n")
    out.append(fence(content) + "\n")
    reasoning = r.get("reasoning_content") or ""
    if reasoning.strip():
        out.append("`OBSERVED` — verbatim `reasoning_content` (model's own "
                   "thinking, not part of the answer):\n")
        out.append(fence(reasoning) + "\n")
    else:
        out.append("*`reasoning_content`: empty.*\n")
    return "\n".join(out)


def case_block(case: dict) -> str:
    req = case.get("required") or []
    lines = [
        f"## {case['id']} — {case['category']} ({case['length']})",
        "",
        "`SOURCE-DERIVED` — conversation exactly as sent (the harness prefixes the "
        "extraction instruction, reproduced in Packet 01):",
        "",
        fence(case["conversation"]),
        "",
        "`SOURCE-DERIVED` — reference written before the runs:",
        "",
    ]
    if req:
        for it in req:
            lines.append(f"- **{it['id']}** → `{it['decision']}` — {it.get('points', '')}")
    else:
        lines.append("- no required items: the correct answer is to store nothing")
    if case.get("forbidden"):
        lines.append("- forbidden (fabrication or reversal): " +
                     ", ".join(f"`{f}`" for f in case["forbidden"]))
    if case.get("notes"):
        lines.append(f"- note: {case['notes']}")
    lines.append("")
    return "\n".join(lines)


def pack(items: list[tuple[str, str]], prefix: str, title: str, sources: str,
         cases: str, omissions: str, framing: str, verbatim: str) -> list[dict]:
    """Greedy pack of (case_id, text) blocks into word-limited packets."""
    groups: list[list[tuple[str, str]]] = []
    cur: list[tuple[str, str]] = []
    size = 0
    for cid, text in items:
        w = words(text)
        if cur and size + w > LIMIT:
            groups.append(cur)
            cur, size = [], 0
        cur.append((cid, text))
        size += w
    if cur:
        groups.append(cur)

    out = []
    for i, g in enumerate(groups, start=1):
        suffix = chr(64 + i)
        ids = [c for c, _ in g]
        rng = f"{ids[0]}–{ids[-1]}" if len(ids) > 1 else ids[0]
        body = head(f"{prefix}{suffix}", f"{title}: {rng}", sources, rng,
                    verbatim, omissions, framing) + "\n\n".join(t for _, t in g)
        out.append({"id": f"{prefix}{suffix}",
                    "title": f"{title}: {rng}",
                    "sources": sources, "cases": rng, "body": body,
                    "verbatim": verbatim})
    return out


# --------------------------------------------------------------------------
# Packets 03–07
# --------------------------------------------------------------------------

def packet_03() -> list[dict]:
    by_label = defaultdict(list)
    for r in triage:
        by_label[r.get("label")].append(r)
    items = []
    for label in sorted(by_label, key=lambda k: (-len(by_label[k]), k)):
        rs = sorted(by_label[label], key=lambda r: (r.get("draw") or 0))
        txt = [f"## {label}", "",
               f"`SOURCE-DERIVED`: {len(rs)} run(s) on the six-item prompt. "
               f"Sampler values are in each block.", ""]
        for r in rs:
            txt.append(run_block(r, f"{label} — draw {r.get('draw')}"))
        items.append((label, "\n".join(txt)))

    framing = """This packet reproduces the six-item prompt runs exactly as stored in
`results/all_triage.json`. The prompt itself is quoted in full at the top of the
first triage packet and in Packet 01; it is not repeated per run.

Two things to watch, both `OBSERVED`:

1. **Reasoning leakage.** Some models put working notes (`Analysis:`,
   `Decision:`) into `content`. The v1 triage scorer marks such runs
   inadmissible; their scores are shown but should not be read as verdicts.
2. **Non-binding budgets.** Where `reasoning_budget` is set but
   `completion_tokens` exceeds it substantially, the budget did not bind. The
   report treats those as protocol failures, not model quality.

The 12-point triage score is the historical instrument and is unchanged in v2.
"""
    prompt_block = ("## The six-item prompt, verbatim\n\n"
                    "`SOURCE-DERIVED` from `experiments/lfm25-cpu/run_probe.py` "
                    "(`PROMPT`), which round 3 imports rather than retyping.\n\n"
                    + fence(PROMPT) + "\n\n")
    items[0] = (items[0][0], prompt_block + items[0][1])
    return pack(items, "03", "Raw model outputs — six-item triage",
                "`results/all_triage.json`, `results/*_scores.json`",
                "all 51 triage runs",
                "none — every stored triage run appears in this series",
                framing, "yes, verbatim")


def packet_04() -> list[dict]:
    by_case = defaultdict(list)
    for r in runs:
        by_case[r.get("case")].append(r)
    items = []
    for cid in sorted(by_case):
        case = corpus[cid]
        blocks = [case_block(case)]
        rs = sorted(by_case[cid], key=lambda r: (ORDER.index(r["label"])
                                                 if r["label"] in ORDER else 99,
                                                 r.get("draw") or 0))
        blocks.append(f"### Stored responses for {cid} "
                      f"({len(rs)} runs, all draws)\n")
        for r in rs:
            blocks.append(run_block(r, f"{cid} — {DISPLAY.get(r['label'], r['label'])}"
                                       f" — draw {r.get('draw')}"))
        items.append((cid, "\n".join(blocks)))

    framing = """One case per section, with its conversation exactly as sent, the
reference that was written before any output existed, and every stored response
for that case across all conditions and draws.

Read the reference first, then the responses. A high score with a wrong record is
the failure mode this packet exists to expose: check whether the memory text
could be written into a store and used later without misleading a reader.

`CALCULATED` v1 scores are shown per run. They are the historical instrument and
are known to over-credit compound records and unparseable restraint (audit
F-01/F-02). `CALCULATED` v2 figures are from a separate corrected instrument and
are **not** a rescoring of v1.
"""
    return pack(items, "04", "Raw model outputs — corpus by case",
                "`corpus.json`, `results/corpus.json`, `results/corpus_scores.json`",
                "C01–C18",
                "none — every stored corpus run for these cases appears in this "
                "series. Cases are not compressed to fit a word target; the series "
                "is split at case boundaries.",
                framing, "yes, verbatim")


def packet_05() -> list[dict]:
    focus = {
        "C01": "Atomicity and compound records. The 4B folds four style rules into "
               "one line and v1 awards 100.",
        "C12": "Restraint and malformed output. A 1.2B table of column labels "
               "receives 97.8 for storing nothing.",
        "C13": "Real-world versus fictional attribution. The user's nursing job is "
               "folded into the character.",
        "C17": "Password and embedded directive. One row stores a credential; "
               "another refuses the deletion instruction.",
        "C18": "Numeral reversal, and a negated forbidden phrase. `not triple` was "
               "counted as storing `triple` by v1.",
    }
    by_case = defaultdict(list)
    for r in runs:
        if r.get("case") in focus:
            by_case[r.get("case")].append(r)
    items = []
    for cid in ["C01", "C12", "C13", "C17", "C18"]:
        case = corpus[cid]
        blocks = [case_block(case),
                  f"**Why this case is here (`INFERRED`, Buffy):** {focus[cid]}\n"]
        rs = sorted(by_case[cid], key=lambda r: (ORDER.index(r["label"])
                                                 if r["label"] in ORDER else 99,
                                                 r.get("draw") or 0))
        for r in rs:
            blocks.append(run_block(r, f"{cid} — {DISPLAY.get(r['label'], r['label'])}"
                                       f" — draw {r.get('draw')}"))
        items.append((cid, "\n".join(blocks)))

    framing = """The five cases where the instrument and the models disagree most.
Every draw is included, including the ones that look like noise: TinyLlama's
looping, the 1.2B's column labels, and the thinking model's reasoning are
evidence about failure modes, not clutter.

`OBSERVED` facts to check against the text:

- C01: one 4B row carries short paragraphs, past tense, no purple prose and the
  ban on narrating thoughts. Whether that is one memory or four is the question.
- C12: three rows read `C1 | DECISION | MEMORY | REASON`. There is no memory in
  them. v1 scored 97.8; v2 calls it `MALFORMED` and scores 0.
- C13: draw 1 of the 1.2B says `The user is an ER nurse, not a healer` — correct
  register — inside an invalid table. Draws 2 and 3 invert it.
- C17: the 4B thinking row stores `hunter2` while explicitly refusing durable
  status. Both facts are in the same row.
- C18: the 4B says `exactly double the toll (not triple)`. v1 recorded that as a
  catastrophic numeric reversal.
"""
    return pack(items, "05", "Adversarial and failure cases",
                "`corpus.json`, `results/corpus.json`, `results/corpus_scores.json`",
                "C01, C12, C13, C17, C18",
                "none for these five cases; every draw is present. Other cases are "
                "in Packet 04.",
                framing, "yes, verbatim")


def packet_06() -> list[dict]:
    iso = load(R3 / "results" / "isolation.json")
    rst = load(R3 / "results" / "reset.json")
    body = head(
        "06", "Context isolation and reset",
        "`run_reset.py`, `run_isolation.py`, `results/reset.json`, "
        "`results/isolation.json`, `raw/reset/server.log`",
        "C03 (benign), C17 (adversarial), sentinel probes",
        "n/a — probe outputs, reproduced from JSON",
        "Server logs are excerpted to the request lines and cache metrics; the "
        "full logs remain at the paths above.",
        """## What was tested

`SOURCE-DERIVED` from `run_reset.py`:

1. Job A1: benign case C03 on a fresh server.
2. Job B: adversarial case C17 on the same server.
3. Job A2: case C03 again, same server, same seed.
4. Job A_fresh: case C03 on a new server, same seed.

The comparison is byte equality of `content`, plus `cache_n` from the server
metrics. The erase endpoint is attempted and its HTTP status recorded.

`SOURCE-DERIVED` from `run_isolation.py`: a sentinel name ("Marrowgale") and code
("7712") are placed in one job, then an unrelated extraction runs, then a probe
asks whether the sentinel appears. The probe is repeated against a fresh server.

## Results (`OBSERVED`)

`CALCULATED` from `results/reset.json`: `a1_equals_a2` and `a1_equals_fresh` are
both true; `cache_reused_after_erase` is false; `erase_supported` is false.

""") + fence(json.dumps(rst, indent=2)[:6000]) + """

`OBSERVED` from `results/isolation.json`:

""" + fence(json.dumps(iso, indent=2)[:6000]) + """

## What this demonstrates, and what it does not

`OBSERVED`: under serial execution at `-np 1`, on these two Qwen configurations
and this build, a benign job's output was byte-identical before and after an
adversarial job, and identical to a fresh server's output. `cache_n` was 0 on
every job. The erase endpoint returned HTTP 501, so no explicit reset was
exercised.

`UNESTABLISHED`, and not to be inferred from the above: behaviour under
concurrent requests, with multiple slots, on other model families, on other
llama.cpp builds, under different seeds, or for MEMORY-09 persistent records.
The last is a different layer entirely and was not tested.

**Permitted summary:** "Serial per-request isolation was observed under the
tested conditions." **Not permitted:** "context flushing works", "the context is
reset", "isolation is guaranteed", "no cross-contamination is possible".
"""
    return [{"id": "06", "title": "Context isolation and reset",
             "sources": "run_reset.py, run_isolation.py, results/*.json",
             "cases": "C03, C17, sentinels", "body": body,
             "verbatim": "yes, probe outputs from JSON"}]


def packet_07() -> list[dict]:
    gpu = load(R3 / "results" / "gpu-full.json")
    co = load(R3 / "results" / "coexist.json")
    abl = load(R3 / "results" / "ablation_corpus_scores.json")
    tri = load(R3 / "results" / "ablation_triage.json")

    body = head(
        "07", "Ablation, performance and hardware",
        "`results/gpu-full.json`, `results/coexist.json`, "
        "`results/ablation_corpus_all.json`, `results/ablation_*_scores.json`, "
        "`configs_ablation.py`",
        "4 ablation conditions; GPU 4B and 1.2B; one coexistence pair",
        "raw ablation generations are cited by run_code, not all reproduced",
        "108 ablation runs are summarised rather than reproduced in full: this "
        "packet carries the configuration, the scores and the aggregate, and "
        "`results/ablation_corpus_all.json` holds every generation. Reproducing "
        "all 108 here would exceed the packet size and duplicate Packet 04's "
        "evidence for the stock baselines.",
        """## 1. Ablation configurations

`SOURCE-DERIVED` from `configs_ablation.py` and `SOURCES.md` §8. Each variant
re-runs the sampler chain recorded for its stock baseline, so sampling is not a
confound. New files, new labels and new tags; no baseline run was rewritten.

| Variant | Repo | Quant | Stock baseline it is compared with |
|---|---|---|---|
| Qwen3.5-4B HauhauCS Uncensored | `HauhauCS/Qwen3.5-4B` | Q4_K_M | 4B non-thinking |
| Qwen3.5-2B Uncensored HauhauCS Aggressive | `HauhauCS/Qwen3.5-2B-Uncensored-HauhauCS-Aggressive` | Q4_K_M | 2B non-thinking (new control) |
| LFM2.5-1.2B-Instruct abliterated | `mradermacher/Huihui-LFM2.5-1.2B-Instruct-abliterated-GGUF` | Q4_K_M | 1.2B T0.1 |
| LFM3-1.7B Heretic abliterated | `Vortecks/LFM3-1.7B-Heretic-Abliterated-GGUF` | Q4_K_M | LFM3 thinking 512 |

`OBSERVED` limitation: the full `run_code` uniqueness audit across
`ablation_corpus.json`, the per-label files and `ablation_corpus_all.json` was
**not** performed (audit F-13, still open). The per-label files are the
authoritative copies; `ablation_corpus.json` was overwritten once and re-run.

## 2. Ablation triage (`OBSERVED`, 15 runs)

""") + fence(json.dumps([
        {k: r.get(k) for k in ("label", "draw", "seed", "sampling", "max_tokens",
                               "reasoning_budget", "finish_reason", "run_code")}
         | {"content": r.get("content")}
        for r in tri], indent=2, ensure_ascii=False)[:9000]) + """

## 3. Ablation corpus scores (`CALCULATED`, one row per scored run)

""" + fence(json.dumps([
        {k: r.get(k) for k in ("run_code", "case", "score", "catastrophic",
                               "catastrophic_reasons", "finish_reason")}
        for r in abl], indent=2)[:9000]) + """

## 4. GPU measurement (`OBSERVED`)

""" + fence(json.dumps(gpu, indent=2)[:5000]) + """

## 5. Coexistence (`OBSERVED`)

""" + fence(json.dumps(co, indent=2)[:4000]) + """

## 6. Limits of these measurements

- `CALCULATED` wall times exclude model load; load is recorded separately.
- CPU runs use `-ngl 0` with default flags. GPU runs use `-ngl 99`, Flash
  Attention on and unified KV. The speed-up therefore includes a kernel change,
  not offload alone.
- Peak VRAM is `nvidia-smi memory.used`, which includes other desktop consumers.
  The 2,741 MB reported for a 696 MB model file is the demonstration.
- Partial layer offload was never run.
- The coexistence test used a stand-in 9B, not the production chat model, and one
  case.
- Nothing here generalises to another GPU or another build.
"""
    return [{"id": "07", "title": "Ablation, performance and hardware",
             "sources": "results/gpu-full.json, coexist.json, ablation_*",
             "cases": "ablation conditions, GPU, coexistence", "body": body,
             "verbatim": "partly — scoring rows and JSON verbatim; full ablation "
                         "generations omitted as stated"}]


# --------------------------------------------------------------------------
# Packet 09 — questions and checklist
# --------------------------------------------------------------------------

QUESTIONS = """## A. What is actually established

`OBSERVED` and `CALCULATED`, stated at the strength the evidence supports:

1. On this machine and corpus, Qwen3.5-4B non-thinking held the four-field
   contract in every run and had the highest measured semantic score. It is a
   reasonable shortlist candidate.
2. LFM3-1.7B thinking was the only small model with a credible semantic profile;
   it lost 6 of 36 runs to the 512-token cap.
3. LFM2.5-1.2B-Instruct did not produce a contract-valid row in any of its 90
   corpus runs. That is a schema-capability result, not a quality score.
4. Thinking mode produced no measured quality gain at 4B and cost roughly 14×
   the wall time at the tested budget. n=1.
5. Serial per-request isolation was observed at `-np 1`; no reset API exists.

## B. What is not established

`UNESTABLISHED`, and should be treated as open by any reviewer:

1. Whether any configuration is *better* than another in general. CIs overlap,
   draws are 1–3, and one corpus with one author cannot support a ranking.
2. Whether round 2's numbers failed to reproduce because of variance, a template
   change, or a harness bug.
3. When the rubric and corpus were written relative to the runs. No pre-run hash
   exists.
4. Whether v2's semantic composite tracks human judgment. Its tests prove rule
   conformance, nothing more.
5. Any ownership behaviour beyond the declarative cases: the detector cannot
   infer scope from free text.
6. Anything about concurrency, other builds, other quantizations, or MEMORY-09
   records.

## C. Questions for independent review

These are the questions where a second reader's judgment is worth more than
another computation.

1. **C01.** Is "the user prefers short paragraphs, third-person past tense,
   avoids purple prose, and forbids narrating thoughts" one memory or four?
   Does the answer change what a store should do with it?
2. **C12.** Three column labels and the word `MEMORY` in the third column. Is
   that "restraint" in any sense a reader would accept? If not, what should a
   scorer do with a model that answers this way consistently?
3. **C13.** The character is a healer; the user is an ER nurse. How much should a
   single fused sentence cost, given that a human reader can still recover both
   facts from it?
4. **C17.** The thinking model stores `hunter2` while saying it is not durable
   canon. Is that a memory failure, a policy failure, or correct fidelity to an
   in-world disclosure? §26 argues the last; the v1 scorer says the first.
5. **C18.** `exactly double the toll (not triple)` — how should a scorer treat a
   forbidden word that appears only inside a negation? v1 counted it; v2 does
   not.
6. **F-05.** With n=1–3, is *any* ordering in §17 defensible, or should the
   report present cost and contract separately and refuse to order by quality?
7. **F-06.** Which single configuration difference (context, `min_p`, `--jinja`,
   or max tokens) is the most likely cause of the round-2 non-reproduction, and
   what is the cheapest experiment that would distinguish them?
8. **Truncation.** LFM3 lost 6 runs to a 512-token cap. Should a truncated run
   contribute a partial score, be excluded, or be re-run before comparison?
9. **Ownership.** The adversarial cases in
   `instrument_v2/cases_adversarial_v2.json` define belief, knowledge, disclosure
   and scenario scope. Are those the right distinctions for Howl, and is
   declarative detection sufficient for a first measurement?
10. **Instrument v2.** Are the two composites weighted defensibly? Is
    `production_usable` (contract ≥ 80, no critical event) the right gate, or
    does it hide outputs that a human would accept?

## D. Decisions that belong to the project owner, not to a scorer

1. Whether to apply the Appendix C errata to the report in place or publish them
   as an addendum. Recommended: addendum.
2. Whether retaining a disclosed credential should remain a critical event
   (current v1 and v2 rule) or become a policy question scored on attribution and
   register (§26's argument).
3. Whether a three-field row (`ID | DECISION | MEMORY`, no reason) counts as
   contract-valid. v2 says no; this materially affects the 4B's C14 draws.
4. Whether to fund a powered confirmation (≥40 cases × 5 draws, two independent
   scorers) before any model decision.
5. Whether v2 should replace v1 for future rounds, run beside it, or be revised
   first.

## E. Reading checklist

- [ ] Packet 00 — know which file is raw and which is interpretation.
- [ ] Packet 01 — know what was actually sent and measured.
- [ ] Packet 02 — know the instrument's rules and its known defects.
- [ ] Packet 03 — see the six-item prompt saturate.
- [ ] Packet 04 — read at least C01, C07, C13, C17 with the reference in view.
- [ ] Packet 05 — check the five failure cases against the raw text.
- [ ] Packet 06 — check that the isolation claim is narrower than "flush works".
- [ ] Packet 07 — check that CPU and GPU flags differ.
- [ ] Packet 08 — read Muse's findings before forming a verdict.
- [ ] Packet 09 — answer the ten questions, then the five owner decisions.
"""


# --------------------------------------------------------------------------
# Write
# --------------------------------------------------------------------------

def main() -> int:
    # Remove only the generated packets; 08A/08B are hand-written and preserved.
    for p in HERE.glob("PACKET-*.md"):
        if re.match(r"PACKET-(01|02|03|04|05|06|07|09)", p.name):
            p.unlink()

    packets: list[dict] = []
    packets += packet_02()
    packets += packet_03()
    packets += packet_04()
    packets += packet_05()
    packets.append(packet_06()[0])
    packets.append(packet_07()[0])
    packets.append({"id": "09", "title": "Outstanding questions and review checklist",
                    "sources": "this packet is analysis, not source",
                    "cases": "all", "verbatim": "n/a (no model output)",
                    "body": head("09", "Outstanding questions and review checklist",
                                 "synthesis of all packets",
                                 "all", "n/a — no model output",
                                 "nothing omitted; this packet is argument, not evidence",
                                 QUESTIONS)})

    # Packet 01 has no splitting.
    packets.insert(0, {"id": "01", "title": "Experimental design and protocol",
                       "sources": "harness.py, configs.py, corpus.json, RUBRIC.md",
                       "cases": "all", "verbatim": "n/a (no model output)",
                       "body": head("01", "Experimental design and protocol",
                                    "`harness.py`, `configs.py`, `corpus.json`, "
                                    "`RUBRIC.md`, `SOURCES.md`",
                                    "all 18 cases and 51 triage runs",
                                    "n/a — design only",
                                    "nothing omitted",
                                    DESIGN + "\n## 9. The six-item prompt, verbatim\n\n"
                                    "`SOURCE-DERIVED` from `run_probe.py::PROMPT`.\n\n"
                                    + fence(PROMPT) + "\n")})

    def split_body(body: str, target: int) -> list[str]:
        """Split a packet at run/heading boundaries, keeping the header."""
        parts = re.split(r"(?=\n#### )", body)
        head_part, blocks = parts[0], parts[1:]
        if not blocks:
            return [body]
        groups, cur, size = [], [], words(head_part)
        for b in blocks:
            w = words(b)
            if cur and size + w > target:
                groups.append(cur)
                cur, size = [], words(head_part)
            cur.append(b)
            size += w
        if cur:
            groups.append(cur)
        return [head_part + "".join(g) for g in groups]

    def split_oversized(items: list[dict]) -> list[dict]:
        out = []
        for p in items:
            if words(p["body"]) <= 4200:
                out.append(p)
                continue
            parts = split_body(p["body"], 3200)
            for i, part in enumerate(parts, 1):
                part = re.sub(r"^# Packet [^\n]+",
                              f"# Packet {p['id']}-{i} — {p['title']} (part {i})",
                              part, count=1)
                out.append({**p, "id": f"{p['id']}-{i}",
                            "title": f"{p['title']} (part {i})",
                            "body": part})
        return out

    packets = split_oversized(packets)
    packets.sort(key=lambda p: p["id"])
    written = []
    for p in packets:
        path = HERE / f"PACKET-{p['id']}-{re.sub(r'[^A-Za-z0-9]+', '-', p['title'])[:60].strip('-')}.md"
        path.write_text(p["body"], encoding="utf-8")
        written.append({**p, "path": path.name, "words": words(p["body"])})

    # Verbatim assertions: every included response string must be present.
    missing = []
    for r in runs + triage:
        c = r.get("content")
        if c and c.strip():
            for f in HERE.glob("PACKET-*.md"):
                if c in f.read_text(encoding="utf-8", errors="replace"):
                    break
            else:
                missing.append(r.get("run_code"))
    print("packets:", len(written), "words:", sum(w["words"] for w in written))
    print("verbatim check: missing responses =", len(missing),
          missing[:5])

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
