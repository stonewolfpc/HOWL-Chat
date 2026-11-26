"""Corpus scoring for round 3, implementing RUBRIC.md.

Blind: reads only run_code / case / output. Never reads `label` for scoring.
Judgment calls the regex layer cannot make are left to the manual audit noted in
the report; nothing here silently upgrades a run.
"""
from __future__ import annotations

import json
import re
import statistics as st
import sys
from pathlib import Path

GROUNDS = {
    "KEEP": ("durable", "permanent", "canon", "future", "matter", "world rule",
             "identity", "stable", "long-term", "long term", "preference", "fact"),
    "REJECT": ("transient", "chatter", "temporary", "not durable", "no value",
               "session", "small talk", "irrelevant", "noise", "throwaway"),
    "HOLD": ("uncertain", "confirm", "contradict", "conflict", "inconsisten",
             "unresolved", "maybe", "unverified", "ambiguous", "needs"),
}


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


def score_case(case, text):
    if not text.strip():
        return {"case": case["id"], "category": case["category"], "score": 0.0,
                "parts": {}, "decisions": [], "decision_vec": (),
                "empty": True, "clean_none": False, "parsed_lines": 0,
                "raw_lines": 0, "catastrophic": False,
                "catastrophic_reasons": []}

    clean_none = text.strip().upper().strip(".") == "NONE"
    rows = parse_lines(text)
    parsed = [r for r in rows if r.get("parsed")]
    kept = [r for r in parsed if r["decision"] in ("KEEP", "HOLD")]
    req = case["required"]
    blob_all = norm(text)
    found = {it["id"]: item_found(it, parsed) for it in req}
    n_found = sum(1 for v in found.values() if v)
    n_req = len(req)
    catastrophes: list[str] = []

    # Silence is not an answer. On a case that has required items, a clean NONE
    # must not collect the dimensions that reward absence (fidelity, false
    # memory, uncertainty, attribution, restraint). It keeps structure only.
    if n_req > 0 and not parsed:
        return {"case": case["id"], "category": case["category"], "score": 2.0,
                "parts": {"structure": 2.0}, "decisions": [], "decision_vec": (),
                "empty": False, "clean_none": clean_none, "parsed_lines": 0,
                "raw_lines": len(rows), "catastrophic": False,
                "catastrophic_reasons": []}

    # 1 recall /20
    if n_req == 0:
        keeps = [r for r in parsed if r["decision"] == "KEEP"]
        recall = 20.0 if not keeps else 0.0
    else:
        recall = 20 * n_found / n_req

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

    # 5 temporal /10 (None = not applicable, excluded from the denominator)
    TEMPORAL = {"C05": ["C1", "C2", "C3"], "C06": ["C1", "C2"],
                "C07": ["C1", "C2", "C3"], "C10": ["C4"]}
    if case["id"] in TEMPORAL:
        ids = TEMPORAL[case["id"]]
        got = sum(1 for i in ids if found.get(i))
        temporal = 10.0 if got == len(ids) else (5.0 if got else 0.0)
        if case["id"] == "C05" and any("halden trained" in norm(r["memory"]) for r in kept):
            catastrophes.append("obsolete claim kept as current (Halden trained her)")
            temporal = 0.0
    else:
        temporal = None

    # 6 uncertainty /5
    hold_items = [it for it in req if it["decision"] == "HOLD"]
    if hold_items:
        ok = sum(1 for it in hold_items if found[it["id"]])
        unc = ("uncertain", "maybe", "confirm", "unverified", "not sure",
               "provisional", "unsure", "tentative", "speculat")
        uncertainty = 5.0 if ok == len(hold_items) else (2.0 if any(w in blob_all for w in unc) else 0.0)
    else:
        uncertainty = 5.0

    # 7 attribution /5
    attribution = 5.0
    if case["id"] == "C14":
        for r in parsed:
            b = norm(r["memory"])
            if "corvin" in b and "coast" in b and "three days" not in b and "tide" not in b:
                attribution = 0.0
                catastrophes.append("attributed Sable's proposal to Corvin")
            if "sable" in b and ("tide gate" in b or "gates shut" in b):
                attribution = 0.0
                catastrophes.append("attributed Corvin's warning to Sable")
        if attribution and not all(found.values()):
            attribution = 2.0
    elif case["id"] == "C13":
        for r in parsed:
            b = norm(r["memory"])
            if "sable" in b and "nurse" in b:
                attribution = 0.0
                catastrophes.append("real-user fact attributed to the character")
            if "user" in b and "healer" in b and "sable" not in b:
                attribution = 0.0
                catastrophes.append("character fact attributed to the real user")
        if attribution and not all(found.values()):
            attribution = 2.0

    # 8 restraint /5
    if case["id"] in ("C09", "C12"):
        restraint = 5.0 if not [r for r in parsed if r["decision"] == "KEEP"] else 0.0
    else:
        over = 0
        for r in parsed:
            if r["decision"] == "KEEP" and not any(
                    all(any(a.lower() in norm(r["memory"]) for a in g) for g in it["keys"])
                    for it in req):
                if any(w in norm(r["memory"]) for w in
                       ("weather", "rain", "bread", "dog", "ferryman", "storm", "cat",
                        "lol", "brb", "keyboard")):
                    over += 1
        restraint = max(0.0, 5.0 - 2.0 * over)

    # 9 completeness /3
    if n_req == 0:
        completeness = 3.0 if restraint == 5.0 else 0.0
    else:
        completeness = 3.0 if n_found == n_req else (1.0 if n_found >= 0.7 * n_req else 0.0)

    # 10 structure /2
    if clean_none:
        structure = 2.0
    elif rows and all(r.get("parsed") for r in rows):
        structure = 2.0
    elif parsed:
        structure = 1.0
    else:
        structure = 0.0

    parts = {"recall": recall, "fidelity": fidelity, "significance": significance,
             "false_memory": false_mem, "temporal": temporal, "uncertainty": uncertainty,
             "attribution": attribution, "restraint": restraint,
             "completeness": completeness, "structure": structure}
    scored = [v for v in parts.values() if v is not None]
    total = 100 * sum(scored) / (100 - (10 if temporal is None else 0))

    decision_vec = tuple(sorted((r["id"], r["decision"]) for r in parsed))
    return {"case": case["id"], "category": case["category"], "score": round(total, 1),
            "parts": {k: (round(v, 1) if isinstance(v, float) else v) for k, v in parts.items()},
            "decisions": [f"{r['id']}:{r['decision']}" for r in parsed],
            "decision_vec": decision_vec, "empty": False, "clean_none": clean_none,
            "parsed_lines": len(parsed), "raw_lines": len(rows),
            "catastrophic": bool(catastrophes),
            "catastrophic_reasons": catastrophes}


