"""Legacy 12-point scoring of the original six-item triage prompt.

Same rule as rounds 1-2: per item 2 = expected decision AND a reason naming the
right ground; 1 = expected decision with absent/vague reason; 0 = wrong
decision. Max 12. Format compliance and invented content are recorded but not
scored into the 12. Non-blind: round-1/2 scores were already published per
model.
"""
import json, re, sys
from pathlib import Path

KEY = {"M1": ("KEEP", ("canon", "durable", "identity", "training", "trained", "fact", "stated")),
       "M2": ("REJECT", ("transient", "chatter", "small talk", "temporary", "brb", "not durable", "session", "no value")),
       "M3": ("KEEP", ("durable", "trait", "character", "fact", "identity", "cannot swim", "mage")),
       "M4": ("HOLD", ("contradict", "conflict", "inconsisten", "cannot both", "needs confirm", "unresolved", "disagree")),
       "M5": ("KEEP", ("preference", "durable", "style", "affect", "future", "stable", "reply")),
       "M6": ("HOLD", ("uncertain", "maybe", "confirm", "unverified", "not sure", "provisional", "canon status"))}

def parse(text):
    rows = {}
    for line in text.splitlines():
        line = line.strip().strip("*#` ")
        if not line or line.upper() == "NONE":
            continue
        m = re.match(r"^\[?\s*(M[1-6])\s*\]?\s*[|:\-–]\s*(?:DECISION\s*[|:\-]\s*)?(KEEP|REJECT|HOLD)\b(.*)$", line, re.I)
        if m:
            rows.setdefault(m.group(1).upper(), []).append((m.group(2).upper(), m.group(3).lower()))
    out = {}
    for k, v in rows.items():
        out[k] = v[0]  # first line for an item counts
    return out, rows

def strict_lines(text):
    """Only lines that are a clean M#/decision row, not leaked working notes."""
    out = []
    for line in text.splitlines():
        s = line.strip()
        if s.startswith(('*', '-', '•')):
            continue
        if re.match(r"^M[1-6]\s*\|\s*(KEEP|REJECT|HOLD)\s*\|", s, re.I):
            out.append(s)
    return out


def score(rec):
    text = rec.get("content") or ""
    parsed, all_rows = parse(text)
    strict = strict_lines(text)
    # Leaked working notes look like bullet rows that carry 'Analysis:' or
    # 'Decision:'. A bare word like 'wait' in a reason is not evidence.
    leaked = any(line.strip().startswith(("*", "-", "•")) and
                 ("analysis:" in line.lower() or "decision:" in line.lower())
                 for line in text.splitlines())
    truncated = rec.get("finish_reason") == "length"
    # Admissible = the model stopped on its own and produced at least one
    # parseable verdict. Incomplete answers keep their partial score (missing
    # items are 0); only cap-truncated or reasoning-leaked runs are excluded.
    admissible = (not truncated) and len(parsed) >= 1 and not leaked
    total = 0
    detail = {}
    for item, (want, grounds) in KEY.items():
        got = parsed.get(item)
        if got is None:
            detail[item] = {"decision": None, "points": 0}
            continue
        dec, reason = got
        if dec != want:
            pts = 0
        else:
            pts = 2 if any(g in reason for g in grounds) else 1
        total += pts
        detail[item] = {"decision": dec, "points": pts, "reason": reason[:160]}
    lines = [l for l in text.splitlines() if re.search(r"M[1-6]", l)]
    invented = []
    for bad in ("sable's mother", "corvin is her father", "sable is a water mage", "sable has a sister", "sable is 30"):
        if bad in text.lower():
            invented.append(bad)
    return {"score12": total, "items": detail,
            "item_lines": len(lines), "strict_lines": len(strict),
            "format_ok": len(lines) == 6 and len(parsed) == 6,
            "truncated": truncated, "reasoning_leaked_into_content": leaked,
            "admissible": admissible,
            "invented": invented,
            "empty": not text.strip(),
            "finish_reason": rec.get("finish_reason"),
            "completion_tokens": rec.get("completion_tokens"),
            "wall_s": rec.get("wall_s"), "reasoning_chars": len(rec.get("reasoning_content") or ""),
            "run_code": rec.get("run_code"), "label": rec.get("label")}

if __name__ == "__main__":
    src = Path(sys.argv[1]); dst = Path(sys.argv[2])
    runs = json.loads(src.read_text(encoding="utf-8"))
    rows = [score(r) for r in runs if r.get("kind") == "triage" and "error" not in r]
    dst.write_text(json.dumps(rows, indent=2), encoding="utf-8")
    import statistics as st
    agg, bad = {}, {}
    for r in rows:
        if r["admissible"]:
            agg.setdefault(r["label"], []).append(r["score12"])
        else:
            bad.setdefault(r["label"], []).append(
                (r["score12"], "truncated" if r["truncated"] else "leaked/malformed"))
    for k, v in sorted(agg.items(), key=lambda kv: -st.mean(kv[1])):
        print(f"{k:28s} mean={st.mean(v):5.2f}  draws={v}")
    for k, v in bad.items():
        print(f"{k:28s} INADMISSIBLE {v}")
