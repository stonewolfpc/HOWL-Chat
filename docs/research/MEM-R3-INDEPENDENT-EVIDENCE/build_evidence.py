"""Build the Round 3 independent-evidence reading packet series.

Rules this script enforces:

* Only stored generation records are used. Score files are never read.
* Nothing is scored, ranked, classified or interpreted.
* Sampler values, context settings, token counts and hardware measurements are
  written ONLY to REFERENCE-EXECUTION-METADATA.md, never to a reading packet.
* Every stored `content` and `reasoning_content` string is written as bytes and
  then re-checked as an exact byte substring of the packet that carries it.
* Duplicate copies of a run are compared; disagreement is reported, not hidden.

Run from this directory:  python build_evidence.py
"""

from __future__ import annotations

import ast
import hashlib
import json
import re
from collections import Counter, defaultdict
from pathlib import Path

HERE = Path(__file__).resolve().parent
ROOT = HERE.parent.parent.parent
R3 = ROOT / "experiments" / "mem-extraction" / "round3"
RESULTS = R3 / "results"

LIMIT = 3400          # words per packet before splitting at a run boundary
HARD_LIMIT = 4200     # report packets above this rather than hide them

# ---------------------------------------------------------------------------
# Prompt provenance (read from the code, not retyped)
# ---------------------------------------------------------------------------

HARNESS_SRC = (R3 / "harness.py").read_text(encoding="utf-8", errors="strict")


def _constant(path: Path, name: str) -> str:
    tree = ast.parse(path.read_text(encoding="utf-8", errors="strict"))
    for node in tree.body:
        if isinstance(node, ast.Assign) and any(
                getattr(t, "id", "") == name for t in node.targets):
            return ast.literal_eval(node.value)
    raise SystemExit(f"{name} not found in {path}")


EXTRACTION_INSTRUCTION = _constant(R3 / "harness.py", "EXTRACTION_INSTRUCTION")
TRIAGE_PROMPT = _constant(ROOT / "experiments" / "lfm25-cpu" / "run_probe.py", "PROMPT")

COMPOSE_LINE = next(
    i for i, line in enumerate(HARNESS_SRC.splitlines(), 1)
    if 'user = EXTRACTION_INSTRUCTION' in line)

# ---------------------------------------------------------------------------
# Load every generation record, deduplicated by run_code
# ---------------------------------------------------------------------------

CORPUS_DOC = json.loads((R3 / "corpus.json").read_bytes().decode("utf-8"))
CASES = {c["id"]: c for c in CORPUS_DOC["cases"]}


def load_list(path: Path):
    try:
        data = json.loads(path.read_bytes().decode("utf-8"))
    except Exception:
        return []
    return data if isinstance(data, list) else []


FILES = sorted(p.name for p in RESULTS.glob("*.json"))
GEN: dict[str, list[tuple[str, int, dict]]] = defaultdict(list)
for name in FILES:
    rel = f"results/{name}"
    for idx, rec in enumerate(load_list(RESULTS / name)):
        if isinstance(rec, dict) and rec.get("run_code") and "content" in rec:
            GEN[rec["run_code"]].append((rel, idx, rec))

assert GEN, "no generation records found"


def rank(item: tuple[str, int, dict]) -> tuple:
    rel, idx, rec = item
    tag = rec.get("tag")
    if tag and rel == f"results/{tag}.json":
        return (0, len(rel), idx)
    if rel == "results/corpus.json" and rec.get("kind") == "corpus":
        return (1, 0, idx)
    if rel == "results/all_triage.json":
        return (2, 0, idx)
    if rel == "results/ablation_corpus_all.json":
        return (3, 0, idx)
    return (4, len(rel), idx)


RUNS = {}
for code, copies in GEN.items():
    copies_sorted = sorted(copies, key=rank)
    rel, idx, rec = copies_sorted[0]
    others = sorted({c[0] for c in copies} - {rel})
    RUNS[code] = {"rec": rec, "src": rel, "idx": idx, "copies": others,
                  "n_copies": len(copies)}

# Consistency of duplicate generation copies.
DISAGREE = []
for code, copies in GEN.items():
    seen = {(c[2].get("content"), c[2].get("reasoning_content")) for c in copies}
    if len(seen) > 1:
        DISAGREE.append(code)

# ---------------------------------------------------------------------------
# Corpus prompt reconstruction (declared, not presented as recorded)
# ---------------------------------------------------------------------------

COMPOSE_EXPR = 'user = EXTRACTION_INSTRUCTION + "\\n\\nConversation:\\n" + case["conversation"]'


def corpus_prompt(case_id: str) -> str:
    return (EXTRACTION_INSTRUCTION + "\n\nConversation:\n"
            + CASES[case_id]["conversation"])


# ---------------------------------------------------------------------------
# Rendering helpers
# ---------------------------------------------------------------------------

