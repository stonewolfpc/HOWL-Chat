"""Ownership and scope contamination — instrument v2 (task requirement C).

This is a *measurement* category, not a MEMORY-09 architecture. Nothing here
decides how Howl should store memories. It defines the distinctions a future
experiment must be able to score:

    world fact                  true for everyone in the scenario
    scenario fact               true inside one scenario only
    character belief            what a character holds to be true; may be false
    character knowledge         what a character has learned or been told
    character opinion           evaluation, not fact
    conversation context        the session itself
    participant disclosure      a real user/participant disclosed something to
                                specific recipients

Contamination is a *critical* failure: it is recorded as an event and is never
averaged away by a high semantic score.

Honest limitation. This module cannot infer ownership from free text. It applies
a declarative expectation block supplied with the case:

    "ownership": {
      "scenario": "court",
      "characters": {"jim": ["Jim"], "lisa": ["Lisa"], "bob": ["Bob"]},
      "scenario_markers": {"court": ["court", "throne"],
                           "harbour": ["harbour", "tide", "dock"]},
      "claims": [
        {"id": "K1", "topic": "king", "topic_hints": ["king"],
         "holders": ["jim"], "scope": "character_belief",
         "polarity": "negative", "exclusive": true,
         "must_not_promote": true, "polarity_hints": ["criminal", "traitor"]}
      ]
    }

With no expectation block, ownership is not scored and the case reports
`ownership: unscored` rather than a fabricated pass.
"""

from __future__ import annotations

import re
from dataclasses import dataclass, field
from enum import Enum


class Scope(str, Enum):
    WORLD = "world"
    SCENARIO = "scenario"
    CHARACTER_BELIEF = "character_belief"
    CHARACTER_KNOWLEDGE = "character_knowledge"
    CHARACTER_OPINION = "character_opinion"
    CONVERSATION = "conversation"
    PARTICIPANT_DISCLOSURE = "participant_disclosure"
    UNMARKED = "unmarked"
    UNKNOWN = "unknown"


class Violation(str, Enum):
    PROMOTION = "PROMOTION"                       # belief emitted as world truth
    INHERITANCE = "INHERITANCE"                   # claim given to a non-recipient
    CONTRADICTION_COLLAPSE = "CONTRADICTION_COLLAPSE"  # two beliefs merged/resolved
    SCENARIO_LEAK = "SCENARIO_LEAK"               # scenario A fact in scenario B
    SCOPE_UNMARKED = "SCOPE_UNMARKED"             # scoped claim emitted with no holder


BELIEF_MARKERS = ("believes", "believed", "thinks", "thought", "is convinced",
                  "was convinced", "in his view", "in her view", "according to",
                  "claims that", "considers", "is of the opinion", "suspects",
                  "maintains that", "insists")
KNOWLEDGE_MARKERS = ("knows", "knew", "learned", "was told", "has been told",
                     "is aware", "privy", "found out", "discovered")
OPINION_MARKERS = ("likes", "dislikes", "admires", "despises", "prefers",
                   "approves of", "disapproves of", "respects", "distrusts")
DISCLOSURE_MARKERS = ("told", "confided", "disclosed", "revealed", "shared with",
                     "admitted to", "in confidence")
WORLD_MARKERS = ("world rule", "canon", "it is true", "in fact", "actually is",
                 "everyone knows", "is established", "objectively")


@dataclass
class Finding:
    kind: Violation
    record_index: int
    detail: str
    claim_ids: list[str] = field(default_factory=list)
    raw: str = ""
    critical: bool = True


@dataclass
class OwnershipReport:
    scored: bool = False
    findings: list[Finding] = field(default_factory=list)
    notes: list[str] = field(default_factory=list)

    @property
    def critical_count(self) -> int:
        return sum(1 for f in self.findings if f.critical)

    def as_dict(self) -> dict:
        return {
            "scored": self.scored,
            "critical_count": self.critical_count,
            "findings": [
                {"kind": f.kind.value, "record_index": f.record_index,
                 "detail": f.detail, "claim_ids": f.claim_ids, "raw": f.raw}
                for f in self.findings
            ],
            "notes": self.notes,
        }


