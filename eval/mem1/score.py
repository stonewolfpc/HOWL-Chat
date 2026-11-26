"""M1 and M2 (Amend. 1 §2–§3), plus the temporal veto.

M1 — answer quality, scoreable identically in every arm. Credit depends only on
the response and the probe's key, so A0's value is a real value, not a
structural zero.

M2 — evidence-supported answers. Requires all of: `CORRECT`, a selected record
that supports the keyed answer, and the supporting span absent from the window.
Its **applicability differs by arm**, and that difference is represented as data
(`"N/A"`), never as 0:

    A1  applicable
    D   applicable but expected ~0 — non-zero means decoys leak (instrument failure)
    A0  undefined: no records exist

`composite_score` exists only to refuse: a single number that mixes M1 and M2 is
a prohibited aggregate (Amend. 1 §4).
"""
from __future__ import annotations

import re

from . import (ABSTAIN_APPROPRIATE, ABSTAIN_INAPPROPRIATE, CORRECT, INCORRECT,
               MALFORMED, M2_UNDEFINED_ARMS)
from .errors import ProhibitedAggregateError, TemporalVetoError
from .model import ACTIVE, SUPERSEDED

_DEFAULT_ABSTAIN = ("i don't know", "i do not know", "not established", "no record",
                    "cannot say", "i can't say", "no information")


def norm(text: str) -> str:
    return re.sub(r"\s+", " ", (text or "").strip().rstrip(".!?").lower()).strip()


def looks_abstained(response: str, probe) -> bool:
    lowered = (response or "").lower()
    terms = tuple(probe.abstain_terms) or _DEFAULT_ABSTAIN
    return any(t in lowered for t in terms)


def is_malformed(response: str) -> bool:
    text = (response or "").strip()
    return not text or len(text.split()) < 2


def classify(response: str, probe) -> str:
    """Exactly one M1 category, by rule, for any arm."""
    if is_malformed(response):
        return MALFORMED

    answered = False
    for term in probe.acceptable or ((probe.key_value,) if probe.key_value else ()):
        if term and norm(term) in norm(response):
            answered = True
            break

    if looks_abstained(response, probe):
        return ABSTAIN_APPROPRIATE if not probe.answerable else ABSTAIN_INAPPROPRIATE
    if probe.answerable:
        if answered:
            for bad in probe.contradiction_terms:
                if bad and norm(bad) in norm(response):
                    return INCORRECT
            return CORRECT
        return INCORRECT
    # Unanswerable probe: any asserted answer is wrong, by construction.
    return INCORRECT


def credit(category: str, probe) -> int:
    """Amend. 1 §2 scoring rule, symmetric for every arm."""
    if probe.answerable:
        return 1 if category == CORRECT else 0
    return 1 if category == ABSTAIN_APPROPRIATE else 0


def support_check(store, selected_ids: list[str], probe) -> bool:
    """Does a selected record support the keyed answer? Declared, deterministic.

    For a retraction probe the supporting evidence is an **explicit retraction event
    record** that references the retracted record. That keeps "retracted" distinct
    from "never established" (kind `absence`) and from "not known" (no record).
    """
    if probe.retraction_of:
        want = f"retraction:{probe.retraction_of}"
        for rid in selected_ids:
            rec = store.get(rid)
            if (rec.status == ACTIVE and rec.kind == "event"
                    and rec.predicate == want and rec.retracts):
                return True
        return False

    if not probe.key_predicate or not probe.key_value:
        return False
    want = norm(probe.key_value)
    for rid in selected_ids:
        rec = store.get(rid)
        if rec.predicate != probe.key_predicate or rec.status != ACTIVE:
            continue
        value = norm(rec.value)
        if value == want or want in value or value in want:
            return True
    return False


def m2_status(arm: str, probe, category: str, supported: bool, span_absent: bool) -> str:
    """Applicability first, then credit. A0 is `N/A`, never 0."""
    if arm in M2_UNDEFINED_ARMS:
        return "N/A"
    if not probe.answerable:
        return "N/A"
    if category == CORRECT and supported and span_absent:
        return "SUPPORTED"
    return "NOT_SUPPORTED"


def assert_temporal(store, probe, selected_ids: list[str]) -> None:
    """Temporal correctness veto (plan §4.3).

    * a `current` probe that selected a superseded/revoked record: stale as current;
    * a `historical` probe whose superseded value is presented by an active record:
      history collapsed into the present.
    """
    if probe.temporal not in ("current", "historical", "none"):
        raise ValueError(f"unknown temporal mode {probe.temporal!r}")
    if probe.temporal == "none":
        return
    for rid in selected_ids:
        rec = store.get(rid)
        if probe.temporal == "current" and rec.status in (SUPERSEDED, "revoked"):
            raise TemporalVetoError(
                f"{probe.probe_id}: selected {rid} is {rec.status} but the probe asks for current truth"
            )
        if probe.temporal == "historical" and rec.status == ACTIVE and probe.key_value:
            if norm(rec.value) == norm(probe.key_value):
                raise TemporalVetoError(
                    f"{probe.probe_id}: historical value presented by active record {rid}"
                )


def m1_for_conversation(probes: list, receipts: list) -> float:
    """Mean credit over the conversation's 13 probes (equal probe weight)."""
    by_id = {r.probe_id: r for r in receipts}
    missing = [p.probe_id for p in probes if p.probe_id not in by_id]
    if missing:
        raise ProhibitedAggregateError(f"cannot score a partial conversation; missing {missing}")
    total = sum(by_id[p.probe_id].m1_credit for p in probes)
    return total / len(probes)


def composite_score(*_args, **_kwargs):
    """Refused by design: no composite memory score exists in this design."""
    raise ProhibitedAggregateError(
        "composite scores are prohibited: report M1 and M2 separately (Amend. 1 §4)"
    )


__all__ = ["norm", "looks_abstained", "is_malformed", "classify", "credit", "support_check",
           "m2_status", "assert_temporal", "m1_for_conversation", "composite_score"]
