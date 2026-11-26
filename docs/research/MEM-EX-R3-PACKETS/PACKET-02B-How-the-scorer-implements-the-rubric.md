# Packet 02B — How the scorer implements the rubric

**Series:** Round 3 memory-extraction evidence, source-preserving reading packet
**Prepared:** 2026-11-19 · **For:** independent review (ChatGPT)
**Original sources:** `score_corpus.py` (295 lines) and `score_triage.py` (106 lines), excerpts
**Case / run identifiers covered:** scoring functions, not case-specific
**Raw model output in this packet:** none — no model output
**Omitted:** The excerpts below are the parser, the fidelity/false-memory block and the
NONE early return. The remainder of `score_corpus.py` is item matching,
significance, temporal, attribution, restraint, completeness and structure;
it is cited by function name rather than reproduced to keep this packet
readable. The full file is 295 lines and is the authoritative source.

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

## 1. The parser

`SOURCE-DERIVED`, `score_corpus.py::parse_lines`. A line is a record only if it
splits into at least three fields on `|` or `;` and the second field is exactly
`KEEP`, `REJECT` or `HOLD`. Anything else becomes `{"parsed": False}`.

Key consequence, and the root of audit finding F-02: a table whose second column
is the literal word `DECISION` is not a record, so it contributes no `KEEP`. On a
case with no required items, "no KEEP" was scored as restraint.



---
```text
def parse_lines(text: str) -> list[dict]:
    rows = []
    for raw in text.splitlines():
        line = raw.strip().strip("*#` ")
        if not line or line.upper().startswith("NONE"):
            continue
        parts = [p.strip() for p in re.split(r"\s*[|;]\s*", line)]
        if len(parts) < 3:
            rows.append({"raw": line, "parsed": False})
            continue
        cid = parts[0]
        dec = parts[1].upper().strip("[] ")
        if dec not in ("KEEP", "REJECT", "HOLD"):
            rows.append({"raw": line, "parsed": False})
            continue
        rows.append({"raw": line, "parsed": True, "id": cid, "decision": dec,
                     "memory": parts[2], "reason": " | ".join(parts[3:])})
    return rows


def norm(s: str) -> str:
    return re.sub(r"\s+", " ", s.lower())


def _words(s):
    return re.findall(r"[a-z0-9']+", s.lower())


def _alt_matches(alt, words):
    """Prefix match tolerant of plurals/participles ('narrating' vs 'narrate').
    Phrases require every word, in any position."""
    parts = _words(alt)
    if not parts:
        return False
    for p in parts:
        if p == "no" or p == "not":
            if not any(w in ("no", "not", "never") for w in words):
                return False
            continue
        stem = p[: max(4, len(p) - 2)]
        if not any(w.startswith(stem) for w in words):
            return False
    return True


def item_found(item, rows):
    for r in rows:
        if not r.get("parsed") or r["decision"] != item["decision"]:
            continue
        blob = norm(r["memory"] + " " + r["reason"])
        words = _words(blob)
        if all(any(_alt_matches(alt, words) for alt in g) for g in item["keys"]):
            return r
    return None



```
## 2. Fidelity and false memory

`SOURCE-DERIVED`. A forbidden string inside a `KEEP` or `HOLD` memory zeroes
fidelity (20 points) **and** subtracts the entire false-memory dimension
(15 points) **and** appends a catastrophic reason. That is three consequences
from one string, which is audit finding F-04.

```text
    # 2 fidelity /20 — a forbidden claim inside a kept line is a preserved
    # falsehood and earns no fidelity credit for that case.
    fidelity = 20.0
    for bad in case.get("forbidden", []):
        if any(bad.lower() in norm(r["memory"]) for r in kept):
            fidelity = 0.0
            catastrophes.append(f"kept forbidden claim: {bad}")
        elif bad.lower() in blob_all:
            fidelity -= 6
    fidelity = max(0.0, fidelity)

    # 3 significance /15
    if n_req == 0:
        significance = 15.0 if not [r for r in parsed if r["decision"] == "KEEP"] else 0.0
    else:
        pts = 0.0
        for it in req:
            r = found[it["id"]]
            if r:
                pts += 3
                if any(g in norm(r["reason"]) for g in GROUNDS[it["decision"]]):
                    pts += 2
        significance = 15 * pts / (5 * n_req)

    # 4 false memory /15
    false_mem = 15.0
    for bad in case.get("forbidden", []):
        if any(bad.lower() in norm(r["memory"]) for r in kept):
            false_mem -= 15
    false_mem = max(0.0, false_mem)


```
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