def fence(text: str, lang: str = "text") -> str:
    longest = max([len(m) for m in re.findall(r"`+", text)] or [0])
    bar = "`" * max(3, longest + 1)
    return f"{bar}{lang}\n{text}\n{bar}"


def words(s: str) -> int:
    return len(s.split())


def strip_fences(text: str) -> str:
    """Remove fenced blocks, honouring a closing fence of at least the
    opening fence's length. Model output may itself contain short fences."""
    out, open_len = [], 0
    for line in text.splitlines():
        s = line.strip()
        if open_len:
            if s.startswith("`" * open_len) and set(s) == {"`"}:
                open_len = 0
            continue
        m = re.match(r"^(`{3,})[^`]*$", s)
        if m:
            open_len = len(m.group(1))
            continue
        out.append(line)
    return "\n".join(out)


def ident(rec: dict) -> list[str]:
    lines = [f'* Stored model identifier: `{rec.get("label")}`',
             f'* Model file: `{rec.get("file")}`']
    if rec.get("base"):
        lines.append(f'* Stored base name: `{rec["base"]}`')
    if rec.get("tag"):
        lines.append(f'* Batch tag: `{rec["tag"]}`')
    return lines


def run_block(code: str, info: dict, heading: str | None = None) -> str:
    rec = info["rec"]
    draw = rec.get("draw")
    case = rec.get("case")
    label = rec.get("label")
    head = heading or f"`{code}` — `{label}`, draw {draw}"
    out = [f"#### {head}", ""]
    out += ident(rec)
    if case:
        out.append(f"* Case: {case}")
    if draw is not None:
        out.append(f"* Draw: {draw}")
    out.append(f'* Source: `{info["src"]}`, array index {info["idx"]} (0-based)')
    if info["copies"]:
        out.append("* Also present in: " + ", ".join(f"`{c}`" for c in info["copies"]))
    out.append("")

    content = rec.get("content")
    if content == "":
        out.append("`OBSERVED` — stored `content` is the empty string "
                   "(no characters between the fences):")
        out.append(fence("") + "\n")
    else:
        out.append("`OBSERVED` — stored `content`, verbatim:")
        out.append("")
        out.append(fence(content) + "\n")

    reasoning = rec.get("reasoning_content")
    if reasoning:
        out.append("`OBSERVED` — stored `reasoning_content`, verbatim. This is a "
                   "separate stored field, not part of `content`:")
        out.append("")
        out.append(fence(reasoning) + "\n")
    elif reasoning == "":
        out.append("*Stored `reasoning_content` is the empty string.*\n")
    else:
        out.append("*No `reasoning_content` field is present in this record.*\n")
    return "\n".join(out)


def packet_head(pid: str, title: str, evidence: str, sources: str,
                coverage: str) -> str:
    return f"""# Packet {pid} — {title}

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19 · **For:** independent review
**Evidence in this packet:** {evidence}
**Primary sources:** {sources}
**Coverage:** {coverage}

**Excluded from this series by design:** v1 and v2 scores, scorer statuses and
event labels, model rankings, expected answers, derived quality judgments, and
any reviewer's commentary. Sampler values, context settings, token counts,
timings and hardware measurements are not reproduced here; they are indexed by
`run_code` in `REFERENCE-EXECUTION-METADATA.md`.

**Labels:** `OBSERVED` = present in the stored record. `SOURCE-DERIVED` = read
from a source file; not itself a model output. `UNESTABLISHED` = not established
by the available evidence.

Everything inside a fence is model output exactly as stored. Wording outside
fences is limited to provenance, definitions and pointers.

---
"""


# ---------------------------------------------------------------------------
# Packet 00 — scope and conventions
# ---------------------------------------------------------------------------