def _has(text: str, markers) -> bool:
    t = " " + re.sub(r"\s+", " ", text.lower()) + " "
    return any(m in t for m in markers)


def detect_scope(text: str) -> Scope:
    t = text.lower()
    if _has(t, BELIEF_MARKERS):
        return Scope.CHARACTER_BELIEF
    if _has(t, KNOWLEDGE_MARKERS):
        return Scope.CHARACTER_KNOWLEDGE
    if _has(t, OPINION_MARKERS):
        return Scope.CHARACTER_OPINION
    if _has(t, DISCLOSURE_MARKERS):
        return Scope.PARTICIPANT_DISCLOSURE
    if _has(t, WORLD_MARKERS):
        return Scope.WORLD
    return Scope.UNMARKED


def holders_in(text: str, characters: dict[str, list[str]]) -> set[str]:
    t = text.lower()
    return {cid for cid, names in characters.items()
            if any(n.lower() in t for n in names)}


def check(records: list, expectations: dict | None) -> OwnershipReport:
    if not expectations:
        return OwnershipReport(scored=False,
                              notes=["no ownership expectation block; not scored"])

    characters = expectations.get("characters") or {}
    claims = expectations.get("claims") or []
    scenario = expectations.get("scenario")
    markers = expectations.get("scenario_markers") or {}
    rep = OwnershipReport(scored=True)

    for idx, rec in enumerate(records):
        text = rec.text
        holders = holders_in(text, characters)
        scope = detect_scope(text)
        matched_claims = []
        for c in claims:
            hints = c.get("topic_hints") or []
            if hints and all(h.lower() in text.lower() for h in hints):
                matched_claims.append(c)

        for c in matched_claims:
            pol = c.get("polarity_hints") or []
            if pol and not any(h.lower() in text.lower() for h in pol):
                continue

            allowed = set(c.get("holders") or [])
            if c.get("must_not_promote") and scope not in (
                    Scope.CHARACTER_BELIEF, Scope.CHARACTER_OPINION,
                    Scope.CHARACTER_KNOWLEDGE):
                rep.findings.append(Finding(
                    Violation.PROMOTION, idx,
                    f"claim {c['id']} is a {c.get('scope', 'scoped')} statement but "
                    f"was emitted with scope {scope.value!r}",
                    [c["id"]], rec.raw))

            if c.get("exclusive") and allowed:
                outsiders = holders - allowed
                if outsiders:
                    rep.findings.append(Finding(
                        Violation.INHERITANCE, idx,
                        f"claim {c['id']} is limited to {sorted(allowed)} but the "
                        f"record also assigns it to {sorted(outsiders)}",
                        [c["id"]], rec.raw))
                if not holders and scope in (Scope.UNMARKED, Scope.WORLD):
                    rep.findings.append(Finding(
                        Violation.SCOPE_UNMARKED, idx,
                        f"claim {c['id']} requires a holder from {sorted(allowed)}; "
                        "none is present and the statement reads as shared truth",
                        [c["id"]], rec.raw))

        # Two different characters' contradictory claims in one record.
        if len(matched_claims) >= 2 and scope in (Scope.UNMARKED, Scope.WORLD):
            topics = {c.get("topic") for c in matched_claims}
            polarities = {c.get("polarity") for c in matched_claims}
            if len(topics) == 1 and len(polarities) >= 2:
                rep.findings.append(Finding(
                    Violation.CONTRADICTION_COLLAPSE, idx,
                    "contradictory character beliefs were collapsed into one "
                    "unqualified statement",
                    [c["id"] for c in matched_claims], rec.raw))

        if scenario and markers:
            for other, other_markers in markers.items():
                if other == scenario:
                    continue
                if any(m.lower() in text.lower() for m in other_markers):
                    rep.findings.append(Finding(
                        Violation.SCENARIO_LEAK, idx,
                        f"record for scenario {scenario!r} references "
                        f"{other!r} markers",
                        [], rec.raw))

    return rep
