"""Layer 4a — retrieval.

Identity-aware selection with a real lexical branch, a budget, and an explicit
abstention decision. Semantic retrieval and reranking are Gate B and are ABSENT
here; this module never pretends to have them. `branches` on every result says
which branch produced the hit, so a semantic branch that never fires cannot hide.
"""
from __future__ import annotations

from dataclasses import dataclass, field
import math
import re

from eval.s3.records import Record, ACTIVE

_STOP = set("""a an the and or of to in on at for with is was were be been it its
i you he she they we my your his her their me him them this that these those
who what when where why how do does did have has had not no yes if then than
about into over under after before as by from""".split())


def terms(text: str) -> list[str]:
    return [w for w in re.findall(r"[a-z0-9']+", text.lower()) if w and w not in _STOP]


@dataclass
class Scored:
    record_id: str
    score: float
    branch: str
    why: str


@dataclass
class RetrievalResult:
    query: str
    owner_id: str
    candidates: list[Scored] = field(default_factory=list)
    selected: list[Scored] = field(default_factory=list)
    rejected: list[Scored] = field(default_factory=list)
    abstained: bool = False
    reason: str = ""
    branches_used: list[str] = field(default_factory=lambda: ["lexical"])

    def as_dict(self) -> dict:
        return {"query": self.query, "owner_id": self.owner_id,
                "candidates": [c.__dict__ for c in self.candidates],
                "selected": [s.__dict__ for s in self.selected],
                "rejected": [s.__dict__ for s in self.rejected],
                "abstained": self.abstained, "reason": self.reason,
                "branches_used": self.branches_used}


class LexicalIndex:
    """BM25 over record text. Small, deterministic, auditable."""

    def __init__(self, records: list[Record]) -> None:
        self.records = records
        self.docs = {r.record_id: terms(f"{r.text} {r.value}") for r in records}
        self.avgdl = (sum(len(d) for d in self.docs.values()) / len(self.docs)) if self.docs else 0.0
        df: dict[str, int] = {}
        for d in self.docs.values():
            for t in set(d):
                df[t] = df.get(t, 0) + 1
        n = max(1, len(self.docs))
        self.idf = {t: math.log(1 + (n - c + 0.5) / (c + 0.5)) for t, c in df.items()}

    def score(self, record_id: str, q: list[str]) -> tuple[float, list[str]]:
        doc = self.docs[record_id]
        if not doc:
            return 0.0, []
        dl = len(doc)
        total, matched = 0.0, []
        for t in q:
            tf = doc.count(t)
            if not tf:
                continue
            matched.append(t)
            denom = tf + 1.5 * (1 - 0.75 + 0.75 * dl / (self.avgdl or 1))
            total += self.idf.get(t, 1.0) * (tf * 2.5 / denom)
        return total, matched


def retrieve(query: str, store, owner_id: str, max_records: int = 6,
             max_chars: int = 1200, min_score: float = 1.0) -> RetrievalResult:
    """Select current records owned by `owner_id` relevant to `query`.

    Abstention is a decision, not a failure: no overlap, or a score below
    `min_score`, returns `abstained=True` with a reason.
    """
    # Ownership filter first: a record that belongs to another character is not a
    # candidate, so it can neither be selected nor leaked.
    pool = [r for r in store.current(owner_id=owner_id)]
    res = RetrievalResult(query=query, owner_id=owner_id)
    if not pool:
        res.abstained, res.reason = True, "no_records_for_owner"
        return res

    index = LexicalIndex(pool)
    q = terms(query)
    if not q:
        res.abstained, res.reason = True, "empty_query"
        return res

    for r in pool:
        s, matched = index.score(r.record_id, q)
        why = f"terms={','.join(matched)}" if matched else "no_term_overlap"
        res.candidates.append(Scored(r.record_id, round(s, 4), "lexical", why))
    res.candidates.sort(key=lambda c: c.score, reverse=True)

    budget, used = [], 0
    for c in res.candidates:
        rec = store.get(c.record_id)
        if c.score < min_score:
            res.rejected.append(c)
            continue
        if len(budget) >= max_records or used + len(rec.text) > max_chars:
            res.rejected.append(c)
            continue
        budget.append(c)
        used += len(rec.text)

    if not budget:
        res.abstained = True
        res.reason = "below_min_score" if res.candidates else "no_candidates"
        return res
    res.selected = budget
    res.reason = "selected"
    return res
