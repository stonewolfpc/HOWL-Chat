"""Corrected scoring for instrument v2.

Two composites, never one
-------------------------
    semantic_score   what the model understood: content, fidelity, significance,
                     temporal handling, attribution, ownership/scope, uncertainty,
                     restraint.
    contract_score   whether the output is usable: parseability, four-field
                     schema, atomicity, termination, uniqueness.

A model that states the right fact in a broken format can therefore score high
on semantic and low on contract; the `production_usable` flag stays False. This
is the F-03 separation. It does not pretend malformed output is deployable.

One incident, one deduction (F-04)
----------------------------------
Every failure is an *event* with exactly one `primary_dimension`. Only that
dimension is reduced. Other affected dimensions are listed under
`secondary_impacts` for visibility and are not deducted a second time. Events
are always reported separately from both composites, so a high score cannot hide
an ownership or credential event.
"""

from __future__ import annotations

import re
import statistics as st
from dataclasses import dataclass, field

import atomicity
import matching
import outcomes
import ownership


WEIGHTS_REQUIRED = {
    "recall": 30.0,
    "fidelity": 18.0,
    "significance": 12.0,
    "temporal": 10.0,
    "attribution": 10.0,
    "ownership": 15.0,
    "uncertainty": 5.0,
}

WEIGHTS_RESTRAINT = {
    "restraint": 60.0,
    "fidelity": 25.0,
    "ownership": 15.0,
}

CONTRACT_WEIGHTS = {
    "parseability": 30.0,
    "schema": 25.0,
    "atomicity": 25.0,
    "termination": 10.0,
    "uniqueness": 10.0,
}

GROUNDS = {
    "KEEP": ("durable", "permanent", "canon", "future", "matter", "world rule",
             "identity", "stable", "long-term", "long term", "preference", "fact"),
    "REJECT": ("transient", "chatter", "temporary", "not durable", "no value",
               "session", "small talk", "irrelevant", "noise", "throwaway"),
    "HOLD": ("uncertain", "confirm", "contradict", "conflict", "inconsisten",
             "unresolved", "maybe", "unverified", "ambiguous", "needs"),
}


@dataclass
class Event:
    code: str
    primary_dimension: str | None
    detail: str
    severity: str = "high"
    secondary_impacts: list[str] = field(default_factory=list)
    policy_sensitive: bool = False
    raw: str = ""

    def as_dict(self) -> dict:
        return {
            "code": self.code,
            "severity": self.severity,
            "primary_dimension": self.primary_dimension,
            "secondary_impacts": self.secondary_impacts,
            "policy_sensitive": self.policy_sensitive,
            "detail": self.detail,
            "raw": self.raw,
        }


def _pct(parts: dict, weights: dict) -> float | None:
    num = 0.0
    den = 0.0
    for k, w in weights.items():
        v = parts.get(k)
        if v is None:
            continue
        num += v
        den += w
    return round(100.0 * num / den, 1) if den else None


def _ground_hit(record, decision: str) -> bool:
    t = record.text.lower()
    return any(g in t for g in GROUNDS.get(decision, ()))


_NEGATORS = ("not ", "never ", "no ", "rather than ", "instead of ",
             "without ", "n't ", "isn't ", "aren't ", "doesn't ", "don't ")


def forbidden_occurrences(text: str, phrase: str) -> list[dict]:
    """Every occurrence of a forbidden phrase, with negation context.

    'double the toll (not triple)' mentions the prohibited word in order to
    deny it. Counting that as a stored falsehood is a false positive, and one
    that the v1 scorer made on C18 (run jjf67xzt).
    """
    low, p = (text or "").lower(), phrase.lower()
    out, i = [], low.find(p)
    while i != -1:
        window = low[max(0, i - 24):i]
        out.append({"index": i, "negated": any(n in window for n in _NEGATORS)})
        i = low.find(p, i + max(1, len(p)))
    return out


