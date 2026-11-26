"""Frozen pre-v2 eligibility logic. TEST ORACLE ONLY — never production retrieval.

`retrieve.eligible` was rewritten from intertwined context/owner branches into an
explicit visibility table. Replacing branch logic with a table is not automatically
behaviour-preserving, so the old branches are preserved here verbatim and compared
against the table over a constructed matrix. This module must not be imported by
production paths, and its body must not be edited to make the table agree with it.

The oracle describes v1 semantics: visibility keyed on `scope`, access keyed on the
requester, and no awareness of ownership kind, activation, or requested scenario.
That is exactly what makes it a useful control: for `owner`, `restricted` and
`world` scopes the table must reproduce it exactly.
"""
from __future__ import annotations

from .model import ACTIVE

CONTEXTS = ("character_perspective", "world_state")


def reference_eligible(store, owner_id: str, context: str):
    """v1 branch logic, copied unchanged. Reference implementation, not retrieval."""
    if context not in CONTEXTS:
        raise ValueError(f"unknown retrieval context: {context!r}")
    if not owner_id:
        return []

    out = []
    for rec in store.current():
        if rec.status != ACTIVE:
            continue
        if context == "character_perspective":
            if rec.owner_id == owner_id and rec.scope in ("owner", "restricted"):
                if rec.scope == "restricted" and owner_id not in rec.audience:
                    continue
                out.append(rec)
        else:  # world_state
            if rec.scope == "world":
                out.append(rec)
            elif rec.scope == "restricted" and owner_id in rec.audience:
                out.append(rec)
    return out


def reference_ids(store, owner_id: str, context: str) -> list[str]:
    return [r.record_id for r in reference_eligible(store, owner_id, context)]


__all__ = ["CONTEXTS", "reference_eligible", "reference_ids"]
