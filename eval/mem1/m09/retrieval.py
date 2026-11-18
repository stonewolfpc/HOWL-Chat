"""Retrieval contexts for the v9 engine.

Three call surfaces, and no fourth:

  * ``perspective`` — one character: its own claims, accounts and beliefs. A private
    record owned by someone else is never a candidate, so it cannot be selected
    (rule 5 / H2 / G3).
  * ``world`` — authoritative records answering world truth for one domain+target.
    Identity is preserved: one entry per proposition, a conflict returned as an
    unresolved conflict, and ``multiple`` when several identities stand unrelated
    (rules 41, J14, U-1). No record is ever selected as settled on the caller's behalf.
  * ``belief_facts`` — the same shape over ``belief_ascription`` identities only. It
    answers "what does <character> believe" from canon; it never returns a character's
    private belief record (C-2, O-5, rule 56).

Eligibility precedes matching: the candidate pool is built by ownership and class
first. A record that is not eligible is not scored, so ranking can never leak it.

Both canon views consume one snapshot taken inside ``Engine.authoritative`` and add no
read of their own, so no field of a response comes from a later state than any other
(S-1, S-2).
"""
from __future__ import annotations

import re

from .errors import RetrievalError

CONTEXTS = ("perspective", "world", "belief_facts")

_STOP = set("""a an the and or of to in on at for with is was were be been it its
i you he she they we my your his her their me him them this that these those
who what when where why how do does did have has had not no yes if then than
about into over under after before as by from""".split())


def terms(text: str) -> list[str]:
    return [w for w in re.findall(r"[a-z0-9']+", text.lower()) if w and w not in _STOP]


def _score(text: str, query: list[str]) -> int:
    doc = terms(text)
    return sum(1 for t in query if t in doc)


def perspective(engine, character: str, query: str, limit: int = 6) -> dict:
    """What one character can retrieve. Own records only; abstain when nothing matches."""
    q = terms(query)
    pool = [r for r in engine.all_records()
            if r.owner_character == character and r.status != "revoked"]
    if not pool:
        return {"status": "abstain", "reason": "no_owned_records", "records": []}
    if not q:
        return {"status": "abstain", "reason": "empty_query", "records": []}
    # Canonical tie-break by record id: score ties must not be resolved by dict
    # insertion order, or the same state could answer differently after replay (F7).
    scored = sorted(((_score(r.text, q), r) for r in pool),
                    key=lambda p: (-p[0], p[1].record_id))
    hits = [r for s, r in scored if s > 0][:limit]
    if not hits:
        return {"status": "abstain", "reason": "no_match", "records": []}
    return {"status": "selected", "reason": "match", "records": hits}


def world(engine, domain: tuple[str, str], target: str) -> dict:
    """Canon answer for one applicability, in the per-proposition envelope (REV 6 §4).

    A pass-through: the envelope, the conflict state and the refusal shapes are all
    decided inside the engine, where the gate and the snapshot live. This function adds
    no selection and no second read.
    """
    return engine.authoritative(domain, target, view="world")


def belief_facts(engine, domain: tuple[str, str], target: str) -> dict:
    """Canon answer over belief ascriptions only (REV 6 §5).

    It returns authoritative records *about* a character -- a fact about that
    character -- never the character's own belief record, which stays private,
    unlinked and unenumerable.
    """
    return engine.authoritative(domain, target, view="belief_facts")


def dispatch(engine, context: str, **kwargs) -> dict:
    """Route to a context. An unknown context returns nothing (F4).

    F4 expects an empty result, not an exception: the requirement is that an unknown
    context cannot fail open into a broader pool. Returning an explicit empty answer
    satisfies that and stays callable from a caller that cannot catch.
    """
    if context == "perspective":
        return perspective(engine, kwargs["character"], kwargs["query"])
    if context == "world":
        return world(engine, tuple(kwargs["domain"]), kwargs["target"])
    if context == "belief_facts":
        return belief_facts(engine, tuple(kwargs["domain"]), kwargs["target"])
    return {"status": "abstain", "reason": f"unknown_context:{context}", "records": []}


__all__ = ["CONTEXTS", "terms", "perspective", "world", "belief_facts", "dispatch"]