def packet_00() -> dict:
    n_corpus = sum(1 for i in RUNS.values() if i["rec"].get("kind") == "corpus")
    n_triage = sum(1 for i in RUNS.values() if i["rec"].get("kind") == "triage")
    body = f"""# Packet 00 — Scope, exclusions, and how to read this series

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19 · **For:** independent review
**Runs covered:** {len(RUNS)} stored generations ({n_corpus} corpus-prompt,
{n_triage} six-item-prompt)

## 1. What this series is

A source-preserving view of the Round 3 model evidence, assembled above and
independently of every scoring instrument. It answers four questions per run:
what was sent, to which model, what came back, and where the record lives.

It does not score, classify, rank or interpret. No file in this series contains
a model recommendation.

## 2. What is in the primary reading set

For every one of the {len(RUNS)} stored generations:

1. the prompt, or an explicit statement that the exact request was not recorded
   (see §4 and Packet 01);
2. the stored model identifier and model file, exactly as recorded;
3. the stored response, byte-identical to the JSON value, including empty,
   malformed, repetitive and truncated output;
4. `run_code`, case, draw, batch tag and the exact source file with its array
   index.

Corpus runs are grouped by case (Packets 02A–02R) because the case conversation
is the variable a reader needs in view. Six-item runs are grouped by configuration
(Packets 03A–03M) because the prompt is identical across them.

## 3. What is deliberately absent

- **Scores of any kind**, including the v1 rubric and the separate v2
  recalculation.
- **Scorer statuses and event labels**, and any classification of a response as
  valid, malformed, truncated or catastrophic.
- **Rankings or comparative quality claims.**
- **Expected answers.** The corpus's `required`/`forbidden` blocks are the
  grading key. The model never saw them, so they are not model evidence and are
  not reproduced. Their location is given in §6 so a reviewer can obtain them
  independently.
- **Commentary by Buffy or by Muse**, including the earlier packet series'
  interpretation.
- **Sampler and execution metadata** in the packets. It is indexed separately,
  by `run_code`, in `REFERENCE-EXECUTION-METADATA.md`.
- **Any retrospective correction** to a prompt or response. Nothing here has
  been rewritten.

## 4. The prompt is the one incomplete component

No per-run request body was ever written to disk. The result files store the
response, token counts and a `prompt_tokens` figure, but not the prompt text or a
prompt hash.

What exists is the code path that built each request and the corpus text it read.
Packet 01 reproduces the extraction instruction, the exact composition
expression, and the six-item prompt, and states precisely which parts are
`SOURCE-DERIVED` reconstruction. No reconstruction is presented as a recorded
request.

**Consequence for a reviewer:** the response side of every run is
`OBSERVED`; the prompt side is `SOURCE-DERIVED` for every run. Prompt bytes
identical to what the server received are `UNESTABLISHED` throughout this series.

## 5. Evidence labels

| Label | Meaning in this series |
|---|---|
| `OBSERVED` | Present in a stored record: a response string, a recorded identifier. |
| `SOURCE-DERIVED` | Read from a source file — the instruction text, the composition expression, the case conversation. |
| `UNESTABLISHED` | Not established by the available evidence; do not infer it. |

`INFERRED`, `CALCULATED` and `SUPERSEDED` are not used: this series contains no
inference, no arithmetic over runs, and no replacement of earlier claims.

## 6. Where everything else lives

| For | File | Key |
|---|---|---|
| Every run, its packet and its source | `INDEX.md` | `run_code` |
| Sampler chain, context, token counts, timings, memory | `REFERENCE-EXECUTION-METADATA.md` | `run_code`, `label` |
| Model file, size, SHA-256, repository string where recorded | `REFERENCE-MODEL-IDENTITY.md` | `label`, `file` |
| Grading key (not reproduced here) | `experiments/mem-extraction/round3/corpus.json`, each case's `required`, `forbidden`, `notes` | case id |
| Findings and gaps | `REPORT.md` | §1–§6 |

## 7. Reading order

1. **01** — the prompt record, including what is and is not recoverable.
2. **02A–02R** — one case at a time: what was given, then every response.
3. **03A–03M** — the six-item prompt across configurations.
4. `REFERENCE-EXECUTION-METADATA.md` only when a specific claim needs checking.
5. `REPORT.md` for coverage and integrity statements.

## 8. What this series does not establish

- That the reconstructed prompt equals the bytes the server received.
- That the recorded configuration values were successfully applied.
- That any response is correct, useful or safe.
- Anything about models not present in the stored records.
"""
    return {"id": "00", "title": "Scope, exclusions, and how to read this series",
            "body": body, "runs": [], "sources": "—", "evidence": "guide only"}


# ---------------------------------------------------------------------------
# Packet 01 — the prompt record
# ---------------------------------------------------------------------------