def stored_forbidden(records: list, forbidden: list[str]) -> tuple[list[str], list[str]]:
    """Return (stored, negated_only) phrase lists for KEEP/HOLD rows.

    One incident, one event: callers emit a single event for the whole returned
    `stored` list rather than one per phrase.
    """
    stored: list[str] = []
    negated: list[str] = []
    for bad in forbidden:
        occs = []
        for r in records:
            if r.decision in ("KEEP", "HOLD"):
                occs.extend(forbidden_occurrences(r.text, bad))
        if not occs:
            continue
        if any(not o["negated"] for o in occs):
            stored.append(bad)
        else:
            negated.append(bad)
    return stored, negated


def score_run(record: dict, case: dict) -> dict:
    """Score one run against one case. Never raises on malformed input."""
    parsed = outcomes.classify_run(record)

    base = {
        "case": case.get("id"),
        "run_code": record.get("run_code"),
        "label": record.get("label"),
        "draw": record.get("draw"),
        "status": parsed.status.value,
        "finish_reason": record.get("finish_reason"),
        "truncated": parsed.truncated,
        "quality_eligible": parsed.quality_eligible,
        "notes": parsed.notes,
        "events": [],
        "semantic": None,
        "contract": None,
        "production_usable": False,
        "parts": {},
        "contract_parts": {},
        "atomicity": [],
        "ownership": {"scored": False},
        "decisions": parsed.decisions(),
    }

    if parsed.status == outcomes.Status.EXECUTION_FAILURE:
        base["events"] = [Event("EXECUTION_FAILURE", None,
                                "; ".join(parsed.notes), "critical").as_dict()]
        return base

    required = case.get("required") or []
    n_req = len(required)
    matches = matching.match_items(parsed.records, required) if n_req else {}
    atom = atomicity.judge_case(parsed.records, matches, case)
    own = ownership.check(parsed.records, case.get("ownership"))

    events: list[Event] = []
    dims: dict[str, float | None] = {}
    cdims: dict[str, float | None] = {}

    absence_ok = parsed.status in outcomes.VALID

    # ---- semantic -------------------------------------------------------
    if n_req:
        # Absence-based credit only: a format-invalid row that *does* carry the
        # fact keeps content credit (F-03). What it never gets is usability.
        if not absence_ok and not parsed.records:
            dims.update(recall=0.0, fidelity=0.0, significance=0.0,
                        temporal=0.0, attribution=0.0, ownership=0.0,
                        uncertainty=0.0)
        else:
            credit = sum(atom.credit_for(m["record_index"])
                         for m in matches.values())
            dims["recall"] = WEIGHTS_REQUIRED["recall"] * credit / n_req

            dims["fidelity"] = WEIGHTS_REQUIRED["fidelity"]
            stored, negated_only = stored_forbidden(
                parsed.records, case.get("forbidden") or [])
            if stored:
                dims["fidelity"] = 0.0
                holder = next((r for r in parsed.records
                               if r.decision in ("KEEP", "HOLD")), None)
                events.append(Event(
                    "FORBIDDEN_CLAIM_STORED", "fidelity",
                    "stored forbidden claim(s): " + ", ".join(sorted(stored)),
                    "critical", ["false_memory"], True,
                    holder.raw if holder else ""))
            else:
                joined = " ".join(r.text for r in parsed.records).lower()
                for bad in case.get("forbidden") or []:
                    occs = forbidden_occurrences(joined, bad)
                    if any(not o["negated"] for o in occs):
                        dims["fidelity"] = max(0.0, dims["fidelity"] - 4.0)

            if not absence_ok and not matches:
                dims["fidelity"] = 0.0

            sig = 0.0
            for item in required:
                m = matches.get(item["id"])
                if not m:
                    continue
                share = 1.0 / n_req
                if m["decision_correct"]:
                    sig += 0.6 * share
                    if _ground_hit(parsed.records[m["record_index"]],
                                   item.get("decision") or "KEEP"):
                        sig += 0.4 * share
            dims["significance"] = WEIGHTS_REQUIRED["significance"] * sig

            tids = case.get("temporal_ids")
            if tids:
                got = sum(1 for t in tids
                          if t in matches and matches[t]["decision_correct"])
                dims["temporal"] = WEIGHTS_REQUIRED["temporal"] * got / len(tids)
            else:
                dims["temporal"] = None

            checks = case.get("attribution_forbidden") or []
            if checks or case.get("ownership"):
                dims["attribution"] = WEIGHTS_REQUIRED["attribution"]
                for chk in checks:
                    toks = [t.lower() for t in chk.get("tokens") or []]
                    hit = next((r for r in parsed.records
                                if toks and all(t in r.text.lower() for t in toks)), None)
                    if hit:
                        dims["attribution"] = 0.0
                        events.append(Event(
                            "ATTRIBUTION_ERROR", "attribution",
                            chk.get("why", f"attribution collapsed: {toks}"),
                            "critical", [], False, hit.raw))
            else:
                dims["attribution"] = None

            hold_items = [i for i in required if i.get("decision") == "HOLD"]
            if hold_items:
                ok = sum(1 for i in hold_items
                         if i["id"] in matches and matches[i["id"]]["decision_correct"])
                dims["uncertainty"] = (WEIGHTS_REQUIRED["uncertainty"]
                                       * ok / len(hold_items))
            else:
                dims["uncertainty"] = None

            # Silence on a case that has required items is not restraint. It
            # earns no semantic credit at all (historical amendment 1, kept).
            if not parsed.records:
                for k in list(dims):
                    dims[k] = 0.0 if dims[k] is not None else None
    else:
        # A case whose correct answer is to store nothing.
        absent = parsed.status in outcomes.RESTRAINT_VALID and not parsed.kept
        dims["restraint"] = WEIGHTS_RESTRAINT["restraint"] if absent else 0.0
        if not absence_ok:
            dims["restraint"] = 0.0
        dims["fidelity"] = WEIGHTS_RESTRAINT["fidelity"] if absence_ok else 0.0
        stored, _negated = stored_forbidden(
            parsed.records, case.get("forbidden") or [])
        if stored:
            dims["fidelity"] = 0.0
            holder = next((r for r in parsed.records
                           if r.decision in ("KEEP", "HOLD")), None)
            events.append(Event(
                "FORBIDDEN_CLAIM_STORED", "fidelity",
                "stored forbidden claim(s): " + ", ".join(sorted(stored)),
                "critical", ["restraint"], True,
                holder.raw if holder else ""))

    # ---- ownership ------------------------------------------------------
    if own.scored:
        dims["ownership"] = (0.0 if not absence_ok else
                             max(0.0, 15.0 - 15.0 * min(1, own.critical_count)))
        for f in own.findings:
            events.append(Event(f.kind.value, "ownership", f.detail,
                                "critical", ["semantic"], False, f.raw))

    # ---- contract -------------------------------------------------------
    if parsed.status in (outcomes.Status.VALID_NONE, outcomes.Status.VALID_EMPTY):
        # An explicit empty declaration conforms by construction; there are no
        # rows to check for schema, atomicity or uniqueness.
        cdims.update({"parseability": CONTRACT_WEIGHTS["parseability"],
                      "schema": None, "atomicity": None,
                      "uniqueness": None,
                      "termination": 0.0 if parsed.truncated
                      else CONTRACT_WEIGHTS["termination"]})
    else:
        cdims["parseability"] = (CONTRACT_WEIGHTS["parseability"]
                                 if parsed.status in outcomes.VALID else 0.0)
        rows = parsed.records
        if rows:
            cdims["schema"] = CONTRACT_WEIGHTS["schema"] * (
                sum(1 for r in rows if r.contract_valid) / len(rows))
        else:
            cdims["schema"] = 0.0
        if n_req:
            cdims["atomicity"] = CONTRACT_WEIGHTS["atomicity"] * atom.structure_pass_rate
        else:
            cdims["atomicity"] = None
        cdims["termination"] = (0.0 if (parsed.truncated or not absence_ok)
                                else CONTRACT_WEIGHTS["termination"])
        seen = [re.sub(r"\s+", " ", r.memory.lower()).strip()
                for r in rows if r.memory]
        cdims["uniqueness"] = (CONTRACT_WEIGHTS["uniqueness"]
                               if seen and len(seen) == len(set(seen)) else 0.0)

    for v in atom.violations:
        events.append(Event("ATOMICITY_" + v["verdict"], "atomicity",
                            v["evidence"], "medium", ["recall"], False, v["raw"]))

    # ---- composites -----------------------------------------------------
    weights = WEIGHTS_REQUIRED if n_req else WEIGHTS_RESTRAINT
    semantic = _pct(dims, weights)
    contract = _pct(cdims, CONTRACT_WEIGHTS)

    critical = any(e.severity == "critical" for e in events)
    base.update({
        "semantic": semantic,
        "contract": contract,
        "parts": {k: (round(v, 1) if isinstance(v, float) else v)
                  for k, v in dims.items()},
        "contract_parts": {k: (round(v, 1) if isinstance(v, float) else v)
                           for k, v in cdims.items()},
        "atomicity": [
            {"record_index": i, "verdict": j.verdict.value,
             "credit": j.credit, "ids": j.distinct_ids,
             "confidence": j.confidence}
            for i, j in sorted(atom.judgements.items())],
        "atomicity_review": atom.review,
        "ownership": own.as_dict(),
        "events": [e.as_dict() for e in events],
        "production_usable": bool(
            contract is not None and contract >= 80 and semantic is not None
            and (parsed.status in outcomes.DECISIONAL
                 or (not n_req and parsed.status in outcomes.RESTRAINT_VALID))
            and not parsed.truncated and not critical),
    })
    return base