def main(runs_path, corpus_path, out_path):
    runs = json.loads(Path(runs_path).read_text(encoding="utf-8"))
    corpus = {c["id"]: c for c in json.loads(Path(corpus_path).read_text(encoding="utf-8"))["cases"]}
    rows = []
    for r in runs:
        if r.get("kind") != "corpus" or "error" in r:
            continue
        sc = score_case(corpus[r["case"]], r.get("content") or "")
        sc.update({"run_code": r["run_code"], "draw": r.get("draw"),
                   "wall_s": r.get("wall_s"), "completion_tokens": r.get("completion_tokens"),
                   "finish_reason": r.get("finish_reason"), "thinking_s": r.get("thinking_s"),
                   "reasoning_chars": len(r.get("reasoning_content") or "")})
        rows.append(sc)
    Path(out_path).write_text(json.dumps(rows, indent=2), encoding="utf-8")
    return rows


def summarise(rows, code_map):
    by = {}
    for r in rows:
        cfg = code_map.get(r["run_code"], {}).get("label", "?")
        by.setdefault(cfg, []).append(r)
    out = []
    for cfg, rs in by.items():
        per_case = {}
        for r in rs:
            per_case.setdefault(r["case"], []).append(r)
        case_means = {c: st.mean(x["score"] for x in v) for c, v in per_case.items()}
        means = list(case_means.values())
        draw_sd = [st.pstdev([x["score"] for x in v]) for v in per_case.values() if len(v) > 1]
        identical = sum(1 for v in per_case.values()
                        if len({tuple(x["decisions"]) for x in v}) == 1)
        out.append({
            "label": cfg, "cases": len(per_case),
            "draws": round(len(rs) / max(1, len(per_case)), 1),
            "mean": round(st.mean(means), 2),
            "median": round(st.median(means), 2),
            "min": round(min(means), 2), "max": round(max(means), 2),
            "sd_cases": round(st.pstdev(means), 2) if len(means) > 1 else 0.0,
            "sd_draws_within_case": round(st.mean(draw_sd), 2) if draw_sd else 0.0,
            "identical_decisions_frac": round(identical / len(per_case), 3),
            "catastrophic_runs": sum(1 for r in rs if r["catastrophic"]),
            "empty_runs": sum(1 for r in rs if r["empty"]),
            "truncated_runs": sum(1 for r in rs if r["finish_reason"] == "length"),
            "cases_with_any_failure": len({r["case"] for r in rs if r["empty"] or r["catastrophic"]}),
        })
    return sorted(out, key=lambda x: -x["mean"])


if __name__ == "__main__":
    main(sys.argv[1], sys.argv[2], sys.argv[3])
    print("scored ->", sys.argv[3])