def packet_01() -> dict:
    harness_bytes = (R3 / "harness.py").read_bytes()
    harness_hash = hashlib.sha256(harness_bytes).hexdigest()
    corpus_bytes = (R3 / "corpus.json").read_bytes()
    corpus_hash = hashlib.sha256(corpus_bytes).hexdigest()
    probe_bytes = (ROOT / "experiments" / "lfm25-cpu" / "run_probe.py").read_bytes()
    probe_hash = hashlib.sha256(probe_bytes).hexdigest()

    body = f"""# Packet 01 — The prompt record

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19
**Evidence in this packet:** the stored instruction text, the recorded composition
expression, and the six-item prompt
**Primary sources:** `experiments/mem-extraction/round3/harness.py`;
`experiments/lfm25-cpu/run_probe.py`; `experiments/mem-extraction/round3/corpus.json`

**Status of the prompt evidence:** `SOURCE-DERIVED`, not `OBSERVED`. No request
body, prompt string or prompt hash was stored per run. See §5.

---

## 1. Corpus prompt: the instruction, verbatim from source

`SOURCE-DERIVED` — `harness.py`, constant `EXTRACTION_INSTRUCTION`, reproduced
exactly:

{fence(EXTRACTION_INSTRUCTION)}

## 2. Corpus prompt: how the user message was composed

`SOURCE-DERIVED` — `harness.py` line {COMPOSE_LINE}, inside `run_corpus`:

{fence(COMPOSE_EXPR, 'python')}

`SOURCE-DERIVED` — the request builder then sends a `messages` array. For corpus
specifications no system message is defined, so the array holds one `user`
message: the composed string above. No corpus record stores a `system` field.

## 3. Corpus prompt: the case text

`SOURCE-DERIVED` — `corpus.json`. Each case's `conversation` string is reproduced
in the packet for that case, inside a fence, exactly as stored. The composition
inserts it after the literal separator `\\n\\nConversation:\\n`.

## 4. Six-item prompt, verbatim from source

`SOURCE-DERIVED` — `experiments/lfm25-cpu/run_probe.py`, constant `PROMPT`,
imported by the round-3 harness rather than retyped. It is the complete user
message for every six-item run:

{fence(TRIAGE_PROMPT)}

Two configurations additionally sent a system message. The value is stored per
run and is reproduced in the packet that carries the run: `smollm3-3b-think-rb1024`
sent `/think`; `smollm3-3b-nonthink` sent `/no_think`. No other six-item run
records a system message.

## 5. What is recovered, and what is not

`SOURCE-DERIVED`, verified against the files:

| Item | Value |
|---|---|
| `harness.py` SHA-256 | `{harness_hash}` |
| `corpus.json` SHA-256 | `{corpus_hash}` |
| `run_probe.py` SHA-256 | `{probe_hash}` |
| `harness.py` modification time | 2026-11-18 22:51 |
| `corpus.json` modification time | 2026-11-18 22:56 |
| Earliest stored generation | `results/stageA.json`, 2026-11-18 23:00 |

A hash taken now proves the file's current content. It does not prove the file
was unchanged when the requests were sent; the modification times are consistent
with that, and nothing stronger is claimed.

`UNESTABLISHED` for every run:

1. That the composed string equals the bytes the server received. No request
   body was logged.
2. That the chat template rendered `messages` identically for each model family.
   Templates are applied server-side and were not dumped.
3. That the instruction text was the same at the moment of each request. Only one
   revision exists in the workspace, and its modification time precedes the runs.

What is recorded per run and could be used to test a reconstruction later:
`prompt_tokens` ({sum(1 for i in RUNS.values() if i["rec"].get("prompt_tokens") is not None)}
of {len(RUNS)} runs have a value). Comparing it to a re-tokenisation requires each
model's tokenizer; that cross-check was not performed here and is not claimed.

## 6. Consequence for the packets

Packets 02A–02R print, for each case, the full reconstructed user message once,
labelled `SOURCE-DERIVED` reconstruction. Packets 03A–03M do not repeat the
six-item prompt per run; they point here, because the prompt is a single constant
and repetition would add no evidence.
"""
    return {"id": "01", "title": "The prompt record", "body": body, "runs": [],
            "sources": "harness.py, run_probe.py, corpus.json",
            "evidence": "instruction text, composition expression, six-item prompt"}


# ---------------------------------------------------------------------------
# Packets 02 — corpus by case
# ---------------------------------------------------------------------------

def corpus_packets() -> list[dict]:
    by_case = defaultdict(list)
    for code, info in RUNS.items():
        rec = info["rec"]
        if rec.get("kind") == "corpus" and rec.get("case"):
            by_case[rec["case"]].append(code)

    packets = []
    order = sorted(by_case, key=lambda c: int(c[1:]))
    for i, cid in enumerate(order):
        suffix = chr(65 + i)
        case = CASES[cid]
        given = corpus_prompt(cid)
        intro = f"""## Case {cid} — `{case.get('category')}`, {case.get('length')}

### What the model was given

`SOURCE-DERIVED` reconstruction, not a recorded request body. It applies the
composition expression in Packet 01 §2 to the stored case text. The instruction
is Packet 01 §1.

{fence(given)}

### Stored responses for this case

"""
        blocks = [(code, run_block(code, RUNS[code]))
                  for code in sorted(by_case[cid],
                                     key=lambda c: (RUNS[c]["rec"].get("label") or "",
                                                    RUNS[c]["rec"].get("draw") or 0,
                                                    c))]

        # greedy split at run boundaries
        groups, cur, size = [], [], words(intro)
        for code, b in blocks:
            w = words(b)
            if cur and size + w > LIMIT:
                groups.append(cur)
                cur, size = [], 0
            cur.append((code, b))
            size += w
        if cur:
            groups.append(cur)

        for j, group in enumerate(groups, 1):
            pid = f"02{suffix}" + (f"-{j}" if len(groups) > 1 else "")
            title = (f"Corpus case {cid}: responses"
                     + (f" (part {j} of {len(groups)})" if len(groups) > 1 else ""))
            count = (f"{len(group)} of {len(by_case[cid])} stored runs for this case, "
                     "sorted by stored model identifier and then draw. Every stored run "
                     "for the case appears in this packet series.\n")
            body = (packet_head(
                pid, title,
                f"the reconstructed prompt for one case and {len(group)} stored responses",
                "`corpus.json`, `results/*.json`", f"{cid}")
                + "\n" + intro + count)
            body += "\n" + "\n".join(b for _, b in group)
            packets.append({"id": pid, "title": title, "body": body,
                            "runs": [c for c, _ in group],
                            "sources": "corpus.json, results/*.json",
                            "evidence": f"{cid} prompt reconstruction and responses"})
    return packets


