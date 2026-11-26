"""Atomicity evaluation — instrument v2 (audit finding F-01).

v1 asked whether the required keywords appeared *somewhere on a line*. One row
carrying four unrelated rules therefore earned full recall (C01 = 100). v2 asks
a second question: is the record a legitimate memory unit?

Verdicts
--------
ATOMIC              one asserted proposition about one subject.
LEGITIMATE_COMPOUND several predicates that are jointly true of one event or
                    subject, or a declared joint group in the case reference
                    ("Halden trained her as a cartographer for nine years").
MERGED_DISTINCT     two or more *separately required* facts from different
                    requirement groups in one record, coordinated by a clause
                    boundary. Recoverable but lossy: half credit per item.
UNRELATED           a required fact padded with material from another topic.
UNCERTAIN           heuristic could not decide; no penalty, flagged for review.

No arbitrary splitting rule is imposed: nothing requires a model to make one
row per fact. The rule is only that a record that merges *distinct required
facts* is not credited as if those facts had been kept apart.

The heuristic is lexical (clause boundaries + a verb lexicon + the case's own
`compound_ok` declarations). It cannot parse language, so every MERGED_DISTINCT
carries a confidence value and is listed for review.
"""

from __future__ import annotations

import re
from dataclasses import dataclass, field
from enum import Enum


class Verdict(str, Enum):
    ATOMIC = "ATOMIC"
    LEGITIMATE_COMPOUND = "LEGITIMATE_COMPOUND"
    MERGED_DISTINCT = "MERGED_DISTINCT"
    UNRELATED = "UNRELATED"
    UNCERTAIN = "UNCERTAIN"


#: Recall multiplier applied to a required item found in a record of each kind.
CREDIT = {
    Verdict.ATOMIC: 1.0,
    Verdict.LEGITIMATE_COMPOUND: 1.0,
    Verdict.MERGED_DISTINCT: 0.5,
    Verdict.UNRELATED: 0.5,
    Verdict.UNCERTAIN: 1.0,
}

_VERBS = (
    "is was are were has have had hates hate likes like prefers prefer wants want "
    "cannot can't can could will would must should keeps keep carries carry pays pay "
    "expires expire doubles double trained train trains grew grow lives live works work "
    "says say said told tells believes believe thinks think knows know remembers "
    "avoid avoids refused refuse fears fear sets set burns burn forbids forbid "
    "isn't aren't doesn't don't never"
).split()

_CLAUSE_SPLIT = re.compile(
    r"\s*(?:;|\u2014|--|,\s*(?:and|but|while|whereas)\b|\s+(?:and|but|while|whereas)\s+(?=(?:also\b|then\b|she\b|he\b|they\b|it\b|his\b|her\b|their\b|the\b|its\b|[A-Z][a-z]+)))",
    re.IGNORECASE,
)


@dataclass
class Judgement:
    verdict: Verdict
    confidence: float
    distinct_ids: list[str] = field(default_factory=list)
    clause_count: int = 1
    predicate_count: int = 0
    evidence: str = ""

    @property
    def credit(self) -> float:
        return CREDIT[self.verdict]


@dataclass
class CaseAtomicity:
    judgements: dict[int, Judgement] = field(default_factory=dict)
    violations: list[dict] = field(default_factory=list)
    review: list[dict] = field(default_factory=list)

    def credit_for(self, record_index: int) -> float:
        j = self.judgements.get(record_index)
        return j.credit if j else 0.0

    @property
    def structure_pass_rate(self) -> float:
        if not self.judgements:
            return 0.0
        ok = sum(1 for j in self.judgements.values()
                 if j.verdict in (Verdict.ATOMIC, Verdict.LEGITIMATE_COMPOUND))
        return ok / len(self.judgements)


def clauses(text: str) -> list[str]:
    parts = [p.strip() for p in _CLAUSE_SPLIT.split(text or "") if p.strip()]
    return parts or [(text or "").strip()]


def predicate_count(text: str) -> int:
    toks = re.findall(r"[a-z']+", (text or "").lower())
    return sum(1 for t in toks if t in _VERBS)


def _compound_ok(case: dict, ids: list[str]) -> bool:
    groups = case.get("compound_ok") or []
    s = set(ids)
    return any(s and s <= set(g) for g in groups)


def judge_record(text: str, distinct_ids: list[str], case: dict) -> Judgement:
    cl = clauses(text)
    pc = predicate_count(text)
    ids = sorted(set(distinct_ids))
    noise = [m for m in (case.get("noise_markers") or [])
             if m.lower() in (text or "").lower()]

    if len(ids) >= 2:
        if _compound_ok(case, ids):
            return Judgement(Verdict.LEGITIMATE_COMPOUND, 0.9, ids, len(cl), pc,
                             f"{len(ids)} required facts declared joint for this case")
        conf = 0.9 if len(cl) >= 2 else 0.6
        return Judgement(Verdict.MERGED_DISTINCT, conf, ids, len(cl), pc,
                         f"{len(ids)} separately required facts in one record "
                         f"across {len(cl)} clause(s)")

    if ids and noise:
        return Judgement(Verdict.UNRELATED, 0.7, ids, len(cl), pc,
                         f"required fact padded with off-topic marker(s): {noise}")

    if len(cl) >= 4 and pc >= 4:
        return Judgement(Verdict.UNCERTAIN, 0.4, ids, len(cl), pc,
                         "long multi-predicate record with one matched fact; "
                         "needs human reading")
    return Judgement(Verdict.ATOMIC, 0.85 if len(cl) == 1 else 0.6, ids, len(cl), pc, "")


def judge_case(records: list, matches: dict[str, dict], case: dict) -> CaseAtomicity:
    """Judge every record that carries at least one matched required item."""
    ids_by_record: dict[int, list[str]] = {}
    for item_id, m in matches.items():
        ids_by_record.setdefault(m["record_index"], []).append(item_id)

    out = CaseAtomicity()
    for idx, ids in ids_by_record.items():
        rec = records[idx]
        j = judge_record(rec.text, ids, case)
        out.judgements[idx] = j
        if j.verdict == Verdict.MERGED_DISTINCT:
            out.violations.append({
                "record_index": idx,
                "verdict": j.verdict.value,
                "ids": j.distinct_ids,
                "confidence": j.confidence,
                "evidence": j.evidence,
                "raw": rec.raw,
            })
        elif j.verdict == Verdict.UNRELATED:
            out.violations.append({
                "record_index": idx,
                "verdict": j.verdict.value,
                "ids": j.distinct_ids,
                "confidence": j.confidence,
                "evidence": j.evidence,
                "raw": rec.raw,
            })
        elif j.verdict == Verdict.UNCERTAIN:
            out.review.append({"record_index": idx, "evidence": j.evidence,
                               "raw": rec.raw})
    return out