def aggregate(results: list[dict]) -> dict:
    """Per-configuration summary. Truncated runs are counted, flagged, and also
    reported separately; they are never silently dropped (F-08)."""
    by: dict[str, list[dict]] = {}
    for r in results:
        by.setdefault(r.get("label") or "?", []).append(r)

    out = {}
    for label, rows in by.items():
        sem = [r["semantic"] for r in rows
               if r["semantic"] is not None and r["quality_eligible"]]
        sem_all = [r["semantic"] for r in rows if r["semantic"] is not None]
        con = [r["contract"] for r in rows if r["contract"] is not None]
        per_case: dict[str, list[float]] = {}
        for r in rows:
            if r["semantic"] is not None and r["quality_eligible"]:
                per_case.setdefault(r["case"], []).append(r["semantic"])
        case_means = [st.mean(v) for v in per_case.values()]
        draws = len(rows) / max(1, len({r["case"] for r in rows}))
        out[label] = {
            "runs": len(rows),
            "draws_per_case": round(draws, 1),
            "cases": len({r["case"] for r in rows}),
            "semantic_mean_eligible": round(st.mean(sem), 1) if sem else None,
            "semantic_mean_all_runs": round(st.mean(sem_all), 1) if sem_all else None,
            "case_mean_of_means": round(st.mean(case_means), 1) if case_means else None,
            "contract_mean": round(st.mean(con), 1) if con else None,
            "eligible_runs": len(sem),
            "truncated_runs": sum(1 for r in rows if r["truncated"]),
            "invalid_runs": sum(1 for r in rows
                                if r["status"] not in
                                {"VALID_KEEP", "VALID_HOLD", "VALID_REJECT",
                                 "VALID_NONE", "VALID_EMPTY"}),
            "critical_events": sum(1 for r in rows
                                   if any(e["severity"] == "critical"
                                          for e in r["events"])),
            "production_usable_runs": sum(1 for r in rows if r["production_usable"]),
            "single_draw_observation": draws < 2,
            "underpowered": draws < 5,
        }
    return out