# ---------------------------------------------------------------------------
# Packets 03 — six-item prompt by configuration
# ---------------------------------------------------------------------------

def triage_packets() -> list[dict]:
    by_label = defaultdict(list)
    for code, info in RUNS.items():
        rec = info["rec"]
        if rec.get("kind") == "triage":
            by_label[rec.get("label")].append(code)

    # Flatten to one item per run so a long reasoning block can be split off
    # without dragging its whole configuration into an oversized packet.
    items = []  # (label, code, block, opens_label)
    for label in sorted(by_label):
        codes = sorted(by_label[label], key=lambda c: (RUNS[c]["rec"].get("draw") or 0, c))
        sysmsg = {RUNS[c]["rec"].get("system") for c in codes}
        sysmsg.discard(None)
        for k, code in enumerate(codes):
            text = run_block(code, RUNS[code])
            if k == 0:
                head = f"## `{label}` — {len(codes)} stored run(s)\n"
                if sysmsg:
                    head += ("\nRecorded system message: "
                             + ", ".join(f"`{s}`" for s in sorted(sysmsg)) + "\n")
                text = head + "\n" + text
            items.append((label, code, text, k == 0))

    groups, cur, size = [], [], 0
    for label, code, text, opens in items:
        w = words(text)
        if cur and size + w > LIMIT:
            groups.append(cur)
            cur, size = [], 0
        cur.append((label, code, text, opens))
        size += w
    if cur:
        groups.append(cur)

    out = []
    for i, group in enumerate(groups):
        suffix = chr(65 + i)
        labels = sorted({g[0] for g in group})
        cov = (f"`{labels[0]}`" if len(labels) == 1
               else f"`{labels[0]}` … `{labels[-1]}`")
        body = packet_head(
            f"03{suffix}", "Six-item prompt: responses by configuration",
            f"{len(group)} stored response(s) for {len(labels)} configuration(s)",
            "`results/all_triage.json`, `results/stageB.json`, "
            "`results/stageC.json`, `results/tinyllama.json`, "
            "`results/gpu-full*.json`, `results/budget.json`, "
            "`results/ablation_triage.json`",
            cov)
        body += ("\nThe user message is the constant reproduced in Packet 01 §4 and is "
                 "identical for every run below. Draws differ only in seed and, where "
                 "recorded, system message.\n\n")
        rendered = []
        prev = None
        for label, code, text, opens in group:
            if label != prev and not opens:
                text = f"## `{label}` (continued)\n\n" + text
            rendered.append(text)
            prev = label
        body += "\n".join(rendered)
        out.append({"id": f"03{suffix}",
                    "title": "Six-item prompt: responses by configuration",
                    "body": body,
                    "runs": [g[1] for g in group],
                    "sources": "results/*.json",
                    "evidence": "six-item prompt responses"})
    return out


# ---------------------------------------------------------------------------
# Reference indexes
# ---------------------------------------------------------------------------

