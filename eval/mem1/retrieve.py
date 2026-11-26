"""Context-explicit retrieval (plan §3.1).

Two contexts, and no third:

  * `character_perspective` — **default and fail-closed**. A character sees its
    own records and restricted records it is in the audience of. Nothing else.
  * `world_state` — world-scoped records, plus restricted records the viewer is
    in the audience of.

Eligibility precedes ranking, structurally: the candidate pool is built by
ownership/scope first, then scored. A record that belongs to someone else is
never a candidate, so it can neither be selected nor leak — and if one somehow
appears in the selection, that is an ownership veto, not a ranking curiosity.
"""
from __future__ import annotations

import math
import re
from dataclasses import dataclass, field

from .errors import DesignViolationError, OwnershipVetoError
from .model import ACTIVE, SCOPES, record_is_active

CONTEXTS = ("character_perspective", "world_state")
DEFAULT_CONTEXT = "character_perspective"

# Visibility table: which scopes each context may see at all. Ownership is not
# consulted here; who owns a record and where it may be seen are separate axes.
# Every scope in `SCOPES` must appear in at least one context, and every entry must
# be a declared scope: an omitted scope would silently hide a whole class of
# records, so the table is asserted total rather than trusted.
VISIBILITY = {
    "character_perspective": ("owner", "restricted", "chat_state", "scenario_state"),
    "world_state": ("world", "restricted", "chat_state", "scenario_state"),
}


def check_visibility_table() -> None:
    """The table must cover every scope and name no scope it does not declare."""
    listed: set[str] = set()
    for context, scopes in VISIBILITY.items():
        if context not in CONTEXTS:
            raise DesignViolationError(f"visibility table names unknown context {context!r}")
        for scope in scopes:
            if scope not in SCOPES:
                raise DesignViolationError(
                    f"visibility table names undeclared scope {scope!r} in {context}")
            listed.add(scope)
    missing = [s for s in SCOPES if s not in listed]
    if missing:
        raise DesignViolationError(
            f"visibility table omits declared scopes {missing}: those records would "
            "be unreachable in every context")


check_visibility_table()

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
    context: str
    candidates: list[Scored] = field(default_factory=list)
    selected: list[Scored] = field(default_factory=list)
    rejected: list[Scored] = field(default_factory=list)
    abstained: bool = False
    reason: str = ""
    branches_used: list[str] = field(default_factory=lambda: ["lexical"])

    def ids(self) -> list[str]:
        return [s.record_id for s in self.selected]

    def as_dict(self) -> dict:
        return {"query": self.query, "owner_id": self.owner_id, "context": self.context,
                "candidates": [vars(c) for c in self.candidates],
                "selected": [vars(s) for s in self.selected],
                "rejected": [vars(s) for s in self.rejected],
                "abstained": self.abstained, "reason": self.reason,
                "branches_used": self.branches_used}


def eligible(store, owner_id: str, context: str, *, active_state, requested_scenario):
    """Eligibility filter. Fail-closed: an unknown perspective sees nothing.

    Four independent questions, never one field answering several:

      * visibility — does this context allow this scope (`VISIBILITY`)?
      * access — does this requester pass the scope's own condition?
      * scenario context — is this the scenario that was asked from?
      * activation — is the owning world/scenario active (`record_is_active`)?

    State is supplied explicitly and has no default: `active_state=None` (or an
    undeclared state) gates nothing (D1), but omitting the argument is a TypeError.
    A call site that forgets the state must fail loudly rather than silently retrieve
    as if nothing were declared. Scenario records are visible only when
    `requested_scenario` matches their owning scenario; with no requested scenario
    they are never returned, so one scenario cannot see another's records even when
    both share a parent world.
    """
    if context not in CONTEXTS:
        raise ValueError(f"unknown retrieval context: {context!r}")
    if not owner_id:
        return []

    allowed = VISIBILITY[context]
    out = []
    for rec in store.current():
        if rec.status != ACTIVE:
            continue
        if rec.scope not in allowed:
            continue
        if rec.scope == "owner":
            if rec.owner_id != owner_id:
                continue
        elif rec.scope == "restricted":
            if owner_id not in rec.audience:
                continue
            if context == "character_perspective" and rec.owner_id != owner_id:
                continue
        elif rec.scope == "scenario_state":
            if not requested_scenario or rec.owner_id != requested_scenario:
                continue
        # `world` and `chat_state` carry no requester condition of their own.
        if not record_is_active(rec, active_state):
            continue
        out.append(rec)
    return out


