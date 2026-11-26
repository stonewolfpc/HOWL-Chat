"""Gate B — hybrid retrieval: branches, fusion, rerank, attribution.

Each branch is scored alone first and reports its own hits. Fusion is a separate
step, and the reranker is a further step. No stage can claim a contribution it
did not make: `branches` on every selected hit names the branches that actually
produced a non-zero score, and an unavailable branch is recorded as unavailable.
"""
from __future__ import annotations

from dataclasses import dataclass, field

from eval.s3.records import Record
from eval.s3.retrieve import LexicalIndex, terms
from eval.s3.vectors import SubwordIndex, SEMANTIC_STATUS

RRF_K = 60
MIN_FUSED = 0.012
MIN_SUBWORD = 0.45


@dataclass
class BranchHit:
    record_id: str
    branch: str
    score: float
    rank: int


@dataclass
class HybridResult:
    query: str
    pool_size: int
    owner_filter: str | None
    per_branch: dict[str, list[BranchHit]] = field(default_factory=dict)
    fused: list[BranchHit] = field(default_factory=list)
    reranked: list[BranchHit] = field(default_factory=list)
    selected: list[dict] = field(default_factory=list)
    abstained: bool = False
    reason: str = ""
    branches_available: list[str] = field(default_factory=list)
    branches_unavailable: dict[str, str] = field(default_factory=dict)

    def as_dict(self) -> dict:
        return {
            "query": self.query, "pool_size": self.pool_size, "owner_filter": self.owner_filter,
            "per_branch": {b: [h.__dict__ for h in hs] for b, hs in self.per_branch.items()},
            "fused": [h.__dict__ for h in self.fused],
            "reranked": [h.__dict__ for h in self.reranked],
            "selected": self.selected, "abstained": self.abstained, "reason": self.reason,
            "branches_available": self.branches_available,
            "branches_unavailable": self.branches_unavailable,
        }


def pool(store, owner_id: str | None, include_history: bool = False) -> list[Record]:
    """Candidate pool. owner_id=None deliberately returns everything so that
    pre-filter contamination can be measured instead of assumed away."""
    recs = store.all() if include_history else [r for r in store.all() if r.status == "active"]
    if owner_id is not None:
        recs = [r for r in recs if r.owner_id == owner_id]
    return recs


def _rank(scores: dict[str, float]) -> list[BranchHit]:
    ordered = sorted(scores.items(), key=lambda kv: kv[1], reverse=True)
    return [BranchHit(rid, "", round(s, 4), i + 1) for i, (rid, s) in enumerate(ordered) if s > 0]


# Calibrated, not tuned until green. Measured on the 7-record baseline with
# nomic-embed-text-v1.5: correct hits scored 0.656 and 0.750; irrelevant queries
# topped out at 0.580. 0.60 sits in the observed gap. This is CORPUS-SPECIFIC and
# must be recalibrated when the record set grows — it is not a model constant.
# Note it does NOT rescue P5: "Who trained the horse?" scores 0.722 on the
# trainer record, above any threshold that keeps P1 working.
MIN_SEMANTIC = 0.60


def run_branches(query: str, recs: list[Record], semantic=None) -> dict[str, list[BranchHit]]:
    texts = {r.record_id: f"{r.text} {r.value}" for r in recs}
    lex = LexicalIndex(recs)
    q = terms(query)
    lex_scores = {}
    for r in recs:
        s, _ = lex.score(r.record_id, q)
        if s:
            lex_scores[r.record_id] = s
    sw = SubwordIndex.build(texts)
    sw_scores = {rid: s for rid, s in sw.score(query).items() if s >= MIN_SUBWORD}
    sem_scores: dict[str, float] = {}
    if semantic is not None:
        docs = semantic.score_documents(texts)
        sem_scores = {rid: s for rid, s in semantic.score_query(query, docs).items()
                      if s >= MIN_SEMANTIC}
    out = {"lexical": _rank(lex_scores), "subword": _rank(sw_scores),
           "semantic": _rank(sem_scores)}
    for b in out:
        for h in out[b]:
            h.branch = b
    return out


def fuse(per_branch: dict[str, list[BranchHit]]) -> list[BranchHit]:
    """Reciprocal rank fusion. Only branches with hits contribute, and only those
    are named on the fused hit."""
    acc: dict[str, float] = {}
    who: dict[str, set[str]] = {}
    for branch, hits in per_branch.items():
        for h in hits:
            acc[h.record_id] = acc.get(h.record_id, 0.0) + 1.0 / (RRF_K + h.rank)
            who.setdefault(h.record_id, set()).add(branch)
    ordered = sorted(acc.items(), key=lambda kv: kv[1], reverse=True)
    return [BranchHit(rid, "+".join(sorted(who[rid])), round(s, 6), i + 1)
            for i, (rid, s) in enumerate(ordered)]


def rerank(query: str, fused: list[BranchHit], store) -> list[BranchHit]:
    """Heuristic reranker — NOT a cross-encoder. It rewards query-term coverage and
    record kinds that answer questions, and penalises stale-state history."""
    q = set(terms(query))
    out: list[BranchHit] = []
    for h in fused:
        rec = store.get(h.record_id)
        covered = len(q & set(terms(f"{rec.text} {rec.value}"))) / max(1, len(q))
        kind_prior = {"fact": 0.05, "absence": 0.05, "commitment": 0.02}.get(rec.kind, 0.0)
        score = h.score + 0.10 * covered + kind_prior - (0.02 if rec.status != "active" else 0.0)
        out.append(BranchHit(h.record_id, h.branch, round(score, 6), 0))
    out.sort(key=lambda h: h.score, reverse=True)
    for i, h in enumerate(out, start=1):
        h.rank = i
    return out


def hybrid(query: str, store, owner_id: str | None, max_records: int = 6,
           max_chars: int = 1200, include_history: bool = False, semantic=None) -> HybridResult:
    recs = pool(store, owner_id, include_history)
    res = HybridResult(query=query, pool_size=len(recs), owner_filter=owner_id)
    if semantic is None:
        res.branches_unavailable = {"semantic": SEMANTIC_STATUS["reason"]}
        res.branches_available = ["lexical", "subword"]
    else:
        res.branches_available = ["lexical", "subword", "semantic"]
    if not recs:
        res.abstained, res.reason = True, "empty_pool"
        return res

    res.per_branch = run_branches(query, recs, semantic=semantic)
    res.fused = fuse(res.per_branch)
    res.reranked = rerank(query, res.fused, store)

    if not res.reranked or res.reranked[0].score < MIN_FUSED:
        res.abstained, res.reason = True, "below_min_fused"
        return res

    used = 0
    for h in res.reranked:
        rec = store.get(h.record_id)
        if len(res.selected) >= max_records or used + len(rec.text) > max_chars:
            break
        res.selected.append({"record_id": rec.record_id, "branches": h.branch,
                             "fused_score": h.score, "value": rec.value,
                             "owner_id": rec.owner_id, "status": rec.status})
        used += len(rec.text)
    res.reason = "selected"
    return res