def reference_metadata() -> tuple[str, dict]:
    per_label = {}
    for code, info in RUNS.items():
        rec = info["rec"]
        label = rec.get("label")
        d = per_label.setdefault(label, {"kinds": Counter(), "sampling": Counter(),
                                         "max_tokens": Counter(), "budget": Counter(),
                                         "ngl": Counter(), "ctx": Counter(),
                                         "system": Counter(), "command": 0,
                                         "n": 0, "logs": set()})
        d["n"] += 1
        d["kinds"][rec.get("kind")] += 1
        d["sampling"][json.dumps(rec.get("sampling"), sort_keys=True)] += 1
        d["max_tokens"][rec.get("max_tokens")] += 1
        d["budget"][rec.get("reasoning_budget")] += 1
        d["ngl"][rec.get("ngl")] += 1
        d["ctx"][rec.get("ctx")] += 1
        d["system"][rec.get("system")] += 1
        if rec.get("command"):
            d["command"] += 1
        tag = rec.get("tag")
        if tag:
            d["logs"].add(f"raw/{tag}-{label}/server.log")
        elif rec.get("kind") == "corpus":
            d["logs"].add(f"raw/corpus-{label}/server.log")

    lines = ["""# Reference index — execution metadata

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19

This file is the separate reference index required by the evidence-preparation
brief. It exists so that sampler parameters, context settings, token counts,
timings and hardware measurements stay out of the reading packets while
remaining independently inspectable.

Nothing here is evidence that a configuration was applied. A recorded value is a
statement about what the harness intended to send and what the server reported;
it is not a measurement of the server's effective sampler chain.

## 1. Where each value is stored

| Value | JSON key | Present in |
|---|---|---|
| Sampler chain | `sampling` (object) | all generation records |
| Output cap | `max_tokens` | all generation records |
| Reasoning budget | `reasoning_budget` | all generation records (may be null) |
| Reasoning effort / flag | `reasoning_effort`, `reasoning_flag` | six-item records |
| Context size | `ctx` | six-item records only |
| GPU layers | `ngl` | all generation records |
| Seed | `seed` | all generation records |
| Server launch command | `command` | six-item records only |
| Prompt / completion tokens | `prompt_tokens`, `completion_tokens` | all generation records |
| Timings | `ttft_s`, `thinking_s`, `answer_s`, `wall_s`, `load_s` | all generation records |
| Throughput | `prefill_ms`, `decode_ms`, `decode_tok_s`, `prefill_tok_s` | six-item records |
| Peak memory | `peak_rss_mb`, `peak_vram_mb`, `peak_gpu_util` | all generation records |
| Server log | `raw/<tag>-<label>/server.log` | files on disk |

`OBSERVED` gaps in the records themselves:

- Corpus records store no `ctx` and no `command`. Context size is therefore not
  recorded per corpus run; the corpus launch logs under `raw/corpus-<label>/`
  contain the server banners and are the place to check it.
- Corpus records store no `tag`; the batch is recoverable only from which result
  file holds the record (see `INDEX.md`).
- Timings are wall-clock around the HTTP request and exclude model load, which is
  stored separately in `load_s` where recorded.

## 2. Per-configuration recorded values

One row per stored model identifier. Where a value varies within a label, every
distinct value is listed with a count.
"""]
    for label in sorted(per_label):
        d = per_label[label]
        def vals(counter):
            return ", ".join(f"`{k}`×{v}" for k, v in sorted(
                counter.items(), key=lambda kv: str(kv[0])))
        lines.append(f"### `{label}`\n")
        lines.append(f"* Runs: {d['n']} ({', '.join(f'{k} {v}' for k, v in d['kinds'].items())})")
        lines.append(f"* Sampling: {vals(d['sampling'])}")
        lines.append(f"* Max tokens: {vals(d['max_tokens'])}")
        lines.append(f"* Reasoning budget: {vals(d['budget'])}")
        lines.append(f"* `ngl`: {vals(d['ngl'])}")
        lines.append(f"* `ctx`: {vals(d['ctx'])}")
        lines.append(f"* System message: {vals(d['system'])}")
        lines.append(f"* `command` recorded: {d['command']} of {d['n']} runs")
        if d["logs"]:
            lines.append("* Server log(s): " + ", ".join(f"`{p}`" for p in sorted(d["logs"])))
        lines.append("")

    lines.append("## 3. Per-run values\n")
    lines.append("| run_code | label | case | draw | seed | finish | prompt_tok | compl_tok | cap | budget | wall_s | rss_mb | vram_mb | source |")
    lines.append("|---|---|---|---:|---:|---|---:|---:|---:|---:|---:|---:|---:|---|")
    for code in sorted(RUNS):
        r = RUNS[code]["rec"]
        lines.append("| `{c}` | `{l}` | {case} | {d} | {s} | `{f}` | {pt} | {ct} | {cap} | {b} | {w} | {rss} | {vram} | `{src}`[{idx}] |".format(
            c=code, l=r.get("label"), case=r.get("case") or "—", d=r.get("draw"),
            s=r.get("seed"), f=r.get("finish_reason"), pt=r.get("prompt_tokens"),
            ct=r.get("completion_tokens"), cap=r.get("max_tokens"),
            b=r.get("reasoning_budget"), w=r.get("wall_s"),
            rss=r.get("peak_rss_mb"), vram=r.get("peak_vram_mb"),
            src=RUNS[code]["src"], idx=RUNS[code]["idx"]))
    lines.append("""
## 4. What these values do not show

- Whether the server applied the requested sampler chain. llama.cpp accepts the
  values; no echo of the effective chain is stored.
- Whether `max_tokens`, `reasoning_budget` or `ctx` bound a particular run. The
  recorded `finish_reason` and `completion_tokens` are the only direct indicators.
- Whether measured memory belongs to the model alone. `peak_vram_mb` is
  `nvidia-smi memory.used`, which includes other consumers on the desktop.
""")
    stats = {"labels": len(per_label)}
    return "\n".join(lines), stats