def _score(query_terms: list[str], pool) -> list[Scored]:
    docs = {r.record_id: terms(f"{r.text} {r.value}") for r in pool}
    n = max(1, len(docs))
    df: dict[str, int] = {}
    for d in docs.values():
        for t in set(d):
            df[t] = df.get(t, 0) + 1
    idf = {t: math.log(1 + (n - c + 0.5) / (c + 0.5)) for t, c in df.items()}
    scored: list[Scored] = []
    for r in pool:
        doc = docs[r.record_id]
        if not doc:
            scored.append(Scored(r.record_id, 0.0, "lexical", "empty_record"))
            continue
        matched, total = [], 0.0
        for t in query_terms:
            tf = doc.count(t)
            if tf:
                matched.append(t)
                total += idf.get(t, 1.0) * (tf * 2.5 / (tf + 1.5 * (0.25 + 0.75 * len(doc) / 6)))
        why = f"terms={','.join(matched)}" if matched else "no_term_overlap"
        scored.append(Scored(r.record_id, round(total, 4), "lexical", why))
    scored.sort(key=lambda s: s.score, reverse=True)
    return scored


def retrieve(query: str, store, owner_id: str, context: str = DEFAULT_CONTEXT,
             max_records: int = 6, max_chars: int = 1200,
             min_score: float = 0.0, *, active_state,
             requested_scenario) -> RetrievalResult:
    """Select relevant, eligible records; abstain when nothing overlaps.

    `min_score` is a declared floor fixed before the run (plan §2.6), not a tuned
    knob. The default 0.0 means "any verified content-term overlap": ranking still
    uses BM25, but selection does not depend on its absolute magnitude, which
    collapses when the eligible pool is small (a one-record pool would otherwise
    reject a perfect match).
    """
    res = RetrievalResult(query=query, owner_id=owner_id, context=context)
    pool = eligible(store, owner_id, context, active_state=active_state,
                    requested_scenario=requested_scenario)
    if not pool:
        res.abstained, res.reason = True, "no_eligible_records"
        return res
    q = terms(query)
    if not q:
        res.abstained, res.reason = True, "empty_query"
        return res

    res.candidates = _score(q, pool)
    budget, used = [], 0
    for cand in res.candidates:
        rec = store.get(cand.record_id)
        if cand.score <= min_score:
            res.rejected.append(cand)
            continue
        if len(budget) >= max_records or used + len(rec.text) > max_chars:
            res.rejected.append(cand)
            continue
        # Veto check at the point of selection, not later.
        if rec.scope == "owner" and rec.owner_id != owner_id:
            raise OwnershipVetoError(
                f"{cand.record_id} owned by {rec.owner_id} selected for {owner_id}"
            )
        if rec.status != ACTIVE:
            raise OwnershipVetoError(f"{cand.record_id} is {rec.status} but was selected")
        budget.append(cand)
        used += len(rec.text)

    if not budget:
        res.abstained = True
        res.reason = "below_min_score" if res.candidates else "no_candidates"
        return res
    res.selected = budget
    res.reason = "selected"
    return res


def assert_selection_owned(store, owner_id: str, selected_ids: list[str]) -> None:
    """Belt and braces: any foreign record in a selection is a veto."""
    for rid in selected_ids:
        rec = store.get(rid)
        if rec.scope == "owner" and rec.owner_id != owner_id:
            raise OwnershipVetoError(f"foreign record {rid} in selection for {owner_id}")


__all__ = ["CONTEXTS", "DEFAULT_CONTEXT", "VISIBILITY", "check_visibility_table",
           "Scored", "RetrievalResult", "eligible", "retrieve", "assert_selection_owned"]