def reference_identity() -> str:
    manifest = json.loads((RESULTS / "model_manifest.json").read_bytes().decode("utf-8"))
    by_file = {m["file"]: m for m in manifest}
    per_file = defaultdict(lambda: {"labels": set(), "kinds": set(),
                                    "sizes": set(), "presets": set(), "notes": set()})
    for info in RUNS.values():
        r = info["rec"]
        d = per_file[r.get("file")]
        d["labels"].add(r.get("label"))
        d["kinds"].add(r.get("kind"))
        if r.get("file_mb"):
            d["sizes"].add(r["file_mb"])
        if r.get("preset"):
            d["presets"].add(r["preset"])
        if r.get("notes"):
            d["notes"].add(r["notes"])

    lines = ["""# Reference index — model identity

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19

Every reading packet shows the stored identifier and the model file verbatim.
This index adds the recorded file size, the SHA-256 from
`results/model_manifest.json`, and the repository string where one was recorded.

`OBSERVED` limitations:

- Quantisation is not a separate field. It is visible only inside the file name
  (`Q4_K_M`, `QAD-Q4_0`) and is reproduced there, not normalised.
- Repository identifiers were recorded only for ablation conditions, in the
  `notes` field. For baseline conditions the repository is not stored with the
  run.
- `model_manifest.json` was written at 02:12 on 2026-11-19, after the runs it
  describes. A hash taken afterwards cannot prove the file was unchanged during
  the experiment; it identifies the file present when the manifest was written.
"""]
    lines.append("| Model file | Identifier(s) | Kind(s) | Size on disk (bytes) | SHA-256 | Repository or preset string recorded |")
    lines.append("|---|---|---|---:|---|---|")
    for f in sorted(per_file):
        d = per_file[f]
        m = by_file.get(f)
        size = m["bytes"] if m else "not in manifest"
        sha = f"`{m['sha256']}`" if m else "not in manifest"
        repo = "; ".join(sorted(d["presets"] | d["notes"])) or "not recorded"
        lines.append("| `{}` | {} | {} | {} | {} | {} |".format(
            f, ", ".join(f"`{x}`" for x in sorted(d["labels"])),
            ", ".join(sorted(d["kinds"])), size, sha,
            repo.replace("|", "/")))
    lines.append("""
## Labels used more than once

`qwen35-2b-nonthink` is stored as the identifier for both six-item runs and
corpus runs, in different files and with different output caps. The identifier
alone is therefore not unique across prompt types; `kind` and the source file
disambiguate. The same is true of `abl-2b-nonthink`, `abl-4b-nonthink` and
`abl-lfm120`, which appear in both the corpus and six-item sets.
""")
    return "\n".join(lines)


# ---------------------------------------------------------------------------
# Index and statistics
# ---------------------------------------------------------------------------

def make_index(packets: list[dict]) -> str:
    run_packet = {}
    for p in packets:
        for code in p["runs"]:
            run_packet[code] = p["id"]

    lines = ["""# Master index — every stored run, its packet, its source

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19

One row per unique `run_code`. `Source` is the file and 0-based array index of
the copy used in the packet; `Other copies` lists further result files that hold
the same generation (compared for agreement; see `REPORT.md` §4).

`Kind` is `corpus` (extraction instruction plus one case conversation) or
`triage` (the six-item prompt).
"""]
    corpus_pkts = [p for p in packets if p["id"].startswith("02")]
    triage_pkts = [p for p in packets if p["id"].startswith("03")]
    lines.append("## Packet summary\n")
    lines.append("| Packet | Title | Runs | Primary sources |")
    lines.append("|---|---|---:|---|")
    for p in packets:
        lines.append(f"| {p['id']} | {p['title']} | {len(p['runs'])} | {p['sources']} |")
    lines.append(f"\n**Total:** {len(packets)} packets covering {len(run_packet)} unique runs "
                 f"({len(corpus_pkts)} corpus-case packets, {len(triage_pkts)} six-item packets).\n")

    lines.append("## Every run\n")
    lines.append("| run_code | kind | case | label | draw | packet | source | other copies |")
    lines.append("|---|---|---|---|---:|---|---|---|")
    for code in sorted(RUNS, key=lambda c: (RUNS[c]["rec"].get("kind") or "",
                                            RUNS[c]["rec"].get("case") or "",
                                            RUNS[c]["rec"].get("label") or "",
                                            RUNS[c]["rec"].get("draw") or 0, c)):
        r = RUNS[code]["rec"]
        info = RUNS[code]
        lines.append("| `{}` | {} | {} | `{}` | {} | {} | `{}`[{}] | {} |".format(
            code, r.get("kind") or "—", r.get("case") or "—", r.get("label"),
            r.get("draw"), run_packet.get(code, "MISSING"), info["src"], info["idx"],
            ", ".join(f"`{c}`" for c in info["copies"]) or "—"))
    return "\n".join(lines)


def main() -> int:
    packets = [packet_00(), packet_01()]
    packets += corpus_packets()
    packets += triage_packets()
    packets.sort(key=lambda p: p["id"])

    written = []
    for p in packets:
        name = f"PACKET-{p['id']}-{re.sub(r'[^A-Za-z0-9]+', '-', p['title'])[:60].strip('-')}.md"
        path = HERE / name
        path.write_bytes(p["body"].encode("utf-8"))
        written.append((p, path))

    (HERE / "INDEX.md").write_bytes(make_index(packets).encode("utf-8"))
    meta, _ = reference_metadata()
    (HERE / "REFERENCE-EXECUTION-METADATA.md").write_bytes(meta.encode("utf-8"))
    (HERE / "REFERENCE-MODEL-IDENTITY.md").write_bytes(reference_identity().encode("utf-8"))

    # ---- verbatim assertions -------------------------------------------
    missing, missing_reasoning = [], []
    for code, info in RUNS.items():
        rec = info["rec"]
        pkt = next((path for p, path in written if code in p["runs"]), None)
        if pkt is None:
            missing.append((code, "no packet"))
            continue
        blob = pkt.read_bytes()
        content = rec.get("content") or ""
        if content.encode("utf-8") not in blob:
            missing.append((code, "content"))
        reasoning = rec.get("reasoning_content") or ""
        if reasoning and reasoning.encode("utf-8") not in blob:
            missing_reasoning.append(code)

    over = [(p["id"], words(p["body"])) for p, _ in written
            if words(p["body"]) > HARD_LIMIT]

    bad_attr = []
    for p, path in written:
        blob = path.read_bytes()
        for code in p["runs"]:
            if f"`{code}`".encode() not in blob:
                bad_attr.append((p["id"], code))
    total_words = sum(words(p["body"]) for p, _ in written)

    # ---- scan for excluded material outside fences ---------------------
    banned = ["v1 score", "v2 `", "MALFORMED", "VALID_KEEP", "VALID_NONE",
              "VALID_HOLD", "catastrophic", "FORBIDDEN_CLAIM", "event code"]
    hits = []
    for p, path in written:
        if not p["id"].startswith(("02", "03")):
            continue  # the guide and prompt record name the exclusions
        text = strip_fences(p["body"])
        low = text.lower()
        for token in banned:
            if token.lower() in low:
                hits.append((p["id"], token))

    stats = {
        "runs": len(RUNS),
        "corpus_runs": sum(1 for i in RUNS.values() if i["rec"].get("kind") == "corpus"),
        "triage_runs": sum(1 for i in RUNS.values() if i["rec"].get("kind") == "triage"),
        "labels": len({i["rec"].get("label") for i in RUNS.values()}),
        "files": len({i["rec"].get("file") for i in RUNS.values()}),
        "packets": len(written),
        "words": total_words,
        "missing_content": missing,
        "missing_reasoning": missing_reasoning,
        "oversize": over,
        "bad_attribution": bad_attr,
        "banned_hits": hits,
        "disagree": DISAGREE,
        "duplicate_codes": sum(1 for i in RUNS.values() if i["n_copies"] > 1),
        "max_copies": max(i["n_copies"] for i in RUNS.values()),
        "empty_content": [c for c, i in RUNS.items() if i["rec"].get("content") == ""],
        "empty_reasoning": sum(1 for i in RUNS.values() if i["rec"].get("reasoning_content") == ""),
        "absent_reasoning_field": sum(1 for i in RUNS.values() if "reasoning_content" not in i["rec"]),
        "finish": dict(Counter(i["rec"].get("finish_reason") for i in RUNS.values())),
        "length_by_label": dict(Counter(i["rec"].get("label") for i in RUNS.values()
                                     if i["rec"].get("finish_reason") == "length")),
        "labels_both_kinds": sorted(l for l in {i["rec"].get("label") for i in RUNS.values()}
                                 if len({i["rec"].get("kind") for i in RUNS.values()
                                         if i["rec"].get("label") == l}) > 1),
        "no_prompt_tokens": [c for c, i in RUNS.items() if i["rec"].get("prompt_tokens") is None],
        "no_ctx": sum(1 for i in RUNS.values() if "ctx" not in i["rec"]),
        "no_command": sum(1 for i in RUNS.values() if "command" not in i["rec"]),
        "system": dict(Counter(str(i["rec"].get("system")) for i in RUNS.values()
                               if i["rec"].get("system"))),
        "packet_words": [{"id": p["id"], "runs": len(p["runs"]), "words": words(p["body"])}
                         for p, _ in written],
    }
    (HERE / "build_stats.json").write_bytes(
        json.dumps(stats, indent=2, default=str).encode("utf-8"))

    print(json.dumps({k: v for k, v in stats.items() if k != "packet_words"}, indent=2)[:4000])
    print("packets:", len(written), "words:", total_words)
    print("missing content:", len(missing), "missing reasoning:", len(missing_reasoning))
    print("oversize:", over)
    print("banned hits:", hits)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
